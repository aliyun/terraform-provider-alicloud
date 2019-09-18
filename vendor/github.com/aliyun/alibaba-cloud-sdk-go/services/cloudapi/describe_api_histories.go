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

// DescribeApiHistories invokes the cloudapi.DescribeApiHistories API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapihistories.html
func (client *Client) DescribeApiHistories(request *DescribeApiHistoriesRequest) (response *DescribeApiHistoriesResponse, err error) {
	response = CreateDescribeApiHistoriesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeApiHistoriesWithChan invokes the cloudapi.DescribeApiHistories API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapihistories.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiHistoriesWithChan(request *DescribeApiHistoriesRequest) (<-chan *DescribeApiHistoriesResponse, <-chan error) {
	responseChan := make(chan *DescribeApiHistoriesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeApiHistories(request)
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

// DescribeApiHistoriesWithCallback invokes the cloudapi.DescribeApiHistories API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapihistories.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiHistoriesWithCallback(request *DescribeApiHistoriesRequest, callback func(response *DescribeApiHistoriesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeApiHistoriesResponse
		var err error
		defer close(result)
		response, err = client.DescribeApiHistories(request)
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

// DescribeApiHistoriesRequest is the request struct for api DescribeApiHistories
type DescribeApiHistoriesRequest struct {
	*requests.RpcRequest
	StageName     string `position:"Query" name:"StageName"`
	GroupId       string `position:"Query" name:"GroupId"`
	PageNumber    string `position:"Query" name:"PageNumber"`
	ApiName       string `position:"Query" name:"ApiName"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
	PageSize      string `position:"Query" name:"PageSize"`
	ApiId         string `position:"Query" name:"ApiId"`
}

// DescribeApiHistoriesResponse is the response struct for api DescribeApiHistories
type DescribeApiHistoriesResponse struct {
	*responses.BaseResponse
	RequestId   string                            `json:"RequestId" xml:"RequestId"`
	TotalCount  int                               `json:"TotalCount" xml:"TotalCount"`
	PageSize    int                               `json:"PageSize" xml:"PageSize"`
	PageNumber  int                               `json:"PageNumber" xml:"PageNumber"`
	ApiHisItems ApiHisItemsInDescribeApiHistories `json:"ApiHisItems" xml:"ApiHisItems"`
}

// CreateDescribeApiHistoriesRequest creates a request to invoke DescribeApiHistories API
func CreateDescribeApiHistoriesRequest() (request *DescribeApiHistoriesRequest) {
	request = &DescribeApiHistoriesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeApiHistories", "apigateway", "openAPI")
	return
}

// CreateDescribeApiHistoriesResponse creates a response to parse from DescribeApiHistories response
func CreateDescribeApiHistoriesResponse() (response *DescribeApiHistoriesResponse) {
	response = &DescribeApiHistoriesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
