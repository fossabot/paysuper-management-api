package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Version, get from ldflags
	appVersion = "unknown"
	// Cmd version command
	Cmd = &cobra.Command{
		Use:           "version",
		Short:         "Application version",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(appVersion)
		},
	}
)

// Version returns version of application
func Version() string {
	return appVersion
}
