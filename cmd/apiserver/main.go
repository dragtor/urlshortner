package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dragtor/urlshortner/services/shortner"
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
		http.Error(w, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	urlmeta, err := api.ShortnerService.CreateShortUrl(requestData.Url)
	if err != nil {
		http.Error(w, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	response := ResponseData{
		ShortUrl: urlmeta.GetShortUrl(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Failed to write JSON response:", err)
	}

}

func (api *APIHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["encoded-url"]
	urlmeta, err := api.ShortnerService.GetSourceUrlForShortUrl(shortUrl)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, urlmeta.GetSourceUrl(), http.StatusSeeOther)
}

func (api *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"message": "Ok"}
	json.NewEncoder(w).Encode(data)
}

func (api *APIHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := api.ShortnerService.GetMetrics(3)
	if err != nil {
		http.Error(w, "Failed to fetch metrics", http.StatusInternalServerError)
	}
	var resp MetricsEndpointResponseData
	resp.Data = make([]*DomainMetrics, 0)
	for _, m := range metrics {
		resp.Data = append(resp.Data, &DomainMetrics{Domain: m.GetDomain(), Count: m.GetCount()})
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Failed to write JSON response:", err)
	}

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
