package bssopenapi

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

// QueryOrders invokes the bssopenapi.QueryOrders API synchronously
// api document: https://help.aliyun.com/api/bssopenapi/queryorders.html
func (client *Client) QueryOrders(request *QueryOrdersRequest) (response *QueryOrdersResponse, err error) {
	response = CreateQueryOrdersResponse()
	err = client.DoAction(request, response)
	return
}

// QueryOrdersWithChan invokes the bssopenapi.QueryOrders API asynchronously
// api document: https://help.aliyun.com/api/bssopenapi/queryorders.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryOrdersWithChan(request *QueryOrdersRequest) (<-chan *QueryOrdersResponse, <-chan error) {
	responseChan := make(chan *QueryOrdersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryOrders(request)
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

// QueryOrdersWithCallback invokes the bssopenapi.QueryOrders API asynchronously
// api document: https://help.aliyun.com/api/bssopenapi/queryorders.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryOrdersWithCallback(request *QueryOrdersRequest, callback func(response *QueryOrdersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryOrdersResponse
		var err error
		defer close(result)
		response, err = client.QueryOrders(request)
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

// QueryOrdersRequest is the request struct for api QueryOrders
type QueryOrdersRequest struct {
	*requests.RpcRequest
	ProductCode      string           `position:"Query" name:"ProductCode"`
	SubscriptionType string           `position:"Query" name:"SubscriptionType"`
	PageSize         requests.Integer `position:"Query" name:"PageSize"`
	PaymentStatus    string           `position:"Query" name:"PaymentStatus"`
	CreateTimeStart  string           `position:"Query" name:"CreateTimeStart"`
	PageNum          requests.Integer `position:"Query" name:"PageNum"`
	OwnerId          requests.Integer `position:"Query" name:"OwnerId"`
	CreateTimeEnd    string           `position:"Query" name:"CreateTimeEnd"`
	ProductType      string           `position:"Query" name:"ProductType"`
	OrderType        string           `position:"Query" name:"OrderType"`
}

// QueryOrdersResponse is the response struct for api QueryOrders
type QueryOrdersResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateQueryOrdersRequest creates a request to invoke QueryOrders API
func CreateQueryOrdersRequest() (request *QueryOrdersRequest) {
	request = &QueryOrdersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "QueryOrders", "", "")
	return
}

// CreateQueryOrdersResponse creates a response to parse from QueryOrders response
func CreateQueryOrdersResponse() (response *QueryOrdersResponse) {
	response = &QueryOrdersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
