package cli

import (
	"context"
	"errors"
	"path"

	pb "github.com/Hookey/go-sync/api/pb"
	"github.com/spf13/cobra"
)

func get(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 || len(args) > 2 {
		return errors.New("`get` requires `src` and/or `dst` arguments")
	}

	src, err := validatePath(args[0])
	if err != nil {
		return
	}

	// Default `dst` to the base segment of the source path; use the second argument if provided.
	dst := path.Base(src)
	if len(args) == 2 {
		dst = args[1]
	}

	_, err = c.Get(context.Background(), &pb.GetRequest{Src: src, Dst: dst})

	return
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [flags] <source> [<target>]",
	Short: "Download a file",
	RunE:  get,
}

func init() {
	RootCmd.AddCommand(getCmd)
}
