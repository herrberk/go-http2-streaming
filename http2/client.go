package http2

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

type Client struct {
	client *http.Client
}

func (c *Client) Dial() {
	// Adds TLS cert-key pair
	certs, err := tls.LoadX509KeyPair("./http2/certs/key.crt", "./http2/certs/key.key")
	if err != nil {
		log.Fatal(err)
	}

	t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{certs},
			InsecureSkipVerify: true,
		},
	}
	c.client = &http.Client{Transport: t}
}

func (c *Client) Post(data []byte) {
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "localhost:10000",
			Path:   "/",
		},
		Header: http.Header{},
		Body:   ioutil.NopCloser(bytes.NewReader(data)),
	}

	// Sends the request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode == 500 {
		return
	}

	defer resp.Body.Close()

	bufferedReader := bufio.NewReader(resp.Body)

	buffer := make([]byte, 4*1024)

	var totalBytesReceived int

	// Reads the response
	for {
		len, err := bufferedReader.Read(buffer)
		if len > 0 {
			totalBytesReceived += len
			log.Println(len, "bytes received")
			// Prints received data
			// log.Println(string(buffer[:len]))
		}

		if err != nil {
			if err == io.EOF {
				// Last chunk received
				// log.Println(err)
			}
			break
		}
	}
	log.Println("Total Bytes Sent:", len(data))
	log.Println("Total Bytes Received:", totalBytesReceived)
}
