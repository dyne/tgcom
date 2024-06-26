package server

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

func listFilesInDir(directory string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	var fileList strings.Builder
	for _, file := range files {
		fileList.WriteString(file.Name() + "\n")
	}

	return fileList.String(), nil
}

func StartServer() {
	srv, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					command := sess.Command()
					if len(command) < 2 {
						wish.Println(sess, "Usage tgcom <directory>")
						next(sess)
					}
					dir := command[1]
					files, err := listFilesInDir(dir)
					if err != nil {
						wish.Printf(sess, "Error listing files: %v\n", err)
						next(sess)
					}
					wish.Println(sess, files)
					next(sess)
				}
			},
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	log.Info("Starting SSH server", "port", "2222")
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
	}
}
