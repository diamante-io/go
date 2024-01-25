package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/support/errors"
)

type maybeDuration struct {
	maybeDuration time.Duration
	maybeError    error
}

func main() {
	// Friendbot must be running as a local server. Get Friendbot URL from CL.
	fbURL := flag.String("url", "http://0.0.0.0:8000/", "URL of friendbot")
	numRequests := flag.Int("requests", 500, "number of requests")
	flag.Parse()
	durationChannel := make(chan maybeDuration, *numRequests)
	for i := 0; i < *numRequests; i++ {
		kp, err := keypair.Random()
		if err != nil {
			panic(err)
		}
		address := kp.Address()
		go makeFriendbotRequest(address, *fbURL, durationChannel)

		time.Sleep(time.Duration(500) * time.Millisecond)
	}
	durations := []maybeDuration{}
	for i := 0; i < *numRequests; i++ {
		durations = append(durations, <-durationChannel)
	}
	close(durationChannel)
	log.Printf("Got %d times with average %s", *numRequests, mean(durations))
}

func makeFriendbotRequest(address, fbURL string, durationChannel chan maybeDuration) {
	start := time.Now()
	formData := url.Values{
		"addr": {address},
	}
	resp, err := http.PostForm(fbURL, formData)
	if err != nil {
		log.Printf("Got post error: %s", err)
		durationChannel <- maybeDuration{maybeError: errors.Wrap(err, "posting form")}
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Got decode error: %s", err)
		durationChannel <- maybeDuration{maybeError: errors.Wrap(err, "decoding json")}
	}
	timeTrack(start, "makeFriendbotRequest", durationChannel)
}

func timeTrack(start time.Time, name string, durationChannel chan maybeDuration) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
	durationChannel <- maybeDuration{maybeDuration: elapsed}
}

func mean(durations []maybeDuration) time.Duration {
	var total time.Duration
	count := 0
	for _, d := range durations {
		if d.maybeError != nil {
			continue
		}
		total += d.maybeDuration
		count++
	}
	return total / time.Duration(count)
}
