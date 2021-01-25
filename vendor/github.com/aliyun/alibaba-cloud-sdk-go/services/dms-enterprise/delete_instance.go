package dms_enterprise

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

// DeleteInstance invokes the dms_enterprise.DeleteInstance API synchronously
func (client *Client) DeleteInstance(request *DeleteInstanceRequest) (response *DeleteInstanceResponse, err error) {
	response = CreateDeleteInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteInstanceWithChan invokes the dms_enterprise.DeleteInstance API asynchronously
func (client *Client) DeleteInstanceWithChan(request *DeleteInstanceRequest) (<-chan *DeleteInstanceResponse, <-chan error) {
	responseChan := make(chan *DeleteInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteInstance(request)
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

// DeleteInstanceWithCallback invokes the dms_enterprise.DeleteInstance API asynchronously
func (client *Client) DeleteInstanceWithCallback(request *DeleteInstanceRequest, callback func(response *DeleteInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteInstanceResponse
		var err error
		defer close(result)
		response, err = client.DeleteInstance(request)
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

// DeleteInstanceRequest is the request struct for api DeleteInstance
type DeleteInstanceRequest struct {
	*requests.RpcRequest
	Tid  requests.Integer `position:"Query" name:"Tid"`
	Sid  string           `position:"Query" name:"Sid"`
	Port requests.Integer `position:"Query" name:"Port"`
	Host string           `position:"Query" name:"Host"`
}

// DeleteInstanceResponse is the response struct for api DeleteInstance
type DeleteInstanceResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	Success      bool   `json:"Success" xml:"Success"`
	ErrorMessage string `json:"ErrorMessage" xml:"ErrorMessage"`
	ErrorCode    string `json:"ErrorCode" xml:"ErrorCode"`
}

// CreateDeleteInstanceRequest creates a request to invoke DeleteInstance API
func CreateDeleteInstanceRequest() (request *DeleteInstanceRequest) {
	request = &DeleteInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dms-enterprise", "2018-11-01", "DeleteInstance", "dmsenterprise", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDeleteInstanceResponse creates a response to parse from DeleteInstance response
func CreateDeleteInstanceResponse() (response *DeleteInstanceResponse) {
	response = &DeleteInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
