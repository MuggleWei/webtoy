package webtoy_base

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type MessageRsp struct {
	Code   int         `json:"code,omitempty"`
	ErrMsg string      `json:"msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func HttpResponse(w http.ResponseWriter, res *MessageRsp) error {
	b, err := json.Marshal(*res)
	if err != nil {
		log.Warning("failed marshal response model")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Warning("failed http write: %v", err.Error())
	}
	return err
}

// http service registry post
func HttpSRPost(serviceName, urlPath string, transport http.RoundTripper, req interface{}, rspData interface{}) (*MessageRsp, error) {
	log.Debugf("http service registry post: %v %v", serviceName, urlPath)

	srClient := GetSrClientComponent().Client
	addr, err := srClient.ClientLB.GetService(serviceName)
	if err != nil {
		errMsg := fmt.Sprintf("failed get service %v address", serviceName)
		log.Errorf("%v", errMsg)
		return nil, errors.New(errMsg)
	}

	url := "http://" + addr + urlPath

	b, err := HttpClientPost(url, transport, req)
	if err != nil {
		log.Errorf("%v", err.Error())
		return nil, err
	}

	var rsp MessageRsp
	err = json.Unmarshal(b, &rsp)
	if err != nil {
		log.Errorf("failed unmarshal service response: %v, %v", err.Error(), string(b))
		return nil, err
	}

	if rspData != nil {
		decodeConfig := mapstructure.DecoderConfig{TagName: "json", Result: rspData}
		decoder, err := mapstructure.NewDecoder(&decodeConfig)
		if err != nil {
			log.Errorf("failed new response data decoder")
			return nil, err
		}

		err = decoder.Decode(rsp.Data)
		if err != nil {
			log.Errorf("failed decode response data: %v", err.Error())
			return nil, err
		}

		rsp.Data = rspData
	}

	return &rsp, nil
}

func HttpClientPost(url string, transport http.RoundTripper, req interface{}) ([]byte, error) {
	b, err := json.Marshal(req)
	if err != nil {
		log.Errorf("failed marshal req message")
		return nil, err
	}

	return HttpClientPostBytes(url, transport, b)
}

func HttpClientPostBytes(url string, transport http.RoundTripper, b []byte) ([]byte, error) {
	log.Debugf("HttpClientPostBytes: %v", url)

	client := &http.Client{Transport: transport}
	rsp, err := client.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Errorf("failed HttpTransportGet %v", err.Error())
		panic(err)
	}

	if rsp.StatusCode != http.StatusOK {
		log.Errorf("service return status code not equal 200")
		return nil, errors.New("upstream serivce return " + fmt.Sprint(rsp.StatusCode))
	}

	defer rsp.Body.Close()

	return io.ReadAll(rsp.Body)
}

func HttpTransportGet(url string, transport http.RoundTripper, w http.ResponseWriter) ([]byte, error) {
	log.Debugf("HttpTransportGet: %v", url)

	client := &http.Client{Transport: transport}
	rsp, err := client.Get(url)
	if err != nil {
		log.Errorf("failed HttpTransportGet %v", err.Error())
		panic(err)
	}
	defer rsp.Body.Close()

	w.Header().Set("Content-Type", rsp.Header.Get("Content-Type"))
	for _, cookie := range rsp.Cookies() {
		http.SetCookie(w, cookie)
	}

	return io.ReadAll(rsp.Body)
}

func HttpTransportPost(url string, transport http.RoundTripper, obj interface{}, w http.ResponseWriter) ([]byte, error) {
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return HttpTransportPostBytes(url, transport, body, w)
}

func HttpTransportPostBytes(url string, transport http.RoundTripper, body []byte, w http.ResponseWriter) ([]byte, error) {
	log.Debugf("HttpTransportPostBytes: %v", url)

	client := &http.Client{Transport: transport}
	rsp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("failed HttpTransportPostBytes %v", err.Error())
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		log.Errorf("service return status code not equal 200")
		return nil, err
	}

	defer rsp.Body.Close()

	w.Header().Set("Content-Type", rsp.Header.Get("Content-Type"))
	for _, cookie := range rsp.Cookies() {
		http.SetCookie(w, cookie)
	}

	return io.ReadAll(rsp.Body)
}
