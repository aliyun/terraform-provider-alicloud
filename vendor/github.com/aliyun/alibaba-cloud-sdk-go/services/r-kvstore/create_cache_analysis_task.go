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

// CreateCacheAnalysisTask invokes the r_kvstore.CreateCacheAnalysisTask API synchronously
// api document: https://help.aliyun.com/api/r-kvstore/createcacheanalysistask.html
func (client *Client) CreateCacheAnalysisTask(request *CreateCacheAnalysisTaskRequest) (response *CreateCacheAnalysisTaskResponse, err error) {
	response = CreateCreateCacheAnalysisTaskResponse()
	err = client.DoAction(request, response)
	return
}

// CreateCacheAnalysisTaskWithChan invokes the r_kvstore.CreateCacheAnalysisTask API asynchronously
// api document: https://help.aliyun.com/api/r-kvstore/createcacheanalysistask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCacheAnalysisTaskWithChan(request *CreateCacheAnalysisTaskRequest) (<-chan *CreateCacheAnalysisTaskResponse, <-chan error) {
	responseChan := make(chan *CreateCacheAnalysisTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateCacheAnalysisTask(request)
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

// CreateCacheAnalysisTaskWithCallback invokes the r_kvstore.CreateCacheAnalysisTask API asynchronously
// api document: https://help.aliyun.com/api/r-kvstore/createcacheanalysistask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCacheAnalysisTaskWithCallback(request *CreateCacheAnalysisTaskRequest, callback func(response *CreateCacheAnalysisTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateCacheAnalysisTaskResponse
		var err error
		defer close(result)
		response, err = client.CreateCacheAnalysisTask(request)
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

// CreateCacheAnalysisTaskRequest is the request struct for api CreateCacheAnalysisTask
type CreateCacheAnalysisTaskRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// CreateCacheAnalysisTaskResponse is the response struct for api CreateCacheAnalysisTask
type CreateCacheAnalysisTaskResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateCacheAnalysisTaskRequest creates a request to invoke CreateCacheAnalysisTask API
func CreateCreateCacheAnalysisTaskRequest() (request *CreateCacheAnalysisTaskRequest) {
	request = &CreateCacheAnalysisTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "CreateCacheAnalysisTask", "redisa", "openAPI")
	return
}

// CreateCreateCacheAnalysisTaskResponse creates a response to parse from CreateCacheAnalysisTask response
func CreateCreateCacheAnalysisTaskResponse() (response *CreateCacheAnalysisTaskResponse) {
	response = &CreateCacheAnalysisTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
