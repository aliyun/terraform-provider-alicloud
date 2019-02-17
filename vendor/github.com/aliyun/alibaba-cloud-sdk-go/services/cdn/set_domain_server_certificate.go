package cdn

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

// SetDomainServerCertificate invokes the cdn.SetDomainServerCertificate API synchronously
// api document: https://help.aliyun.com/api/cdn/setdomainservercertificate.html
func (client *Client) SetDomainServerCertificate(request *SetDomainServerCertificateRequest) (response *SetDomainServerCertificateResponse, err error) {
	response = CreateSetDomainServerCertificateResponse()
	err = client.DoAction(request, response)
	return
}

// SetDomainServerCertificateWithChan invokes the cdn.SetDomainServerCertificate API asynchronously
// api document: https://help.aliyun.com/api/cdn/setdomainservercertificate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetDomainServerCertificateWithChan(request *SetDomainServerCertificateRequest) (<-chan *SetDomainServerCertificateResponse, <-chan error) {
	responseChan := make(chan *SetDomainServerCertificateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetDomainServerCertificate(request)
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

// SetDomainServerCertificateWithCallback invokes the cdn.SetDomainServerCertificate API asynchronously
// api document: https://help.aliyun.com/api/cdn/setdomainservercertificate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetDomainServerCertificateWithCallback(request *SetDomainServerCertificateRequest, callback func(response *SetDomainServerCertificateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetDomainServerCertificateResponse
		var err error
		defer close(result)
		response, err = client.SetDomainServerCertificate(request)
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

// SetDomainServerCertificateRequest is the request struct for api SetDomainServerCertificate
type SetDomainServerCertificateRequest struct {
	*requests.RpcRequest
	PrivateKey              string           `position:"Query" name:"PrivateKey"`
	ForceSet                string           `position:"Query" name:"ForceSet"`
	ServerCertificateStatus string           `position:"Query" name:"ServerCertificateStatus"`
	ServerCertificate       string           `position:"Query" name:"ServerCertificate"`
	SecurityToken           string           `position:"Query" name:"SecurityToken"`
	CertType                string           `position:"Query" name:"CertType"`
	CertName                string           `position:"Query" name:"CertName"`
	DomainName              string           `position:"Query" name:"DomainName"`
	OwnerId                 requests.Integer `position:"Query" name:"OwnerId"`
	Region                  string           `position:"Query" name:"Region"`
}

// SetDomainServerCertificateResponse is the response struct for api SetDomainServerCertificate
type SetDomainServerCertificateResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetDomainServerCertificateRequest creates a request to invoke SetDomainServerCertificate API
func CreateSetDomainServerCertificateRequest() (request *SetDomainServerCertificateRequest) {
	request = &SetDomainServerCertificateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "SetDomainServerCertificate", "", "")
	return
}

// CreateSetDomainServerCertificateResponse creates a response to parse from SetDomainServerCertificate response
func CreateSetDomainServerCertificateResponse() (response *SetDomainServerCertificateResponse) {
	response = &SetDomainServerCertificateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
