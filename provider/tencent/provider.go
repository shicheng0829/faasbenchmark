package tencent

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/gddo/httputil/header"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/syncedtrace"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
)

type Tencent struct {
	region      string
	name        string
	credentials string
}

func credsPath() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("getting absolute project path")
	}
	fileDir := filepath.Dir(filename)
	return filepath.Join(fileDir, "..", "..", "credentials", "tencent"), nil
}

func New() (*Tencent, error) {
	name := "tencent"

	credsFile, err := credsPath()
	if err != nil {
		return nil, errors.WithMessage(err, "getting tencent credentials path")
	}

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "ap-guangzhou"

	return &Tencent{region: region, name: name, credentials: credsFile}, nil
}

func (tencent *Tencent) Name() string {
	return tencent.name
}

func getRegion(session *session.Session) (string, error) {
	metaClient := ec2metadata.New(session)
	region, err := metaClient.Region()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (tencent *Tencent) buildGFuncInvokeReq(funcName string, projectId string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}

	// https://YOUR_REGION-YOUR_PROJECT_ID.cloudfunctions.net/FUNCTION_NAME?sleep={time}
	// https://service-keypade6-1256474564.gz.apigw.tencentcs.com/release/?sleep=500
	//fmt.Println("projectID:",projectId)
	//https://service-34fsrgdc-1256474564.gz.apigw.tencentcs.com/release/memstress-dev-memstress4
	funcUrl.Scheme = "https"
	//funcUrl.Host = "service-34fsrgdc-1256474564.gz.apigw.tencentcs.com"
	//funcUrl.Host = "service-7t1pqv44-1256474564.gz.apigw.tencentcs.com"
	//funcUrl.Path = "release/memstress-dev-memstress4"
	//funcUrl.Path = "release/sleepfunc-dev-sleep"
	funcUrl.Host = "service-88noweis-1256474564.gz.apigw.tencentcs.com"
	//funcUrl.Host = fmt.Sprintf("%s-%s.alicloudapi.com", projectId, aliyun.region)
	//funcUrl.Path = path.Join(funcUrl.Path, funcName)
	funcUrl.Path = funcName[6:] + "/"
	//fmt.Println(funcUrl.Path)
	//fmt.Println(funcUrl.String())
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

func (tencent *Tencent) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) (func(uniqueId string) (*http.Request, error)) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return tencent.buildGFuncInvokeReq(function.Handler(),stack.Project(), qParams, &localHeaders, body)
	}
}

func (tencent *Tencent) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
