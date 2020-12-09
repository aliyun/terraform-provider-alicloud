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

// DescribeDataImportPreCheckResult invokes the drds.DescribeDataImportPreCheckResult API synchronously
func (client *Client) DescribeDataImportPreCheckResult(request *DescribeDataImportPreCheckResultRequest) (response *DescribeDataImportPreCheckResultResponse, err error) {
	response = CreateDescribeDataImportPreCheckResultResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDataImportPreCheckResultWithChan invokes the drds.DescribeDataImportPreCheckResult API asynchronously
func (client *Client) DescribeDataImportPreCheckResultWithChan(request *DescribeDataImportPreCheckResultRequest) (<-chan *DescribeDataImportPreCheckResultResponse, <-chan error) {
	responseChan := make(chan *DescribeDataImportPreCheckResultResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDataImportPreCheckResult(request)
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

// DescribeDataImportPreCheckResultWithCallback invokes the drds.DescribeDataImportPreCheckResult API asynchronously
func (client *Client) DescribeDataImportPreCheckResultWithCallback(request *DescribeDataImportPreCheckResultRequest, callback func(response *DescribeDataImportPreCheckResultResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDataImportPreCheckResultResponse
		var err error
		defer close(result)
		response, err = client.DescribeDataImportPreCheckResult(request)
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

// DescribeDataImportPreCheckResultRequest is the request struct for api DescribeDataImportPreCheckResult
type DescribeDataImportPreCheckResultRequest struct {
	*requests.RpcRequest
	TaskId requests.Integer `position:"Query" name:"TaskId"`
}

// DescribeDataImportPreCheckResultResponse is the response struct for api DescribeDataImportPreCheckResult
type DescribeDataImportPreCheckResultResponse struct {
	*responses.BaseResponse
	RequestId      string         `json:"RequestId" xml:"RequestId"`
	Success        bool           `json:"Success" xml:"Success"`
	PreCheckResult PreCheckResult `json:"PreCheckResult" xml:"PreCheckResult"`
}

// CreateDescribeDataImportPreCheckResultRequest creates a request to invoke DescribeDataImportPreCheckResult API
func CreateDescribeDataImportPreCheckResultRequest() (request *DescribeDataImportPreCheckResultRequest) {
	request = &DescribeDataImportPreCheckResultRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2019-01-23", "DescribeDataImportPreCheckResult", "Drds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDataImportPreCheckResultResponse creates a response to parse from DescribeDataImportPreCheckResult response
func CreateDescribeDataImportPreCheckResultResponse() (response *DescribeDataImportPreCheckResultResponse) {
	response = &DescribeDataImportPreCheckResultResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
