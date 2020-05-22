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

// DescribePortMaxConns invokes the ddoscoo.DescribePortMaxConns API synchronously
// api document: https://help.aliyun.com/api/ddoscoo/describeportmaxconns.html
func (client *Client) DescribePortMaxConns(request *DescribePortMaxConnsRequest) (response *DescribePortMaxConnsResponse, err error) {
	response = CreateDescribePortMaxConnsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribePortMaxConnsWithChan invokes the ddoscoo.DescribePortMaxConns API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describeportmaxconns.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribePortMaxConnsWithChan(request *DescribePortMaxConnsRequest) (<-chan *DescribePortMaxConnsResponse, <-chan error) {
	responseChan := make(chan *DescribePortMaxConnsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribePortMaxConns(request)
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

// DescribePortMaxConnsWithCallback invokes the ddoscoo.DescribePortMaxConns API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describeportmaxconns.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribePortMaxConnsWithCallback(request *DescribePortMaxConnsRequest, callback func(response *DescribePortMaxConnsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribePortMaxConnsResponse
		var err error
		defer close(result)
		response, err = client.DescribePortMaxConns(request)
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

// DescribePortMaxConnsRequest is the request struct for api DescribePortMaxConns
type DescribePortMaxConnsRequest struct {
	*requests.RpcRequest
	EndTime         requests.Integer `position:"Query" name:"EndTime"`
	StartTime       requests.Integer `position:"Query" name:"StartTime"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SourceIp        string           `position:"Query" name:"SourceIp"`
	InstanceIds     *[]string        `position:"Query" name:"InstanceIds"  type:"Repeated"`
}

// DescribePortMaxConnsResponse is the response struct for api DescribePortMaxConns
type DescribePortMaxConnsResponse struct {
	*responses.BaseResponse
	RequestId    string             `json:"RequestId" xml:"RequestId"`
	PortMaxConns []PortMaxConnsItem `json:"PortMaxConns" xml:"PortMaxConns"`
}

// CreateDescribePortMaxConnsRequest creates a request to invoke DescribePortMaxConns API
func CreateDescribePortMaxConnsRequest() (request *DescribePortMaxConnsRequest) {
	request = &DescribePortMaxConnsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribePortMaxConns", "ddoscoo", "openAPI")
	return
}

// CreateDescribePortMaxConnsResponse creates a response to parse from DescribePortMaxConns response
func CreateDescribePortMaxConnsResponse() (response *DescribePortMaxConnsResponse) {
	response = &DescribePortMaxConnsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
