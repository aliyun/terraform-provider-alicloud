package cloudapi

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

// CreateAccessControlList invokes the cloudapi.CreateAccessControlList API synchronously
func (client *Client) CreateAccessControlList(request *CreateAccessControlListRequest) (response *CreateAccessControlListResponse, err error) {
	response = CreateCreateAccessControlListResponse()
	err = client.DoAction(request, response)
	return
}

// CreateAccessControlListWithChan invokes the cloudapi.CreateAccessControlList API asynchronously
func (client *Client) CreateAccessControlListWithChan(request *CreateAccessControlListRequest) (<-chan *CreateAccessControlListResponse, <-chan error) {
	responseChan := make(chan *CreateAccessControlListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateAccessControlList(request)
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

// CreateAccessControlListWithCallback invokes the cloudapi.CreateAccessControlList API asynchronously
func (client *Client) CreateAccessControlListWithCallback(request *CreateAccessControlListRequest, callback func(response *CreateAccessControlListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateAccessControlListResponse
		var err error
		defer close(result)
		response, err = client.CreateAccessControlList(request)
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

// CreateAccessControlListRequest is the request struct for api CreateAccessControlList
type CreateAccessControlListRequest struct {
	*requests.RpcRequest
	AclName          string `position:"Query" name:"AclName"`
	AddressIPVersion string `position:"Query" name:"AddressIPVersion"`
	SecurityToken    string `position:"Query" name:"SecurityToken"`
}

// CreateAccessControlListResponse is the response struct for api CreateAccessControlList
type CreateAccessControlListResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateAccessControlListRequest creates a request to invoke CreateAccessControlList API
func CreateCreateAccessControlListRequest() (request *CreateAccessControlListRequest) {
	request = &CreateAccessControlListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "CreateAccessControlList", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateAccessControlListResponse creates a response to parse from CreateAccessControlList response
func CreateCreateAccessControlListResponse() (response *CreateAccessControlListResponse) {
	response = &CreateAccessControlListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
