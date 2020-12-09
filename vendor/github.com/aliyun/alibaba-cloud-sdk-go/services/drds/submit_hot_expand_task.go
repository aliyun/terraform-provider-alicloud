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

// SubmitHotExpandTask invokes the drds.SubmitHotExpandTask API synchronously
func (client *Client) SubmitHotExpandTask(request *SubmitHotExpandTaskRequest) (response *SubmitHotExpandTaskResponse, err error) {
	response = CreateSubmitHotExpandTaskResponse()
	err = client.DoAction(request, response)
	return
}

// SubmitHotExpandTaskWithChan invokes the drds.SubmitHotExpandTask API asynchronously
func (client *Client) SubmitHotExpandTaskWithChan(request *SubmitHotExpandTaskRequest) (<-chan *SubmitHotExpandTaskResponse, <-chan error) {
	responseChan := make(chan *SubmitHotExpandTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SubmitHotExpandTask(request)
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

// SubmitHotExpandTaskWithCallback invokes the drds.SubmitHotExpandTask API asynchronously
func (client *Client) SubmitHotExpandTaskWithCallback(request *SubmitHotExpandTaskRequest, callback func(response *SubmitHotExpandTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SubmitHotExpandTaskResponse
		var err error
		defer close(result)
		response, err = client.SubmitHotExpandTask(request)
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

// SubmitHotExpandTaskRequest is the request struct for api SubmitHotExpandTask
type SubmitHotExpandTaskRequest struct {
	*requests.RpcRequest
	Mapping              *[]SubmitHotExpandTaskMapping              `position:"Query" name:"Mapping"  type:"Repeated"`
	TaskDesc             string                                     `position:"Query" name:"TaskDesc"`
	SupperAccountMapping *[]SubmitHotExpandTaskSupperAccountMapping `position:"Query" name:"SupperAccountMapping"  type:"Repeated"`
	ExtendedMapping      *[]SubmitHotExpandTaskExtendedMapping      `position:"Query" name:"ExtendedMapping"  type:"Repeated"`
	TaskName             string                                     `position:"Query" name:"TaskName"`
	DrdsInstanceId       string                                     `position:"Query" name:"DrdsInstanceId"`
	InstanceDbMapping    *[]SubmitHotExpandTaskInstanceDbMapping    `position:"Query" name:"InstanceDbMapping"  type:"Repeated"`
	DbName               string                                     `position:"Query" name:"DbName"`
}

// SubmitHotExpandTaskMapping is a repeated param struct in SubmitHotExpandTaskRequest
type SubmitHotExpandTaskMapping struct {
	DbShardColumn string `name:"DbShardColumn"`
	TbShardColumn string `name:"TbShardColumn"`
	ShardTbValue  string `name:"ShardTbValue"`
	HotDbName     string `name:"HotDbName"`
	ShardDbValue  string `name:"ShardDbValue"`
	HotTableName  string `name:"HotTableName"`
	LogicTable    string `name:"LogicTable"`
}

// SubmitHotExpandTaskSupperAccountMapping is a repeated param struct in SubmitHotExpandTaskRequest
type SubmitHotExpandTaskSupperAccountMapping struct {
	InstanceName   string `name:"InstanceName"`
	SupperAccount  string `name:"SupperAccount"`
	SupperPassword string `name:"SupperPassword"`
}

// SubmitHotExpandTaskExtendedMapping is a repeated param struct in SubmitHotExpandTaskRequest
type SubmitHotExpandTaskExtendedMapping struct {
	SrcInstanceId string `name:"SrcInstanceId"`
	SrcDb         string `name:"SrcDb"`
}

// SubmitHotExpandTaskInstanceDbMapping is a repeated param struct in SubmitHotExpandTaskRequest
type SubmitHotExpandTaskInstanceDbMapping struct {
	DbList       string `name:"DbList"`
	InstanceName string `name:"InstanceName"`
}

// SubmitHotExpandTaskResponse is the response struct for api SubmitHotExpandTask
type SubmitHotExpandTaskResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
}

// CreateSubmitHotExpandTaskRequest creates a request to invoke SubmitHotExpandTask API
func CreateSubmitHotExpandTaskRequest() (request *SubmitHotExpandTaskRequest) {
	request = &SubmitHotExpandTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2019-01-23", "SubmitHotExpandTask", "Drds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateSubmitHotExpandTaskResponse creates a response to parse from SubmitHotExpandTask response
func CreateSubmitHotExpandTaskResponse() (response *SubmitHotExpandTaskResponse) {
	response = &SubmitHotExpandTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
