package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/Hookey/go-sync/api/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var c pb.APIClient

func validatePath(p string) (path string, err error) {
	path = p

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	path = strings.TrimSuffix(path, "/")

	return
}

func initClient(cmd *cobra.Command, args []string) (err error) {
	address, _ := cmd.Flags().GetString("addr")
	conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = pb.NewAPIClient(conn)
	return err
}

func closeClient(cmd *cobra.Command, args []string) error {
	return conn.Close()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sync",
	Short: "A command line tool for users and team admins",
	Long: `Use sync to quickly interact with your cloud, upload/download files,
manage your team and more. It is easy, scriptable and works on all platforms!`,
	SilenceUsage:       true,
	PersistentPreRunE:  initClient,
	PersistentPostRunE: closeClient,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().String("addr", "localhost:50051", "API address")
}
