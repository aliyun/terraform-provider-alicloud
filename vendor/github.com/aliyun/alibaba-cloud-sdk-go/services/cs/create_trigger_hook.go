package cs

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

// CreateTriggerHook invokes the cs.CreateTriggerHook API synchronously
// api document: https://help.aliyun.com/api/cs/createtriggerhook.html
func (client *Client) CreateTriggerHook(request *CreateTriggerHookRequest) (response *CreateTriggerHookResponse, err error) {
	response = CreateCreateTriggerHookResponse()
	err = client.DoAction(request, response)
	return
}

// CreateTriggerHookWithChan invokes the cs.CreateTriggerHook API asynchronously
// api document: https://help.aliyun.com/api/cs/createtriggerhook.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateTriggerHookWithChan(request *CreateTriggerHookRequest) (<-chan *CreateTriggerHookResponse, <-chan error) {
	responseChan := make(chan *CreateTriggerHookResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateTriggerHook(request)
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

// CreateTriggerHookWithCallback invokes the cs.CreateTriggerHook API asynchronously
// api document: https://help.aliyun.com/api/cs/createtriggerhook.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateTriggerHookWithCallback(request *CreateTriggerHookRequest, callback func(response *CreateTriggerHookResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateTriggerHookResponse
		var err error
		defer close(result)
		response, err = client.CreateTriggerHook(request)
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

// CreateTriggerHookRequest is the request struct for api CreateTriggerHook
type CreateTriggerHookRequest struct {
	*requests.RoaRequest
}

// CreateTriggerHookResponse is the response struct for api CreateTriggerHook
type CreateTriggerHookResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateTriggerHookRequest creates a request to invoke CreateTriggerHook API
func CreateCreateTriggerHookRequest() (request *CreateTriggerHookRequest) {
	request = &CreateTriggerHookRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("CS", "2015-12-15", "CreateTriggerHook", "/hook/trigger", "", "")
	request.Method = requests.PUT
	return
}

// CreateCreateTriggerHookResponse creates a response to parse from CreateTriggerHook response
func CreateCreateTriggerHookResponse() (response *CreateTriggerHookResponse) {
	response = &CreateTriggerHookResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
