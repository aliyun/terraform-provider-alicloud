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

// DeleteQosCar invokes the smartag.DeleteQosCar API synchronously
// api document: https://help.aliyun.com/api/smartag/deleteqoscar.html
func (client *Client) DeleteQosCar(request *DeleteQosCarRequest) (response *DeleteQosCarResponse, err error) {
	response = CreateDeleteQosCarResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteQosCarWithChan invokes the smartag.DeleteQosCar API asynchronously
// api document: https://help.aliyun.com/api/smartag/deleteqoscar.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteQosCarWithChan(request *DeleteQosCarRequest) (<-chan *DeleteQosCarResponse, <-chan error) {
	responseChan := make(chan *DeleteQosCarResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteQosCar(request)
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

// DeleteQosCarWithCallback invokes the smartag.DeleteQosCar API asynchronously
// api document: https://help.aliyun.com/api/smartag/deleteqoscar.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteQosCarWithCallback(request *DeleteQosCarRequest, callback func(response *DeleteQosCarResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteQosCarResponse
		var err error
		defer close(result)
		response, err = client.DeleteQosCar(request)
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

// DeleteQosCarRequest is the request struct for api DeleteQosCar
type DeleteQosCarRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	QosId                string           `position:"Query" name:"QosId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	QosCarId             string           `position:"Query" name:"QosCarId"`
}

// DeleteQosCarResponse is the response struct for api DeleteQosCar
type DeleteQosCarResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteQosCarRequest creates a request to invoke DeleteQosCar API
func CreateDeleteQosCarRequest() (request *DeleteQosCarRequest) {
	request = &DeleteQosCarRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DeleteQosCar", "smartag", "openAPI")
	return
}

// CreateDeleteQosCarResponse creates a response to parse from DeleteQosCar response
func CreateDeleteQosCarResponse() (response *DeleteQosCarResponse) {
	response = &DeleteQosCarResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
