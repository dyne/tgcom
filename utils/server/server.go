package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/dyne/tgcom/utils/tui"
	"github.com/dyne/tgcom/utils/tui/modelutils"
)

const (
	envHostKey = "_TGCOM_HOSTKEY"
)

var (
	pathTgcom   = filepath.Join(xdg.DataHome, "tgcom")
	pathHostKey = filepath.Join(pathTgcom, "hostkey")
	teaOptions  = []tea.ProgramOption{tea.WithAltScreen(), tea.WithOutput(os.Stderr)}
	dir         string
)

func StartServer(port string) {
	withHostKey := wish.WithHostKeyPath(pathHostKey)
	if pem, ok := os.LookupEnv(envHostKey); ok {
		withHostKey = wish.WithHostKeyPEM([]byte(pem))
	}
	srv, err := wish.NewServer(
		wish.WithAddress(":"+port),
		wish.WithMiddleware(
			bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, _ := s.Pty()
				// Initialize the file selector model with the directory argument
				model := tui.Model{
					State:         "FileSelection",
					FilesSelector: modelutils.InitialModel(dir, pty.Window.Height-5), // Initialize the FilesSelector model with window height
				}
				if model.Error != nil {
					wish.Println(s, model.Error.Error())
					return nil, nil
				}
				return model, teaOptions
			}),
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					command := s.Command()
					if len(command) < 2 {
						wish.Println(s, "Usage tgcom <directory>")
						next(s)
						return
					}
					dir = command[1]
					next(s)
				}
			},
			activeterm.Middleware(),
			lm.Middleware(),
		),

		withHostKey,
	)
	if err != nil {
		log.Fatalf("could not create server: %s", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("starting server: %s", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server returned an error: %s", err)
		}
	}()

	<-done
	log.Println("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %s", err)
	}
}
