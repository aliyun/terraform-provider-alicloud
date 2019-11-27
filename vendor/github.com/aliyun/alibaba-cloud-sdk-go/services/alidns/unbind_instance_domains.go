package alidns

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

// UnbindInstanceDomains invokes the alidns.UnbindInstanceDomains API synchronously
// api document: https://help.aliyun.com/api/alidns/unbindinstancedomains.html
func (client *Client) UnbindInstanceDomains(request *UnbindInstanceDomainsRequest) (response *UnbindInstanceDomainsResponse, err error) {
	response = CreateUnbindInstanceDomainsResponse()
	err = client.DoAction(request, response)
	return
}

// UnbindInstanceDomainsWithChan invokes the alidns.UnbindInstanceDomains API asynchronously
// api document: https://help.aliyun.com/api/alidns/unbindinstancedomains.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnbindInstanceDomainsWithChan(request *UnbindInstanceDomainsRequest) (<-chan *UnbindInstanceDomainsResponse, <-chan error) {
	responseChan := make(chan *UnbindInstanceDomainsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UnbindInstanceDomains(request)
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

// UnbindInstanceDomainsWithCallback invokes the alidns.UnbindInstanceDomains API asynchronously
// api document: https://help.aliyun.com/api/alidns/unbindinstancedomains.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnbindInstanceDomainsWithCallback(request *UnbindInstanceDomainsRequest, callback func(response *UnbindInstanceDomainsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UnbindInstanceDomainsResponse
		var err error
		defer close(result)
		response, err = client.UnbindInstanceDomains(request)
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

// UnbindInstanceDomainsRequest is the request struct for api UnbindInstanceDomains
type UnbindInstanceDomainsRequest struct {
	*requests.RpcRequest
	DomainNames  string `position:"Query" name:"DomainNames"`
	InstanceId   string `position:"Query" name:"InstanceId"`
	UserClientIp string `position:"Query" name:"UserClientIp"`
	Lang         string `position:"Query" name:"Lang"`
}

// UnbindInstanceDomainsResponse is the response struct for api UnbindInstanceDomains
type UnbindInstanceDomainsResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	SuccessCount int    `json:"SuccessCount" xml:"SuccessCount"`
	FailedCount  int    `json:"FailedCount" xml:"FailedCount"`
}

// CreateUnbindInstanceDomainsRequest creates a request to invoke UnbindInstanceDomains API
func CreateUnbindInstanceDomainsRequest() (request *UnbindInstanceDomainsRequest) {
	request = &UnbindInstanceDomainsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "UnbindInstanceDomains", "alidns", "openAPI")
	return
}

// CreateUnbindInstanceDomainsResponse creates a response to parse from UnbindInstanceDomains response
func CreateUnbindInstanceDomainsResponse() (response *UnbindInstanceDomainsResponse) {
	response = &UnbindInstanceDomainsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
