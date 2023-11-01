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

// DeleteSceneDefensePolicy invokes the ddoscoo.DeleteSceneDefensePolicy API synchronously
func (client *Client) DeleteSceneDefensePolicy(request *DeleteSceneDefensePolicyRequest) (response *DeleteSceneDefensePolicyResponse, err error) {
	response = CreateDeleteSceneDefensePolicyResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteSceneDefensePolicyWithChan invokes the ddoscoo.DeleteSceneDefensePolicy API asynchronously
func (client *Client) DeleteSceneDefensePolicyWithChan(request *DeleteSceneDefensePolicyRequest) (<-chan *DeleteSceneDefensePolicyResponse, <-chan error) {
	responseChan := make(chan *DeleteSceneDefensePolicyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteSceneDefensePolicy(request)
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

// DeleteSceneDefensePolicyWithCallback invokes the ddoscoo.DeleteSceneDefensePolicy API asynchronously
func (client *Client) DeleteSceneDefensePolicyWithCallback(request *DeleteSceneDefensePolicyRequest, callback func(response *DeleteSceneDefensePolicyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteSceneDefensePolicyResponse
		var err error
		defer close(result)
		response, err = client.DeleteSceneDefensePolicy(request)
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

// DeleteSceneDefensePolicyRequest is the request struct for api DeleteSceneDefensePolicy
type DeleteSceneDefensePolicyRequest struct {
	*requests.RpcRequest
	SourceIp string `position:"Query" name:"SourceIp"`
	PolicyId string `position:"Query" name:"PolicyId"`
}

// DeleteSceneDefensePolicyResponse is the response struct for api DeleteSceneDefensePolicy
type DeleteSceneDefensePolicyResponse struct {
	*responses.BaseResponse
	Success   bool   `json:"Success" xml:"Success"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteSceneDefensePolicyRequest creates a request to invoke DeleteSceneDefensePolicy API
func CreateDeleteSceneDefensePolicyRequest() (request *DeleteSceneDefensePolicyRequest) {
	request = &DeleteSceneDefensePolicyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DeleteSceneDefensePolicy", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDeleteSceneDefensePolicyResponse creates a response to parse from DeleteSceneDefensePolicy response
func CreateDeleteSceneDefensePolicyResponse() (response *DeleteSceneDefensePolicyResponse) {
	response = &DeleteSceneDefensePolicyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
