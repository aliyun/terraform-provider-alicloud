package ecs

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

// DescribeDisks invokes the ecs.DescribeDisks API synchronously
// api document: https://help.aliyun.com/api/ecs/describedisks.html
func (client *Client) DescribeDisks(request *DescribeDisksRequest) (response *DescribeDisksResponse, err error) {
	response = CreateDescribeDisksResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDisksWithChan invokes the ecs.DescribeDisks API asynchronously
// api document: https://help.aliyun.com/api/ecs/describedisks.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDisksWithChan(request *DescribeDisksRequest) (<-chan *DescribeDisksResponse, <-chan error) {
	responseChan := make(chan *DescribeDisksResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDisks(request)
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

// DescribeDisksWithCallback invokes the ecs.DescribeDisks API asynchronously
// api document: https://help.aliyun.com/api/ecs/describedisks.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDisksWithCallback(request *DescribeDisksRequest, callback func(response *DescribeDisksResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDisksResponse
		var err error
		defer close(result)
		response, err = client.DescribeDisks(request)
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

// DescribeDisksRequest is the request struct for api DescribeDisks
type DescribeDisksRequest struct {
	*requests.RpcRequest
	ResourceOwnerId               requests.Integer    `position:"Query" name:"ResourceOwnerId"`
	Filter2Value                  string              `position:"Query" name:"Filter.2.Value"`
	AutoSnapshotPolicyId          string              `position:"Query" name:"AutoSnapshotPolicyId"`
	DiskName                      string              `position:"Query" name:"DiskName"`
	DeleteAutoSnapshot            requests.Boolean    `position:"Query" name:"DeleteAutoSnapshot"`
	ResourceGroupId               string              `position:"Query" name:"ResourceGroupId"`
	DiskChargeType                string              `position:"Query" name:"DiskChargeType"`
	LockReason                    string              `position:"Query" name:"LockReason"`
	Filter1Key                    string              `position:"Query" name:"Filter.1.Key"`
	Tag                           *[]DescribeDisksTag `position:"Query" name:"Tag"  type:"Repeated"`
	EnableAutoSnapshot            requests.Boolean    `position:"Query" name:"EnableAutoSnapshot"`
	DryRun                        requests.Boolean    `position:"Query" name:"DryRun"`
	Filter1Value                  string              `position:"Query" name:"Filter.1.Value"`
	Portable                      requests.Boolean    `position:"Query" name:"Portable"`
	OwnerId                       requests.Integer    `position:"Query" name:"OwnerId"`
	AdditionalAttributes          *[]string           `position:"Query" name:"AdditionalAttributes"  type:"Repeated"`
	InstanceId                    string              `position:"Query" name:"InstanceId"`
	ZoneId                        string              `position:"Query" name:"ZoneId"`
	Status                        string              `position:"Query" name:"Status"`
	SnapshotId                    string              `position:"Query" name:"SnapshotId"`
	PageNumber                    requests.Integer    `position:"Query" name:"PageNumber"`
	PageSize                      requests.Integer    `position:"Query" name:"PageSize"`
	DiskIds                       string              `position:"Query" name:"DiskIds"`
	DeleteWithInstance            requests.Boolean    `position:"Query" name:"DeleteWithInstance"`
	ResourceOwnerAccount          string              `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                  string              `position:"Query" name:"OwnerAccount"`
	EnableAutomatedSnapshotPolicy requests.Boolean    `position:"Query" name:"EnableAutomatedSnapshotPolicy"`
	Filter2Key                    string              `position:"Query" name:"Filter.2.Key"`
	DiskType                      string              `position:"Query" name:"DiskType"`
	EnableShared                  requests.Boolean    `position:"Query" name:"EnableShared"`
	Encrypted                     requests.Boolean    `position:"Query" name:"Encrypted"`
	Category                      string              `position:"Query" name:"Category"`
	KMSKeyId                      string              `position:"Query" name:"KMSKeyId"`
}

// DescribeDisksTag is a repeated param struct in DescribeDisksRequest
type DescribeDisksTag struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// DescribeDisksResponse is the response struct for api DescribeDisks
type DescribeDisksResponse struct {
	*responses.BaseResponse
	RequestId  string               `json:"RequestId" xml:"RequestId"`
	TotalCount int                  `json:"TotalCount" xml:"TotalCount"`
	PageNumber int                  `json:"PageNumber" xml:"PageNumber"`
	PageSize   int                  `json:"PageSize" xml:"PageSize"`
	Disks      DisksInDescribeDisks `json:"Disks" xml:"Disks"`
}

// CreateDescribeDisksRequest creates a request to invoke DescribeDisks API
func CreateDescribeDisksRequest() (request *DescribeDisksRequest) {
	request = &DescribeDisksRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeDisks", "ecs", "openAPI")
	return
}

// CreateDescribeDisksResponse creates a response to parse from DescribeDisks response
func CreateDescribeDisksResponse() (response *DescribeDisksResponse) {
	response = &DescribeDisksResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
