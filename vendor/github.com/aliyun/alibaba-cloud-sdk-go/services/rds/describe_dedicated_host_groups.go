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

// DescribeDedicatedHostGroups invokes the rds.DescribeDedicatedHostGroups API synchronously
// api document: https://help.aliyun.com/api/rds/describededicatedhostgroups.html
func (client *Client) DescribeDedicatedHostGroups(request *DescribeDedicatedHostGroupsRequest) (response *DescribeDedicatedHostGroupsResponse, err error) {
	response = CreateDescribeDedicatedHostGroupsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDedicatedHostGroupsWithChan invokes the rds.DescribeDedicatedHostGroups API asynchronously
// api document: https://help.aliyun.com/api/rds/describededicatedhostgroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDedicatedHostGroupsWithChan(request *DescribeDedicatedHostGroupsRequest) (<-chan *DescribeDedicatedHostGroupsResponse, <-chan error) {
	responseChan := make(chan *DescribeDedicatedHostGroupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDedicatedHostGroups(request)
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

// DescribeDedicatedHostGroupsWithCallback invokes the rds.DescribeDedicatedHostGroups API asynchronously
// api document: https://help.aliyun.com/api/rds/describededicatedhostgroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDedicatedHostGroupsWithCallback(request *DescribeDedicatedHostGroupsRequest, callback func(response *DescribeDedicatedHostGroupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDedicatedHostGroupsResponse
		var err error
		defer close(result)
		response, err = client.DescribeDedicatedHostGroups(request)
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

// DescribeDedicatedHostGroupsRequest is the request struct for api DescribeDedicatedHostGroups
type DescribeDedicatedHostGroupsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	DedicatedHostGroupId string           `position:"Query" name:"DedicatedHostGroupId"`
}

// DescribeDedicatedHostGroupsResponse is the response struct for api DescribeDedicatedHostGroups
type DescribeDedicatedHostGroupsResponse struct {
	*responses.BaseResponse
	RequestId           string              `json:"RequestId" xml:"RequestId"`
	DedicatedHostGroups DedicatedHostGroups `json:"DedicatedHostGroups" xml:"DedicatedHostGroups"`
}

// CreateDescribeDedicatedHostGroupsRequest creates a request to invoke DescribeDedicatedHostGroups API
func CreateDescribeDedicatedHostGroupsRequest() (request *DescribeDedicatedHostGroupsRequest) {
	request = &DescribeDedicatedHostGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeDedicatedHostGroups", "", "")
	return
}

// CreateDescribeDedicatedHostGroupsResponse creates a response to parse from DescribeDedicatedHostGroups response
func CreateDescribeDedicatedHostGroupsResponse() (response *DescribeDedicatedHostGroupsResponse) {
	response = &DescribeDedicatedHostGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
