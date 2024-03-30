package webui

import (
	"net/http"

	"github.com/rakyll/statik/fs"
)

func NewFs() http.FileSystem {

	binFs, err := fs.NewWithNamespace("webui")
	if err == nil && binFs != nil {
		return binFs
	}

	return nil
}
