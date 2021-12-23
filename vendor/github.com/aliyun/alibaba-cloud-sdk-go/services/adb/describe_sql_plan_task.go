package adb

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

// DescribeSQLPlanTask invokes the adb.DescribeSQLPlanTask API synchronously
func (client *Client) DescribeSQLPlanTask(request *DescribeSQLPlanTaskRequest) (response *DescribeSQLPlanTaskResponse, err error) {
	response = CreateDescribeSQLPlanTaskResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSQLPlanTaskWithChan invokes the adb.DescribeSQLPlanTask API asynchronously
func (client *Client) DescribeSQLPlanTaskWithChan(request *DescribeSQLPlanTaskRequest) (<-chan *DescribeSQLPlanTaskResponse, <-chan error) {
	responseChan := make(chan *DescribeSQLPlanTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSQLPlanTask(request)
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

// DescribeSQLPlanTaskWithCallback invokes the adb.DescribeSQLPlanTask API asynchronously
func (client *Client) DescribeSQLPlanTaskWithCallback(request *DescribeSQLPlanTaskRequest, callback func(response *DescribeSQLPlanTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSQLPlanTaskResponse
		var err error
		defer close(result)
		response, err = client.DescribeSQLPlanTask(request)
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

// DescribeSQLPlanTaskRequest is the request struct for api DescribeSQLPlanTask
type DescribeSQLPlanTaskRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ProcessId            string           `position:"Query" name:"ProcessId"`
	StageId              string           `position:"Query" name:"StageId"`
}

// DescribeSQLPlanTaskResponse is the response struct for api DescribeSQLPlanTask
type DescribeSQLPlanTaskResponse struct {
	*responses.BaseResponse
	RequestId string        `json:"RequestId" xml:"RequestId"`
	TaskList  []SqlPlanTask `json:"TaskList" xml:"TaskList"`
}

// CreateDescribeSQLPlanTaskRequest creates a request to invoke DescribeSQLPlanTask API
func CreateDescribeSQLPlanTaskRequest() (request *DescribeSQLPlanTaskRequest) {
	request = &DescribeSQLPlanTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("adb", "2019-03-15", "DescribeSQLPlanTask", "ads", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSQLPlanTaskResponse creates a response to parse from DescribeSQLPlanTask response
func CreateDescribeSQLPlanTaskResponse() (response *DescribeSQLPlanTaskResponse) {
	response = &DescribeSQLPlanTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
