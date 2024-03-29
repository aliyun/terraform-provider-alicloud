package ddoscoo

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

// DescribeDomainStatusCodeList invokes the ddoscoo.DescribeDomainStatusCodeList API synchronously
func (client *Client) DescribeDomainStatusCodeList(request *DescribeDomainStatusCodeListRequest) (response *DescribeDomainStatusCodeListResponse, err error) {
	response = CreateDescribeDomainStatusCodeListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainStatusCodeListWithChan invokes the ddoscoo.DescribeDomainStatusCodeList API asynchronously
func (client *Client) DescribeDomainStatusCodeListWithChan(request *DescribeDomainStatusCodeListRequest) (<-chan *DescribeDomainStatusCodeListResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainStatusCodeListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainStatusCodeList(request)
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

// DescribeDomainStatusCodeListWithCallback invokes the ddoscoo.DescribeDomainStatusCodeList API asynchronously
func (client *Client) DescribeDomainStatusCodeListWithCallback(request *DescribeDomainStatusCodeListRequest, callback func(response *DescribeDomainStatusCodeListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainStatusCodeListResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainStatusCodeList(request)
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

// DescribeDomainStatusCodeListRequest is the request struct for api DescribeDomainStatusCodeList
type DescribeDomainStatusCodeListRequest struct {
	*requests.RpcRequest
	StartTime       requests.Integer `position:"Query" name:"StartTime"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SourceIp        string           `position:"Query" name:"SourceIp"`
	QueryType       string           `position:"Query" name:"QueryType"`
	EndTime         requests.Integer `position:"Query" name:"EndTime"`
	Domain          string           `position:"Query" name:"Domain"`
	Interval        requests.Integer `position:"Query" name:"Interval"`
}

// DescribeDomainStatusCodeListResponse is the response struct for api DescribeDomainStatusCodeList
type DescribeDomainStatusCodeListResponse struct {
	*responses.BaseResponse
	RequestId      string       `json:"RequestId" xml:"RequestId"`
	StatusCodeList []StatusCode `json:"StatusCodeList" xml:"StatusCodeList"`
}

// CreateDescribeDomainStatusCodeListRequest creates a request to invoke DescribeDomainStatusCodeList API
func CreateDescribeDomainStatusCodeListRequest() (request *DescribeDomainStatusCodeListRequest) {
	request = &DescribeDomainStatusCodeListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeDomainStatusCodeList", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainStatusCodeListResponse creates a response to parse from DescribeDomainStatusCodeList response
func CreateDescribeDomainStatusCodeListResponse() (response *DescribeDomainStatusCodeListResponse) {
	response = &DescribeDomainStatusCodeListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
