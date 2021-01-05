package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/litmuschaos/test-tools/pkg/log"
)

var (
	timeSumsMu sync.RWMutex
	timeSums   int64
	qps        string
	prevReq    int
	sum        int
	totalCount string
)

func main() {
	log.Info("[Status]: Starting QPS provider...")
	go runDataLoop()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timeSumsMu.RLock()
		defer timeSumsMu.RUnlock()
		fmt.Fprint(w, qps)
	})
	http.ListenAndServe(":8080", nil)
}

//// Start the goroutine that will sum the current time
// once per second.
// Create a handler that will read-lock the mutext and
// write the summed time to the client
func runDataLoop() {
	queue := list.New()
	timeInt := os.Getenv("TIME") //The time period in second to calculate mean value of QPS.
	url := os.Getenv("URL")      //URL of endpoint metics
	route := os.Getenv("ROUTE")  //route is a endpoint of application to get qps
	flag.Parse()

	timeOfInt, _ := strconv.Atoi(timeInt)

	for {
		// Within an infinite loop, lock the mutex and
		// increment our value, then sleep for 1 second until
		// the next time we need to get a value.
		start := time.Now()
		timeSums += time.Now().Unix()
		prevReq = 0
		for i := 1; ; i++ {

			req, err := GetRequests(url, route)
			if err != nil {
				qps = strconv.Itoa(0)
				log.Errorf("Failed to fetch responce, err: %v", err)
			}

			second := int(time.Now().Sub(start).Seconds())
			reqs, _ := strconv.Atoi(req)
			sum = reqs

			if second < timeOfInt+1 {
				qps = string(strconv.Itoa(reqs - prevReq))
				queue.PushBack(reqs)
			} else {
				front := queue.Front()
				queue.Remove(front)
				queue.PushBack(reqs)
				sum -= front.Value.(int)
				qps = string(strconv.Itoa(sum / timeOfInt))
			}
			prevReq = reqs

			log.Infof("[Status]: Current total requests : ", req)
			log.Infof("[Status]: Current QPS value is : ", qps)
		}
	}
}

//GetRequests will fetch the responce from metrics and calculate the total requests from front-end of sock-shop.
func GetRequests(url string, route string) (string, error) {

	time.Sleep(1 * time.Second)
	response, err := http.Get(url)
	if err != nil {
		qps = strconv.Itoa(0)
		log.Errorf("Failed to fetch responce, err: %v", err)
		return "", err
	} else {
		defer response.Body.Close()
		metric, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("Failed to read responce, err: %v", err)
			return "", err
		}

		metrics := string(metric)

		value := strings.Split(metrics, "\n")
		flag := 0
		for i := 0; i < len(value) && flag == 0; i++ {

			if true == strings.Contains(string(value[i]), `request_duration_seconds_count{service="front-end",method="get",route="`+route+`",status_code="200`) {
				val := strings.Split(value[i], " ")
				totalCount = (val[1])
				flag = 1
			}
		}

	}
	return totalCount, nil
}