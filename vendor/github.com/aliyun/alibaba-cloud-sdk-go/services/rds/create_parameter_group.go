package rds

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

// CreateParameterGroup invokes the rds.CreateParameterGroup API synchronously
func (client *Client) CreateParameterGroup(request *CreateParameterGroupRequest) (response *CreateParameterGroupResponse, err error) {
	response = CreateCreateParameterGroupResponse()
	err = client.DoAction(request, response)
	return
}

// CreateParameterGroupWithChan invokes the rds.CreateParameterGroup API asynchronously
func (client *Client) CreateParameterGroupWithChan(request *CreateParameterGroupRequest) (<-chan *CreateParameterGroupResponse, <-chan error) {
	responseChan := make(chan *CreateParameterGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateParameterGroup(request)
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

// CreateParameterGroupWithCallback invokes the rds.CreateParameterGroup API asynchronously
func (client *Client) CreateParameterGroupWithCallback(request *CreateParameterGroupRequest, callback func(response *CreateParameterGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateParameterGroupResponse
		var err error
		defer close(result)
		response, err = client.CreateParameterGroup(request)
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

// CreateParameterGroupRequest is the request struct for api CreateParameterGroup
type CreateParameterGroupRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	EngineVersion        string           `position:"Query" name:"EngineVersion"`
	Engine               string           `position:"Query" name:"Engine"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ParameterGroupName   string           `position:"Query" name:"ParameterGroupName"`
	Parameters           string           `position:"Query" name:"Parameters"`
	ParameterGroupDesc   string           `position:"Query" name:"ParameterGroupDesc"`
}

// CreateParameterGroupResponse is the response struct for api CreateParameterGroup
type CreateParameterGroupResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateParameterGroupRequest creates a request to invoke CreateParameterGroup API
func CreateCreateParameterGroupRequest() (request *CreateParameterGroupRequest) {
	request = &CreateParameterGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "CreateParameterGroup", "rds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateParameterGroupResponse creates a response to parse from CreateParameterGroup response
func CreateCreateParameterGroupResponse() (response *CreateParameterGroupResponse) {
	response = &CreateParameterGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
