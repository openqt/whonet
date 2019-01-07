package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"runtime"
)

var (
	AppVersion string
	GoVersion  string
	GitVersion string
	BuildTime  string
)

// versionCmd represents the appVersion command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Version",
	Long:  `Print version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("AppVersion:     %s\n", AppVersion)
		fmt.Printf("GitVersion:     %s\n", GitVersion)
		fmt.Printf("GoCompiler:     %s\n", GoVersion)
		fmt.Printf("Build Time:     %s\n", BuildTime)
		fmt.Printf("Go Version:     %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
