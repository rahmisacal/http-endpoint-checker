package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
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
	req, _ := http.NewRequest("GET", os.Getenv("endpoint"), nil)
	startTime := time.Now()

	//ssl handshake
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	responseTime := time.Since(startTime)
	if err != nil {
		fmt.Println("Errors on response client ", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Errors on response body close ", err.Error())
		}
	}(res.Body)

	currentTime := time.Now()
	expectedStatusCode, err := strconv.Atoi(os.Getenv("ExpectedStatusCode"))
	if err != nil {
		fmt.Println("Error on load ExpectedStatusCode from env:", err.Error())
		expectedStatusCode = 200
	}
	mailBody := "Description of endpoint url details..\n" + "Endpoint: " + os.Getenv("endpoint") + "\n" + "ResponseStatusCode: " + strconv.Itoa(res.StatusCode) + "\n" + "ResponseTime: " + responseTime.String()
	if res.StatusCode == expectedStatusCode  {
		fmt.Println(currentTime.Format("2006.01.02 15:04:05") + " " + os.Getenv("endpoint") + " is working fine. ResponseCode: "+ strconv.Itoa(res.StatusCode) + " ResponseTime: " + responseTime.String())
	} else {
		fmt.Println(currentTime.Format("2006.01.02 15:04:05") + " " + os.Getenv("endpoint") + " not response! ResponseCode: "+ strconv.Itoa(res.StatusCode) + " ResponseTime: " + responseTime.String())
		SendMail("smtp.example.com.tr:25", "notify@smtp.example.com.tr", os.Getenv("Cluster") + " HTTP Egress Check", mailBody, []string{"example@example.com.tr"})
	}
}

func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}