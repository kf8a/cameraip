package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"strings"
	"time"
  "net"
)

var (
	ping_counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "camera_pings",
		Help: "Number of times the stats server was pinged.",
	})
)

type camera struct {
	IP       string    `json:"ip"`
	DateTime time.Time `json:"time"`
}

func status(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	enc.Encode(cameras)
}

func handler(cameras map[string]camera, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera_id := vars["camera"]
	match, _ := regexp.MatchString("g\\d+SALT-4e2816a6aa799eb76d1a9ff7265d5371", camera_id)

	if match {
		camera_id = strings.Replace(camera_id, "SALT-4e2816a6aa799eb76d1a9ff7265d5371", "", 1)
		ping_counter.Add(1)
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// ip := strings.Split(r.RemoteAddr, ":")[0]

		c := camera{IP: ip, DateTime: time.Now()}

		cameras[camera_id] = c
	} else {

	}
}

func init() {
	prometheus.MustRegister(ping_counter)
}

func main() {

	cameras := make(map[string]camera)

	r := mux.NewRouter()

	r.Handle("/metrics", prometheus.Handler())

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		status(cameras, w, r)
	})
	r.HandleFunc("/{camera}", func(w http.ResponseWriter, r *http.Request) {
		handler(cameras, w, r)
	})

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
