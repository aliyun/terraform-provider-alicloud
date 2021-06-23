package cr

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

// GetImageLayer invokes the cr.GetImageLayer API synchronously
func (client *Client) GetImageLayer(request *GetImageLayerRequest) (response *GetImageLayerResponse, err error) {
	response = CreateGetImageLayerResponse()
	err = client.DoAction(request, response)
	return
}

// GetImageLayerWithChan invokes the cr.GetImageLayer API asynchronously
func (client *Client) GetImageLayerWithChan(request *GetImageLayerRequest) (<-chan *GetImageLayerResponse, <-chan error) {
	responseChan := make(chan *GetImageLayerResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetImageLayer(request)
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

// GetImageLayerWithCallback invokes the cr.GetImageLayer API asynchronously
func (client *Client) GetImageLayerWithCallback(request *GetImageLayerRequest, callback func(response *GetImageLayerResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetImageLayerResponse
		var err error
		defer close(result)
		response, err = client.GetImageLayer(request)
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

// GetImageLayerRequest is the request struct for api GetImageLayer
type GetImageLayerRequest struct {
	*requests.RoaRequest
	RepoNamespace string `position:"Path" name:"RepoNamespace"`
	RepoName      string `position:"Path" name:"RepoName"`
	Tag           string `position:"Path" name:"Tag"`
}

// GetImageLayerResponse is the response struct for api GetImageLayer
type GetImageLayerResponse struct {
	*responses.BaseResponse
}

// CreateGetImageLayerRequest creates a request to invoke GetImageLayer API
func CreateGetImageLayerRequest() (request *GetImageLayerRequest) {
	request = &GetImageLayerRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("cr", "2016-06-07", "GetImageLayer", "/repos/[RepoNamespace]/[RepoName]/tags/[Tag]/layers", "acr", "openAPI")
	request.Method = requests.GET
	return
}

// CreateGetImageLayerResponse creates a response to parse from GetImageLayer response
func CreateGetImageLayerResponse() (response *GetImageLayerResponse) {
	response = &GetImageLayerResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
