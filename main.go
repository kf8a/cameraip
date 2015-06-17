package main

import (
	// "fmt"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type camera struct {
	IP       string
	DateTime time.Time
}

func status(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	for k, v := range cameras {
		fmt.Fprintf(w, "camera %s: last ip: %s on %s\n", k, v.IP, v.DateTime)
	}
}

func handler(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera_id := vars["camera"]

	ip := strings.Split(r.RemoteAddr, ":")[0]

	c := camera{IP: ip, DateTime: time.Now()}

	cameras[camera_id] = c
}

func main() {
	cameras := make(map[string]camera)

	r := mux.NewRouter()

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		status(cameras, w, r)
	})
	r.HandleFunc("/{camera}", func(w http.ResponseWriter, r *http.Request) {
		handler(cameras, w, r)
	})

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
