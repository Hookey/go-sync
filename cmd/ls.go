package cmd

import (
	"context"
	"log"
	"time"

	pb "example.com/sync/api/pb"
	"github.com/spf13/cobra"
)

func ls(cmd *cobra.Command, args []string) (err error) {
	path := ""
	if len(args) > 0 {
		if path, err = validatePath(args[0]); err != nil {
			return err
		}
	}

	/*arg := files.NewListFolderArg(path)
	arg.Recursive, _ = cmd.Flags().GetBool("recurse")
	arg.IncludeDeleted, _ = cmd.Flags().GetBool("include-deleted")
	onlyDeleted, _ := cmd.Flags().GetBool("only-deleted")
	arg.IncludeDeleted = arg.IncludeDeleted || onlyDeleted
	long, _ := cmd.Flags().GetBool("long")*/

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r, err := c.Ls(ctx, &pb.LsRequest{Path: path}); err != nil {
		log.Fatalf("could not ls: %v", err)
	} else {
		log.Printf("ls: %s", r.GetLsmessage())
	}

	return err
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [flags] [<path>]",
	Short: "List files and folders",
	Example: `  sync ls / # Or just 'ls'
  sync ls /some-folder # Or 'ls some-folder'
  sync ls /some-folder/some-file.pdf
  sync ls -l`,
	RunE: ls,
}

func init() {
	RootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("long", "l", false, "Long listing")
	lsCmd.Flags().BoolP("recurse", "R", false, "Recursively list all subfolders")
	lsCmd.Flags().BoolP("include-deleted", "d", false, "Include deleted files")
	lsCmd.Flags().BoolP("only-deleted", "D", false, "Only show deleted files")
}
