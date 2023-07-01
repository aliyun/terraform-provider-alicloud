package cloudapi

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

// DescribePluginApis invokes the cloudapi.DescribePluginApis API synchronously
func (client *Client) DescribePluginApis(request *DescribePluginApisRequest) (response *DescribePluginApisResponse, err error) {
	response = CreateDescribePluginApisResponse()
	err = client.DoAction(request, response)
	return
}

// DescribePluginApisWithChan invokes the cloudapi.DescribePluginApis API asynchronously
func (client *Client) DescribePluginApisWithChan(request *DescribePluginApisRequest) (<-chan *DescribePluginApisResponse, <-chan error) {
	responseChan := make(chan *DescribePluginApisResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribePluginApis(request)
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

// DescribePluginApisWithCallback invokes the cloudapi.DescribePluginApis API asynchronously
func (client *Client) DescribePluginApisWithCallback(request *DescribePluginApisRequest, callback func(response *DescribePluginApisResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribePluginApisResponse
		var err error
		defer close(result)
		response, err = client.DescribePluginApis(request)
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

// DescribePluginApisRequest is the request struct for api DescribePluginApis
type DescribePluginApisRequest struct {
	*requests.RpcRequest
	Method        string           `position:"Query" name:"Method"`
	PluginId      string           `position:"Query" name:"PluginId"`
	GroupId       string           `position:"Query" name:"GroupId"`
	Description   string           `position:"Query" name:"Description"`
	PageNumber    requests.Integer `position:"Query" name:"PageNumber"`
	Path          string           `position:"Query" name:"Path"`
	ApiName       string           `position:"Query" name:"ApiName"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	ApiId         string           `position:"Query" name:"ApiId"`
}

// DescribePluginApisResponse is the response struct for api DescribePluginApis
type DescribePluginApisResponse struct {
	*responses.BaseResponse
	PageNumber  int                             `json:"PageNumber" xml:"PageNumber"`
	RequestId   string                          `json:"RequestId" xml:"RequestId"`
	PageSize    int                             `json:"PageSize" xml:"PageSize"`
	TotalCount  int                             `json:"TotalCount" xml:"TotalCount"`
	ApiSummarys ApiSummarysInDescribePluginApis `json:"ApiSummarys" xml:"ApiSummarys"`
}

// CreateDescribePluginApisRequest creates a request to invoke DescribePluginApis API
func CreateDescribePluginApisRequest() (request *DescribePluginApisRequest) {
	request = &DescribePluginApisRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribePluginApis", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribePluginApisResponse creates a response to parse from DescribePluginApis response
func CreateDescribePluginApisResponse() (response *DescribePluginApisResponse) {
	response = &DescribePluginApisResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
