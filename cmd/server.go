package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/dyne/tgcom/utils/server"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the SSH server",
	Long:  `Start the SSH server that allows remote interactions with tgcom.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

func init() {
	// Register the server command
	rootCmd.AddCommand(serverCmd)
}

func executeRemoteCommand(remotePath string) {
	parts := strings.SplitN(remotePath, "@", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid format. Usage: tgcom -w user@remote:/path/folder")
		os.Exit(1)
	}

	userHost := parts[0]
	pathParts := strings.SplitN(parts[1], ":", 2)
	if len(pathParts) != 2 {
		fmt.Println("Invalid format. Usage: tgcom -w user@remote:/path/folder")
		os.Exit(1)
	}

	host := pathParts[0]
	dir := pathParts[1]

	sshCmd := "ssh"
	sshArgs := []string{"-t", "-p", "2222", userHost + "@" + host, "tgcom", dir}

	// Start SSH command with PTY
	if err := startSSHWithPTY(sshCmd, sshArgs); err != nil {
		log.Fatalf("Error starting SSH with PTY: %v", err)
	}
}

func startSSHWithPTY(cmd string, args []string) error {
	// Create SSH command
	sshCommand := exec.Command(cmd, args...)

	// Start PTY
	ptmx, err := pty.Start(sshCommand)
	if err != nil {
		return fmt.Errorf("failed to start PTY: %w", err)
	}
	defer ptmx.Close()

	// Set terminal attributes
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to make terminal raw: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Resize PTY to current terminal size
	if err := resizePTY(ptmx); err != nil {
		return fmt.Errorf("failed to resize PTY: %w", err)
	}

	// Forward input to PTY
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	// Forward output from PTY
	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	// Handle PTY signals and resize
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		for range ch {
			if err := resizePTY(ptmx); err != nil {
				log.Printf("Error resizing PTY: %v", err)
			}
		}
	}()

	// Wait for SSH command to finish
	if err := sshCommand.Wait(); err != nil {
		return fmt.Errorf("SSH command failed: %w", err)
	}

	// Wait a bit before exiting to ensure all output is processed
	time.Sleep(100 * time.Millisecond)

	return nil
}

func resizePTY(ptmx *os.File) error {
	size, err := pty.GetsizeFull(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}
	if err := pty.Setsize(ptmx, size); err != nil {
		return fmt.Errorf("failed to set terminal size: %w", err)
	}
	return nil
}
