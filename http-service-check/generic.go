package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	DoEvery(15*time.Second, CheckStatusEvery)
}

func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func CheckStatusEvery(t time.Time) {
	req, err := http.NewRequest("GET", os.Getenv("endpoint"), nil)
	if err != nil {
		fmt.Printf("ERROR : " + time.Now().String() + "Request was not generated.\n")
	}
	startTime := time.Now()

	//ssl handshake
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Errors on response client ", err.Error())
	}
	responseTime := time.Since(startTime)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Errors on response body close ", err.Error())
		}
	}(res.Body)

	currentTime := time.Now()
	expectedStatusCode, err := strconv.Atoi(os.Getenv("ExpectedStatusCode"))
	if err != nil {
		fmt.Println("Error on load ExpectedStatusCode from env: ", err.Error())
		expectedStatusCode = 200
	}
	if res.StatusCode == expectedStatusCode  {
		fmt.Println(currentTime.Format("2006.01.02 15:04:05") + " " + os.Getenv("endpoint") + " is working fine. ResponseCode: " + strconv.Itoa(res.StatusCode) + " ResponseTime: " + responseTime.String())
	} else {
		fmt.Println(currentTime.Format("2006.01.02 15:04:05") + " " + os.Getenv("endpoint") + " is not working. ResponseCode: " + strconv.Itoa(res.StatusCode) + " ResponseTime: " + responseTime.String())
	}
}