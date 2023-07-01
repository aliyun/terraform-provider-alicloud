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

// QueryAvailableInstances invokes the bssopenapi.QueryAvailableInstances API synchronously
func (client *Client) QueryAvailableInstances(request *QueryAvailableInstancesRequest) (response *QueryAvailableInstancesResponse, err error) {
	response = CreateQueryAvailableInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// QueryAvailableInstancesWithChan invokes the bssopenapi.QueryAvailableInstances API asynchronously
func (client *Client) QueryAvailableInstancesWithChan(request *QueryAvailableInstancesRequest) (<-chan *QueryAvailableInstancesResponse, <-chan error) {
	responseChan := make(chan *QueryAvailableInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryAvailableInstances(request)
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

// QueryAvailableInstancesWithCallback invokes the bssopenapi.QueryAvailableInstances API asynchronously
func (client *Client) QueryAvailableInstancesWithCallback(request *QueryAvailableInstancesRequest, callback func(response *QueryAvailableInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryAvailableInstancesResponse
		var err error
		defer close(result)
		response, err = client.QueryAvailableInstances(request)
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

// QueryAvailableInstancesRequest is the request struct for api QueryAvailableInstances
type QueryAvailableInstancesRequest struct {
	*requests.RpcRequest
	ProductCode      string           `position:"Query" name:"ProductCode"`
	SubscriptionType string           `position:"Query" name:"SubscriptionType"`
	PageNum          requests.Integer `position:"Query" name:"PageNum"`
	ProductType      string           `position:"Query" name:"ProductType"`
	PageSize         requests.Integer `position:"Query" name:"PageSize"`
	RenewStatus      string           `position:"Query" name:"RenewStatus"`
	OwnerId          requests.Integer `position:"Query" name:"OwnerId"`
	EndTimeStart     string           `position:"Query" name:"EndTimeStart"`
	CreateTimeEnd    string           `position:"Query" name:"CreateTimeEnd"`
	InstanceIDs      string           `position:"Query" name:"InstanceIDs"`
	EndTimeEnd       string           `position:"Query" name:"EndTimeEnd"`
	CreateTimeStart  string           `position:"Query" name:"CreateTimeStart"`
	Region           string           `position:"Query" name:"Region"`
}

// QueryAvailableInstancesResponse is the response struct for api QueryAvailableInstances
type QueryAvailableInstancesResponse struct {
	*responses.BaseResponse
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateQueryAvailableInstancesRequest creates a request to invoke QueryAvailableInstances API
func CreateQueryAvailableInstancesRequest() (request *QueryAvailableInstancesRequest) {
	request = &QueryAvailableInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "QueryAvailableInstances", "bssopenapi", "openAPI")
	request.Method = requests.POST
	return
}

// CreateQueryAvailableInstancesResponse creates a response to parse from QueryAvailableInstances response
func CreateQueryAvailableInstancesResponse() (response *QueryAvailableInstancesResponse) {
	response = &QueryAvailableInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
