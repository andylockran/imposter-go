package httpserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/imposter-project/imposter-go/internal/adapter"
	"github.com/imposter-project/imposter-go/internal/config"
	"github.com/imposter-project/imposter-go/internal/handler"
)

// httpServer represents the HTTP server configuration.
type httpServer struct {
	Addr      string
	ConfigDir string
	Configs   []config.Config
}

// StartServer initialises and starts the HTTP server.
func StartServer() {
	var configDirArg string
	if len(os.Args) >= 2 {
		configDirArg = os.Args[1]
	}

	imposterConfig, configDir, configs := adapter.InitialiseImposter(configDirArg)

	// Initialise and start the server with multiple configs
	srv := newServer(imposterConfig, configDir, configs)
	srv.start(imposterConfig)
}

// newServer creates a new instance of httpServer.
func newServer(imposterConfig *config.ImposterConfig, configDir string, configs []config.Config) *httpServer {
	return &httpServer{
		Addr:      ":" + imposterConfig.ServerPort,
		ConfigDir: configDir,
		Configs:   configs,
	}
}

// start begins listening for HTTP requests and handles them.
func (s *httpServer) start(imposterConfig *config.ImposterConfig) {
	fmt.Printf("Server is listening on %s...\n", s.Addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleRequest(w, r, s.ConfigDir, s.Configs, imposterConfig)
	})

	http.HandleFunc("/system/store/", handler.HandleStoreRequest)

	if err := http.ListenAndServe(s.Addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
