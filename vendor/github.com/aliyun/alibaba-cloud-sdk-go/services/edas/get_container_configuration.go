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

// GetContainerConfiguration invokes the edas.GetContainerConfiguration API synchronously
func (client *Client) GetContainerConfiguration(request *GetContainerConfigurationRequest) (response *GetContainerConfigurationResponse, err error) {
	response = CreateGetContainerConfigurationResponse()
	err = client.DoAction(request, response)
	return
}

// GetContainerConfigurationWithChan invokes the edas.GetContainerConfiguration API asynchronously
func (client *Client) GetContainerConfigurationWithChan(request *GetContainerConfigurationRequest) (<-chan *GetContainerConfigurationResponse, <-chan error) {
	responseChan := make(chan *GetContainerConfigurationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetContainerConfiguration(request)
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

// GetContainerConfigurationWithCallback invokes the edas.GetContainerConfiguration API asynchronously
func (client *Client) GetContainerConfigurationWithCallback(request *GetContainerConfigurationRequest, callback func(response *GetContainerConfigurationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetContainerConfigurationResponse
		var err error
		defer close(result)
		response, err = client.GetContainerConfiguration(request)
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

// GetContainerConfigurationRequest is the request struct for api GetContainerConfiguration
type GetContainerConfigurationRequest struct {
	*requests.RoaRequest
	AppId   string `position:"Query" name:"AppId"`
	GroupId string `position:"Query" name:"GroupId"`
}

// GetContainerConfigurationResponse is the response struct for api GetContainerConfiguration
type GetContainerConfigurationResponse struct {
	*responses.BaseResponse
	Code                   int                    `json:"Code" xml:"Code"`
	Message                string                 `json:"Message" xml:"Message"`
	RequestId              string                 `json:"RequestId" xml:"RequestId"`
	ContainerConfiguration ContainerConfiguration `json:"ContainerConfiguration" xml:"ContainerConfiguration"`
}

// CreateGetContainerConfigurationRequest creates a request to invoke GetContainerConfiguration API
func CreateGetContainerConfigurationRequest() (request *GetContainerConfigurationRequest) {
	request = &GetContainerConfigurationRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "GetContainerConfiguration", "/pop/v5/app/container_config", "edas", "openAPI")
	request.Method = requests.GET
	return
}

// CreateGetContainerConfigurationResponse creates a response to parse from GetContainerConfiguration response
func CreateGetContainerConfigurationResponse() (response *GetContainerConfigurationResponse) {
	response = &GetContainerConfigurationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
