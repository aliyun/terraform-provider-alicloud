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

// ModifyAutoSnapshotPolicy invokes the ecs.ModifyAutoSnapshotPolicy API synchronously
func (client *Client) ModifyAutoSnapshotPolicy(request *ModifyAutoSnapshotPolicyRequest) (response *ModifyAutoSnapshotPolicyResponse, err error) {
	response = CreateModifyAutoSnapshotPolicyResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyAutoSnapshotPolicyWithChan invokes the ecs.ModifyAutoSnapshotPolicy API asynchronously
func (client *Client) ModifyAutoSnapshotPolicyWithChan(request *ModifyAutoSnapshotPolicyRequest) (<-chan *ModifyAutoSnapshotPolicyResponse, <-chan error) {
	responseChan := make(chan *ModifyAutoSnapshotPolicyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyAutoSnapshotPolicy(request)
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

// ModifyAutoSnapshotPolicyWithCallback invokes the ecs.ModifyAutoSnapshotPolicy API asynchronously
func (client *Client) ModifyAutoSnapshotPolicyWithCallback(request *ModifyAutoSnapshotPolicyRequest, callback func(response *ModifyAutoSnapshotPolicyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyAutoSnapshotPolicyResponse
		var err error
		defer close(result)
		response, err = client.ModifyAutoSnapshotPolicy(request)
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

// ModifyAutoSnapshotPolicyRequest is the request struct for api ModifyAutoSnapshotPolicy
type ModifyAutoSnapshotPolicyRequest struct {
	*requests.RpcRequest
	DataDiskPolicyEnabled             requests.Boolean `position:"Query" name:"DataDiskPolicyEnabled"`
	ResourceOwnerId                   requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DataDiskPolicyRetentionDays       requests.Integer `position:"Query" name:"DataDiskPolicyRetentionDays"`
	SystemDiskPolicyRetentionLastWeek requests.Boolean `position:"Query" name:"SystemDiskPolicyRetentionLastWeek"`
	SystemDiskPolicyRetentionDays     requests.Integer `position:"Query" name:"SystemDiskPolicyRetentionDays"`
	DataDiskPolicyTimePeriod          requests.Integer `position:"Query" name:"DataDiskPolicyTimePeriod"`
	ResourceOwnerAccount              string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                      string           `position:"Query" name:"OwnerAccount"`
	SystemDiskPolicyTimePeriod        requests.Integer `position:"Query" name:"SystemDiskPolicyTimePeriod"`
	OwnerId                           requests.Integer `position:"Query" name:"OwnerId"`
	DataDiskPolicyRetentionLastWeek   requests.Boolean `position:"Query" name:"DataDiskPolicyRetentionLastWeek"`
	SystemDiskPolicyEnabled           requests.Boolean `position:"Query" name:"SystemDiskPolicyEnabled"`
}

// ModifyAutoSnapshotPolicyResponse is the response struct for api ModifyAutoSnapshotPolicy
type ModifyAutoSnapshotPolicyResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyAutoSnapshotPolicyRequest creates a request to invoke ModifyAutoSnapshotPolicy API
func CreateModifyAutoSnapshotPolicyRequest() (request *ModifyAutoSnapshotPolicyRequest) {
	request = &ModifyAutoSnapshotPolicyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "ModifyAutoSnapshotPolicy", "", "")
	request.Method = requests.POST
	return
}

// CreateModifyAutoSnapshotPolicyResponse creates a response to parse from ModifyAutoSnapshotPolicy response
func CreateModifyAutoSnapshotPolicyResponse() (response *ModifyAutoSnapshotPolicyResponse) {
	response = &ModifyAutoSnapshotPolicyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
