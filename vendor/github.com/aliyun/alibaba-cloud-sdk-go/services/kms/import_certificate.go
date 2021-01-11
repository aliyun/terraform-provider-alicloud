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

// ImportCertificate invokes the kms.ImportCertificate API synchronously
func (client *Client) ImportCertificate(request *ImportCertificateRequest) (response *ImportCertificateResponse, err error) {
	response = CreateImportCertificateResponse()
	err = client.DoAction(request, response)
	return
}

// ImportCertificateWithChan invokes the kms.ImportCertificate API asynchronously
func (client *Client) ImportCertificateWithChan(request *ImportCertificateRequest) (<-chan *ImportCertificateResponse, <-chan error) {
	responseChan := make(chan *ImportCertificateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ImportCertificate(request)
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

// ImportCertificateWithCallback invokes the kms.ImportCertificate API asynchronously
func (client *Client) ImportCertificateWithCallback(request *ImportCertificateRequest, callback func(response *ImportCertificateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ImportCertificateResponse
		var err error
		defer close(result)
		response, err = client.ImportCertificate(request)
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

// ImportCertificateRequest is the request struct for api ImportCertificate
type ImportCertificateRequest struct {
	*requests.RpcRequest
	PKCS12Blob string `position:"Query" name:"PKCS12Blob"`
	Passphrase string `position:"Query" name:"Passphrase"`
}

// ImportCertificateResponse is the response struct for api ImportCertificate
type ImportCertificateResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	CertificateId string `json:"CertificateId" xml:"CertificateId"`
	Arn           string `json:"Arn" xml:"Arn"`
}

// CreateImportCertificateRequest creates a request to invoke ImportCertificate API
func CreateImportCertificateRequest() (request *ImportCertificateRequest) {
	request = &ImportCertificateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "ImportCertificate", "kms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateImportCertificateResponse creates a response to parse from ImportCertificate response
func CreateImportCertificateResponse() (response *ImportCertificateResponse) {
	response = &ImportCertificateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
