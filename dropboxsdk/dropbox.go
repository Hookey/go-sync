package dropboxsdk

import (
	"fmt"
	"strings"

	"github.com/Hookey/go-sync/core"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
)

var _ core.Storage = (*Dropbox)(nil)

type Dropbox struct {
	Config dropbox.Config
	//projectID  string
	//accessID   string
	//privateKey []byte
	//client     *storage.Client
}

func NewEngine(cfg dropbox.Config) *Dropbox {
	return &Dropbox{Config: cfg}
}

func validatePath(p string) (path string, err error) {
	path = p

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	path = strings.TrimSuffix(path, "/")

	return
}
