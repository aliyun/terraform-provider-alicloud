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

// GetQosAttribute invokes the smartag.GetQosAttribute API synchronously
func (client *Client) GetQosAttribute(request *GetQosAttributeRequest) (response *GetQosAttributeResponse, err error) {
	response = CreateGetQosAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// GetQosAttributeWithChan invokes the smartag.GetQosAttribute API asynchronously
func (client *Client) GetQosAttributeWithChan(request *GetQosAttributeRequest) (<-chan *GetQosAttributeResponse, <-chan error) {
	responseChan := make(chan *GetQosAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetQosAttribute(request)
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

// GetQosAttributeWithCallback invokes the smartag.GetQosAttribute API asynchronously
func (client *Client) GetQosAttributeWithCallback(request *GetQosAttributeRequest, callback func(response *GetQosAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetQosAttributeResponse
		var err error
		defer close(result)
		response, err = client.GetQosAttribute(request)
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

// GetQosAttributeRequest is the request struct for api GetQosAttribute
type GetQosAttributeRequest struct {
	*requests.RpcRequest
	QosId string `position:"Query"`
}

// GetQosAttributeResponse is the response struct for api GetQosAttribute
type GetQosAttributeResponse struct {
	*responses.BaseResponse
	RequestId               string                       `json:"RequestId" xml:"RequestId"`
	ErrorConfigSmartAGCount int                          `json:"ErrorConfigSmartAGCount" xml:"ErrorConfigSmartAGCount"`
	QosName                 string                       `json:"QosName" xml:"QosName"`
	QosDescription          string                       `json:"QosDescription" xml:"QosDescription"`
	QosPolicies             []QosPolicyInGetQosAttribute `json:"QosPolicies" xml:"QosPolicies"`
	QosCars                 []QosCar                     `json:"QosCars" xml:"QosCars"`
}

// CreateGetQosAttributeRequest creates a request to invoke GetQosAttribute API
func CreateGetQosAttributeRequest() (request *GetQosAttributeRequest) {
	request = &GetQosAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "GetQosAttribute", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateGetQosAttributeResponse creates a response to parse from GetQosAttribute response
func CreateGetQosAttributeResponse() (response *GetQosAttributeResponse) {
	response = &GetQosAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
