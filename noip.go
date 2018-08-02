package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"strconv"
)

func runUpdate() {
	username := os.Getenv("USER_NAME")
	if username == "" {
		panic("USER_NAME NOT PASSED!")
	}

	password := os.Getenv("USER_PASSWD")
	if password == "" {
		panic("USER_PASSWD NOT PASSED!")
	}

	hostname := os.Getenv("HOST_NAME")
	if hostname == "" {
		panic("HOST_NAME NOT PASSED!")
	}

	ip := os.Getenv("IP")
	if ip == "" {
		fmt.Println("WARNING: IP not passed. Using external IP.")
		resp, err := http.Get("http://myexternalip.com/raw")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		ip = strings.Replace(string(body), "\n", "", -1)
		fmt.Printf("External IP: %s\n", ip)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", username, password, hostname, ip), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", fmt.Sprintf("Go NoIP DDNS Update Client by bertof from %s", os.Getenv("hostname")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}

func main() {
	intervalString := os.Getenv("INTERVAL")
	if intervalString == "" {
		intervalString = "3600"
	}

	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		panic(err)
	}

	runUpdate()
	for range time.NewTicker(time.Duration(interval) * time.Second).C {
		runUpdate()
	}
}
