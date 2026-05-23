package b3i

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/b3i/pkg/device"
)

var (
	deviceIP   string
	password   string
	insecure   bool
)

var rootCmd = &cobra.Command{
	Use:   "b3i",
	Short: "B3i is a CLI tool for managing BB10 and PlayBook applications",
	Long:  `B3i (BlackBerry Bar installer) allows you to install, uninstall, and manage BAR files on your BB10 or PlayBook device.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&deviceIP, "device", "d", "", "Device IP address")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Device password")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", true, "Allow insecure TLS connections (default true for dev devices)")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(launchCmd)
	rootCmd.AddCommand(terminateCmd)
	rootCmd.AddCommand(serveCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		apps, err := client.ListApps()
		if err != nil {
			return err
		}
		for _, app := range apps {
			fmt.Printf("%s\t%s\t%s\t%s\n", app.ID, app.Name, app.Version, app.Status)
		}
		return nil
	},
}

var installCmd = &cobra.Command{
	Use:   "install [bar-file]",
	Short: "Install a BAR file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		fmt.Printf("Installing %s...\n", args[0])
		return client.InstallApp(args[0])
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [app-id]",
	Short: "Uninstall an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		return client.UninstallApp(args[0])
	},
}

var launchCmd = &cobra.Command{
	Use:   "launch [app-id]",
	Short: "Launch an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		return client.ManageApp(args[0], "launch")
	},
}

var terminateCmd = &cobra.Command{
	Use:   "terminate [app-id]",
	Short: "Terminate an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		return client.ManageApp(args[0], "terminate")
	},
}

func getClient() (*device.Client, error) {
	if deviceIP == "" {
		deviceIP = os.Getenv("B3I_DEVICE")
	}
	if password == "" {
		password = os.Getenv("B3I_PASSWORD")
	}
	if deviceIP == "" {
		return nil, fmt.Errorf("device IP is required (use --device or B3I_DEVICE env var)")
	}

	client := device.NewClient(deviceIP, password, insecure)
	if password != "" {
		if err := client.Login(); err != nil {
			return nil, fmt.Errorf("failed to login: %w", err)
		}
	}
	return client, nil
}
