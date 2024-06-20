package sshutils

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

type Model struct {
	Username  string
	Host      string
	Directory string
	Output    string
	Err       error
}

func ParseSSHURL(sshURL string) (username, host, path string, err error) {
	// Trim any leading or trailing whitespace
	sshURL = strings.TrimSpace(sshURL)

	// Split URL by '/'
	parts := strings.Split(sshURL, "/")

	// URL should have at least 4 parts (protocol, empty string, user@host, path)
	if len(parts) < 4 {
		err = fmt.Errorf("invalid SSH URL format")
		return
	}

	// Extract user@host part
	userHost := parts[2]
	usernameHost := strings.Split(userHost, "@")
	if len(usernameHost) != 2 {
		err = fmt.Errorf("invalid SSH URL format")
		return
	}
	username = usernameHost[0]
	host = usernameHost[1]

	// Construct path from remaining parts
	path = strings.Join(parts[3:], "/")

	return
}

func HandleSSH(sshURL string) error {
	username, host, path, err := ParseSSHURL(sshURL)
	if err != nil {
		return fmt.Errorf("invalid SSH URL: %w", err)
	}

	p := tea.NewProgram(Model{
		Username:  username,
		Host:      host,
		Directory: path,
	})

	if _, cmdErr := p.Run(); cmdErr != nil {
		return cmdErr
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return sshCmd(m.Username, m.Host, m.Directory)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sshOutputMsg:
		return Model{Username: m.Username, Host: m.Host, Directory: m.Directory, Output: msg.Output, Err: msg.Err}, nil
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %s\n", m.Err.Error())
	}
	return fmt.Sprintf("%s\nPress 'q' to quit...\n", m.Output)
}

type sshOutputMsg struct {
	Output string
	Err    error
}

func sshCmd(username, host, directory string) tea.Cmd {
	return func() tea.Msg {
		// Validate directory path
		if err := validateDirectory(directory); err != nil {
			return sshOutputMsg{Err: err}
		}

		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
					key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519")
					if err != nil {
						return nil, err
					}
					signer, err := ssh.ParsePrivateKey(key)
					if err != nil {
						return nil, err
					}
					return []ssh.Signer{signer}, nil
				}),
			},
			HostKeyCallback: verifyHostKey,
		}

		client, err := ssh.Dial("tcp", host+":22", config)
		if err != nil {
			return sshOutputMsg{Err: err}
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			return sshOutputMsg{Err: err}
		}
		defer session.Close()

		var stdout, stderr bytes.Buffer
		session.Stdout = &stdout
		session.Stderr = &stderr

		cmd := fmt.Sprintf("ls %s", directory)
		err = session.Run(cmd)
		if err != nil {
			return sshOutputMsg{Err: err}
		}

		return sshOutputMsg{Output: stdout.String()}
	}
}

func validateDirectory(directory string) error {
	// Trim any leading or trailing whitespace
	directory = strings.TrimSpace(directory)

	if directory == "" {
		return fmt.Errorf("directory path is empty")
	}

	// Check if directory path contains any forbidden characters
	forbiddenChars := "<>|&;$(){}[]`\\\"'"
	for _, char := range forbiddenChars {
		if strings.Contains(directory, string(char)) {
			return fmt.Errorf("directory path contains forbidden character: %s", string(char))
		}
	}
	return nil
}

func verifyHostKey(hostname string, remote net.Addr, receivedKey ssh.PublicKey) error {
	host, _, err := net.SplitHostPort(hostname)
	if err == nil {
		hostname = host
	}

	hostname = strings.ToLower(hostname)

	knownHosts, err := fetchKnownHosts()
	if err != nil {
		return err
	}

	for _, knownKeys := range knownHosts[hostname] {
		if bytes.Equal(knownKeys.Marshal(), receivedKey.Marshal()) {
			return nil
		}
	}

	return fmt.Errorf("unrecognized host key for %s", hostname)
}

func fetchKnownHosts() (map[string][]ssh.PublicKey, error) {
	file, err := os.Open(os.Getenv("HOME") + "/.ssh/known_hosts")
	if err != nil {
		return nil, fmt.Errorf("failed to open known_hosts file: %v", err)
	}
	defer file.Close()

	knownHosts := make(map[string][]ssh.PublicKey)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		hostname := fields[0]
		publicKeyBytes := strings.Join(fields[1:], " ")

		publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKeyBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key for %s: %v", hostname, err)
		}

		knownHosts[hostname] = append(knownHosts[hostname], publicKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning known_hosts file: %v", err)
	}

	return knownHosts, nil
}
