package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type APIHandler struct {
}

func (api *APIHandler) PostURLShortnerHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *APIHandler) GetURLShortnerHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *APIHandler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["encoded-url"]
	var redirectUrl string
	if url == "shubham" {
		redirectUrl = "https://facebook.com"
	} else {
		redirectUrl = "https://google.com"
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (api *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"message": "Ok"}
	json.NewEncoder(w).Encode(data)
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	router := mux.NewRouter()

	api := APIHandler{}

	router.HandleFunc("/api/healthcheck", api.HealthCheck).Methods("GET")
	router.HandleFunc("/api/shorturl", api.GetURLShortnerHandler).Methods("GET")
	router.HandleFunc("/api/shorturl", api.PostURLShortnerHandler).Methods("POST")
	router.HandleFunc("/{encoded-url}", api.RedirectURLHandler).Methods("GET")

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
