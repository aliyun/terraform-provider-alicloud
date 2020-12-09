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

// UninstallKibanaPlugin invokes the elasticsearch.UninstallKibanaPlugin API synchronously
func (client *Client) UninstallKibanaPlugin(request *UninstallKibanaPluginRequest) (response *UninstallKibanaPluginResponse, err error) {
	response = CreateUninstallKibanaPluginResponse()
	err = client.DoAction(request, response)
	return
}

// UninstallKibanaPluginWithChan invokes the elasticsearch.UninstallKibanaPlugin API asynchronously
func (client *Client) UninstallKibanaPluginWithChan(request *UninstallKibanaPluginRequest) (<-chan *UninstallKibanaPluginResponse, <-chan error) {
	responseChan := make(chan *UninstallKibanaPluginResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UninstallKibanaPlugin(request)
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

// UninstallKibanaPluginWithCallback invokes the elasticsearch.UninstallKibanaPlugin API asynchronously
func (client *Client) UninstallKibanaPluginWithCallback(request *UninstallKibanaPluginRequest, callback func(response *UninstallKibanaPluginResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UninstallKibanaPluginResponse
		var err error
		defer close(result)
		response, err = client.UninstallKibanaPlugin(request)
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

// UninstallKibanaPluginRequest is the request struct for api UninstallKibanaPlugin
type UninstallKibanaPluginRequest struct {
	*requests.RoaRequest
	InstanceId  string `position:"Path" name:"InstanceId"`
	ClientToken string `position:"Query" name:"clientToken"`
}

// UninstallKibanaPluginResponse is the response struct for api UninstallKibanaPlugin
type UninstallKibanaPluginResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Result    []string `json:"Result" xml:"Result"`
}

// CreateUninstallKibanaPluginRequest creates a request to invoke UninstallKibanaPlugin API
func CreateUninstallKibanaPluginRequest() (request *UninstallKibanaPluginRequest) {
	request = &UninstallKibanaPluginRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "UninstallKibanaPlugin", "/openapi/instances/[InstanceId]/kibana-plugins/actions/uninstall", "elasticsearch", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUninstallKibanaPluginResponse creates a response to parse from UninstallKibanaPlugin response
func CreateUninstallKibanaPluginResponse() (response *UninstallKibanaPluginResponse) {
	response = &UninstallKibanaPluginResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
