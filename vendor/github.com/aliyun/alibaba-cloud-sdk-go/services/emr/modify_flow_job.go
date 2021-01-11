package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyFlowJob invokes the emr.ModifyFlowJob API synchronously
func (client *Client) ModifyFlowJob(request *ModifyFlowJobRequest) (response *ModifyFlowJobResponse, err error) {
	response = CreateModifyFlowJobResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyFlowJobWithChan invokes the emr.ModifyFlowJob API asynchronously
func (client *Client) ModifyFlowJobWithChan(request *ModifyFlowJobRequest) (<-chan *ModifyFlowJobResponse, <-chan error) {
	responseChan := make(chan *ModifyFlowJobResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyFlowJob(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ModifyFlowJobWithCallback invokes the emr.ModifyFlowJob API asynchronously
func (client *Client) ModifyFlowJobWithCallback(request *ModifyFlowJobRequest, callback func(response *ModifyFlowJobResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyFlowJobResponse
		var err error
		defer close(result)
		response, err = client.ModifyFlowJob(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ModifyFlowJobRequest is the request struct for api ModifyFlowJob
type ModifyFlowJobRequest struct {
	*requests.RpcRequest
	RetryPolicy       string                       `position:"Query" name:"RetryPolicy"`
	RunConf           string                       `position:"Query" name:"RunConf"`
	Description       string                       `position:"Query" name:"Description"`
	ParamConf         string                       `position:"Query" name:"ParamConf"`
	ResourceList      *[]ModifyFlowJobResourceList `position:"Query" name:"ResourceList"  type:"Repeated"`
	FailAct           string                       `position:"Query" name:"FailAct"`
	Mode              string                       `position:"Query" name:"Mode"`
	MonitorConf       string                       `position:"Query" name:"MonitorConf"`
	Id                string                       `position:"Query" name:"Id"`
	MaxRetry          requests.Integer             `position:"Query" name:"MaxRetry"`
	AlertConf         string                       `position:"Query" name:"AlertConf"`
	ProjectId         string                       `position:"Query" name:"ProjectId"`
	EnvConf           string                       `position:"Query" name:"EnvConf"`
	MaxRunningTimeSec requests.Integer             `position:"Query" name:"MaxRunningTimeSec"`
	ClusterId         string                       `position:"Query" name:"ClusterId"`
	Params            string                       `position:"Query" name:"Params"`
	CustomVariables   string                       `position:"Query" name:"CustomVariables"`
	RetryInterval     requests.Integer             `position:"Query" name:"RetryInterval"`
	Name              string                       `position:"Query" name:"Name"`
}

// ModifyFlowJobResourceList is a repeated param struct in ModifyFlowJobRequest
type ModifyFlowJobResourceList struct {
	Path  string `name:"Path"`
	Alias string `name:"Alias"`
}

// ModifyFlowJobResponse is the response struct for api ModifyFlowJob
type ModifyFlowJobResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      bool   `json:"Data" xml:"Data"`
}

// CreateModifyFlowJobRequest creates a request to invoke ModifyFlowJob API
func CreateModifyFlowJobRequest() (request *ModifyFlowJobRequest) {
	request = &ModifyFlowJobRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "ModifyFlowJob", "emr", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyFlowJobResponse creates a response to parse from ModifyFlowJob response
func CreateModifyFlowJobResponse() (response *ModifyFlowJobResponse) {
	response = &ModifyFlowJobResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
