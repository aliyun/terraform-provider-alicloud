package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
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
	listFunctionsRequest.Limit = tea.Int32(20)
	runtime := &util.RuntimeOptions{}
	functionNames := make([]string, 0)
	nextToken := ""

	for {
		if nextToken != "" {
			listFunctionsRequest.NextToken = &nextToken
		}
		_response, _err := client.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
		if _err != nil {
			return "", _err
		}

		if _response.Body.Functions == nil || len(_response.Body.Functions) < 1 {
			break
		}

		for _, fc := range _response.Body.Functions {
			functionNames = append(functionNames, *fc.FunctionName)
		}

		nextToken = ""
		if _response.Body.NextToken != nil {
			nextToken = *_response.Body.NextToken
		}

		if nextToken == "" {
			break
		}
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

func _checkInvocationIsExist(client *fc_open20210406.Client, serviceName, invocationId string) (exist bool, _functionName string, _err error) {
	listFunctionsHeaders := &fc_open20210406.ListFunctionsHeaders{}
	listFunctionsRequest := &fc_open20210406.ListFunctionsRequest{}
	runtime := &util.RuntimeOptions{}
	_response, _err := client.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
	if _err != nil {
		return false, "", _err
	}

	var returnError error
	for _, fc := range _response.Body.Functions {
		functionName := *fc.FunctionName
		getStatefulAsyncInvocationHeaders := &fc_open20210406.GetStatefulAsyncInvocationHeaders{}
		getStatefulAsyncInvocationRequest := &fc_open20210406.GetStatefulAsyncInvocationRequest{}
		_resp, _err := client.GetStatefulAsyncInvocationWithOptions(tea.String(serviceName), tea.String(functionName), tea.String(invocationId), getStatefulAsyncInvocationRequest, getStatefulAsyncInvocationHeaders, runtime)
		if _err != nil {
			if strings.Contains(_err.Error(), "StatefulAsyncInvocationNotFound") {
				continue
			}
			returnError = _err
			log.Printf("[ERROR] getting invocation %s failed. Error:%v", invocationId, _err)
		} else if *_resp.Body.InvocationId == invocationId {
			return true, functionName, nil
		}
	}
	return false, "", returnError
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
		if strings.Contains(_err.Error(), "StatefulAsyncInvocationAlreadyExists") {
			log.Printf("the invocation %s has been existed in the function: %s", invocationId, functionName)
		} else {
			return _err
		}
	}

	getStatefulAsyncInvocationHeaders := &fc_open20210406.GetStatefulAsyncInvocationHeaders{}
	getStatefulAsyncInvocationRequest := &fc_open20210406.GetStatefulAsyncInvocationRequest{}

	for true {
		_response, _err := client.GetStatefulAsyncInvocationWithOptions(tea.String(serviceName), tea.String(functionName), tea.String(invocationId), getStatefulAsyncInvocationRequest, getStatefulAsyncInvocationHeaders, runtime)
		if _err != nil {
			return _err
		}
		if fmt.Sprint(*_response.Body.EndTime) == "0" {
			time.Sleep(5 * time.Second)
			continue
		}
		if *_response.Body.Status != "Succeeded" {
			return fmt.Errorf(*_response.Body.InvocationErrorMessage)
		}
		return nil
	}
	return nil
}

func _fetchRunRawLog(ossBucketRegion, ossBucketName, ossObjectPath, accessKeyId, accessKeySecret string) (_err error) {

	endpoint := fmt.Sprintf("https://oss-%s.aliyuncs.com", ossBucketRegion)
	var options []oss.ClientOption
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, options...)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(ossBucketName)
	if err != nil {
		return err
	}
	targetObjectPath := "terraform.run.raw.log"
	sourceObjectPath := ossObjectPath + "/" + targetObjectPath

	lastOffset := 0
	logFileIsExisted := false
	for true {
		if !logFileIsExisted {
			if ok, err := bucket.IsObjectExist(sourceObjectPath); err != nil {
				fmt.Println("\033[33m[ERROR]\033[0m checking terraform run raw log whether is existed failed. Error:", err, ". Retry...")
				time.Sleep(5 * time.Second)
				continue
			} else if !ok {
				fmt.Println("waiting for terraform run raw log...")
				time.Sleep(10 * time.Second)
				continue
			} else {
				logFileIsExisted = true
			}
		}
		err = bucket.GetObjectToFile(sourceObjectPath, targetObjectPath)
		if err != nil {
			fmt.Println("\033[33m[ERROR]\033[0m getting terraform run raw log failed. Error:", err, ". Retry...")
			time.Sleep(5 * time.Second)
			continue
		}
		runtimeLogFile, err := os.Open(targetObjectPath)
		if err != nil {
			fmt.Println("reading terraform run raw log file failed:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer runtimeLogFile.Close()

		runLogContent := make([]byte, 100000000)
		offset, er := runtimeLogFile.Read(runLogContent)
		if er != nil && fmt.Sprint(er) != "EOF" {
			log.Println("[ERROR] reading run log response failed:", err)
		}
		if offset > lastOffset {
			fmt.Println(string(runLogContent[lastOffset:offset]))
			lastOffset = offset
		}
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
	invocationId := strings.Replace(ossObjectPath, "/", "_", -1)
	diffFuncNames := strings.Trim(strings.TrimSpace(os.Args[9]), ";")
	functionName := ""
	log.Println("oss object path:", ossObjectPath)
	log.Println("trace id:", invocationId)
	if exist, fcName, err := _checkInvocationIsExist(client, serviceName, invocationId); err != nil {
		log.Printf("[ERROR] checking invocation %s failed. Error: %v", invocationId, err)
	} else if exist {
		log.Printf("the invocation %s has been existed in the function %s", invocationId, fcName)
		os.Exit(0)
	}
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
	go _fetchRunRawLog(ossBucketRegion, ossBucketName, ossObjectPath, accessKey, secretKey)
	if err := _invokeFunction(client, serviceName, functionName, invocationId, ossBucketRegion, ossBucketName, ossObjectPath, diffFuncNames); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
