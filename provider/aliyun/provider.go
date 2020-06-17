package aliyun

import (
	"bytes"
	"fmt"
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

type Aliyun struct {
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
	return filepath.Join(fileDir, "..", "..", "credentials", "alicloud"), nil
}

func New() (*Aliyun, error) {
	name := "aliyun"

	credsFile, err := credsPath()
	if err != nil {
		return nil, errors.WithMessage(err, "getting aliyun credentials path")
	}

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "cn-shanghai"

	return &Aliyun{region: region, name: name, credentials: credsFile}, nil
}

func (aliyun *Aliyun) Name() string {
	return aliyun.name
}

func getRegion(session *session.Session) (string, error) {
	metaClient := ec2metadata.New(session)
	region, err := metaClient.Region()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (aliyun *Aliyun) buildGFuncInvokeReq(funcName string, projectId string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}

	// https://YOUR_REGION-YOUR_PROJECT_ID.cloudfunctions.net/FUNCTION_NAME?sleep={time}
	// http://69d4ed74258e4ce08eac8edf4f44c000-cn-shanghai.alicloudapi.com/index.handler
	// https://1581223932488159.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/aliyuntestservice-dev/testhandler/
	//projectId = "69d4ed74258e4ce08eac8edf4f44c000"
	funcUrl.Scheme = "https"
	funcUrl.Host = "1581223932488159.cn-shanghai.fc.aliyuncs.com"
	funcUrl.Path = "2016-08-15/proxy/aliyuntestservice-dev/testhandler/"
	//funcUrl.Host = fmt.Sprintf("%s-%s.alicloudapi.com", projectId, aliyun.region)
	//funcUrl.Path = path.Join(funcUrl.Path, funcName)
	//fmt.Println(funcUrl.Path)
	fmt.Println(funcUrl.String())
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
		return aliyun.buildGFuncInvokeReq(function.Handler(),stack.Project(), qParams, &localHeaders, body)
	}
}

func (aliyun *Aliyun) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
