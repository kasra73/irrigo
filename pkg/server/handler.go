package server

import (
	"io"
	"irrigo/pkg/config"
	"net/http"
)

type StartResponse = func(status int, header http.Header)

type Handler interface {
	Handle(req *http.Request, start StartResponse) (io.Reader, error)
}

type Filter interface {
	Handler
}

type FilterFactory interface {
	NewFilter(
		name string,
		globalConfig config.Global,
		mergedConfig config.GeneralSettings,
		localConfig config.Filter,
		next Handler) Filter
}
