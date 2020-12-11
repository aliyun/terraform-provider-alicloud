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

// ScaleoutApplicationWithNewInstances invokes the edas.ScaleoutApplicationWithNewInstances API synchronously
func (client *Client) ScaleoutApplicationWithNewInstances(request *ScaleoutApplicationWithNewInstancesRequest) (response *ScaleoutApplicationWithNewInstancesResponse, err error) {
	response = CreateScaleoutApplicationWithNewInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// ScaleoutApplicationWithNewInstancesWithChan invokes the edas.ScaleoutApplicationWithNewInstances API asynchronously
func (client *Client) ScaleoutApplicationWithNewInstancesWithChan(request *ScaleoutApplicationWithNewInstancesRequest) (<-chan *ScaleoutApplicationWithNewInstancesResponse, <-chan error) {
	responseChan := make(chan *ScaleoutApplicationWithNewInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ScaleoutApplicationWithNewInstances(request)
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

// ScaleoutApplicationWithNewInstancesWithCallback invokes the edas.ScaleoutApplicationWithNewInstances API asynchronously
func (client *Client) ScaleoutApplicationWithNewInstancesWithCallback(request *ScaleoutApplicationWithNewInstancesRequest, callback func(response *ScaleoutApplicationWithNewInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ScaleoutApplicationWithNewInstancesResponse
		var err error
		defer close(result)
		response, err = client.ScaleoutApplicationWithNewInstances(request)
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

// ScaleoutApplicationWithNewInstancesRequest is the request struct for api ScaleoutApplicationWithNewInstances
type ScaleoutApplicationWithNewInstancesRequest struct {
	*requests.RoaRequest
	TemplateVersion    string           `position:"Query" name:"TemplateVersion"`
	TemplateInstanceId string           `position:"Query" name:"TemplateInstanceId"`
	AppId              string           `position:"Query" name:"AppId"`
	GroupId            string           `position:"Query" name:"GroupId"`
	ScalingNum         requests.Integer `position:"Query" name:"ScalingNum"`
	TemplateId         string           `position:"Query" name:"TemplateId"`
	ScalingPolicy      string           `position:"Query" name:"ScalingPolicy"`
}

// ScaleoutApplicationWithNewInstancesResponse is the response struct for api ScaleoutApplicationWithNewInstances
type ScaleoutApplicationWithNewInstancesResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	Code          int    `json:"Code" xml:"Code"`
	Message       string `json:"Message" xml:"Message"`
	ChangeOrderId string `json:"ChangeOrderId" xml:"ChangeOrderId"`
}

// CreateScaleoutApplicationWithNewInstancesRequest creates a request to invoke ScaleoutApplicationWithNewInstances API
func CreateScaleoutApplicationWithNewInstancesRequest() (request *ScaleoutApplicationWithNewInstancesRequest) {
	request = &ScaleoutApplicationWithNewInstancesRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "ScaleoutApplicationWithNewInstances", "/pop/v5/scaling/scale_out", "edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateScaleoutApplicationWithNewInstancesResponse creates a response to parse from ScaleoutApplicationWithNewInstances response
func CreateScaleoutApplicationWithNewInstancesResponse() (response *ScaleoutApplicationWithNewInstancesResponse) {
	response = &ScaleoutApplicationWithNewInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
