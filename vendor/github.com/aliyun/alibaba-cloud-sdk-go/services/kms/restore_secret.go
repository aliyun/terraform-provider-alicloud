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

// RestoreSecret invokes the kms.RestoreSecret API synchronously
func (client *Client) RestoreSecret(request *RestoreSecretRequest) (response *RestoreSecretResponse, err error) {
	response = CreateRestoreSecretResponse()
	err = client.DoAction(request, response)
	return
}

// RestoreSecretWithChan invokes the kms.RestoreSecret API asynchronously
func (client *Client) RestoreSecretWithChan(request *RestoreSecretRequest) (<-chan *RestoreSecretResponse, <-chan error) {
	responseChan := make(chan *RestoreSecretResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RestoreSecret(request)
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

// RestoreSecretWithCallback invokes the kms.RestoreSecret API asynchronously
func (client *Client) RestoreSecretWithCallback(request *RestoreSecretRequest, callback func(response *RestoreSecretResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RestoreSecretResponse
		var err error
		defer close(result)
		response, err = client.RestoreSecret(request)
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

// RestoreSecretRequest is the request struct for api RestoreSecret
type RestoreSecretRequest struct {
	*requests.RpcRequest
	SecretName string `position:"Query" name:"SecretName"`
}

// RestoreSecretResponse is the response struct for api RestoreSecret
type RestoreSecretResponse struct {
	*responses.BaseResponse
	RequestId  string `json:"RequestId" xml:"RequestId"`
	SecretName string `json:"SecretName" xml:"SecretName"`
}

// CreateRestoreSecretRequest creates a request to invoke RestoreSecret API
func CreateRestoreSecretRequest() (request *RestoreSecretRequest) {
	request = &RestoreSecretRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "RestoreSecret", "kms-service", "openAPI")
	request.Method = requests.POST
	return
}

// CreateRestoreSecretResponse creates a response to parse from RestoreSecret response
func CreateRestoreSecretResponse() (response *RestoreSecretResponse) {
	response = &RestoreSecretResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
