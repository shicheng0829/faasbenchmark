package aliyun

import (
	"bytes"
	"github.com/golang/gddo/httputil/header"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/syncedtrace"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Aliyun struct {
	region      string
	name        string
	//credentials string
}

//func credsPath() (string, error) {
//	_, filename, _, ok := runtime.Caller(1)
//	if !ok {
//		return "", errors.New("getting absolute project path")
//	}
//	fileDir := filepath.Dir(filename)
//	return filepath.Join(fileDir, "..", "..", "credentials", "alicloud"), nil
//}

func New() (*Aliyun, error) {
	name := "aliyun"

	//credsFile, err := credsPath()
	//if err != nil {
	//	return nil, errors.WithMessage(err, "getting aliyun credentials path")
	//}

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "cn-shanghai"

	return &Aliyun{region: region, name: name}, nil
}

func (aliyun *Aliyun) Name() string {
	return aliyun.name
}

func (aliyun *Aliyun) buildGFuncInvokeReq(funcName string, projectId string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}
	funcUrl.Scheme = "https"
	funcUrl.Host = "1581223932488159.cn-shanghai.fc.aliyuncs.com"
	funcUrl.Path = "2016-08-15/proxy/faasbenchmark/"
	funcUrl.Path = path.Join(funcUrl.Path, funcName)
	funcUrl.Path += "/"
	req, err := http.NewRequest("GET", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = qParams.Encode()

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	return req, nil
}

func (aliyun *Aliyun) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) (func(uniqueId string) (*http.Request, error)) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return aliyun.buildGFuncInvokeReq(function.Name(),stack.Project(), qParams, &localHeaders, body)
	}
}

func (aliyun *Aliyun) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
