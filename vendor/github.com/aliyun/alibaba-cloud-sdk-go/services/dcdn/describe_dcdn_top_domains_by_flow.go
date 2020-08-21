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

// DescribeDcdnTopDomainsByFlow invokes the dcdn.DescribeDcdnTopDomainsByFlow API synchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdntopdomainsbyflow.html
func (client *Client) DescribeDcdnTopDomainsByFlow(request *DescribeDcdnTopDomainsByFlowRequest) (response *DescribeDcdnTopDomainsByFlowResponse, err error) {
	response = CreateDescribeDcdnTopDomainsByFlowResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnTopDomainsByFlowWithChan invokes the dcdn.DescribeDcdnTopDomainsByFlow API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdntopdomainsbyflow.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnTopDomainsByFlowWithChan(request *DescribeDcdnTopDomainsByFlowRequest) (<-chan *DescribeDcdnTopDomainsByFlowResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnTopDomainsByFlowResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnTopDomainsByFlow(request)
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

// DescribeDcdnTopDomainsByFlowWithCallback invokes the dcdn.DescribeDcdnTopDomainsByFlow API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdntopdomainsbyflow.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnTopDomainsByFlowWithCallback(request *DescribeDcdnTopDomainsByFlowRequest, callback func(response *DescribeDcdnTopDomainsByFlowResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnTopDomainsByFlowResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnTopDomainsByFlow(request)
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

// DescribeDcdnTopDomainsByFlowRequest is the request struct for api DescribeDcdnTopDomainsByFlow
type DescribeDcdnTopDomainsByFlowRequest struct {
	*requests.RpcRequest
	StartTime string           `position:"Query" name:"StartTime"`
	Limit     requests.Integer `position:"Query" name:"Limit"`
	EndTime   string           `position:"Query" name:"EndTime"`
	OwnerId   requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDcdnTopDomainsByFlowResponse is the response struct for api DescribeDcdnTopDomainsByFlow
type DescribeDcdnTopDomainsByFlowResponse struct {
	*responses.BaseResponse
	RequestId         string     `json:"RequestId" xml:"RequestId"`
	StartTime         string     `json:"StartTime" xml:"StartTime"`
	EndTime           string     `json:"EndTime" xml:"EndTime"`
	DomainCount       int64      `json:"DomainCount" xml:"DomainCount"`
	DomainOnlineCount int64      `json:"DomainOnlineCount" xml:"DomainOnlineCount"`
	TopDomains        TopDomains `json:"TopDomains" xml:"TopDomains"`
}

// CreateDescribeDcdnTopDomainsByFlowRequest creates a request to invoke DescribeDcdnTopDomainsByFlow API
func CreateDescribeDcdnTopDomainsByFlowRequest() (request *DescribeDcdnTopDomainsByFlowRequest) {
	request = &DescribeDcdnTopDomainsByFlowRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnTopDomainsByFlow", "", "")
	return
}

// CreateDescribeDcdnTopDomainsByFlowResponse creates a response to parse from DescribeDcdnTopDomainsByFlow response
func CreateDescribeDcdnTopDomainsByFlowResponse() (response *DescribeDcdnTopDomainsByFlowResponse) {
	response = &DescribeDcdnTopDomainsByFlowResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
