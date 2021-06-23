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

// GetAppDeployment invokes the edas.GetAppDeployment API synchronously
func (client *Client) GetAppDeployment(request *GetAppDeploymentRequest) (response *GetAppDeploymentResponse, err error) {
	response = CreateGetAppDeploymentResponse()
	err = client.DoAction(request, response)
	return
}

// GetAppDeploymentWithChan invokes the edas.GetAppDeployment API asynchronously
func (client *Client) GetAppDeploymentWithChan(request *GetAppDeploymentRequest) (<-chan *GetAppDeploymentResponse, <-chan error) {
	responseChan := make(chan *GetAppDeploymentResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetAppDeployment(request)
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

// GetAppDeploymentWithCallback invokes the edas.GetAppDeployment API asynchronously
func (client *Client) GetAppDeploymentWithCallback(request *GetAppDeploymentRequest, callback func(response *GetAppDeploymentResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetAppDeploymentResponse
		var err error
		defer close(result)
		response, err = client.GetAppDeployment(request)
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

// GetAppDeploymentRequest is the request struct for api GetAppDeployment
type GetAppDeploymentRequest struct {
	*requests.RoaRequest
	AppId string `position:"Query" name:"AppId"`
}

// GetAppDeploymentResponse is the response struct for api GetAppDeployment
type GetAppDeploymentResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      string `json:"Data" xml:"Data"`
}

// CreateGetAppDeploymentRequest creates a request to invoke GetAppDeployment API
func CreateGetAppDeploymentRequest() (request *GetAppDeploymentRequest) {
	request = &GetAppDeploymentRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "GetAppDeployment", "/pop/v5/oam/app_deployment", "Edas", "openAPI")
	request.Method = requests.GET
	return
}

// CreateGetAppDeploymentResponse creates a response to parse from GetAppDeployment response
func CreateGetAppDeploymentResponse() (response *GetAppDeploymentResponse) {
	response = &GetAppDeploymentResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
