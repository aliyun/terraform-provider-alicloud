package dcdn

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

// PutDcdnKv invokes the dcdn.PutDcdnKv API synchronously
func (client *Client) PutDcdnKv(request *PutDcdnKvRequest) (response *PutDcdnKvResponse, err error) {
	response = CreatePutDcdnKvResponse()
	err = client.DoAction(request, response)
	return
}

// PutDcdnKvWithChan invokes the dcdn.PutDcdnKv API asynchronously
func (client *Client) PutDcdnKvWithChan(request *PutDcdnKvRequest) (<-chan *PutDcdnKvResponse, <-chan error) {
	responseChan := make(chan *PutDcdnKvResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PutDcdnKv(request)
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

// PutDcdnKvWithCallback invokes the dcdn.PutDcdnKv API asynchronously
func (client *Client) PutDcdnKvWithCallback(request *PutDcdnKvRequest, callback func(response *PutDcdnKvResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PutDcdnKvResponse
		var err error
		defer close(result)
		response, err = client.PutDcdnKv(request)
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

// PutDcdnKvRequest is the request struct for api PutDcdnKv
type PutDcdnKvRequest struct {
	*requests.RpcRequest
	Namespace string `position:"Query" name:"Namespace"`
	Value     string `position:"Body" name:"Value"`
	Key       string `position:"Query" name:"Key"`
}

// PutDcdnKvResponse is the response struct for api PutDcdnKv
type PutDcdnKvResponse struct {
	*responses.BaseResponse
	Length    int    `json:"Length" xml:"Length"`
	Value     string `json:"Value" xml:"Value"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreatePutDcdnKvRequest creates a request to invoke PutDcdnKv API
func CreatePutDcdnKvRequest() (request *PutDcdnKvRequest) {
	request = &PutDcdnKvRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "PutDcdnKv", "", "")
	request.Method = requests.POST
	return
}

// CreatePutDcdnKvResponse creates a response to parse from PutDcdnKv response
func CreatePutDcdnKvResponse() (response *PutDcdnKvResponse) {
	response = &PutDcdnKvResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}