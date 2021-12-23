package adb

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

// ModifyAutoRenewAttribute invokes the adb.ModifyAutoRenewAttribute API synchronously
func (client *Client) ModifyAutoRenewAttribute(request *ModifyAutoRenewAttributeRequest) (response *ModifyAutoRenewAttributeResponse, err error) {
	response = CreateModifyAutoRenewAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyAutoRenewAttributeWithChan invokes the adb.ModifyAutoRenewAttribute API asynchronously
func (client *Client) ModifyAutoRenewAttributeWithChan(request *ModifyAutoRenewAttributeRequest) (<-chan *ModifyAutoRenewAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyAutoRenewAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyAutoRenewAttribute(request)
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

// ModifyAutoRenewAttributeWithCallback invokes the adb.ModifyAutoRenewAttribute API asynchronously
func (client *Client) ModifyAutoRenewAttributeWithCallback(request *ModifyAutoRenewAttributeRequest, callback func(response *ModifyAutoRenewAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyAutoRenewAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyAutoRenewAttribute(request)
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

// ModifyAutoRenewAttributeRequest is the request struct for api ModifyAutoRenewAttribute
type ModifyAutoRenewAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Duration             string           `position:"Query" name:"Duration"`
	RenewalStatus        string           `position:"Query" name:"RenewalStatus"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	PeriodUnit           string           `position:"Query" name:"PeriodUnit"`
}

// ModifyAutoRenewAttributeResponse is the response struct for api ModifyAutoRenewAttribute
type ModifyAutoRenewAttributeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyAutoRenewAttributeRequest creates a request to invoke ModifyAutoRenewAttribute API
func CreateModifyAutoRenewAttributeRequest() (request *ModifyAutoRenewAttributeRequest) {
	request = &ModifyAutoRenewAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("adb", "2019-03-15", "ModifyAutoRenewAttribute", "ads", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyAutoRenewAttributeResponse creates a response to parse from ModifyAutoRenewAttribute response
func CreateModifyAutoRenewAttributeResponse() (response *ModifyAutoRenewAttributeResponse) {
	response = &ModifyAutoRenewAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
