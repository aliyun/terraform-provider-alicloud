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

// DescribeDomainSrcFlowData invokes the cdn.DescribeDomainSrcFlowData API synchronously
// api document: https://help.aliyun.com/api/cdn/describedomainsrcflowdata.html
func (client *Client) DescribeDomainSrcFlowData(request *DescribeDomainSrcFlowDataRequest) (response *DescribeDomainSrcFlowDataResponse, err error) {
	response = CreateDescribeDomainSrcFlowDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainSrcFlowDataWithChan invokes the cdn.DescribeDomainSrcFlowData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainsrcflowdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainSrcFlowDataWithChan(request *DescribeDomainSrcFlowDataRequest) (<-chan *DescribeDomainSrcFlowDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainSrcFlowDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainSrcFlowData(request)
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

// DescribeDomainSrcFlowDataWithCallback invokes the cdn.DescribeDomainSrcFlowData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainsrcflowdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainSrcFlowDataWithCallback(request *DescribeDomainSrcFlowDataRequest, callback func(response *DescribeDomainSrcFlowDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainSrcFlowDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainSrcFlowData(request)
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

// DescribeDomainSrcFlowDataRequest is the request struct for api DescribeDomainSrcFlowData
type DescribeDomainSrcFlowDataRequest struct {
	*requests.RpcRequest
	StartTime  string           `position:"Query" name:"StartTime"`
	FixTimeGap string           `position:"Query" name:"FixTimeGap"`
	TimeMerge  string           `position:"Query" name:"TimeMerge"`
	DomainName string           `position:"Query" name:"DomainName"`
	EndTime    string           `position:"Query" name:"EndTime"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
	Interval   string           `position:"Query" name:"Interval"`
}

// DescribeDomainSrcFlowDataResponse is the response struct for api DescribeDomainSrcFlowData
type DescribeDomainSrcFlowDataResponse struct {
	*responses.BaseResponse
	RequestId              string                 `json:"RequestId" xml:"RequestId"`
	DomainName             string                 `json:"DomainName" xml:"DomainName"`
	StartTime              string                 `json:"StartTime" xml:"StartTime"`
	EndTime                string                 `json:"EndTime" xml:"EndTime"`
	DataInterval           string                 `json:"DataInterval" xml:"DataInterval"`
	SrcFlowDataPerInterval SrcFlowDataPerInterval `json:"SrcFlowDataPerInterval" xml:"SrcFlowDataPerInterval"`
}

// CreateDescribeDomainSrcFlowDataRequest creates a request to invoke DescribeDomainSrcFlowData API
func CreateDescribeDomainSrcFlowDataRequest() (request *DescribeDomainSrcFlowDataRequest) {
	request = &DescribeDomainSrcFlowDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeDomainSrcFlowData", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainSrcFlowDataResponse creates a response to parse from DescribeDomainSrcFlowData response
func CreateDescribeDomainSrcFlowDataResponse() (response *DescribeDomainSrcFlowDataResponse) {
	response = &DescribeDomainSrcFlowDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
