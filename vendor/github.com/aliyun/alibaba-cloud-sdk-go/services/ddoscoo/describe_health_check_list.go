package ddoscoo

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

// DescribeHealthCheckList invokes the ddoscoo.DescribeHealthCheckList API synchronously
func (client *Client) DescribeHealthCheckList(request *DescribeHealthCheckListRequest) (response *DescribeHealthCheckListResponse, err error) {
	response = CreateDescribeHealthCheckListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeHealthCheckListWithChan invokes the ddoscoo.DescribeHealthCheckList API asynchronously
func (client *Client) DescribeHealthCheckListWithChan(request *DescribeHealthCheckListRequest) (<-chan *DescribeHealthCheckListResponse, <-chan error) {
	responseChan := make(chan *DescribeHealthCheckListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeHealthCheckList(request)
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

// DescribeHealthCheckListWithCallback invokes the ddoscoo.DescribeHealthCheckList API asynchronously
func (client *Client) DescribeHealthCheckListWithCallback(request *DescribeHealthCheckListRequest, callback func(response *DescribeHealthCheckListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeHealthCheckListResponse
		var err error
		defer close(result)
		response, err = client.DescribeHealthCheckList(request)
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

// DescribeHealthCheckListRequest is the request struct for api DescribeHealthCheckList
type DescribeHealthCheckListRequest struct {
	*requests.RpcRequest
	NetworkRules string `position:"Query" name:"NetworkRules"`
	SourceIp     string `position:"Query" name:"SourceIp"`
}

// DescribeHealthCheckListResponse is the response struct for api DescribeHealthCheckList
type DescribeHealthCheckListResponse struct {
	*responses.BaseResponse
	RequestId       string            `json:"RequestId" xml:"RequestId"`
	HealthCheckList []HealthCheckItem `json:"HealthCheckList" xml:"HealthCheckList"`
}

// CreateDescribeHealthCheckListRequest creates a request to invoke DescribeHealthCheckList API
func CreateDescribeHealthCheckListRequest() (request *DescribeHealthCheckListRequest) {
	request = &DescribeHealthCheckListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeHealthCheckList", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeHealthCheckListResponse creates a response to parse from DescribeHealthCheckList response
func CreateDescribeHealthCheckListResponse() (response *DescribeHealthCheckListResponse) {
	response = &DescribeHealthCheckListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
