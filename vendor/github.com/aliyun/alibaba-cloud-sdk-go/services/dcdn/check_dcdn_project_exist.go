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

// CheckDcdnProjectExist invokes the dcdn.CheckDcdnProjectExist API synchronously
func (client *Client) CheckDcdnProjectExist(request *CheckDcdnProjectExistRequest) (response *CheckDcdnProjectExistResponse, err error) {
	response = CreateCheckDcdnProjectExistResponse()
	err = client.DoAction(request, response)
	return
}

// CheckDcdnProjectExistWithChan invokes the dcdn.CheckDcdnProjectExist API asynchronously
func (client *Client) CheckDcdnProjectExistWithChan(request *CheckDcdnProjectExistRequest) (<-chan *CheckDcdnProjectExistResponse, <-chan error) {
	responseChan := make(chan *CheckDcdnProjectExistResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CheckDcdnProjectExist(request)
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

// CheckDcdnProjectExistWithCallback invokes the dcdn.CheckDcdnProjectExist API asynchronously
func (client *Client) CheckDcdnProjectExistWithCallback(request *CheckDcdnProjectExistRequest, callback func(response *CheckDcdnProjectExistResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CheckDcdnProjectExistResponse
		var err error
		defer close(result)
		response, err = client.CheckDcdnProjectExist(request)
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

// CheckDcdnProjectExistRequest is the request struct for api CheckDcdnProjectExist
type CheckDcdnProjectExistRequest struct {
	*requests.RpcRequest
	ProjectName string           `position:"Query" name:"ProjectName"`
	OwnerId     requests.Integer `position:"Query" name:"OwnerId"`
}

// CheckDcdnProjectExistResponse is the response struct for api CheckDcdnProjectExist
type CheckDcdnProjectExistResponse struct {
	*responses.BaseResponse
	RequestId string  `json:"RequestId" xml:"RequestId"`
	Content   Content `json:"Content" xml:"Content"`
}

// CreateCheckDcdnProjectExistRequest creates a request to invoke CheckDcdnProjectExist API
func CreateCheckDcdnProjectExistRequest() (request *CheckDcdnProjectExistRequest) {
	request = &CheckDcdnProjectExistRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "CheckDcdnProjectExist", "", "")
	request.Method = requests.GET
	return
}

// CreateCheckDcdnProjectExistResponse creates a response to parse from CheckDcdnProjectExist response
func CreateCheckDcdnProjectExistResponse() (response *CheckDcdnProjectExistResponse) {
	response = &CheckDcdnProjectExistResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
