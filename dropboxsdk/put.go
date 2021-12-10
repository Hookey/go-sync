package dropboxsdk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/ioprogress"
)

const singleShotUploadSizeCutoff int64 = 32 * (1 << 20)

type uploadChunk struct {
	data   []byte
	offset uint64
	close  bool
}

func uploadOneChunk(dbx files.Client, args *files.UploadSessionAppendArg, data []byte) error {
	for {
		err := dbx.UploadSessionAppendV2(args, bytes.NewReader(data))
		if err != nil {
			return err
		}

		rl, ok := err.(auth.RateLimitAPIError)
		if !ok {
			return err
		}

		time.Sleep(time.Second * time.Duration(rl.RateLimitError.RetryAfter))
	}
}

func uploadChunked(dbx files.Client, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64, workers int, chunkSize int64, debug bool) (err error) {
	t0 := time.Now()
	startArgs := files.NewUploadSessionStartArg()
	startArgs.SessionType = &files.UploadSessionType{}
	startArgs.SessionType.Tag = files.UploadSessionTypeConcurrent
	res, err := dbx.UploadSessionStart(startArgs, nil)
	if err != nil {
		return
	}
	if debug {
		log.Printf("Start took: %v\n", time.Since(t0))
	}

	t1 := time.Now()
	wg := sync.WaitGroup{}
	workCh := make(chan uploadChunk, workers)
	errCh := make(chan error, 1)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range workCh {
				cursor := files.NewUploadSessionCursor(res.SessionId, chunk.offset)
				args := files.NewUploadSessionAppendArg(cursor)
				args.Close = chunk.close

				t0 := time.Now()
				if err := uploadOneChunk(dbx, args, chunk.data); err != nil {
					errCh <- err
				}
				if debug {
					log.Printf("Chunk upload at offset %d took: %v\n", chunk.offset, time.Since(t0))
				}
			}
		}()
	}

	written := int64(0)
	for written < sizeTotal {
		data, err := ioutil.ReadAll(&io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return err
		}
		expectedLen := chunkSize
		if written+chunkSize > sizeTotal {
			expectedLen = sizeTotal - written
		}
		if len(data) != int(expectedLen) {
			return fmt.Errorf("failed to read %d bytes from source", expectedLen)
		}

		chunk := uploadChunk{
			data:   data,
			offset: uint64(written),
			close:  written+chunkSize >= sizeTotal,
		}

		select {
		case workCh <- chunk:
		case err := <-errCh:
			return err
		}

		written += int64(len(data))
	}

	close(workCh)
	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
	}
	if debug {
		log.Printf("Full upload took: %v\n", time.Since(t1))
	}

	t2 := time.Now()
	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)
	_, err = dbx.UploadSessionFinish(args, nil)
	if debug {
		log.Printf("Finish took: %v\n", time.Since(t2))
	}
	return
}

func (d *Dropbox) Put(src, dst string) (err error) {
	// TODO: workerpool?
	chunkSize := int64(1 << 22)
	workers := 2
	debug := true

	contents, err := os.Open(src)
	if err != nil {
		return
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return
	}

	progressbar := &ioprogress.Reader{
		Reader: contents,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
			return fmt.Sprintf("Uploading %s/%s",
				humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
		}),
		Size: contentsInfo.Size(),
	}

	commitInfo := files.NewCommitInfo(dst)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	ts := time.Now().UTC().Round(time.Second)
	commitInfo.ClientModified = &ts

	dbx := files.New(d.Config)
	if contentsInfo.Size() > singleShotUploadSizeCutoff {
		return uploadChunked(dbx, progressbar, commitInfo, contentsInfo.Size(), workers, chunkSize, debug)
	}

	if _, err = dbx.Upload(commitInfo, progressbar); err != nil {
		return
	}

	return
}
