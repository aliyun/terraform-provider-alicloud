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

// DescribeSystemLog invokes the ddoscoo.DescribeSystemLog API synchronously
func (client *Client) DescribeSystemLog(request *DescribeSystemLogRequest) (response *DescribeSystemLogResponse, err error) {
	response = CreateDescribeSystemLogResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSystemLogWithChan invokes the ddoscoo.DescribeSystemLog API asynchronously
func (client *Client) DescribeSystemLogWithChan(request *DescribeSystemLogRequest) (<-chan *DescribeSystemLogResponse, <-chan error) {
	responseChan := make(chan *DescribeSystemLogResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSystemLog(request)
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

// DescribeSystemLogWithCallback invokes the ddoscoo.DescribeSystemLog API asynchronously
func (client *Client) DescribeSystemLogWithCallback(request *DescribeSystemLogRequest, callback func(response *DescribeSystemLogResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSystemLogResponse
		var err error
		defer close(result)
		response, err = client.DescribeSystemLog(request)
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

// DescribeSystemLogRequest is the request struct for api DescribeSystemLog
type DescribeSystemLogRequest struct {
	*requests.RpcRequest
	StartTime    requests.Integer `position:"Query" name:"StartTime"`
	PageNumber   requests.Integer `position:"Query" name:"PageNumber"`
	SourceIp     string           `position:"Query" name:"SourceIp"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
	EndTime      requests.Integer `position:"Query" name:"EndTime"`
	EntityObject string           `position:"Query" name:"EntityObject"`
	EntityType   requests.Integer `position:"Query" name:"EntityType"`
}

// DescribeSystemLogResponse is the response struct for api DescribeSystemLog
type DescribeSystemLogResponse struct {
	*responses.BaseResponse
	Total     int64           `json:"Total" xml:"Total"`
	RequestId string          `json:"RequestId" xml:"RequestId"`
	SystemLog []SystemLogItem `json:"SystemLog" xml:"SystemLog"`
}

// CreateDescribeSystemLogRequest creates a request to invoke DescribeSystemLog API
func CreateDescribeSystemLogRequest() (request *DescribeSystemLogRequest) {
	request = &DescribeSystemLogRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeSystemLog", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSystemLogResponse creates a response to parse from DescribeSystemLog response
func CreateDescribeSystemLogResponse() (response *DescribeSystemLogResponse) {
	response = &DescribeSystemLogResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
