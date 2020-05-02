package waf_openapi

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

// CreateCertificateByCertificateId invokes the waf_openapi.CreateCertificateByCertificateId API synchronously
// api document: https://help.aliyun.com/api/waf-openapi/createcertificatebycertificateid.html
func (client *Client) CreateCertificateByCertificateId(request *CreateCertificateByCertificateIdRequest) (response *CreateCertificateByCertificateIdResponse, err error) {
	response = CreateCreateCertificateByCertificateIdResponse()
	err = client.DoAction(request, response)
	return
}

// CreateCertificateByCertificateIdWithChan invokes the waf_openapi.CreateCertificateByCertificateId API asynchronously
// api document: https://help.aliyun.com/api/waf-openapi/createcertificatebycertificateid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCertificateByCertificateIdWithChan(request *CreateCertificateByCertificateIdRequest) (<-chan *CreateCertificateByCertificateIdResponse, <-chan error) {
	responseChan := make(chan *CreateCertificateByCertificateIdResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateCertificateByCertificateId(request)
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

// CreateCertificateByCertificateIdWithCallback invokes the waf_openapi.CreateCertificateByCertificateId API asynchronously
// api document: https://help.aliyun.com/api/waf-openapi/createcertificatebycertificateid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCertificateByCertificateIdWithCallback(request *CreateCertificateByCertificateIdRequest, callback func(response *CreateCertificateByCertificateIdResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateCertificateByCertificateIdResponse
		var err error
		defer close(result)
		response, err = client.CreateCertificateByCertificateId(request)
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

// CreateCertificateByCertificateIdRequest is the request struct for api CreateCertificateByCertificateId
type CreateCertificateByCertificateIdRequest struct {
	*requests.RpcRequest
	CertificateId requests.Integer `position:"Query" name:"CertificateId"`
	InstanceId    string           `position:"Query" name:"InstanceId"`
	SourceIp      string           `position:"Query" name:"SourceIp"`
	Domain        string           `position:"Query" name:"Domain"`
	Lang          string           `position:"Query" name:"Lang"`
}

// CreateCertificateByCertificateIdResponse is the response struct for api CreateCertificateByCertificateId
type CreateCertificateByCertificateIdResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	CertificateId int64  `json:"CertificateId" xml:"CertificateId"`
}

// CreateCreateCertificateByCertificateIdRequest creates a request to invoke CreateCertificateByCertificateId API
func CreateCreateCertificateByCertificateIdRequest() (request *CreateCertificateByCertificateIdRequest) {
	request = &CreateCertificateByCertificateIdRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("waf-openapi", "2019-09-10", "CreateCertificateByCertificateId", "waf", "openAPI")
	return
}

// CreateCreateCertificateByCertificateIdResponse creates a response to parse from CreateCertificateByCertificateId response
func CreateCreateCertificateByCertificateIdResponse() (response *CreateCertificateByCertificateIdResponse) {
	response = &CreateCertificateByCertificateIdResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
