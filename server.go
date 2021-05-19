package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

func hello(event events.APIGatewayResponse) (*events.APIGatewayResponse, error) {
	bodyBytes, err := base64.StdEncoding.DecodeString(event.Body)
	if err != nil {
		return nil, err
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(bodyBytes)))
	if err != nil {
		return nil, err
	}
	req.RequestURI = ""
	u, err := url.Parse(event.Headers["url"]) // url must lowercase
	fmt.Println(event.Headers)
	if err != nil {
		return nil, err
	}
	req.URL = u
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	header := make(map[string]string)
	for k := range resp.Header {
		header[k] = resp.Header.Get(k)
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayResponse{
		IsBase64Encoded: true,
		StatusCode:      resp.StatusCode,
		Headers:         header,
		Body:            base64.StdEncoding.EncodeToString(bodyBytes),
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(hello)
}
