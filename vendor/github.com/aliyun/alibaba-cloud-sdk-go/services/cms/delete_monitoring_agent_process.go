package cms

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

// DeleteMonitoringAgentProcess invokes the cms.DeleteMonitoringAgentProcess API synchronously
// api document: https://help.aliyun.com/api/cms/deletemonitoringagentprocess.html
func (client *Client) DeleteMonitoringAgentProcess(request *DeleteMonitoringAgentProcessRequest) (response *DeleteMonitoringAgentProcessResponse, err error) {
	response = CreateDeleteMonitoringAgentProcessResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteMonitoringAgentProcessWithChan invokes the cms.DeleteMonitoringAgentProcess API asynchronously
// api document: https://help.aliyun.com/api/cms/deletemonitoringagentprocess.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteMonitoringAgentProcessWithChan(request *DeleteMonitoringAgentProcessRequest) (<-chan *DeleteMonitoringAgentProcessResponse, <-chan error) {
	responseChan := make(chan *DeleteMonitoringAgentProcessResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteMonitoringAgentProcess(request)
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

// DeleteMonitoringAgentProcessWithCallback invokes the cms.DeleteMonitoringAgentProcess API asynchronously
// api document: https://help.aliyun.com/api/cms/deletemonitoringagentprocess.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteMonitoringAgentProcessWithCallback(request *DeleteMonitoringAgentProcessRequest, callback func(response *DeleteMonitoringAgentProcessResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteMonitoringAgentProcessResponse
		var err error
		defer close(result)
		response, err = client.DeleteMonitoringAgentProcess(request)
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

// DeleteMonitoringAgentProcessRequest is the request struct for api DeleteMonitoringAgentProcess
type DeleteMonitoringAgentProcessRequest struct {
	*requests.RpcRequest
	ProcessName string `position:"Query" name:"ProcessName"`
	InstanceId  string `position:"Query" name:"InstanceId"`
	ProcessId   string `position:"Query" name:"ProcessId"`
}

// DeleteMonitoringAgentProcessResponse is the response struct for api DeleteMonitoringAgentProcess
type DeleteMonitoringAgentProcessResponse struct {
	*responses.BaseResponse
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Success   bool   `json:"Success" xml:"Success"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteMonitoringAgentProcessRequest creates a request to invoke DeleteMonitoringAgentProcess API
func CreateDeleteMonitoringAgentProcessRequest() (request *DeleteMonitoringAgentProcessRequest) {
	request = &DeleteMonitoringAgentProcessRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "DeleteMonitoringAgentProcess", "cms", "openAPI")
	return
}

// CreateDeleteMonitoringAgentProcessResponse creates a response to parse from DeleteMonitoringAgentProcess response
func CreateDeleteMonitoringAgentProcessResponse() (response *DeleteMonitoringAgentProcessResponse) {
	response = &DeleteMonitoringAgentProcessResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
