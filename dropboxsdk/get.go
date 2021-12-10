package dropboxsdk

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/ioprogress"
)

func (d *Dropbox) Get(src, dst string) (err error) {
	// If `dst` is a directory, append the source filename.
	if f, err := os.Stat(dst); err == nil && f.IsDir() {
		dst = path.Join(dst, path.Base(src))
	}

	arg := files.NewDownloadArg(src)

	dbx := files.New(d.Config)
	res, contents, err := dbx.Download(arg)
	if err != nil {
		return
	}
	defer contents.Close()

	f, err := os.Create(dst)
	if err != nil {
		return
	}
	defer f.Close()

	progressbar := &ioprogress.Reader{
		Reader: contents,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
			return fmt.Sprintf("Downloading %s/%s",
				humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
		}),
		Size: int64(res.Size),
	}

	_, err = io.Copy(f, progressbar)
	return
}
