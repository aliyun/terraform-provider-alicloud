package vpc

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

// ModifySslVpnClientCert invokes the vpc.ModifySslVpnClientCert API synchronously
// api document: https://help.aliyun.com/api/vpc/modifysslvpnclientcert.html
func (client *Client) ModifySslVpnClientCert(request *ModifySslVpnClientCertRequest) (response *ModifySslVpnClientCertResponse, err error) {
	response = CreateModifySslVpnClientCertResponse()
	err = client.DoAction(request, response)
	return
}

// ModifySslVpnClientCertWithChan invokes the vpc.ModifySslVpnClientCert API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifysslvpnclientcert.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySslVpnClientCertWithChan(request *ModifySslVpnClientCertRequest) (<-chan *ModifySslVpnClientCertResponse, <-chan error) {
	responseChan := make(chan *ModifySslVpnClientCertResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifySslVpnClientCert(request)
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

// ModifySslVpnClientCertWithCallback invokes the vpc.ModifySslVpnClientCert API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifysslvpnclientcert.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySslVpnClientCertWithCallback(request *ModifySslVpnClientCertRequest, callback func(response *ModifySslVpnClientCertResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifySslVpnClientCertResponse
		var err error
		defer close(result)
		response, err = client.ModifySslVpnClientCert(request)
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

// ModifySslVpnClientCertRequest is the request struct for api ModifySslVpnClientCert
type ModifySslVpnClientCertRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	SslVpnClientCertId   string           `position:"Query" name:"SslVpnClientCertId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Name                 string           `position:"Query" name:"Name"`
}

// ModifySslVpnClientCertResponse is the response struct for api ModifySslVpnClientCert
type ModifySslVpnClientCertResponse struct {
	*responses.BaseResponse
	RequestId          string `json:"RequestId" xml:"RequestId"`
	Name               string `json:"Name" xml:"Name"`
	SslVpnClientCertId string `json:"SslVpnClientCertId" xml:"SslVpnClientCertId"`
}

// CreateModifySslVpnClientCertRequest creates a request to invoke ModifySslVpnClientCert API
func CreateModifySslVpnClientCertRequest() (request *ModifySslVpnClientCertRequest) {
	request = &ModifySslVpnClientCertRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifySslVpnClientCert", "Vpc", "openAPI")
	return
}

// CreateModifySslVpnClientCertResponse creates a response to parse from ModifySslVpnClientCert response
func CreateModifySslVpnClientCertResponse() (response *ModifySslVpnClientCertResponse) {
	response = &ModifySslVpnClientCertResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
