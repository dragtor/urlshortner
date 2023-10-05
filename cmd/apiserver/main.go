package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dragtor/urlshortner/services/shortner"
	httputils "github.com/dragtor/urlshortner/utils"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	ShortnerService *shortner.ShortnerService
}

type RequestData struct {
	Url string `json:"url"`
}
type ResponseData struct {
	ShortUrl string `json:"shortUrl"`
}

type DomainMetrics struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type MetricsEndpointResponseData struct {
	Data []*DomainMetrics `json:"data"`
}

func (api *APIHandler) PostURLShortnerHandler(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		httputils.HTTPResponseData(w, false, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	urlmeta, err := api.ShortnerService.CreateShortUrl(requestData.Url)
	if err != nil {
		httputils.HTTPResponseData(w, false, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := ResponseData{
		ShortUrl: urlmeta.GetShortUrl(),
	}
	httputils.HTTPResponseData(w, true, resp, "", http.StatusOK)
}

func (api *APIHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["encoded-url"]

	urlmeta, err := api.ShortnerService.GetSourceUrlForShortUrl(shortUrl)
	if err != nil {
		httputils.HTTPResponseData(w, false, nil, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, urlmeta.GetSourceUrl(), http.StatusSeeOther)
}

func (api *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputils.HTTPResponseData(w, true, "ok", "", http.StatusOK)
}

const (
	DEFAULT_HEADCOUNT = 3
)

func (api *APIHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	headCount := queryParams.Get("headcount")
	headcnt, err := strconv.Atoi(headCount)
	if err != nil {
		headcnt = DEFAULT_HEADCOUNT
	}
	metrics, err := api.ShortnerService.GetMetrics(headcnt)
	if err != nil {
		httputils.HTTPResponseData(w, true, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp MetricsEndpointResponseData
	resp.Data = make([]*DomainMetrics, 0)
	for _, m := range metrics {
		resp.Data = append(resp.Data, &DomainMetrics{Domain: m.GetDomain(), Count: m.GetCount()})
	}
	httputils.HTTPResponseData(w, true, resp, "", http.StatusOK)
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	router := mux.NewRouter()

	ss := &shortner.ShortnerService{}
	shorterService, err := shortner.NewShortnerService(shortner.WithMemoryUrlRepository(ss))
	if err != nil {
		panic("Failed to setup service")
	}
	api := APIHandler{ShortnerService: shorterService}

	router.HandleFunc("/urlshortnerservice/v1/healthcheck", api.HealthCheck).Methods("GET")
	router.HandleFunc("/urlshortnerservice/v1/url", api.PostURLShortnerHandler).Methods("POST")
	router.HandleFunc("/{encoded-url}", api.RedirectURLHandler).Methods("GET")
	router.HandleFunc("/urlshortnerservice/v1/metrics", api.GetMetrics).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, shutdownCancel := context.WithTimeout(context.Background(), wait)
	defer shutdownCancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Error shutting down server:", err)
	}
	log.Println("Server shutdown complete.")
}
