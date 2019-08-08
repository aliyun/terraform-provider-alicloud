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

// DescribeGtmInstanceAddressPools invokes the alidns.DescribeGtmInstanceAddressPools API synchronously
// api document: https://help.aliyun.com/api/alidns/describegtminstanceaddresspools.html
func (client *Client) DescribeGtmInstanceAddressPools(request *DescribeGtmInstanceAddressPoolsRequest) (response *DescribeGtmInstanceAddressPoolsResponse, err error) {
	response = CreateDescribeGtmInstanceAddressPoolsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeGtmInstanceAddressPoolsWithChan invokes the alidns.DescribeGtmInstanceAddressPools API asynchronously
// api document: https://help.aliyun.com/api/alidns/describegtminstanceaddresspools.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeGtmInstanceAddressPoolsWithChan(request *DescribeGtmInstanceAddressPoolsRequest) (<-chan *DescribeGtmInstanceAddressPoolsResponse, <-chan error) {
	responseChan := make(chan *DescribeGtmInstanceAddressPoolsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeGtmInstanceAddressPools(request)
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

// DescribeGtmInstanceAddressPoolsWithCallback invokes the alidns.DescribeGtmInstanceAddressPools API asynchronously
// api document: https://help.aliyun.com/api/alidns/describegtminstanceaddresspools.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeGtmInstanceAddressPoolsWithCallback(request *DescribeGtmInstanceAddressPoolsRequest, callback func(response *DescribeGtmInstanceAddressPoolsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeGtmInstanceAddressPoolsResponse
		var err error
		defer close(result)
		response, err = client.DescribeGtmInstanceAddressPools(request)
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

// DescribeGtmInstanceAddressPoolsRequest is the request struct for api DescribeGtmInstanceAddressPools
type DescribeGtmInstanceAddressPoolsRequest struct {
	*requests.RpcRequest
	InstanceId   string           `position:"Query" name:"InstanceId"`
	UserClientIp string           `position:"Query" name:"UserClientIp"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
	Lang         string           `position:"Query" name:"Lang"`
	PageNumber   requests.Integer `position:"Query" name:"PageNumber"`
}

// DescribeGtmInstanceAddressPoolsResponse is the response struct for api DescribeGtmInstanceAddressPools
type DescribeGtmInstanceAddressPoolsResponse struct {
	*responses.BaseResponse
	RequestId  string                                     `json:"RequestId" xml:"RequestId"`
	TotalItems int                                        `json:"TotalItems" xml:"TotalItems"`
	TotalPages int                                        `json:"TotalPages" xml:"TotalPages"`
	PageNumber int                                        `json:"PageNumber" xml:"PageNumber"`
	PageSize   int                                        `json:"PageSize" xml:"PageSize"`
	AddrPools  AddrPoolsInDescribeGtmInstanceAddressPools `json:"AddrPools" xml:"AddrPools"`
}

// CreateDescribeGtmInstanceAddressPoolsRequest creates a request to invoke DescribeGtmInstanceAddressPools API
func CreateDescribeGtmInstanceAddressPoolsRequest() (request *DescribeGtmInstanceAddressPoolsRequest) {
	request = &DescribeGtmInstanceAddressPoolsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "DescribeGtmInstanceAddressPools", "Alidns", "openAPI")
	return
}

// CreateDescribeGtmInstanceAddressPoolsResponse creates a response to parse from DescribeGtmInstanceAddressPools response
func CreateDescribeGtmInstanceAddressPoolsResponse() (response *DescribeGtmInstanceAddressPoolsResponse) {
	response = &DescribeGtmInstanceAddressPoolsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
