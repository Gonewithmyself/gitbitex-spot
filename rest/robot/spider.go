package main

import (
	"encoding/json"
	"fmt"
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

type placeOrderRequest struct {
	ClientOid   string  `json:"client_oid"`
	ProductId   string  `json:"productId"`
	Size        float64 `json:"size"`
	Funds       float64 `json:"funds"`
	Price       float64 `json:"price"`
	Side        string  `json:"side"`
	Type        string  `json:"type"`        // [optional] limit or market (default is limit)
	TimeInForce string  `json:"timeInForce"` // [optional] GTC, GTT, IOC, or FOK (default is GTC)
}

func post(arg interface{}) error {
	var (
		req  = fasthttp.AcquireRequest()
		resp = fasthttp.AcquireResponse()
		//result response
	)
	for i := 0; i < len(hds)-1; i += 2 {
		req.Header.Add(hds[i], hds[i+1])
	}

	req.SetRequestURI("http://localhost:8080/api/orders")
	req.Header.SetMethod("POST")
	d, _ := json.Marshal(arg)
	req.SetBody(d)
	req.Header.SetContentLength(len(d))

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	err := fasthttp.Do(req, resp)
	if err != nil {
		fmt.Println((string(resp.Body())), req.Header.String())
	}
	//
	return err
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
	parseHeader()
	rand.Seed(time.Now().Unix())
	spider.once()
	go loopScraw()
}

var hds []string

func parseHeader() {
	src :=
		`Accept: application/json, text/plain, */*
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.9
Connection: keep-alive
Content-Type: application/json;charset=UTF-8
Cookie: csrftoken=nFp2dDJ3hvDHJ9mv9W3tWitGEu5NQnKqElS6sCP2wIL7vLDTeJZhe2s36dHb6vHU; _ga=GA1.1.1939969193.1562124836; accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IjEyNDM2OTk3NkBxcS5jb20iLCJleHBpcmVkQXQiOjE1NzU4ODU1MTMsImlkIjo0MSwicGFzc3dvcmRIYXNoIjoiMTZlZGU4NmFhM2EzMjA1MmM5YjIxOGM3MjA2M2Q5NjgifQ.Brc2HzV2vRdfPgveq9nrLl2nzXed9jkHUM2Y8ymGxHA; io=jqqSbhve8l88TY9TAAA3
HideError: true
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36`

	// src = strings.ReplaceAll(src, "\n", "")
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		hds = append(hds, strings.Split(line, ": ")...)
	}
}
