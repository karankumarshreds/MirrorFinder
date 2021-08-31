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
	http.HandleFunc("/fastest-mirror", fastestMirror)
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
			latency := time.Since(start) / time.Millisecond

			if err != nil {
				log.Fatal(err)
			} 
			urlChan     <- mirrorUrl
			latencyChan <- latency
		}()
	}

	response := response{
		FastestURL: <-urlChan,
		Latency:    <-latencyChan,
	}
	encoder  := json.NewEncoder(w)
	err      := encoder.Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}



// type response struct {
// 	FastestURL string        `json:"fastest_url"`
// 	Latency    time.Duration `json:"latency"`
// }

// func findFastest(urls [48]string) response {
// 	urlChan := make(chan string)
// 	latencyChan := make(chan time.Duration)

// 	for _, url := range urls {
// 		mirrorURL := url
// 		go func() {
// 			log.Println("Started probing: ", mirrorURL)
// 			start := time.Now()
// 			_, err := http.Get(mirrorURL + "/README")
// 			latency := time.Since(start) / time.Millisecond
// 			if err == nil {
// 				urlChan <- mirrorURL
// 				latencyChan <- latency
// 			}
// 			log.Printf("Got the best mirror: %s with latency: %s", mirrorURL, latency)
// 		}()
// 	}
// 	return response{<-urlChan, <-latencyChan}
// }

// func main() {
// 	http.HandleFunc("/fastest-mirror", func(w http.ResponseWriter, r *http.Request) {
// 		response := findFastest(mirrors.MirrorList)
// 		respJSON, _ := json.Marshal(response)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(respJSON)
// 	})
// 	port := ":8000"
// 	server := &http.Server{
// 		Addr:           port,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// 	fmt.Printf("Starting server on port %s\n", port)
// 	log.Fatal(server.ListenAndServe())
// }