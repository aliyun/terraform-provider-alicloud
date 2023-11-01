package ddosbgp

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

// DeleteBlackhole invokes the ddosbgp.DeleteBlackhole API synchronously
func (client *Client) DeleteBlackhole(request *DeleteBlackholeRequest) (response *DeleteBlackholeResponse, err error) {
	response = CreateDeleteBlackholeResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteBlackholeWithChan invokes the ddosbgp.DeleteBlackhole API asynchronously
func (client *Client) DeleteBlackholeWithChan(request *DeleteBlackholeRequest) (<-chan *DeleteBlackholeResponse, <-chan error) {
	responseChan := make(chan *DeleteBlackholeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteBlackhole(request)
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

// DeleteBlackholeWithCallback invokes the ddosbgp.DeleteBlackhole API asynchronously
func (client *Client) DeleteBlackholeWithCallback(request *DeleteBlackholeRequest, callback func(response *DeleteBlackholeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteBlackholeResponse
		var err error
		defer close(result)
		response, err = client.DeleteBlackhole(request)
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

// DeleteBlackholeRequest is the request struct for api DeleteBlackhole
type DeleteBlackholeRequest struct {
	*requests.RpcRequest
	Ip              string `position:"Query" name:"Ip"`
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	InstanceId      string `position:"Query" name:"InstanceId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
}

// DeleteBlackholeResponse is the response struct for api DeleteBlackhole
type DeleteBlackholeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteBlackholeRequest creates a request to invoke DeleteBlackhole API
func CreateDeleteBlackholeRequest() (request *DeleteBlackholeRequest) {
	request = &DeleteBlackholeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddosbgp", "2018-07-20", "DeleteBlackhole", "", "")
	request.Method = requests.POST
	return
}

// CreateDeleteBlackholeResponse creates a response to parse from DeleteBlackhole response
func CreateDeleteBlackholeResponse() (response *DeleteBlackholeResponse) {
	response = &DeleteBlackholeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
