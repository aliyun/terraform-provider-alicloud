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

// DescribeCustomLine invokes the alidns.DescribeCustomLine API synchronously
// api document: https://help.aliyun.com/api/alidns/describecustomline.html
func (client *Client) DescribeCustomLine(request *DescribeCustomLineRequest) (response *DescribeCustomLineResponse, err error) {
	response = CreateDescribeCustomLineResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeCustomLineWithChan invokes the alidns.DescribeCustomLine API asynchronously
// api document: https://help.aliyun.com/api/alidns/describecustomline.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCustomLineWithChan(request *DescribeCustomLineRequest) (<-chan *DescribeCustomLineResponse, <-chan error) {
	responseChan := make(chan *DescribeCustomLineResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeCustomLine(request)
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

// DescribeCustomLineWithCallback invokes the alidns.DescribeCustomLine API asynchronously
// api document: https://help.aliyun.com/api/alidns/describecustomline.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCustomLineWithCallback(request *DescribeCustomLineRequest, callback func(response *DescribeCustomLineResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeCustomLineResponse
		var err error
		defer close(result)
		response, err = client.DescribeCustomLine(request)
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

// DescribeCustomLineRequest is the request struct for api DescribeCustomLine
type DescribeCustomLineRequest struct {
	*requests.RpcRequest
	LineId       requests.Integer `position:"Query" name:"LineId"`
	UserClientIp string           `position:"Query" name:"UserClientIp"`
	Lang         string           `position:"Query" name:"Lang"`
}

// DescribeCustomLineResponse is the response struct for api DescribeCustomLine
type DescribeCustomLineResponse struct {
	*responses.BaseResponse
	RequestId       string      `json:"RequestId" xml:"RequestId"`
	Id              int64       `json:"Id" xml:"Id"`
	Name            string      `json:"Name" xml:"Name"`
	DomainName      string      `json:"DomainName" xml:"DomainName"`
	CreateTime      string      `json:"CreateTime" xml:"CreateTime"`
	CreateTimestamp int64       `json:"CreateTimestamp" xml:"CreateTimestamp"`
	IpSegments      string      `json:"IpSegments" xml:"IpSegments"`
	Code            string      `json:"Code" xml:"Code"`
	IpSegmentList   []IpSegment `json:"IpSegmentList" xml:"IpSegmentList"`
}

// CreateDescribeCustomLineRequest creates a request to invoke DescribeCustomLine API
func CreateDescribeCustomLineRequest() (request *DescribeCustomLineRequest) {
	request = &DescribeCustomLineRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "DescribeCustomLine", "alidns", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeCustomLineResponse creates a response to parse from DescribeCustomLine response
func CreateDescribeCustomLineResponse() (response *DescribeCustomLineResponse) {
	response = &DescribeCustomLineResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
