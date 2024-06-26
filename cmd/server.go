package cmd

import (
	"github.com/dyne/tgcom/utils/server"
	"github.com/spf13/cobra"
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
