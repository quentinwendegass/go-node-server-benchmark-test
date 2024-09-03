package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"time"
)

var HOST = "http://localhost:8080"
var PARALLEL_REQUESTS = 100
var ITERATIONS = 20

func main() {
	sampleJsonBuffer, err := getSampleJson("https://microsoftedge.github.io/Demos/json-dummy-data/5MB-min.json")
	if err != nil {
		panic(err)
	}

	benchmark(sampleJsonBuffer, PARALLEL_REQUESTS, ITERATIONS)
}

func benchmark(data []byte, parallelRequests int, iterations int) {
	durations := []float64{}

	done, exit := heartBeat()

	for j := 0; j < iterations; j++ {
		ch := make(chan time.Duration, parallelRequests)

		for i := 0; i < parallelRequests; i++ {
			go func() {
				start := time.Now()
				makeRequest(data)
				ch <- time.Since(start)
			}()
		}

		for i := 0; i < parallelRequests; i++ {
			duration := <-ch
			durations = append(durations, float64(duration))
		}

		fmt.Print("Iter ", j, ": ")
		printResults(durations)
	}

	close(exit)
	<-done
	fmt.Print("Final: ")
	printResults(durations)
}

func makeRequest(data []byte) error {
	request, err := http.NewRequest("POST", HOST+"/filter", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("unexpected status code")
	}

	_, err = io.ReadAll(response.Body)

	return err
}

func heartBeat() (chan struct{}, chan struct{}) {
	durations := []float64{}
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				start := time.Now()
				makeStatusRequest()
				duration := time.Since(start)
				durations = append(durations, float64(duration))
				fmt.Print("(Heartbeat) Iter ", len(durations)-1, ": ")
				printResults(durations)
			case <-quit:
				fmt.Print("(Heartbeat) Final: ")
				printResults(durations)
				ticker.Stop()
				done <- struct{}{}
				return
			}
		}
	}()

	return done, quit
}

func getSampleJson(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func makeStatusRequest() {
	client := http.Client{}
	client.Get(HOST + "/status")
}

func printResults(data []float64) {
	mean := calcMean(data)
	median := calcMedian(data)
	min := slices.Min(data)
	max := slices.Max(data)

	fmt.Println("Mean:", time.Duration(mean), "- Median:", time.Duration(median), "- Max:", time.Duration(max), "- Min:", time.Duration(min))
}

func calcMedian(data []float64) float64 {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)

	sort.Float64s(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

func calcMean(data []float64) float64 {
	total := 0.0

	for _, v := range data {
		total += v
	}

	return total / float64(len(data))
}
