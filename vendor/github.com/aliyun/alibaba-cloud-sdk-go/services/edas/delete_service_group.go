package edas

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

// DeleteServiceGroup invokes the edas.DeleteServiceGroup API synchronously
func (client *Client) DeleteServiceGroup(request *DeleteServiceGroupRequest) (response *DeleteServiceGroupResponse, err error) {
	response = CreateDeleteServiceGroupResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteServiceGroupWithChan invokes the edas.DeleteServiceGroup API asynchronously
func (client *Client) DeleteServiceGroupWithChan(request *DeleteServiceGroupRequest) (<-chan *DeleteServiceGroupResponse, <-chan error) {
	responseChan := make(chan *DeleteServiceGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteServiceGroup(request)
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

// DeleteServiceGroupWithCallback invokes the edas.DeleteServiceGroup API asynchronously
func (client *Client) DeleteServiceGroupWithCallback(request *DeleteServiceGroupRequest, callback func(response *DeleteServiceGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteServiceGroupResponse
		var err error
		defer close(result)
		response, err = client.DeleteServiceGroup(request)
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

// DeleteServiceGroupRequest is the request struct for api DeleteServiceGroup
type DeleteServiceGroupRequest struct {
	*requests.RoaRequest
	GroupId string `position:"Query" name:"GroupId"`
}

// DeleteServiceGroupResponse is the response struct for api DeleteServiceGroup
type DeleteServiceGroupResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteServiceGroupRequest creates a request to invoke DeleteServiceGroup API
func CreateDeleteServiceGroupRequest() (request *DeleteServiceGroupRequest) {
	request = &DeleteServiceGroupRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "DeleteServiceGroup", "/pop/v5/service/serviceGroups", "Edas", "openAPI")
	request.Method = requests.DELETE
	return
}

// CreateDeleteServiceGroupResponse creates a response to parse from DeleteServiceGroup response
func CreateDeleteServiceGroupResponse() (response *DeleteServiceGroupResponse) {
	response = &DeleteServiceGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
