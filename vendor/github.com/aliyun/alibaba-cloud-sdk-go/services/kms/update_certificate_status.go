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

// UpdateCertificateStatus invokes the kms.UpdateCertificateStatus API synchronously
func (client *Client) UpdateCertificateStatus(request *UpdateCertificateStatusRequest) (response *UpdateCertificateStatusResponse, err error) {
	response = CreateUpdateCertificateStatusResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateCertificateStatusWithChan invokes the kms.UpdateCertificateStatus API asynchronously
func (client *Client) UpdateCertificateStatusWithChan(request *UpdateCertificateStatusRequest) (<-chan *UpdateCertificateStatusResponse, <-chan error) {
	responseChan := make(chan *UpdateCertificateStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateCertificateStatus(request)
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

// UpdateCertificateStatusWithCallback invokes the kms.UpdateCertificateStatus API asynchronously
func (client *Client) UpdateCertificateStatusWithCallback(request *UpdateCertificateStatusRequest, callback func(response *UpdateCertificateStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateCertificateStatusResponse
		var err error
		defer close(result)
		response, err = client.UpdateCertificateStatus(request)
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

// UpdateCertificateStatusRequest is the request struct for api UpdateCertificateStatus
type UpdateCertificateStatusRequest struct {
	*requests.RpcRequest
	CertificateId string `position:"Query" name:"CertificateId"`
	Status        string `position:"Query" name:"Status"`
}

// UpdateCertificateStatusResponse is the response struct for api UpdateCertificateStatus
type UpdateCertificateStatusResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateCertificateStatusRequest creates a request to invoke UpdateCertificateStatus API
func CreateUpdateCertificateStatusRequest() (request *UpdateCertificateStatusRequest) {
	request = &UpdateCertificateStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "UpdateCertificateStatus", "kms-service", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUpdateCertificateStatusResponse creates a response to parse from UpdateCertificateStatus response
func CreateUpdateCertificateStatusResponse() (response *UpdateCertificateStatusResponse) {
	response = &UpdateCertificateStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
