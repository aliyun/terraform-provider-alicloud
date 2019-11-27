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

// DescribeAvailableDedicatedHostZones invokes the rds.DescribeAvailableDedicatedHostZones API synchronously
// api document: https://help.aliyun.com/api/rds/describeavailablededicatedhostzones.html
func (client *Client) DescribeAvailableDedicatedHostZones(request *DescribeAvailableDedicatedHostZonesRequest) (response *DescribeAvailableDedicatedHostZonesResponse, err error) {
	response = CreateDescribeAvailableDedicatedHostZonesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAvailableDedicatedHostZonesWithChan invokes the rds.DescribeAvailableDedicatedHostZones API asynchronously
// api document: https://help.aliyun.com/api/rds/describeavailablededicatedhostzones.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableDedicatedHostZonesWithChan(request *DescribeAvailableDedicatedHostZonesRequest) (<-chan *DescribeAvailableDedicatedHostZonesResponse, <-chan error) {
	responseChan := make(chan *DescribeAvailableDedicatedHostZonesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAvailableDedicatedHostZones(request)
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

// DescribeAvailableDedicatedHostZonesWithCallback invokes the rds.DescribeAvailableDedicatedHostZones API asynchronously
// api document: https://help.aliyun.com/api/rds/describeavailablededicatedhostzones.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAvailableDedicatedHostZonesWithCallback(request *DescribeAvailableDedicatedHostZonesRequest, callback func(response *DescribeAvailableDedicatedHostZonesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAvailableDedicatedHostZonesResponse
		var err error
		defer close(result)
		response, err = client.DescribeAvailableDedicatedHostZones(request)
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

// DescribeAvailableDedicatedHostZonesRequest is the request struct for api DescribeAvailableDedicatedHostZones
type DescribeAvailableDedicatedHostZonesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeAvailableDedicatedHostZonesResponse is the response struct for api DescribeAvailableDedicatedHostZones
type DescribeAvailableDedicatedHostZonesResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Zones     Zones  `json:"Zones" xml:"Zones"`
}

// CreateDescribeAvailableDedicatedHostZonesRequest creates a request to invoke DescribeAvailableDedicatedHostZones API
func CreateDescribeAvailableDedicatedHostZonesRequest() (request *DescribeAvailableDedicatedHostZonesRequest) {
	request = &DescribeAvailableDedicatedHostZonesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeAvailableDedicatedHostZones", "", "")
	return
}

// CreateDescribeAvailableDedicatedHostZonesResponse creates a response to parse from DescribeAvailableDedicatedHostZones response
func CreateDescribeAvailableDedicatedHostZonesResponse() (response *DescribeAvailableDedicatedHostZonesResponse) {
	response = &DescribeAvailableDedicatedHostZonesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
