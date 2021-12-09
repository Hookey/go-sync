package agent

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "syncagent",
	Short: "Connect to Cloud storage",
	Long: `Connect to interact with your cloud storage, upload/download files,
manage your team and more. It is easy, scriptable and works on all platforms!`,
	Example:      `syncagent dropbox`,
	SilenceUsage: true,
	//PersistentPostRunE: runServer,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().String("port", ":50051", "server port")
}
