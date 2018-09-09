package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/resty.v1"
)

const serverURL = "http://localhost:8080"

const (
	start int = iota
	waitRoom
	login
	booking
	bookingReserve
	bookingGA
	payment
	end
)
const speed = 10

type resStatus struct {
	ID      string
	Message string
	Token   string
}

var globalCache *cache.Cache
var c *resty.Client

func init() {
	globalCache = cache.New(0, 0)
	globalCache.Set("fin", 0, 0)
	globalCache.Set("gotticket", 0, 0)
	globalCache.Set("notgetticket", 0, 0)

	c = resty.New()
	c.SetRetryCount(100)
	c.SetRESTMode()
	c.SetCloseConnection(false)
	c.SetRetryWaitTime(time.Second)
	c.SetRetryMaxWaitTime(time.Hour)
	c.SetLogger(ioutil.Discard)
}

func client() {
	// var token string
	var id string
	state := start
	for state != end {
		// fmt.Println("state: ", state)
		switch state {
		case start:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/waitroom/enter")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				id = r.ID
				// token = r.Token
				state = waitRoom
				// fmt.Printf("id:%v state:%v token:%v\n", id, state, token)
				continue
			}
			time.Sleep(time.Second * speed)
		case waitRoom:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/waitroom/exit/" + id)
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = login
				continue
			}
			time.Sleep(time.Second * 30)
		case login:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/login")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = booking
			}
			time.Sleep(time.Second * speed)
		case booking:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/booking")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = bookingReserve
			}
			time.Sleep(time.Second * speed)
		case bookingReserve:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/booking/reserveseat")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = payment
			}
			time.Sleep(time.Second * speed)
		case bookingGA:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/booking/gaseat")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = payment
			}
			time.Sleep(time.Second * speed)
		case payment:
			resp, _ := c.R().SetResult(resStatus{}).Get(serverURL + "/payment")
			r := resp.Result().(*resStatus)
			if r.Message == "OK" {
				state = end
			}
			time.Sleep(time.Second * speed)
		}
	}

	globalCache.IncrementInt("fin", 1)
}

func launchClient(numberOfClient int) {
	for i := 0; i < numberOfClient; i++ {
		go client()
		time.Sleep(time.Millisecond * 10)
	}

}

func main() {
	numberOfClient := 10000

	go launchClient(numberOfClient)

	for {
		time.Sleep(time.Second * 2)
		fin, _ := globalCache.Get("fin")
		gotticket, _ := globalCache.Get("gotticket")
		notgetticket, _ := globalCache.Get("notgetticket")
		fmt.Printf("Users finished: %v\n", fin)
		fmt.Printf("Users got ticket: %v\n", gotticket)
		fmt.Printf("Users not get ticket: %v\n\n", notgetticket)
		if fin == numberOfClient {
			break
		}
	}
}
