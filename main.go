package main

import (
	// "fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type camera struct {
	Name     string    `json:"name"`
	IP       string    `json:"ip"`
	DateTime time.Time `json:"time"`
}

func status(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	enc.Encode(cameras)
	// for k, v := range cameras {

	// 	fmt.Fprintf(w, "camera %s: last ip: %s on %s\n", k, v.IP, v.DateTime)
	// }
}

func handler(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera_id := vars["camera"]
	match, _ := regexp.MatchString("g\\d+", camera_id)

	if match {
		ip := strings.Split(r.RemoteAddr, ":")[0]

		c := camera{IP: ip, DateTime: time.Now()}

		cameras[camera_id] = c
	}
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
