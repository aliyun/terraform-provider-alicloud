package rds

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

// DescribeParameterGroups invokes the rds.DescribeParameterGroups API synchronously
// api document: https://help.aliyun.com/api/rds/describeparametergroups.html
func (client *Client) DescribeParameterGroups(request *DescribeParameterGroupsRequest) (response *DescribeParameterGroupsResponse, err error) {
	response = CreateDescribeParameterGroupsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeParameterGroupsWithChan invokes the rds.DescribeParameterGroups API asynchronously
// api document: https://help.aliyun.com/api/rds/describeparametergroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeParameterGroupsWithChan(request *DescribeParameterGroupsRequest) (<-chan *DescribeParameterGroupsResponse, <-chan error) {
	responseChan := make(chan *DescribeParameterGroupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeParameterGroups(request)
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

// DescribeParameterGroupsWithCallback invokes the rds.DescribeParameterGroups API asynchronously
// api document: https://help.aliyun.com/api/rds/describeparametergroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeParameterGroupsWithCallback(request *DescribeParameterGroupsRequest, callback func(response *DescribeParameterGroupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeParameterGroupsResponse
		var err error
		defer close(result)
		response, err = client.DescribeParameterGroups(request)
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

// DescribeParameterGroupsRequest is the request struct for api DescribeParameterGroups
type DescribeParameterGroupsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeParameterGroupsResponse is the response struct for api DescribeParameterGroups
type DescribeParameterGroupsResponse struct {
	*responses.BaseResponse
	RequestId       string          `json:"RequestId" xml:"RequestId"`
	ParameterGroups ParameterGroups `json:"ParameterGroups" xml:"ParameterGroups"`
}

// CreateDescribeParameterGroupsRequest creates a request to invoke DescribeParameterGroups API
func CreateDescribeParameterGroupsRequest() (request *DescribeParameterGroupsRequest) {
	request = &DescribeParameterGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeParameterGroups", "rds", "openAPI")
	return
}

// CreateDescribeParameterGroupsResponse creates a response to parse from DescribeParameterGroups response
func CreateDescribeParameterGroupsResponse() (response *DescribeParameterGroupsResponse) {
	response = &DescribeParameterGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
