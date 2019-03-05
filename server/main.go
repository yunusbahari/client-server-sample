package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var svcLatencyHist = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "service_latency_seconds",
		Help: "time required to serve request in seconds",
	},
	[]string{"path", "status"},
)

// Hostname is current server hostname
var Hostname string

func init() {
	var err error
	Hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(svcLatencyHist)
}

func main() {
	h := http.NewServeMux()
	avail := true

	h.HandleFunc("/", handler)
	h.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	h.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		failed := false

		defer func() {
			writeMetric("/readiness", time.Since(begin), failed)
		}()

		w.Header().Add("Content-Type", "application/json")
		if avail {
			w.WriteHeader(200)
			w.Write([]byte(`{"status_code": 200}`))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(`{"status_code": 404}`))
			failed = true
		}

	})

	h.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		failed := false

		defer func() {
			writeMetric("/healthz", time.Since(begin), failed)
		}()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"status_code": 200}`))
	})

	grace := make(chan os.Signal)
	signal.Notify(grace, syscall.SIGINT)
	signal.Notify(grace, syscall.SIGTERM)

	go func() {
		sg := <-grace
		fmt.Printf("Received signal %v\n", sg)
		avail = false
		fmt.Println("gracefully shutdown the system")
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
		}
		os.Exit(0)
	}()

	go func() {
		fmt.Println("listening UDP on port 13000")
		conn, _ := net.ListenPacket("udp", ":13000")
		defer conn.Close()

		for {
			b := make([]byte, 2048)
			conn.ReadFrom(b)
			fmt.Print(string(b))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	log.Print("starting http server on port :9999")
	err := http.ListenAndServe(":9999", h)
	if err != nil {
		log.Print(err)
		grace <- os.Interrupt
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	response(w, r)
}

func response(w http.ResponseWriter, r *http.Request) {
	var str, path string
	begin := time.Now()
	failed := false

	defer func() {
		writeMetric(path, time.Since(begin), failed)
		r.Body.Close()
	}()

	d := make([]byte, 2048)
	r.Body.Read(d)
	r.ParseForm()

	path = r.URL.Path
	str = fmt.Sprint("path: ", path,
		"<br/>querystring : ", r.URL.Query(),
		"<br/>server-host : ", Hostname,
		"<br/>request-body: ", string(d),
		"<br/>")

	w.Header().Add("X-Hostname", Hostname)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(str))
}

func writeMetric(path string, t time.Duration, failed bool) {
	status := "ok"
	if failed {
		status = "fail"
	}
	svcLatencyHist.WithLabelValues(path, status).Observe(t.Seconds())
}
