package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:9999"
	}
	sl := os.Getenv("SLEEP_MS")
	sleep, err := strconv.ParseInt(sl, 10, 64)
	if err != nil {
		sleep = 2000
	}

	h := http.NewServeMux()
	h.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"status_code": 200}`))
	})
	go func() {
		err = http.ListenAndServe(":9998", h)
		if err != nil {
			log.Print(err)
		}
	}()

	cl := http.DefaultClient
	cl.Timeout = 100 * time.Millisecond
	req, err := http.NewRequest("GET", host+"/?querystring=querystring_values", nil)
	if err != nil {
		panic(err.Error())
	}

	for {
		for i := 0; i < 4; i++ {
			go doRequest(cl, req)
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

func doRequest(client *http.Client, req *http.Request) {
	body := make([]byte, 2000)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request error: ", err.Error())
		return
	}

	// log.Println("response from: ", resp.Header.Get("X-Hostname"))
	if resp.StatusCode >= 400 {
		log.Println(resp.StatusCode)
	}
	resp.Body.Read(body)
	defer resp.Body.Close()
	log.Printf("%s\n=============================\n", string(body))
}
