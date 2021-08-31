package main

import (
	_"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/karankumarshreds/MirrorFinder/src"
)


type response struct {
	FastestURL  string           `json:"fastestUrl"`
	Latency     time.Duration    `json:"latency"`
}

func main () {
	http.HandleFunc("fastest/mirror", fastestMirror)
	port := ":8000"
	log.Fatal(http.ListenAndServe(port, nil))
}

func fastestMirror(w http.ResponseWriter, r *http.Request) {
	mirrors      := mirrors.MirrorList
	urlChan      := make(chan string)			
	latencyChan  := make(chan time.Duration)

	for _, url := range mirrors {
		mirrorUrl := url
		go func() {
			start   := time.Now()
			_, err  := http.Get(mirrorUrl + "/README")
			latency := time.Now().Sub(start) / time.Millisecond

			if err == nil {
				urlChan     <- mirrorUrl
				latencyChan <- latency
			}
		}()
	}

	response := response{
		FastestURL: <-urlChan,
		Latency: <-latencyChan,
	}
	encoder      := json.NewEncoder(w)
	err := encoder.Encode(response)
	log.Fatal(err)
}

