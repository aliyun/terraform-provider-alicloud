package vpc

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

// RemoveGlobalAccelerationInstanceIp invokes the vpc.RemoveGlobalAccelerationInstanceIp API synchronously
// api document: https://help.aliyun.com/api/vpc/removeglobalaccelerationinstanceip.html
func (client *Client) RemoveGlobalAccelerationInstanceIp(request *RemoveGlobalAccelerationInstanceIpRequest) (response *RemoveGlobalAccelerationInstanceIpResponse, err error) {
	response = CreateRemoveGlobalAccelerationInstanceIpResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveGlobalAccelerationInstanceIpWithChan invokes the vpc.RemoveGlobalAccelerationInstanceIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/removeglobalaccelerationinstanceip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveGlobalAccelerationInstanceIpWithChan(request *RemoveGlobalAccelerationInstanceIpRequest) (<-chan *RemoveGlobalAccelerationInstanceIpResponse, <-chan error) {
	responseChan := make(chan *RemoveGlobalAccelerationInstanceIpResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveGlobalAccelerationInstanceIp(request)
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

// RemoveGlobalAccelerationInstanceIpWithCallback invokes the vpc.RemoveGlobalAccelerationInstanceIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/removeglobalaccelerationinstanceip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveGlobalAccelerationInstanceIpWithCallback(request *RemoveGlobalAccelerationInstanceIpRequest, callback func(response *RemoveGlobalAccelerationInstanceIpResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveGlobalAccelerationInstanceIpResponse
		var err error
		defer close(result)
		response, err = client.RemoveGlobalAccelerationInstanceIp(request)
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

// RemoveGlobalAccelerationInstanceIpRequest is the request struct for api RemoveGlobalAccelerationInstanceIp
type RemoveGlobalAccelerationInstanceIpRequest struct {
	*requests.RpcRequest
	ResourceOwnerId              requests.Integer `position:"Query" name:"ResourceOwnerId"`
	GlobalAccelerationInstanceId string           `position:"Query" name:"GlobalAccelerationInstanceId"`
	ResourceOwnerAccount         string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                 string           `position:"Query" name:"OwnerAccount"`
	OwnerId                      requests.Integer `position:"Query" name:"OwnerId"`
	IpInstanceId                 string           `position:"Query" name:"IpInstanceId"`
}

// RemoveGlobalAccelerationInstanceIpResponse is the response struct for api RemoveGlobalAccelerationInstanceIp
type RemoveGlobalAccelerationInstanceIpResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRemoveGlobalAccelerationInstanceIpRequest creates a request to invoke RemoveGlobalAccelerationInstanceIp API
func CreateRemoveGlobalAccelerationInstanceIpRequest() (request *RemoveGlobalAccelerationInstanceIpRequest) {
	request = &RemoveGlobalAccelerationInstanceIpRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "RemoveGlobalAccelerationInstanceIp", "Vpc", "openAPI")
	return
}

// CreateRemoveGlobalAccelerationInstanceIpResponse creates a response to parse from RemoveGlobalAccelerationInstanceIp response
func CreateRemoveGlobalAccelerationInstanceIpResponse() (response *RemoveGlobalAccelerationInstanceIpResponse) {
	response = &RemoveGlobalAccelerationInstanceIpResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
