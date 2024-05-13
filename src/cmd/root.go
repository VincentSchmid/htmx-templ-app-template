package cmd

import (
	"fmt"
	"os"

	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

var configFile string
var Bun *bun.DB

var rootCmd = &cobra.Command{
	Use:   "htmx-templ-app-template",
	Short: "htmx-templ-app-template is an application to manage a service provider - client realitionship directly",
	Long:  `htmx-templ-app-template is an application to manage a service provider - client realitionship directly`,
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := appconfig.Init(configFile); err != nil {
			return fmt.Errorf("failed to initialize appconfig: %w", err)
		}

		err := logger.Init(appconfig.Config.Debug)
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	cobra.EnableTraverseRunHooks = true
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "Path to the configuration .pkl file.(required)")
	if err := rootCmd.MarkPersistentFlagRequired("config-file"); err != nil {
		fmt.Fprintf(os.Stderr, "error was an error while marking the flag as required '%s'", err)
		os.Exit(1)
	}
}
