package kms

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

// DisableKey invokes the kms.DisableKey API synchronously
// api document: https://help.aliyun.com/api/kms/disablekey.html
func (client *Client) DisableKey(request *DisableKeyRequest) (response *DisableKeyResponse, err error) {
	response = CreateDisableKeyResponse()
	err = client.DoAction(request, response)
	return
}

// DisableKeyWithChan invokes the kms.DisableKey API asynchronously
// api document: https://help.aliyun.com/api/kms/disablekey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DisableKeyWithChan(request *DisableKeyRequest) (<-chan *DisableKeyResponse, <-chan error) {
	responseChan := make(chan *DisableKeyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DisableKey(request)
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

// DisableKeyWithCallback invokes the kms.DisableKey API asynchronously
// api document: https://help.aliyun.com/api/kms/disablekey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DisableKeyWithCallback(request *DisableKeyRequest, callback func(response *DisableKeyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DisableKeyResponse
		var err error
		defer close(result)
		response, err = client.DisableKey(request)
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

// DisableKeyRequest is the request struct for api DisableKey
type DisableKeyRequest struct {
	*requests.RpcRequest
	KeyId string `position:"Query" name:"KeyId"`
}

// DisableKeyResponse is the response struct for api DisableKey
type DisableKeyResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDisableKeyRequest creates a request to invoke DisableKey API
func CreateDisableKeyRequest() (request *DisableKeyRequest) {
	request = &DisableKeyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "DisableKey", "kms-service", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDisableKeyResponse creates a response to parse from DisableKey response
func CreateDisableKeyResponse() (response *DisableKeyResponse) {
	response = &DisableKeyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
