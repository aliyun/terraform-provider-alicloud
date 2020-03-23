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

// DescribeAvailableResource invokes the r_kvstore.DescribeAvailableResource API synchronously
// api document: https://help.aliyun.com/api/r-kvstore/describeavailableresource.html
func (client *Client) DescribeAvailableResource(request *DescribeAvailableResourceRequest) (response *DescribeAvailableResourceResponse, err error) {
	response = CreateDescribeAvailableResourceResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAvailableResourceWithChan invokes the r_kvstore.DescribeAvailableResource API asynchronously
// api document: https://help.aliyun.com/api/r-kvstore/describeavailableresource.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableResourceWithChan(request *DescribeAvailableResourceRequest) (<-chan *DescribeAvailableResourceResponse, <-chan error) {
	responseChan := make(chan *DescribeAvailableResourceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAvailableResource(request)
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

// DescribeAvailableResourceWithCallback invokes the r_kvstore.DescribeAvailableResource API asynchronously
// api document: https://help.aliyun.com/api/r-kvstore/describeavailableresource.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableResourceWithCallback(request *DescribeAvailableResourceRequest, callback func(response *DescribeAvailableResourceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAvailableResourceResponse
		var err error
		defer close(result)
		response, err = client.DescribeAvailableResource(request)
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

// DescribeAvailableResourceRequest is the request struct for api DescribeAvailableResource
type DescribeAvailableResourceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	Engine               string           `position:"Query" name:"Engine"`
	InstanceChargeType   string           `position:"Query" name:"InstanceChargeType"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	Level                string           `position:"Query" name:"Level"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	OrderType            string           `position:"Query" name:"OrderType"`
}

// DescribeAvailableResourceResponse is the response struct for api DescribeAvailableResource
type DescribeAvailableResourceResponse struct {
	*responses.BaseResponse
	RequestId      string         `json:"RequestId" xml:"RequestId"`
	AvailableZones AvailableZones `json:"AvailableZones" xml:"AvailableZones"`
}

// CreateDescribeAvailableResourceRequest creates a request to invoke DescribeAvailableResource API
func CreateDescribeAvailableResourceRequest() (request *DescribeAvailableResourceRequest) {
	request = &DescribeAvailableResourceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "DescribeAvailableResource", "", "")
	return
}

// CreateDescribeAvailableResourceResponse creates a response to parse from DescribeAvailableResource response
func CreateDescribeAvailableResourceResponse() (response *DescribeAvailableResourceResponse) {
	response = &DescribeAvailableResourceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
