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

// DescribeSceneDefensePolicies invokes the ddoscoo.DescribeSceneDefensePolicies API synchronously
func (client *Client) DescribeSceneDefensePolicies(request *DescribeSceneDefensePoliciesRequest) (response *DescribeSceneDefensePoliciesResponse, err error) {
	response = CreateDescribeSceneDefensePoliciesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSceneDefensePoliciesWithChan invokes the ddoscoo.DescribeSceneDefensePolicies API asynchronously
func (client *Client) DescribeSceneDefensePoliciesWithChan(request *DescribeSceneDefensePoliciesRequest) (<-chan *DescribeSceneDefensePoliciesResponse, <-chan error) {
	responseChan := make(chan *DescribeSceneDefensePoliciesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSceneDefensePolicies(request)
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

// DescribeSceneDefensePoliciesWithCallback invokes the ddoscoo.DescribeSceneDefensePolicies API asynchronously
func (client *Client) DescribeSceneDefensePoliciesWithCallback(request *DescribeSceneDefensePoliciesRequest, callback func(response *DescribeSceneDefensePoliciesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSceneDefensePoliciesResponse
		var err error
		defer close(result)
		response, err = client.DescribeSceneDefensePolicies(request)
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

// DescribeSceneDefensePoliciesRequest is the request struct for api DescribeSceneDefensePolicies
type DescribeSceneDefensePoliciesRequest struct {
	*requests.RpcRequest
	Template        string `position:"Query" name:"Template"`
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
	Status          string `position:"Query" name:"Status"`
}

// DescribeSceneDefensePoliciesResponse is the response struct for api DescribeSceneDefensePolicies
type DescribeSceneDefensePoliciesResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Success   bool     `json:"Success" xml:"Success"`
	Policies  []Policy `json:"Policies" xml:"Policies"`
}

// CreateDescribeSceneDefensePoliciesRequest creates a request to invoke DescribeSceneDefensePolicies API
func CreateDescribeSceneDefensePoliciesRequest() (request *DescribeSceneDefensePoliciesRequest) {
	request = &DescribeSceneDefensePoliciesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeSceneDefensePolicies", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeSceneDefensePoliciesResponse creates a response to parse from DescribeSceneDefensePolicies response
func CreateDescribeSceneDefensePoliciesResponse() (response *DescribeSceneDefensePoliciesResponse) {
	response = &DescribeSceneDefensePoliciesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
