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

// CreateAsyncTask invokes the ddoscoo.CreateAsyncTask API synchronously
// api document: https://help.aliyun.com/api/ddoscoo/createasynctask.html
func (client *Client) CreateAsyncTask(request *CreateAsyncTaskRequest) (response *CreateAsyncTaskResponse, err error) {
	response = CreateCreateAsyncTaskResponse()
	err = client.DoAction(request, response)
	return
}

// CreateAsyncTaskWithChan invokes the ddoscoo.CreateAsyncTask API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/createasynctask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateAsyncTaskWithChan(request *CreateAsyncTaskRequest) (<-chan *CreateAsyncTaskResponse, <-chan error) {
	responseChan := make(chan *CreateAsyncTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateAsyncTask(request)
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

// CreateAsyncTaskWithCallback invokes the ddoscoo.CreateAsyncTask API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/createasynctask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateAsyncTaskWithCallback(request *CreateAsyncTaskRequest, callback func(response *CreateAsyncTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateAsyncTaskResponse
		var err error
		defer close(result)
		response, err = client.CreateAsyncTask(request)
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

// CreateAsyncTaskRequest is the request struct for api CreateAsyncTask
type CreateAsyncTaskRequest struct {
	*requests.RpcRequest
	TaskType        requests.Integer `position:"Query" name:"TaskType"`
	TaskParams      string           `position:"Query" name:"TaskParams"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SourceIp        string           `position:"Query" name:"SourceIp"`
	Lang            string           `position:"Query" name:"Lang"`
}

// CreateAsyncTaskResponse is the response struct for api CreateAsyncTask
type CreateAsyncTaskResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateAsyncTaskRequest creates a request to invoke CreateAsyncTask API
func CreateCreateAsyncTaskRequest() (request *CreateAsyncTaskRequest) {
	request = &CreateAsyncTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "CreateAsyncTask", "ddoscoo", "openAPI")
	return
}

// CreateCreateAsyncTaskResponse creates a response to parse from CreateAsyncTask response
func CreateCreateAsyncTaskResponse() (response *CreateAsyncTaskResponse) {
	response = &CreateAsyncTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
