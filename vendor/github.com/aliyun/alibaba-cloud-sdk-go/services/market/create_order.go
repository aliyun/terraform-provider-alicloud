package market

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

// CreateOrder invokes the market.CreateOrder API synchronously
// api document: https://help.aliyun.com/api/market/createorder.html
func (client *Client) CreateOrder(request *CreateOrderRequest) (response *CreateOrderResponse, err error) {
	response = CreateCreateOrderResponse()
	err = client.DoAction(request, response)
	return
}

// CreateOrderWithChan invokes the market.CreateOrder API asynchronously
// api document: https://help.aliyun.com/api/market/createorder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateOrderWithChan(request *CreateOrderRequest) (<-chan *CreateOrderResponse, <-chan error) {
	responseChan := make(chan *CreateOrderResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateOrder(request)
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

// CreateOrderWithCallback invokes the market.CreateOrder API asynchronously
// api document: https://help.aliyun.com/api/market/createorder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateOrderWithCallback(request *CreateOrderRequest, callback func(response *CreateOrderResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateOrderResponse
		var err error
		defer close(result)
		response, err = client.CreateOrder(request)
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

// CreateOrderRequest is the request struct for api CreateOrder
type CreateOrderRequest struct {
	*requests.RpcRequest
	OrderSouce  string `position:"Query" name:"OrderSouce"`
	Commodity   string `position:"Query" name:"Commodity"`
	ClientToken string `position:"Query" name:"ClientToken"`
	OwnerId     string `position:"Query" name:"OwnerId"`
	PaymentType string `position:"Query" name:"PaymentType"`
	OrderType   string `position:"Query" name:"OrderType"`
}

// CreateOrderResponse is the response struct for api CreateOrder
type CreateOrderResponse struct {
	*responses.BaseResponse
	RequestId   string                   `json:"RequestId" xml:"RequestId"`
	OrderId     string                   `json:"OrderId" xml:"OrderId"`
	InstanceIds InstanceIdsInCreateOrder `json:"InstanceIds" xml:"InstanceIds"`
}

// CreateCreateOrderRequest creates a request to invoke CreateOrder API
func CreateCreateOrderRequest() (request *CreateOrderRequest) {
	request = &CreateOrderRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Market", "2015-11-01", "CreateOrder", "yunmarket", "openAPI")
	return
}

// CreateCreateOrderResponse creates a response to parse from CreateOrder response
func CreateCreateOrderResponse() (response *CreateOrderResponse) {
	response = &CreateOrderResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
