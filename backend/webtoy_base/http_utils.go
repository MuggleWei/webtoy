package webtoy_base

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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
		panic(err)
	}
	defer rsp.Body.Close()

	w.Header().Set("Content-Type", rsp.Header.Get("Content-Type"))
	for _, cookie := range rsp.Cookies() {
		http.SetCookie(w, cookie)
	}

	return io.ReadAll(rsp.Body)
}
