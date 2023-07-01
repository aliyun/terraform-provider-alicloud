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

// ScaleInApplication invokes the edas.ScaleInApplication API synchronously
func (client *Client) ScaleInApplication(request *ScaleInApplicationRequest) (response *ScaleInApplicationResponse, err error) {
	response = CreateScaleInApplicationResponse()
	err = client.DoAction(request, response)
	return
}

// ScaleInApplicationWithChan invokes the edas.ScaleInApplication API asynchronously
func (client *Client) ScaleInApplicationWithChan(request *ScaleInApplicationRequest) (<-chan *ScaleInApplicationResponse, <-chan error) {
	responseChan := make(chan *ScaleInApplicationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ScaleInApplication(request)
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

// ScaleInApplicationWithCallback invokes the edas.ScaleInApplication API asynchronously
func (client *Client) ScaleInApplicationWithCallback(request *ScaleInApplicationRequest, callback func(response *ScaleInApplicationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ScaleInApplicationResponse
		var err error
		defer close(result)
		response, err = client.ScaleInApplication(request)
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

// ScaleInApplicationRequest is the request struct for api ScaleInApplication
type ScaleInApplicationRequest struct {
	*requests.RoaRequest
	ForceStatus string `position:"Query" name:"ForceStatus"`
	AppId       string `position:"Query" name:"AppId"`
	EccInfo     string `position:"Query" name:"EccInfo"`
}

// ScaleInApplicationResponse is the response struct for api ScaleInApplication
type ScaleInApplicationResponse struct {
	*responses.BaseResponse
	ChangeOrderId string `json:"ChangeOrderId" xml:"ChangeOrderId"`
	Code          int    `json:"Code" xml:"Code"`
	Message       string `json:"Message" xml:"Message"`
}

// CreateScaleInApplicationRequest creates a request to invoke ScaleInApplication API
func CreateScaleInApplicationRequest() (request *ScaleInApplicationRequest) {
	request = &ScaleInApplicationRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "ScaleInApplication", "/pop/v5/changeorder/co_scale_in", "Edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateScaleInApplicationResponse creates a response to parse from ScaleInApplication response
func CreateScaleInApplicationResponse() (response *ScaleInApplicationResponse) {
	response = &ScaleInApplicationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
