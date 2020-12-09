package drds

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

// ListVersions invokes the drds.ListVersions API synchronously
func (client *Client) ListVersions(request *ListVersionsRequest) (response *ListVersionsResponse, err error) {
	response = CreateListVersionsResponse()
	err = client.DoAction(request, response)
	return
}

// ListVersionsWithChan invokes the drds.ListVersions API asynchronously
func (client *Client) ListVersionsWithChan(request *ListVersionsRequest) (<-chan *ListVersionsResponse, <-chan error) {
	responseChan := make(chan *ListVersionsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListVersions(request)
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

// ListVersionsWithCallback invokes the drds.ListVersions API asynchronously
func (client *Client) ListVersionsWithCallback(request *ListVersionsRequest, callback func(response *ListVersionsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListVersionsResponse
		var err error
		defer close(result)
		response, err = client.ListVersions(request)
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

// ListVersionsRequest is the request struct for api ListVersions
type ListVersionsRequest struct {
	*requests.RpcRequest
	DrdsVer        string `position:"Query" name:"DrdsVer"`
	DrdsInstanceId string `position:"Query" name:"DrdsInstanceId"`
}

// ListVersionsResponse is the response struct for api ListVersions
type ListVersionsResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Success   bool     `json:"Success" xml:"Success"`
	Versions  Versions `json:"versions" xml:"versions"`
}

// CreateListVersionsRequest creates a request to invoke ListVersions API
func CreateListVersionsRequest() (request *ListVersionsRequest) {
	request = &ListVersionsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2019-01-23", "ListVersions", "Drds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateListVersionsResponse creates a response to parse from ListVersions response
func CreateListVersionsResponse() (response *ListVersionsResponse) {
	response = &ListVersionsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
