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

// ModifyWebPreciseAccessSwitch invokes the ddoscoo.ModifyWebPreciseAccessSwitch API synchronously
func (client *Client) ModifyWebPreciseAccessSwitch(request *ModifyWebPreciseAccessSwitchRequest) (response *ModifyWebPreciseAccessSwitchResponse, err error) {
	response = CreateModifyWebPreciseAccessSwitchResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyWebPreciseAccessSwitchWithChan invokes the ddoscoo.ModifyWebPreciseAccessSwitch API asynchronously
func (client *Client) ModifyWebPreciseAccessSwitchWithChan(request *ModifyWebPreciseAccessSwitchRequest) (<-chan *ModifyWebPreciseAccessSwitchResponse, <-chan error) {
	responseChan := make(chan *ModifyWebPreciseAccessSwitchResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyWebPreciseAccessSwitch(request)
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

// ModifyWebPreciseAccessSwitchWithCallback invokes the ddoscoo.ModifyWebPreciseAccessSwitch API asynchronously
func (client *Client) ModifyWebPreciseAccessSwitchWithCallback(request *ModifyWebPreciseAccessSwitchRequest, callback func(response *ModifyWebPreciseAccessSwitchResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyWebPreciseAccessSwitchResponse
		var err error
		defer close(result)
		response, err = client.ModifyWebPreciseAccessSwitch(request)
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

// ModifyWebPreciseAccessSwitchRequest is the request struct for api ModifyWebPreciseAccessSwitch
type ModifyWebPreciseAccessSwitchRequest struct {
	*requests.RpcRequest
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
	Domain          string `position:"Query" name:"Domain"`
	Config          string `position:"Query" name:"Config"`
}

// ModifyWebPreciseAccessSwitchResponse is the response struct for api ModifyWebPreciseAccessSwitch
type ModifyWebPreciseAccessSwitchResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyWebPreciseAccessSwitchRequest creates a request to invoke ModifyWebPreciseAccessSwitch API
func CreateModifyWebPreciseAccessSwitchRequest() (request *ModifyWebPreciseAccessSwitchRequest) {
	request = &ModifyWebPreciseAccessSwitchRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "ModifyWebPreciseAccessSwitch", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyWebPreciseAccessSwitchResponse creates a response to parse from ModifyWebPreciseAccessSwitch response
func CreateModifyWebPreciseAccessSwitchResponse() (response *ModifyWebPreciseAccessSwitchResponse) {
	response = &ModifyWebPreciseAccessSwitchResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
