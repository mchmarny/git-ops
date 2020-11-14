package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	daprd "github.com/dapr/go-sdk/service/http"
)

const (
	staticMessage = "hello SEA"
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// BuildTime will be overritten during build
	BuildTime = "not set"

	// service
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8080")

	templates = template.Must(template.ParseGlob("resource/template/*"))
)

func main() {

	// server mux
	mux := http.NewServeMux()

	// static content
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("resource/static"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	// other handlers
	mux.HandleFunc("/", appHandler)

	// create a Dapr service
	s := daprd.NewServiceWithMux(address, mux)

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"version": AppVersion,
		"message": staticMessage,
		"build":   BuildTime,
	}

	err := templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./resource/static/img/favicon.ico")
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
