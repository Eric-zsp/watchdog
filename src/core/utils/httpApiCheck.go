package utils

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gologs "github.com/cn-joyconn/gologs"
)

const loggerName = "WebApiRequest"

func HttpGet(urlPath string, queryString string, headers map[string]string, timeOut int) (*http.Response, string, error) {
	client := &http.Client{}
	if timeOut != 0 {
		client.Timeout, _ = time.ParseDuration(strconv.Itoa(timeOut) + "s")
	}
	hasParam := false
	if strings.Index(urlPath, "?") != -1 {
		hasParam = true
	}
	if len(queryString) > 0 {
		if hasParam {
			urlPath += "&" + queryString
		} else {
			urlPath += "?" + queryString
		}
		hasParam = true
	}

	request, err := http.NewRequest("GET", urlPath, nil)
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	if err != nil {
		gologs.GetLogger(loggerName).Error("执行HTTP Get请求时，编码查询字符串“" + queryString + "”发生异常！" + err.Error())
	}

	response, err := client.Do(request)
	result := ""
	if err != nil {
		gologs.GetLogger(loggerName).Error("执行HTTP Get请求时，编码查询字符串“" + queryString + "”发生异常！" + err.Error())
	} else {

		defer response.Body.Close()
		respByte, err1 := ioutil.ReadAll(response.Body)
		if err1 != nil {
			err = err1
			gologs.GetLogger(loggerName).Error("执行HTTP Get请求时，读取response结果出错“" + queryString + "”发生异常！" + err1.Error())
		} else {
			result = string(respByte)
		}
	}

	return response, result, err
}

func HttpPostForm(urlPath string, queryParams map[string]string, headers map[string]string, form map[string]string, timeOut int) (*http.Response, string, error) {
	client := &http.Client{}
	if timeOut != 0 {
		client.Timeout, _ = time.ParseDuration(strconv.Itoa(timeOut) + "s")
	}
	hasParam := false
	if strings.Index(urlPath, "?") != -1 {
		hasParam = true
	}
	qpV := url.Values{}
	if len(queryParams) > 0 {
		for k, v := range queryParams {
			qpV.Add(k, v)
		}
		if hasParam {
			urlPath += "&" + qpV.Encode()
		} else {
			urlPath += "?" + qpV.Encode()
		}
		hasParam = true
	}
	qpV = url.Values{}
	if len(form) > 0 {
		for k, v := range form {
			qpV.Add(k, v)
		}
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range form {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	// request, err := http.NewRequest("POST", urlPath, strings.NewReader(qpV.Encode()))
	request, err := http.NewRequest("POST", urlPath, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	if err != nil {
		gologs.GetLogger(loggerName).Error("执行HTTP POST请求时，请求路径“" + urlPath + "”发生异常！" + err.Error())
	}

	result := ""
	response, err := client.Do(request)
	if err != nil {
		gologs.GetLogger(loggerName).Error("执行HTTP POST请求时，请求路径“" + urlPath + "”发生异常！" + err.Error())
	} else {
		defer response.Body.Close()
		respByte, err1 := ioutil.ReadAll(response.Body)
		if err1 != nil {
			err = err1
			gologs.GetLogger(loggerName).Error("执行HTTP Get请求时，读取response结果出错“" + urlPath + "”发生异常！" + err1.Error())
		} else {
			result = string(respByte)
		}
	}

	return response, result, err
}
