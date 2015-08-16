package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	ping_counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "camera_pings",
		Help: "Number of times the stats server was pinged.",
	})
	process_memory = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_TotalMem",
		Help: "Total Memory used",
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
	match, _ := regexp.MatchString("g\\d+", camera_id)

	if match {
		ping_counter.Add(1)
		ip := strings.Split(r.RemoteAddr, ":")[0]

		c := camera{IP: ip, DateTime: time.Now()}

		cameras[camera_id] = c
	}
}

func process_stats() {
	time.Sleep(1000)
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	process_memory.Set(float64(m.Alloc))
}

func main() {
	go process_stats()

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
