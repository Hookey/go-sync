package cli

import (
	"context"
	"errors"
	"path"

	pb "github.com/Hookey/go-sync/api/pb"
	"github.com/spf13/cobra"
)

func put(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 || len(args) > 2 {
		return errors.New("`put` requires `src` and/or `dst` arguments")
	}

	/*chunkSize, err := cmd.Flags().GetInt64("chunksize")
	if err != nil {
		return err
	}
	if chunkSize%(1<<22) != 0 {
		return errors.New("`put` requires chunk size to be multiple of 4MiB")
	}
	workers, err := cmd.Flags().GetInt("workers")
	if err != nil {
		return err
	}
	if workers < 1 {
		workers = 1
	}
	debug, _ := cmd.Flags().GetBool("debug")*/

	src := args[0]

	// Default `dst` to the base segment of the source path; use the second argument if provided.
	dst := "/" + path.Base(src)
	if len(args) == 2 {
		dst, err = validatePath(args[1])
		if err != nil {
			return
		}
	}

	_, err = c.Put(context.Background(), &pb.PutRequest{Src: src, Dst: dst})

	return
}

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put [flags] <source> [<target>]",
	Short: "Upload a single file",
	Long: `Upload a single file
	- If target is not provided puts the file in the root of your cloud directory.
	- If target is provided it must be the desired filename in the cloud (and not a directory).
	`,

	RunE: put,
}

func init() {
	RootCmd.AddCommand(putCmd)
	putCmd.Flags().IntP("workers", "w", 4, "Number of concurrent upload workers to use")
	putCmd.Flags().Int64P("chunksize", "c", 1<<24, "Chunk size to use (should be multiple of 4MiB)")
	putCmd.Flags().BoolP("debug", "d", false, "Print debug timing")
}
