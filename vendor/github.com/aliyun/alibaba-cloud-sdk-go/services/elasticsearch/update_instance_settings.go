package elasticsearch

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

// UpdateInstanceSettings invokes the elasticsearch.UpdateInstanceSettings API synchronously
// api document: https://help.aliyun.com/api/elasticsearch/updateinstancesettings.html
func (client *Client) UpdateInstanceSettings(request *UpdateInstanceSettingsRequest) (response *UpdateInstanceSettingsResponse, err error) {
	response = CreateUpdateInstanceSettingsResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateInstanceSettingsWithChan invokes the elasticsearch.UpdateInstanceSettings API asynchronously
// api document: https://help.aliyun.com/api/elasticsearch/updateinstancesettings.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateInstanceSettingsWithChan(request *UpdateInstanceSettingsRequest) (<-chan *UpdateInstanceSettingsResponse, <-chan error) {
	responseChan := make(chan *UpdateInstanceSettingsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateInstanceSettings(request)
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

// UpdateInstanceSettingsWithCallback invokes the elasticsearch.UpdateInstanceSettings API asynchronously
// api document: https://help.aliyun.com/api/elasticsearch/updateinstancesettings.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateInstanceSettingsWithCallback(request *UpdateInstanceSettingsRequest, callback func(response *UpdateInstanceSettingsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateInstanceSettingsResponse
		var err error
		defer close(result)
		response, err = client.UpdateInstanceSettings(request)
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

// UpdateInstanceSettingsRequest is the request struct for api UpdateInstanceSettings
type UpdateInstanceSettingsRequest struct {
	*requests.RoaRequest
	InstanceId  string `position:"Path" name:"InstanceId"`
	ClientToken string `position:"Query" name:"clientToken"`
}

// UpdateInstanceSettingsResponse is the response struct for api UpdateInstanceSettings
type UpdateInstanceSettingsResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateInstanceSettingsRequest creates a request to invoke UpdateInstanceSettings API
func CreateUpdateInstanceSettingsRequest() (request *UpdateInstanceSettingsRequest) {
	request = &UpdateInstanceSettingsRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "UpdateInstanceSettings", "/openapi/instances/[InstanceId]/instance-settings", "elasticsearch", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUpdateInstanceSettingsResponse creates a response to parse from UpdateInstanceSettings response
func CreateUpdateInstanceSettingsResponse() (response *UpdateInstanceSettingsResponse) {
	response = &UpdateInstanceSettingsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
