package smartag

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

// DescribeNetworkOptimizationSettings invokes the smartag.DescribeNetworkOptimizationSettings API synchronously
func (client *Client) DescribeNetworkOptimizationSettings(request *DescribeNetworkOptimizationSettingsRequest) (response *DescribeNetworkOptimizationSettingsResponse, err error) {
	response = CreateDescribeNetworkOptimizationSettingsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeNetworkOptimizationSettingsWithChan invokes the smartag.DescribeNetworkOptimizationSettings API asynchronously
func (client *Client) DescribeNetworkOptimizationSettingsWithChan(request *DescribeNetworkOptimizationSettingsRequest) (<-chan *DescribeNetworkOptimizationSettingsResponse, <-chan error) {
	responseChan := make(chan *DescribeNetworkOptimizationSettingsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeNetworkOptimizationSettings(request)
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

// DescribeNetworkOptimizationSettingsWithCallback invokes the smartag.DescribeNetworkOptimizationSettings API asynchronously
func (client *Client) DescribeNetworkOptimizationSettingsWithCallback(request *DescribeNetworkOptimizationSettingsRequest, callback func(response *DescribeNetworkOptimizationSettingsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeNetworkOptimizationSettingsResponse
		var err error
		defer close(result)
		response, err = client.DescribeNetworkOptimizationSettings(request)
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

// DescribeNetworkOptimizationSettingsRequest is the request struct for api DescribeNetworkOptimizationSettings
type DescribeNetworkOptimizationSettingsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query"`
	PageNumber           requests.Integer `position:"Query"`
	NetworkOptId         string           `position:"Query"`
	PageSize             requests.Integer `position:"Query"`
	ResourceOwnerAccount string           `position:"Query"`
	OwnerAccount         string           `position:"Query"`
	OwnerId              requests.Integer `position:"Query"`
}

// DescribeNetworkOptimizationSettingsResponse is the response struct for api DescribeNetworkOptimizationSettings
type DescribeNetworkOptimizationSettingsResponse struct {
	*responses.BaseResponse
	TotalCount int      `json:"TotalCount" xml:"TotalCount"`
	RequestId  string   `json:"RequestId" xml:"RequestId"`
	PageSize   int      `json:"PageSize" xml:"PageSize"`
	PageNumber int      `json:"PageNumber" xml:"PageNumber"`
	Settings   Settings `json:"Settings" xml:"Settings"`
}

// CreateDescribeNetworkOptimizationSettingsRequest creates a request to invoke DescribeNetworkOptimizationSettings API
func CreateDescribeNetworkOptimizationSettingsRequest() (request *DescribeNetworkOptimizationSettingsRequest) {
	request = &DescribeNetworkOptimizationSettingsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DescribeNetworkOptimizationSettings", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeNetworkOptimizationSettingsResponse creates a response to parse from DescribeNetworkOptimizationSettings response
func CreateDescribeNetworkOptimizationSettingsResponse() (response *DescribeNetworkOptimizationSettingsResponse) {
	response = &DescribeNetworkOptimizationSettingsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
