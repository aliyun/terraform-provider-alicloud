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

// DescribeDomainTopUrlVisit invokes the cdn.DescribeDomainTopUrlVisit API synchronously
// api document: https://help.aliyun.com/api/cdn/describedomaintopurlvisit.html
func (client *Client) DescribeDomainTopUrlVisit(request *DescribeDomainTopUrlVisitRequest) (response *DescribeDomainTopUrlVisitResponse, err error) {
	response = CreateDescribeDomainTopUrlVisitResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainTopUrlVisitWithChan invokes the cdn.DescribeDomainTopUrlVisit API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomaintopurlvisit.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainTopUrlVisitWithChan(request *DescribeDomainTopUrlVisitRequest) (<-chan *DescribeDomainTopUrlVisitResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainTopUrlVisitResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainTopUrlVisit(request)
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

// DescribeDomainTopUrlVisitWithCallback invokes the cdn.DescribeDomainTopUrlVisit API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomaintopurlvisit.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainTopUrlVisitWithCallback(request *DescribeDomainTopUrlVisitRequest, callback func(response *DescribeDomainTopUrlVisitResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainTopUrlVisitResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainTopUrlVisit(request)
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

// DescribeDomainTopUrlVisitRequest is the request struct for api DescribeDomainTopUrlVisit
type DescribeDomainTopUrlVisitRequest struct {
	*requests.RpcRequest
	StartTime  string           `position:"Query" name:"StartTime"`
	Percent    string           `position:"Query" name:"Percent"`
	DomainName string           `position:"Query" name:"DomainName"`
	EndTime    string           `position:"Query" name:"EndTime"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
	SortBy     string           `position:"Query" name:"SortBy"`
}

// DescribeDomainTopUrlVisitResponse is the response struct for api DescribeDomainTopUrlVisit
type DescribeDomainTopUrlVisitResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	DomainName string     `json:"DomainName" xml:"DomainName"`
	StartTime  string     `json:"StartTime" xml:"StartTime"`
	AllUrlList AllUrlList `json:"AllUrlList" xml:"AllUrlList"`
	Url200List Url200List `json:"Url200List" xml:"Url200List"`
	Url300List Url300List `json:"Url300List" xml:"Url300List"`
	Url400List Url400List `json:"Url400List" xml:"Url400List"`
	Url500List Url500List `json:"Url500List" xml:"Url500List"`
}

// CreateDescribeDomainTopUrlVisitRequest creates a request to invoke DescribeDomainTopUrlVisit API
func CreateDescribeDomainTopUrlVisitRequest() (request *DescribeDomainTopUrlVisitRequest) {
	request = &DescribeDomainTopUrlVisitRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeDomainTopUrlVisit", "", "")
	return
}

// CreateDescribeDomainTopUrlVisitResponse creates a response to parse from DescribeDomainTopUrlVisit response
func CreateDescribeDomainTopUrlVisitResponse() (response *DescribeDomainTopUrlVisitResponse) {
	response = &DescribeDomainTopUrlVisitResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
