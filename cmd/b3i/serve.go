package b3i

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/b3i/pkg/server"
)

var (
	uiPort int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web UI server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting web UI on port %d...\n", uiPort)
		s := server.NewServer(uiPort, deviceIP, password, insecure)
		return s.Start()
	},
}

func init() {
	serveCmd.Flags().IntVarP(&uiPort, "ui-port", "P", 8080, "Port for the web UI")
}
