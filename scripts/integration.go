package main

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"strings"
	"time"
)

func CreateClient(accessKey, secretKey, accountId, fcRegion string) (_result *fc_open20210406.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(secretKey),
	}
	// 访问的域名
	config.Endpoint = tea.String(fmt.Sprintf("%s.%s.fc.aliyuncs.com", accountId, fcRegion))
	_result = &fc_open20210406.Client{}
	_result, _err = fc_open20210406.NewClient(config)
	return _result, _err
}

func _getIdleFunction(client *fc_open20210406.Client, serviceName string) (_functionName string, _err error) {
	listFunctionsHeaders := &fc_open20210406.ListFunctionsHeaders{}
	listFunctionsRequest := &fc_open20210406.ListFunctionsRequest{}
	runtime := &util.RuntimeOptions{}
	functionNames := make([]string, 0)
	_response, _err := client.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
	if _err != nil {
		return "", _err
	}
	for _, fc := range _response.Body.Functions {
		functionNames = append(functionNames, *fc.FunctionName)
	}

	for _, functionName := range functionNames {
		listStatefulAsyncInvocationsHeaders := &fc_open20210406.ListStatefulAsyncInvocationsHeaders{}
		listStatefulAsyncInvocationsRequest := &fc_open20210406.ListStatefulAsyncInvocationsRequest{}
		_response, _err := client.ListStatefulAsyncInvocationsWithOptions(tea.String(serviceName), tea.String(functionName), listStatefulAsyncInvocationsRequest, listStatefulAsyncInvocationsHeaders, runtime)
		if _err != nil {
			return "", _err
		}
		idle := true
		for _, invocation := range _response.Body.Invocations {
			if fmt.Sprint(*invocation.EndTime) == "0" {
				idle = false
				break
			}
		}
		if idle {
			return functionName, nil
		}
	}
	return "", nil
}

func _invokeFunction(client *fc_open20210406.Client, serviceName, functionName, invocationId, ossBucketRegion, ossBucketName, ossObjectPath, diffFuncNames string) (_err error) {
	invokeFunctionHeaders := &fc_open20210406.InvokeFunctionHeaders{
		XFcInvocationType:            tea.String("Async"),
		XFcLogType:                   tea.String("None"),
		XFcStatefulAsyncInvocationId: tea.String(invocationId),
	}
	body := map[string]interface{}{
		"diffFuncNames":   diffFuncNames,
		"ossBucketName":   ossBucketName,
		"ossBucketRegion": ossBucketRegion,
		"ossObjectPath":   ossObjectPath,
	}
	bodyString, err := json.Marshal(body)
	if err != nil {
		return err
	}
	invokeFunctionRequest := &fc_open20210406.InvokeFunctionRequest{
		Qualifier: tea.String("LATEST"),
		Body:      util.ToBytes(tea.String(string(bodyString))),
	}
	runtime := &util.RuntimeOptions{}

	_, _err = client.InvokeFunctionWithOptions(tea.String(serviceName), tea.String(functionName), invokeFunctionRequest, invokeFunctionHeaders, runtime)
	if _err != nil {
		return _err
	}

	endpoint := fmt.Sprintf("https://oss-%s.aliyuncs.com", ossBucketRegion)
	var options []oss.ClientOption
	accessKey, _ := client.GetAccessKeyId()
	secretKey, _ := client.GetAccessKeySecret()
	ossClient, err := oss.New(endpoint, *accessKey, *secretKey, options...)
	if err != nil {
		return err
	}
	bucket, err := ossClient.Bucket(ossBucketName)
	if err != nil {
		return err
	}

	getStatefulAsyncInvocationHeaders := &fc_open20210406.GetStatefulAsyncInvocationHeaders{}
	getStatefulAsyncInvocationRequest := &fc_open20210406.GetStatefulAsyncInvocationRequest{}

	terraformRunLog := "terraform.run.log"
	lastRunLog := ""
	for true {
		_response, _err := client.GetStatefulAsyncInvocationWithOptions(tea.String(serviceName), tea.String(functionName), tea.String(invocationId), getStatefulAsyncInvocationRequest, getStatefulAsyncInvocationHeaders, runtime)
		if _err != nil {
			return _err
		}
		if err := bucket.GetObjectToFile(ossObjectPath+"/"+terraformRunLog, terraformRunLog); err == nil {
			runLog, _ := os.ReadFile(terraformRunLog)
			printRunLog := strings.TrimPrefix(string(runLog), lastRunLog)
			if printRunLog != "" {
				fmt.Println(printRunLog)
				lastRunLog = string(runLog)
			}
		}
		if fmt.Sprint(*_response.Body.EndTime) == "0" {
			time.Sleep(5 * time.Second)
			continue
		}
		if *_response.Body.Status != "Succeeded" {
			log.Println(*_response.Body.InvocationErrorMessage)
			return fmt.Errorf(*_response.Body.InvocationErrorMessage)
		}
		return nil
	}
	return nil
}
func main() {
	accessKey := os.Args[1]
	secretKey := os.Args[2]
	accountId := os.Args[3]
	serviceName := os.Args[4]
	fcRegion := os.Args[5]
	client, _err := CreateClient(accessKey, secretKey, accountId, fcRegion)
	if _err != nil {
		log.Println(_err)
		os.Exit(1)
	}
	ossBucketRegion := strings.TrimSpace(os.Args[6])
	ossBucketName := strings.TrimSpace(os.Args[7])
	ossObjectPath := strings.TrimSpace(os.Args[8])
	invocationId := strings.Replace(strings.TrimPrefix(ossObjectPath, "github-actions"), "/", "_", -1)
	diffFuncNames := strings.Trim(os.Getenv("DIFF_FUNC_NAMES"), ";")
	functionName := ""
	for true {
		if idleFunc, err := _getIdleFunction(client, serviceName); err != nil {
			log.Println("_getIdleFunction got an error:", err)
			os.Exit(1)
		} else if idleFunc != "" {
			functionName = idleFunc
			break
		}
	}

	log.Println("using function", functionName)
	log.Println("trace id:", invocationId)
	if err := _invokeFunction(client, serviceName, functionName, invocationId, ossBucketRegion, ossBucketName, ossObjectPath, diffFuncNames); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
