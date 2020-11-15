package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"gopkg.in/olahol/melody.v1"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

const (
	// demo message
	greetingsMessage = "hello SEA"

	// component args
	queryURL      = "https://swapi.dev/api/people/%d/"
	minQueryParam = 0
	maxQueryParam = 4 //83
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// BuildTime will be overritten during build
	BuildTime = "not set"

	// service
	logger   = log.New(os.Stdout, "", 0)
	address  = getEnvVar("ADDRESS", ":8080")
	schedule = getEnvVar("SCHEDULE", "demo-cron")

	broadcaster = melody.New()
	templates   = template.Must(template.ParseGlob("resource/template/*"))
	queryCache  = make(map[int][]byte)
)

func main() {

	// server mux
	mux := http.NewServeMux()

	// static content
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("resource/static"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	// other handlers
	mux.HandleFunc("/", appHandler)
	mux.HandleFunc("/ws", wsHandler)

	// websocket upgrade
	broadcaster.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// create a Dapr service
	s := daprd.NewServiceWithMux(address, mux)

	// add some input binding handler
	if err := s.AddBindingInvocationHandler(schedule, scheduleHandler); err != nil {
		logger.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

func scheduleHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	logger.Printf("Schedule - Metadata:%v, Data:%v", in.Metadata, in.Data)

	// number in range of the expected starwars characters
	rndChar := rand.Intn(maxQueryParam-minQueryParam) + minQueryParam

	// if in local cache, use it
	if _, ok := queryCache[rndChar]; ok {
		logger.Printf("cache hit: %d", rndChar)
		body := queryCache[rndChar]
		if err := broadcaster.Broadcast(body); err != nil {
			logger.Println(err)
			return nil, errors.Wrapf(err, "error broadcasting: %s", string(body))
		}
		return nil, nil
	}

	// not in cache, query
	url := fmt.Sprintf(queryURL, rndChar)
	logger.Printf("query URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Println(err)
		return nil, errors.Wrapf(err, "error quering: %s", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return nil, errors.Wrapf(err, "error reading body from: %s", url)
	}
	queryCache[rndChar] = body

	// broadcast to UI
	if err := broadcaster.Broadcast(body); err != nil {
		logger.Println(err)
		return nil, errors.Wrapf(err, "error broadcasting: %s", string(body))
	}

	return nil, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	broadcaster.HandleRequest(w, r)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	proto := r.Header.Get("x-forwarded-proto")
	if proto == "" {
		proto = "http"
	}

	data := map[string]string{
		"host":    r.Host,
		"proto":   proto,
		"version": AppVersion,
		"message": greetingsMessage,
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
