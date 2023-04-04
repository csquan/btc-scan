package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var httpCli *http.Client

func init() {
	httpCli = &http.Client{Timeout: 20 * time.Second}
}

func DoRequest(url string, method string, params interface{}) (data []byte, err error) {
	reqData, err := json.Marshal(params)
	if err != nil {
		return
	}
	return DoRequestWithHeaders(url, method, reqData, nil)
}

func DoRequestWithHeaders(url string, method string, reqData []byte, headers map[string]string) (data []byte, err error) {
	body := bytes.NewReader(reqData)
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("content-type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		logrus.Errorf("do http request error:%v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode:%d, url:%s,method:%s", resp.StatusCode, url, method)
		return
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read response body error:%v", err)
		return
	}
	logrus.Infof("DoRequestWithHeaders host:%s path:%s, input:%s, response:%v", req.Host, req.URL.Path, string(reqData), string(data))
	return
}

func Post(requestUrl string, bytesData []byte) (ret string, err error) {
	res, err := http.Post(requestUrl,
		"application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	return *str, nil
}

func Get(url string) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return string(body), nil
}
