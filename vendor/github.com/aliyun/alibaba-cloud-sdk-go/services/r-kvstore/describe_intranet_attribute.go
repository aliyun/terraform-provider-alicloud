package r_kvstore

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

// DescribeIntranetAttribute invokes the r_kvstore.DescribeIntranetAttribute API synchronously
func (client *Client) DescribeIntranetAttribute(request *DescribeIntranetAttributeRequest) (response *DescribeIntranetAttributeResponse, err error) {
	response = CreateDescribeIntranetAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeIntranetAttributeWithChan invokes the r_kvstore.DescribeIntranetAttribute API asynchronously
func (client *Client) DescribeIntranetAttributeWithChan(request *DescribeIntranetAttributeRequest) (<-chan *DescribeIntranetAttributeResponse, <-chan error) {
	responseChan := make(chan *DescribeIntranetAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeIntranetAttribute(request)
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

// DescribeIntranetAttributeWithCallback invokes the r_kvstore.DescribeIntranetAttribute API asynchronously
func (client *Client) DescribeIntranetAttributeWithCallback(request *DescribeIntranetAttributeRequest, callback func(response *DescribeIntranetAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeIntranetAttributeResponse
		var err error
		defer close(result)
		response, err = client.DescribeIntranetAttribute(request)
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

// DescribeIntranetAttributeRequest is the request struct for api DescribeIntranetAttribute
type DescribeIntranetAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceGroupId      string           `position:"Query" name:"ResourceGroupId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	InstanceId           string           `position:"Query" name:"InstanceId"`
}

// DescribeIntranetAttributeResponse is the response struct for api DescribeIntranetAttribute
type DescribeIntranetAttributeResponse struct {
	*responses.BaseResponse
	RequestId           string `json:"RequestId" xml:"RequestId"`
	IntranetBandwidth   int    `json:"IntranetBandwidth" xml:"IntranetBandwidth"`
	ExpireTime          string `json:"ExpireTime" xml:"ExpireTime"`
	BandwidthExpireTime string `json:"BandwidthExpireTime" xml:"BandwidthExpireTime"`
	AutoRenewal         bool   `json:"AutoRenewal" xml:"AutoRenewal"`
}

// CreateDescribeIntranetAttributeRequest creates a request to invoke DescribeIntranetAttribute API
func CreateDescribeIntranetAttributeRequest() (request *DescribeIntranetAttributeRequest) {
	request = &DescribeIntranetAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "DescribeIntranetAttribute", "redisa", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeIntranetAttributeResponse creates a response to parse from DescribeIntranetAttribute response
func CreateDescribeIntranetAttributeResponse() (response *DescribeIntranetAttributeResponse) {
	response = &DescribeIntranetAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
