package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	url = "https://www.huobi.io/-/x/pro/market/overview5?r=o4qf5j"
)

type detail struct {
	Symbol string
	Close  float64
}

type response struct {
	Status string
	Data   []detail
	m      map[string]float64
}

func (r *response) once() {
	var (
		req    = fasthttp.AcquireRequest()
		resp   = fasthttp.AcquireResponse()
		result response
	)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Println("do: ", err)
		return
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil || result.Status != "ok" {
		log.Println("json unmarshal: ", err, result.Status)
		return
	}

	result.m = make(map[string]float64)
	for i := range result.Data {
		result.m[result.Data[i].Symbol] = result.Data[i].Close
	}
	*r = result
}

func (r *response) get(curr string) float64 {
	curr = strings.ReplaceAll(curr, "-", "")
	curr = strings.ToLower(curr)
	return r.m[curr]
}

func loopScraw() {
	for {
		time.Sleep(time.Minute)
		spider.once()
	}
}

var spider response

func init() {

	rand.Seed(time.Now().Unix())
	spider.once()
	go loopScraw()
}
