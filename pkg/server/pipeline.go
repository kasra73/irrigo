package server

import (
	"fmt"
	"io"
	"irrigo/pkg/config"
	"irrigo/pkg/utils"
	"net/http"
	"time"
)

type Pipeline struct {
	first  Handler       // First filter in the pipeline
	config config.Global // Global config
}

func CreatePipeline(globalConfig config.Global, filterFactory FilterFactory) (*Pipeline, error) {
	var prev Handler
	if filterFactory == nil {
		return nil, fmt.Errorf("filter factory is required")
	}
	for i := len(globalConfig.Pipeline) - 1; i >= 0; i-- {
		name := globalConfig.Pipeline[i]
		// Get the local config for this filter
		localConfig, ok := globalConfig.Filters[name]
		if !ok {
			return nil, fmt.Errorf("filter %s not found in local config", name)
		}
		mergedConfig := utils.MergeConfigs(
			globalConfig.Default,
			localConfig.Overrides)

		if localConfig.Type == "plugin" {
			pluginFilter, err := loadPluginFilter(localConfig.PluginPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load plugin filter: %w", err)
			}
			prev = pluginFilter.NewFilter(name, globalConfig, mergedConfig, localConfig, prev)
		} else {
			prev = filterFactory.NewFilter(
				name,
				globalConfig,
				mergedConfig,
				localConfig,
				prev)
		}
	}
	return &Pipeline{first: prev, config: globalConfig}, nil
}

func (p *Pipeline) Start() error {
	sc := p.config.Default.Server
	// Create a new HTTP server with the configuration
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", sc.BindIP, sc.BindPort),
		ReadTimeout:  time.Duration(sc.BindTimeout) * time.Second,
		WriteTimeout: time.Duration(sc.BindTimeout) * time.Second,
		IdleTimeout:  time.Duration(sc.KeepIdle) * time.Second,
	}
	// Start the web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK
		start := func(status int, header http.Header) {
			for k, values := range header {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
			statusCode = status
		}
		appResponse, err := p.first.Handle(r, start)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if appResponse != nil {
			w.WriteHeader(statusCode)
			io.Copy(w, appResponse)
			return
		}
		if statusCode == 200 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(statusCode)
		}
		w.Write([]byte(""))
	})
	return server.ListenAndServe()
}
