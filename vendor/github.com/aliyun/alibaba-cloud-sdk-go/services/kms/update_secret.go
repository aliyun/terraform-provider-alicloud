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

// UpdateSecret invokes the kms.UpdateSecret API synchronously
func (client *Client) UpdateSecret(request *UpdateSecretRequest) (response *UpdateSecretResponse, err error) {
	response = CreateUpdateSecretResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateSecretWithChan invokes the kms.UpdateSecret API asynchronously
func (client *Client) UpdateSecretWithChan(request *UpdateSecretRequest) (<-chan *UpdateSecretResponse, <-chan error) {
	responseChan := make(chan *UpdateSecretResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateSecret(request)
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

// UpdateSecretWithCallback invokes the kms.UpdateSecret API asynchronously
func (client *Client) UpdateSecretWithCallback(request *UpdateSecretRequest, callback func(response *UpdateSecretResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateSecretResponse
		var err error
		defer close(result)
		response, err = client.UpdateSecret(request)
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

// UpdateSecretRequest is the request struct for api UpdateSecret
type UpdateSecretRequest struct {
	*requests.RpcRequest
	Description string `position:"Query" name:"Description"`
	SecretName  string `position:"Query" name:"SecretName"`
}

// UpdateSecretResponse is the response struct for api UpdateSecret
type UpdateSecretResponse struct {
	*responses.BaseResponse
	RequestId  string `json:"RequestId" xml:"RequestId"`
	SecretName string `json:"SecretName" xml:"SecretName"`
}

// CreateUpdateSecretRequest creates a request to invoke UpdateSecret API
func CreateUpdateSecretRequest() (request *UpdateSecretRequest) {
	request = &UpdateSecretRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "UpdateSecret", "kms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUpdateSecretResponse creates a response to parse from UpdateSecret response
func CreateUpdateSecretResponse() (response *UpdateSecretResponse) {
	response = &UpdateSecretResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
