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

// DescribeDomainBasicConfigs invokes the waf_openapi.DescribeDomainBasicConfigs API synchronously
// api document: https://help.aliyun.com/api/waf-openapi/describedomainbasicconfigs.html
func (client *Client) DescribeDomainBasicConfigs(request *DescribeDomainBasicConfigsRequest) (response *DescribeDomainBasicConfigsResponse, err error) {
	response = CreateDescribeDomainBasicConfigsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainBasicConfigsWithChan invokes the waf_openapi.DescribeDomainBasicConfigs API asynchronously
// api document: https://help.aliyun.com/api/waf-openapi/describedomainbasicconfigs.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainBasicConfigsWithChan(request *DescribeDomainBasicConfigsRequest) (<-chan *DescribeDomainBasicConfigsResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainBasicConfigsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainBasicConfigs(request)
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

// DescribeDomainBasicConfigsWithCallback invokes the waf_openapi.DescribeDomainBasicConfigs API asynchronously
// api document: https://help.aliyun.com/api/waf-openapi/describedomainbasicconfigs.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainBasicConfigsWithCallback(request *DescribeDomainBasicConfigsRequest, callback func(response *DescribeDomainBasicConfigsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainBasicConfigsResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainBasicConfigs(request)
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

// DescribeDomainBasicConfigsRequest is the request struct for api DescribeDomainBasicConfigs
type DescribeDomainBasicConfigsRequest struct {
	*requests.RpcRequest
	PageNumber      requests.Integer `position:"Query" name:"PageNumber"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SourceIp        string           `position:"Query" name:"SourceIp"`
	PageSize        requests.Integer `position:"Query" name:"PageSize"`
	Lang            string           `position:"Query" name:"Lang"`
	InstanceId      string           `position:"Query" name:"InstanceId"`
	DomainKey       string           `position:"Query" name:"DomainKey"`
}

// DescribeDomainBasicConfigsResponse is the response struct for api DescribeDomainBasicConfigs
type DescribeDomainBasicConfigsResponse struct {
	*responses.BaseResponse
	RequestId     string         `json:"RequestId" xml:"RequestId"`
	TotalCount    int            `json:"TotalCount" xml:"TotalCount"`
	DomainConfigs []DomainConfig `json:"DomainConfigs" xml:"DomainConfigs"`
}

// CreateDescribeDomainBasicConfigsRequest creates a request to invoke DescribeDomainBasicConfigs API
func CreateDescribeDomainBasicConfigsRequest() (request *DescribeDomainBasicConfigsRequest) {
	request = &DescribeDomainBasicConfigsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("waf-openapi", "2019-09-10", "DescribeDomainBasicConfigs", "waf", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainBasicConfigsResponse creates a response to parse from DescribeDomainBasicConfigs response
func CreateDescribeDomainBasicConfigsResponse() (response *DescribeDomainBasicConfigsResponse) {
	response = &DescribeDomainBasicConfigsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
