package slb

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

// SetCACertificateName invokes the slb.SetCACertificateName API synchronously
// api document: https://help.aliyun.com/api/slb/setcacertificatename.html
func (client *Client) SetCACertificateName(request *SetCACertificateNameRequest) (response *SetCACertificateNameResponse, err error) {
	response = CreateSetCACertificateNameResponse()
	err = client.DoAction(request, response)
	return
}

// SetCACertificateNameWithChan invokes the slb.SetCACertificateName API asynchronously
// api document: https://help.aliyun.com/api/slb/setcacertificatename.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetCACertificateNameWithChan(request *SetCACertificateNameRequest) (<-chan *SetCACertificateNameResponse, <-chan error) {
	responseChan := make(chan *SetCACertificateNameResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetCACertificateName(request)
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

// SetCACertificateNameWithCallback invokes the slb.SetCACertificateName API asynchronously
// api document: https://help.aliyun.com/api/slb/setcacertificatename.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetCACertificateNameWithCallback(request *SetCACertificateNameRequest, callback func(response *SetCACertificateNameResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetCACertificateNameResponse
		var err error
		defer close(result)
		response, err = client.SetCACertificateName(request)
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

// SetCACertificateNameRequest is the request struct for api SetCACertificateName
type SetCACertificateNameRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	CACertificateName    string           `position:"Query" name:"CACertificateName"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	CACertificateId      string           `position:"Query" name:"CACertificateId"`
}

// SetCACertificateNameResponse is the response struct for api SetCACertificateName
type SetCACertificateNameResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetCACertificateNameRequest creates a request to invoke SetCACertificateName API
func CreateSetCACertificateNameRequest() (request *SetCACertificateNameRequest) {
	request = &SetCACertificateNameRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "SetCACertificateName", "slb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateSetCACertificateNameResponse creates a response to parse from SetCACertificateName response
func CreateSetCACertificateNameResponse() (response *SetCACertificateNameResponse) {
	response = &SetCACertificateNameResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
