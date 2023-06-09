package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	url2 "net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type URequest interface {
	DoRequest(url string, method string, data []byte) ([]byte, error)
	Get(url string, data map[string]interface{}) error
	Post(url string, data map[string]interface{}) error
	ReadData() []byte
}
type URequestImpl struct {
	baseUrl string
	headers map[string]string
	timeout time.Duration
	result  interface{}
}

func (uR *URequestImpl) DoRequest(url string, method string, data []byte) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 10}
	var postData io.Reader

	if strings.ToUpper(method) == "POST" {
		postData = bytes.NewBuffer(data)
	} else if strings.ToUpper(method) == "GET" {
		var m map[string]interface{}
		if len(data) > 0 {
			err := json.Unmarshal(data, &m)
			if err != nil {
				return nil, err
			}
		}
		params := url2.Values{}
		for k, v := range m {
			t := reflect.TypeOf(v)
			if t.Kind() == reflect.String {
				params.Add(k, v.(string))
			} else if t.Kind() == reflect.Int {
				params.Add(k, strconv.Itoa(v.(int)))
			} else if t.Kind() == reflect.Float64 {
				params.Add(k, strconv.FormatFloat(v.(float64), 'f', 2, 64))
			} else if t.Kind() == reflect.Float32 {
				params.Add(k, strconv.FormatFloat(v.(float64), 'f', 2, 32))
			} else if t.Kind() == reflect.Bool {
				params.Add(k, strconv.FormatBool(v.(bool)))
			}
		}
		if len(data) > 0 {
			url = url + "?" + params.Encode()
		}
	}
	req, err := http.NewRequest(method, url, postData)
	if err != nil {
		return nil, err
	}
	for k, v := range uR.headers {
		req.Header.Set(k, v)
	}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// Get send get request
func (uR *URequestImpl) Get(url string, data map[string]interface{}) error {
	var m []byte
	if data != nil {
		marshal, err := json.Marshal(data)
		if err != nil {
			return err
		}
		m = marshal
	}
	bUrl, _ := url2.JoinPath(uR.baseUrl, url)
	request, err := uR.DoRequest(bUrl, "GET", m)
	if err != nil {
		return err
	}
	uR.result = request
	return nil
}

// Post send post request
func (uR *URequestImpl) Post(url string, data map[string]interface{}) error {
	var m []byte
	if data != nil {
		marshal, err := json.Marshal(data)
		if err != nil {
			return err
		}
		m = marshal
	}
	bUrl, _ := url2.JoinPath(uR.baseUrl, url)
	request, err := uR.DoRequest(bUrl, "POST", m)
	if err != nil {
		return err
	}
	uR.result = request
	return nil
}

// ReadData read data from response
func (uR *URequestImpl) ReadData() []byte {
	return uR.result.([]byte)
}
func New(baseUrl string, headers map[string]string) URequest {
	return &URequestImpl{
		baseUrl: baseUrl,
		headers: headers,
		timeout: time.Second * 10,
		result:  nil,
	}
}
