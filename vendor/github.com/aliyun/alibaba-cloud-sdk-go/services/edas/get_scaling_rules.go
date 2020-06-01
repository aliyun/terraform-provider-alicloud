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

// GetScalingRules invokes the edas.GetScalingRules API synchronously
// api document: https://help.aliyun.com/api/edas/getscalingrules.html
func (client *Client) GetScalingRules(request *GetScalingRulesRequest) (response *GetScalingRulesResponse, err error) {
	response = CreateGetScalingRulesResponse()
	err = client.DoAction(request, response)
	return
}

// GetScalingRulesWithChan invokes the edas.GetScalingRules API asynchronously
// api document: https://help.aliyun.com/api/edas/getscalingrules.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetScalingRulesWithChan(request *GetScalingRulesRequest) (<-chan *GetScalingRulesResponse, <-chan error) {
	responseChan := make(chan *GetScalingRulesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetScalingRules(request)
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

// GetScalingRulesWithCallback invokes the edas.GetScalingRules API asynchronously
// api document: https://help.aliyun.com/api/edas/getscalingrules.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetScalingRulesWithCallback(request *GetScalingRulesRequest, callback func(response *GetScalingRulesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetScalingRulesResponse
		var err error
		defer close(result)
		response, err = client.GetScalingRules(request)
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

// GetScalingRulesRequest is the request struct for api GetScalingRules
type GetScalingRulesRequest struct {
	*requests.RoaRequest
	Mode    string `position:"Query" name:"Mode"`
	AppId   string `position:"Query" name:"AppId"`
	GroupId string `position:"Query" name:"GroupId"`
}

// GetScalingRulesResponse is the response struct for api GetScalingRules
type GetScalingRulesResponse struct {
	*responses.BaseResponse
	RequestId  string                `json:"RequestId" xml:"RequestId"`
	Code       int                   `json:"Code" xml:"Code"`
	Message    string                `json:"Message" xml:"Message"`
	UpdateTime int64                 `json:"UpdateTime" xml:"UpdateTime"`
	Data       DataInGetScalingRules `json:"Data" xml:"Data"`
}

// CreateGetScalingRulesRequest creates a request to invoke GetScalingRules API
func CreateGetScalingRulesRequest() (request *GetScalingRulesRequest) {
	request = &GetScalingRulesRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "GetScalingRules", "/pop/v5/app/scalingRules", "edas", "openAPI")
	request.Method = requests.GET
	return
}

// CreateGetScalingRulesResponse creates a response to parse from GetScalingRules response
func CreateGetScalingRulesResponse() (response *GetScalingRulesResponse) {
	response = &GetScalingRulesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
