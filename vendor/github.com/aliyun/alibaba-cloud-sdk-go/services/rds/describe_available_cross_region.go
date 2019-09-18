package rds

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

// DescribeAvailableCrossRegion invokes the rds.DescribeAvailableCrossRegion API synchronously
// api document: https://help.aliyun.com/api/rds/describeavailablecrossregion.html
func (client *Client) DescribeAvailableCrossRegion(request *DescribeAvailableCrossRegionRequest) (response *DescribeAvailableCrossRegionResponse, err error) {
	response = CreateDescribeAvailableCrossRegionResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAvailableCrossRegionWithChan invokes the rds.DescribeAvailableCrossRegion API asynchronously
// api document: https://help.aliyun.com/api/rds/describeavailablecrossregion.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableCrossRegionWithChan(request *DescribeAvailableCrossRegionRequest) (<-chan *DescribeAvailableCrossRegionResponse, <-chan error) {
	responseChan := make(chan *DescribeAvailableCrossRegionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAvailableCrossRegion(request)
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

// DescribeAvailableCrossRegionWithCallback invokes the rds.DescribeAvailableCrossRegion API asynchronously
// api document: https://help.aliyun.com/api/rds/describeavailablecrossregion.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableCrossRegionWithCallback(request *DescribeAvailableCrossRegionRequest, callback func(response *DescribeAvailableCrossRegionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAvailableCrossRegionResponse
		var err error
		defer close(result)
		response, err = client.DescribeAvailableCrossRegion(request)
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

// DescribeAvailableCrossRegionRequest is the request struct for api DescribeAvailableCrossRegion
type DescribeAvailableCrossRegionRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
}

// DescribeAvailableCrossRegionResponse is the response struct for api DescribeAvailableCrossRegion
type DescribeAvailableCrossRegionResponse struct {
	*responses.BaseResponse
	RequestId string                                `json:"RequestId" xml:"RequestId"`
	Regions   RegionsInDescribeAvailableCrossRegion `json:"Regions" xml:"Regions"`
}

// CreateDescribeAvailableCrossRegionRequest creates a request to invoke DescribeAvailableCrossRegion API
func CreateDescribeAvailableCrossRegionRequest() (request *DescribeAvailableCrossRegionRequest) {
	request = &DescribeAvailableCrossRegionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeAvailableCrossRegion", "rds", "openAPI")
	return
}

// CreateDescribeAvailableCrossRegionResponse creates a response to parse from DescribeAvailableCrossRegion response
func CreateDescribeAvailableCrossRegionResponse() (response *DescribeAvailableCrossRegionResponse) {
	response = &DescribeAvailableCrossRegionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
