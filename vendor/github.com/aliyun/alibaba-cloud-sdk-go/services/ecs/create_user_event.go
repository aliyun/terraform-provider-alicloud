package ecs

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

// CreateUserEvent invokes the ecs.CreateUserEvent API synchronously
// api document: https://help.aliyun.com/api/ecs/createuserevent.html
func (client *Client) CreateUserEvent(request *CreateUserEventRequest) (response *CreateUserEventResponse, err error) {
	response = CreateCreateUserEventResponse()
	err = client.DoAction(request, response)
	return
}

// CreateUserEventWithChan invokes the ecs.CreateUserEvent API asynchronously
// api document: https://help.aliyun.com/api/ecs/createuserevent.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateUserEventWithChan(request *CreateUserEventRequest) (<-chan *CreateUserEventResponse, <-chan error) {
	responseChan := make(chan *CreateUserEventResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateUserEvent(request)
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

// CreateUserEventWithCallback invokes the ecs.CreateUserEvent API asynchronously
// api document: https://help.aliyun.com/api/ecs/createuserevent.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateUserEventWithCallback(request *CreateUserEventRequest, callback func(response *CreateUserEventResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateUserEventResponse
		var err error
		defer close(result)
		response, err = client.CreateUserEvent(request)
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

// CreateUserEventRequest is the request struct for api CreateUserEvent
type CreateUserEventRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PlanTime             string           `position:"Query" name:"PlanTime"`
	ExpireTime           string           `position:"Query" name:"ExpireTime"`
	ResourceId           string           `position:"Query" name:"ResourceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	EventType            string           `position:"Query" name:"EventType"`
}

// CreateUserEventResponse is the response struct for api CreateUserEvent
type CreateUserEventResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	EventId   string `json:"EventId" xml:"EventId"`
}

// CreateCreateUserEventRequest creates a request to invoke CreateUserEvent API
func CreateCreateUserEventRequest() (request *CreateUserEventRequest) {
	request = &CreateUserEventRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "CreateUserEvent", "ecs", "openAPI")
	return
}

// CreateCreateUserEventResponse creates a response to parse from CreateUserEvent response
func CreateCreateUserEventResponse() (response *CreateUserEventResponse) {
	response = &CreateUserEventResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
