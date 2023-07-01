package dcdn

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

// DescribeDcdnWafDomainDetail invokes the dcdn.DescribeDcdnWafDomainDetail API synchronously
func (client *Client) DescribeDcdnWafDomainDetail(request *DescribeDcdnWafDomainDetailRequest) (response *DescribeDcdnWafDomainDetailResponse, err error) {
	response = CreateDescribeDcdnWafDomainDetailResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnWafDomainDetailWithChan invokes the dcdn.DescribeDcdnWafDomainDetail API asynchronously
func (client *Client) DescribeDcdnWafDomainDetailWithChan(request *DescribeDcdnWafDomainDetailRequest) (<-chan *DescribeDcdnWafDomainDetailResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnWafDomainDetailResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnWafDomainDetail(request)
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

// DescribeDcdnWafDomainDetailWithCallback invokes the dcdn.DescribeDcdnWafDomainDetail API asynchronously
func (client *Client) DescribeDcdnWafDomainDetailWithCallback(request *DescribeDcdnWafDomainDetailRequest, callback func(response *DescribeDcdnWafDomainDetailResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnWafDomainDetailResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnWafDomainDetail(request)
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

// DescribeDcdnWafDomainDetailRequest is the request struct for api DescribeDcdnWafDomainDetail
type DescribeDcdnWafDomainDetailRequest struct {
	*requests.RpcRequest
	DomainName string `position:"Query" name:"DomainName"`
}

// DescribeDcdnWafDomainDetailResponse is the response struct for api DescribeDcdnWafDomainDetail
type DescribeDcdnWafDomainDetailResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Domain    Domain `json:"Domain" xml:"Domain"`
}

// CreateDescribeDcdnWafDomainDetailRequest creates a request to invoke DescribeDcdnWafDomainDetail API
func CreateDescribeDcdnWafDomainDetailRequest() (request *DescribeDcdnWafDomainDetailRequest) {
	request = &DescribeDcdnWafDomainDetailRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnWafDomainDetail", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnWafDomainDetailResponse creates a response to parse from DescribeDcdnWafDomainDetail response
func CreateDescribeDcdnWafDomainDetailResponse() (response *DescribeDcdnWafDomainDetailResponse) {
	response = &DescribeDcdnWafDomainDetailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
