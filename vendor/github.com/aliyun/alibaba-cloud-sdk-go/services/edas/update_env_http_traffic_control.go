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

// UpdateEnvHttpTrafficControl invokes the edas.UpdateEnvHttpTrafficControl API synchronously
func (client *Client) UpdateEnvHttpTrafficControl(request *UpdateEnvHttpTrafficControlRequest) (response *UpdateEnvHttpTrafficControlResponse, err error) {
	response = CreateUpdateEnvHttpTrafficControlResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateEnvHttpTrafficControlWithChan invokes the edas.UpdateEnvHttpTrafficControl API asynchronously
func (client *Client) UpdateEnvHttpTrafficControlWithChan(request *UpdateEnvHttpTrafficControlRequest) (<-chan *UpdateEnvHttpTrafficControlResponse, <-chan error) {
	responseChan := make(chan *UpdateEnvHttpTrafficControlResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateEnvHttpTrafficControl(request)
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

// UpdateEnvHttpTrafficControlWithCallback invokes the edas.UpdateEnvHttpTrafficControl API asynchronously
func (client *Client) UpdateEnvHttpTrafficControlWithCallback(request *UpdateEnvHttpTrafficControlRequest, callback func(response *UpdateEnvHttpTrafficControlResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateEnvHttpTrafficControlResponse
		var err error
		defer close(result)
		response, err = client.UpdateEnvHttpTrafficControl(request)
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

// UpdateEnvHttpTrafficControlRequest is the request struct for api UpdateEnvHttpTrafficControl
type UpdateEnvHttpTrafficControlRequest struct {
	*requests.RoaRequest
	Condition       string `position:"Body" name:"Condition"`
	UrlPath         string `position:"Body" name:"UrlPath"`
	AppId           string `position:"Body" name:"AppId"`
	LabelAdviceName string `position:"Body" name:"LabelAdviceName"`
	PointcutName    string `position:"Body" name:"PointcutName"`
	TriggerPolicy   string `position:"Body" name:"TriggerPolicy"`
}

// UpdateEnvHttpTrafficControlResponse is the response struct for api UpdateEnvHttpTrafficControl
type UpdateEnvHttpTrafficControlResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// CreateUpdateEnvHttpTrafficControlRequest creates a request to invoke UpdateEnvHttpTrafficControl API
func CreateUpdateEnvHttpTrafficControlRequest() (request *UpdateEnvHttpTrafficControlRequest) {
	request = &UpdateEnvHttpTrafficControlRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "UpdateEnvHttpTrafficControl", "/pop/v5/gray/env_http_traffic_control", "Edas", "openAPI")
	request.Method = requests.PUT
	return
}

// CreateUpdateEnvHttpTrafficControlResponse creates a response to parse from UpdateEnvHttpTrafficControl response
func CreateUpdateEnvHttpTrafficControlResponse() (response *UpdateEnvHttpTrafficControlResponse) {
	response = &UpdateEnvHttpTrafficControlResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
