package resourcemanager

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

// CreatePolicyVersion invokes the resourcemanager.CreatePolicyVersion API synchronously
// api document: https://help.aliyun.com/api/resourcemanager/createpolicyversion.html
func (client *Client) CreatePolicyVersion(request *CreatePolicyVersionRequest) (response *CreatePolicyVersionResponse, err error) {
	response = CreateCreatePolicyVersionResponse()
	err = client.DoAction(request, response)
	return
}

// CreatePolicyVersionWithChan invokes the resourcemanager.CreatePolicyVersion API asynchronously
// api document: https://help.aliyun.com/api/resourcemanager/createpolicyversion.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreatePolicyVersionWithChan(request *CreatePolicyVersionRequest) (<-chan *CreatePolicyVersionResponse, <-chan error) {
	responseChan := make(chan *CreatePolicyVersionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreatePolicyVersion(request)
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

// CreatePolicyVersionWithCallback invokes the resourcemanager.CreatePolicyVersion API asynchronously
// api document: https://help.aliyun.com/api/resourcemanager/createpolicyversion.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreatePolicyVersionWithCallback(request *CreatePolicyVersionRequest, callback func(response *CreatePolicyVersionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreatePolicyVersionResponse
		var err error
		defer close(result)
		response, err = client.CreatePolicyVersion(request)
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

// CreatePolicyVersionRequest is the request struct for api CreatePolicyVersion
type CreatePolicyVersionRequest struct {
	*requests.RpcRequest
	SetAsDefault   requests.Boolean `position:"Query" name:"SetAsDefault"`
	PolicyName     string           `position:"Query" name:"PolicyName"`
	PolicyDocument string           `position:"Query" name:"PolicyDocument"`
}

// CreatePolicyVersionResponse is the response struct for api CreatePolicyVersion
type CreatePolicyVersionResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	PolicyVersion PolicyVersion `json:"PolicyVersion" xml:"PolicyVersion"`
}

// CreateCreatePolicyVersionRequest creates a request to invoke CreatePolicyVersion API
func CreateCreatePolicyVersionRequest() (request *CreatePolicyVersionRequest) {
	request = &CreatePolicyVersionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ResourceManager", "2020-03-31", "CreatePolicyVersion", "resourcemanager", "openAPI")
	return
}

// CreateCreatePolicyVersionResponse creates a response to parse from CreatePolicyVersion response
func CreateCreatePolicyVersionResponse() (response *CreatePolicyVersionResponse) {
	response = &CreatePolicyVersionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
