package hbase

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

// DescribeSecurityGroups invokes the hbase.DescribeSecurityGroups API synchronously
// api document: https://help.aliyun.com/api/hbase/describesecuritygroups.html
func (client *Client) DescribeSecurityGroups(request *DescribeSecurityGroupsRequest) (response *DescribeSecurityGroupsResponse, err error) {
	response = CreateDescribeSecurityGroupsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSecurityGroupsWithChan invokes the hbase.DescribeSecurityGroups API asynchronously
// api document: https://help.aliyun.com/api/hbase/describesecuritygroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSecurityGroupsWithChan(request *DescribeSecurityGroupsRequest) (<-chan *DescribeSecurityGroupsResponse, <-chan error) {
	responseChan := make(chan *DescribeSecurityGroupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSecurityGroups(request)
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

// DescribeSecurityGroupsWithCallback invokes the hbase.DescribeSecurityGroups API asynchronously
// api document: https://help.aliyun.com/api/hbase/describesecuritygroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSecurityGroupsWithCallback(request *DescribeSecurityGroupsRequest, callback func(response *DescribeSecurityGroupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSecurityGroupsResponse
		var err error
		defer close(result)
		response, err = client.DescribeSecurityGroups(request)
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

// DescribeSecurityGroupsRequest is the request struct for api DescribeSecurityGroups
type DescribeSecurityGroupsRequest struct {
	*requests.RpcRequest
	ClusterId string `position:"Query" name:"ClusterId"`
}

// DescribeSecurityGroupsResponse is the response struct for api DescribeSecurityGroups
type DescribeSecurityGroupsResponse struct {
	*responses.BaseResponse
	RequestId        string           `json:"RequestId" xml:"RequestId"`
	SecurityGroupIds SecurityGroupIds `json:"SecurityGroupIds" xml:"SecurityGroupIds"`
}

// CreateDescribeSecurityGroupsRequest creates a request to invoke DescribeSecurityGroups API
func CreateDescribeSecurityGroupsRequest() (request *DescribeSecurityGroupsRequest) {
	request = &DescribeSecurityGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("HBase", "2019-01-01", "DescribeSecurityGroups", "", "")
	return
}

// CreateDescribeSecurityGroupsResponse creates a response to parse from DescribeSecurityGroups response
func CreateDescribeSecurityGroupsResponse() (response *DescribeSecurityGroupsResponse) {
	response = &DescribeSecurityGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
