package agent

import (
	"log"
	"net"
	"os"

	pb "example.com/sync/api/pb"
	//"example.com/sync/dropboxsdk"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func runServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "syncagent",
	Short: "Connect to Cloud storage",
	Long: `Connect to interact with your cloud storage, upload/download files,
manage your team and more. It is easy, scriptable and works on all platforms!`,
	Example:            `syncagent dropbox`,
	SilenceUsage:       true,
	PersistentPostRunE: runServer,
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
