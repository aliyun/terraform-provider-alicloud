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

// RemoveSubInstance invokes the r_kvstore.RemoveSubInstance API synchronously
func (client *Client) RemoveSubInstance(request *RemoveSubInstanceRequest) (response *RemoveSubInstanceResponse, err error) {
	response = CreateRemoveSubInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveSubInstanceWithChan invokes the r_kvstore.RemoveSubInstance API asynchronously
func (client *Client) RemoveSubInstanceWithChan(request *RemoveSubInstanceRequest) (<-chan *RemoveSubInstanceResponse, <-chan error) {
	responseChan := make(chan *RemoveSubInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveSubInstance(request)
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

// RemoveSubInstanceWithCallback invokes the r_kvstore.RemoveSubInstance API asynchronously
func (client *Client) RemoveSubInstanceWithCallback(request *RemoveSubInstanceRequest, callback func(response *RemoveSubInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveSubInstanceResponse
		var err error
		defer close(result)
		response, err = client.RemoveSubInstance(request)
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

// RemoveSubInstanceRequest is the request struct for api RemoveSubInstance
type RemoveSubInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ReleaseSubInstance   requests.Boolean `position:"Query" name:"ReleaseSubInstance"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	InstanceId           string           `position:"Query" name:"InstanceId"`
}

// RemoveSubInstanceResponse is the response struct for api RemoveSubInstance
type RemoveSubInstanceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRemoveSubInstanceRequest creates a request to invoke RemoveSubInstance API
func CreateRemoveSubInstanceRequest() (request *RemoveSubInstanceRequest) {
	request = &RemoveSubInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "RemoveSubInstance", "redisa", "openAPI")
	request.Method = requests.POST
	return
}

// CreateRemoveSubInstanceResponse creates a response to parse from RemoveSubInstance response
func CreateRemoveSubInstanceResponse() (response *RemoveSubInstanceResponse) {
	response = &RemoveSubInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
