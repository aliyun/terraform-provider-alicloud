package smartag

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

// DeleteQosPolicy invokes the smartag.DeleteQosPolicy API synchronously
func (client *Client) DeleteQosPolicy(request *DeleteQosPolicyRequest) (response *DeleteQosPolicyResponse, err error) {
	response = CreateDeleteQosPolicyResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteQosPolicyWithChan invokes the smartag.DeleteQosPolicy API asynchronously
func (client *Client) DeleteQosPolicyWithChan(request *DeleteQosPolicyRequest) (<-chan *DeleteQosPolicyResponse, <-chan error) {
	responseChan := make(chan *DeleteQosPolicyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteQosPolicy(request)
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

// DeleteQosPolicyWithCallback invokes the smartag.DeleteQosPolicy API asynchronously
func (client *Client) DeleteQosPolicyWithCallback(request *DeleteQosPolicyRequest, callback func(response *DeleteQosPolicyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteQosPolicyResponse
		var err error
		defer close(result)
		response, err = client.DeleteQosPolicy(request)
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

// DeleteQosPolicyRequest is the request struct for api DeleteQosPolicy
type DeleteQosPolicyRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	QosPolicyId          string           `position:"Query" name:"QosPolicyId"`
	QosId                string           `position:"Query" name:"QosId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DeleteQosPolicyResponse is the response struct for api DeleteQosPolicy
type DeleteQosPolicyResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteQosPolicyRequest creates a request to invoke DeleteQosPolicy API
func CreateDeleteQosPolicyRequest() (request *DeleteQosPolicyRequest) {
	request = &DeleteQosPolicyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DeleteQosPolicy", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDeleteQosPolicyResponse creates a response to parse from DeleteQosPolicy response
func CreateDeleteQosPolicyResponse() (response *DeleteQosPolicyResponse) {
	response = &DeleteQosPolicyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
