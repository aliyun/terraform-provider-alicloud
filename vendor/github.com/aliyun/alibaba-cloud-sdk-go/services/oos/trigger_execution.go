package oos

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

// TriggerExecution invokes the oos.TriggerExecution API synchronously
// api document: https://help.aliyun.com/api/oos/triggerexecution.html
func (client *Client) TriggerExecution(request *TriggerExecutionRequest) (response *TriggerExecutionResponse, err error) {
	response = CreateTriggerExecutionResponse()
	err = client.DoAction(request, response)
	return
}

// TriggerExecutionWithChan invokes the oos.TriggerExecution API asynchronously
// api document: https://help.aliyun.com/api/oos/triggerexecution.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TriggerExecutionWithChan(request *TriggerExecutionRequest) (<-chan *TriggerExecutionResponse, <-chan error) {
	responseChan := make(chan *TriggerExecutionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.TriggerExecution(request)
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

// TriggerExecutionWithCallback invokes the oos.TriggerExecution API asynchronously
// api document: https://help.aliyun.com/api/oos/triggerexecution.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TriggerExecutionWithCallback(request *TriggerExecutionRequest, callback func(response *TriggerExecutionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *TriggerExecutionResponse
		var err error
		defer close(result)
		response, err = client.TriggerExecution(request)
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

// TriggerExecutionRequest is the request struct for api TriggerExecution
type TriggerExecutionRequest struct {
	*requests.RpcRequest
	ClientToken string `position:"Query" name:"ClientToken"`
	Type        string `position:"Query" name:"Type"`
	Content     string `position:"Query" name:"Content"`
	ExecutionId string `position:"Query" name:"ExecutionId"`
}

// TriggerExecutionResponse is the response struct for api TriggerExecution
type TriggerExecutionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateTriggerExecutionRequest creates a request to invoke TriggerExecution API
func CreateTriggerExecutionRequest() (request *TriggerExecutionRequest) {
	request = &TriggerExecutionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("oos", "2019-06-01", "TriggerExecution", "", "")
	request.Method = requests.POST
	return
}

// CreateTriggerExecutionResponse creates a response to parse from TriggerExecution response
func CreateTriggerExecutionResponse() (response *TriggerExecutionResponse) {
	response = &TriggerExecutionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
