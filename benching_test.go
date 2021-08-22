package benching_test

import (
	"GoProject/types"
	"net/http"
	"sync"
	"testing"
)

var waitGroup sync.WaitGroup

func BenchmarkSample(b *testing.B) {
	var statuses chan types.Statuses = make(chan types.Statuses, b.N)
	for i := 0; i < b.N; i++ {
		waitGroup.Add(1)
		go GetRequest(i, b, statuses)
	}
	waitGroup.Wait()
	for i := 0; i < len(statuses)-1; i++ {
		status, ok := <-statuses
		if !ok {
			break
		}
		b.Logf("# %d Status: %+v", status.Number, status.Status)
	}
}

func GetRequest(i int, b *testing.B, statuses chan<- types.Statuses) {
	defer waitGroup.Done()
	client := http.Client{}
	resp, err := client.Get("http://127.0.0.1:8010/json/hackers")
	if err != nil {
		b.Logf("Error in TEST: %+s", err)
		return
	}

	if resp.StatusCode == 200 {
		statuses <- types.Statuses{Number: i, Status: true}
	} else {
		statuses <- types.Statuses{Number: i, Status: false}
	}
}
