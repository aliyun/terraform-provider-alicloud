package drds

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

// DrdsApiValidateOrder invokes the drds.DrdsApiValidateOrder API synchronously
func (client *Client) DrdsApiValidateOrder(request *DrdsApiValidateOrderRequest) (response *DrdsApiValidateOrderResponse, err error) {
	response = CreateDrdsApiValidateOrderResponse()
	err = client.DoAction(request, response)
	return
}

// DrdsApiValidateOrderWithChan invokes the drds.DrdsApiValidateOrder API asynchronously
func (client *Client) DrdsApiValidateOrderWithChan(request *DrdsApiValidateOrderRequest) (<-chan *DrdsApiValidateOrderResponse, <-chan error) {
	responseChan := make(chan *DrdsApiValidateOrderResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DrdsApiValidateOrder(request)
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

// DrdsApiValidateOrderWithCallback invokes the drds.DrdsApiValidateOrder API asynchronously
func (client *Client) DrdsApiValidateOrderWithCallback(request *DrdsApiValidateOrderRequest, callback func(response *DrdsApiValidateOrderResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DrdsApiValidateOrderResponse
		var err error
		defer close(result)
		response, err = client.DrdsApiValidateOrder(request)
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

// DrdsApiValidateOrderRequest is the request struct for api DrdsApiValidateOrder
type DrdsApiValidateOrderRequest struct {
	*requests.RpcRequest
	Data string `position:"Query" name:"data"`
}

// DrdsApiValidateOrderResponse is the response struct for api DrdsApiValidateOrder
type DrdsApiValidateOrderResponse struct {
	*responses.BaseResponse
	RequestId string `json:"requestId" xml:"requestId"`
	Code      string `json:"code" xml:"code"`
	Msg       string `json:"msg" xml:"msg"`
	Data      bool   `json:"data" xml:"data"`
	Success   bool   `json:"success" xml:"success"`
}

// CreateDrdsApiValidateOrderRequest creates a request to invoke DrdsApiValidateOrder API
func CreateDrdsApiValidateOrderRequest() (request *DrdsApiValidateOrderRequest) {
	request = &DrdsApiValidateOrderRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2015-04-13", "DrdsApiValidateOrder", "Drds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDrdsApiValidateOrderResponse creates a response to parse from DrdsApiValidateOrder response
func CreateDrdsApiValidateOrderResponse() (response *DrdsApiValidateOrderResponse) {
	response = &DrdsApiValidateOrderResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
