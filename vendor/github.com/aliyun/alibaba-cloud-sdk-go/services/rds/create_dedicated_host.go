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

// CreateDedicatedHost invokes the rds.CreateDedicatedHost API synchronously
// api document: https://help.aliyun.com/api/rds/creatededicatedhost.html
func (client *Client) CreateDedicatedHost(request *CreateDedicatedHostRequest) (response *CreateDedicatedHostResponse, err error) {
	response = CreateCreateDedicatedHostResponse()
	err = client.DoAction(request, response)
	return
}

// CreateDedicatedHostWithChan invokes the rds.CreateDedicatedHost API asynchronously
// api document: https://help.aliyun.com/api/rds/creatededicatedhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDedicatedHostWithChan(request *CreateDedicatedHostRequest) (<-chan *CreateDedicatedHostResponse, <-chan error) {
	responseChan := make(chan *CreateDedicatedHostResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateDedicatedHost(request)
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

// CreateDedicatedHostWithCallback invokes the rds.CreateDedicatedHost API asynchronously
// api document: https://help.aliyun.com/api/rds/creatededicatedhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDedicatedHostWithCallback(request *CreateDedicatedHostRequest, callback func(response *CreateDedicatedHostResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateDedicatedHostResponse
		var err error
		defer close(result)
		response, err = client.CreateDedicatedHost(request)
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

// CreateDedicatedHostRequest is the request struct for api CreateDedicatedHost
type CreateDedicatedHostRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	HostName             string           `position:"Query" name:"HostName"`
	DedicatedHostGroupId string           `position:"Query" name:"DedicatedHostGroupId"`
	Period               string           `position:"Query" name:"Period"`
	HostClass            string           `position:"Query" name:"HostClass"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	UsedTime             string           `position:"Query" name:"UsedTime"`
	VSwitchId            string           `position:"Query" name:"VSwitchId"`
	AutoRenew            string           `position:"Query" name:"AutoRenew"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	PayType              string           `position:"Query" name:"PayType"`
}

// CreateDedicatedHostResponse is the response struct for api CreateDedicatedHost
type CreateDedicatedHostResponse struct {
	*responses.BaseResponse
	RequestId        string           `json:"RequestId" xml:"RequestId"`
	OrderId          int64            `json:"OrderId" xml:"OrderId"`
	DedicateHostList DedicateHostList `json:"DedicateHostList" xml:"DedicateHostList"`
}

// CreateCreateDedicatedHostRequest creates a request to invoke CreateDedicatedHost API
func CreateCreateDedicatedHostRequest() (request *CreateDedicatedHostRequest) {
	request = &CreateDedicatedHostRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "CreateDedicatedHost", "Rds", "openAPI")
	return
}

// CreateCreateDedicatedHostResponse creates a response to parse from CreateDedicatedHost response
func CreateCreateDedicatedHostResponse() (response *CreateDedicatedHostResponse) {
	response = &CreateDedicatedHostResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
