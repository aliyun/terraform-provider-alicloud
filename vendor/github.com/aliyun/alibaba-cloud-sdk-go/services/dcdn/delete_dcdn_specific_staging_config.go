package dcdn

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

// DeleteDcdnSpecificStagingConfig invokes the dcdn.DeleteDcdnSpecificStagingConfig API synchronously
// api document: https://help.aliyun.com/api/dcdn/deletedcdnspecificstagingconfig.html
func (client *Client) DeleteDcdnSpecificStagingConfig(request *DeleteDcdnSpecificStagingConfigRequest) (response *DeleteDcdnSpecificStagingConfigResponse, err error) {
	response = CreateDeleteDcdnSpecificStagingConfigResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteDcdnSpecificStagingConfigWithChan invokes the dcdn.DeleteDcdnSpecificStagingConfig API asynchronously
// api document: https://help.aliyun.com/api/dcdn/deletedcdnspecificstagingconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDcdnSpecificStagingConfigWithChan(request *DeleteDcdnSpecificStagingConfigRequest) (<-chan *DeleteDcdnSpecificStagingConfigResponse, <-chan error) {
	responseChan := make(chan *DeleteDcdnSpecificStagingConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteDcdnSpecificStagingConfig(request)
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

// DeleteDcdnSpecificStagingConfigWithCallback invokes the dcdn.DeleteDcdnSpecificStagingConfig API asynchronously
// api document: https://help.aliyun.com/api/dcdn/deletedcdnspecificstagingconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDcdnSpecificStagingConfigWithCallback(request *DeleteDcdnSpecificStagingConfigRequest, callback func(response *DeleteDcdnSpecificStagingConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteDcdnSpecificStagingConfigResponse
		var err error
		defer close(result)
		response, err = client.DeleteDcdnSpecificStagingConfig(request)
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

// DeleteDcdnSpecificStagingConfigRequest is the request struct for api DeleteDcdnSpecificStagingConfig
type DeleteDcdnSpecificStagingConfigRequest struct {
	*requests.RpcRequest
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	ConfigId      string           `position:"Query" name:"ConfigId"`
}

// DeleteDcdnSpecificStagingConfigResponse is the response struct for api DeleteDcdnSpecificStagingConfig
type DeleteDcdnSpecificStagingConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteDcdnSpecificStagingConfigRequest creates a request to invoke DeleteDcdnSpecificStagingConfig API
func CreateDeleteDcdnSpecificStagingConfigRequest() (request *DeleteDcdnSpecificStagingConfigRequest) {
	request = &DeleteDcdnSpecificStagingConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DeleteDcdnSpecificStagingConfig", "", "")
	request.Method = requests.POST
	return
}

// CreateDeleteDcdnSpecificStagingConfigResponse creates a response to parse from DeleteDcdnSpecificStagingConfig response
func CreateDeleteDcdnSpecificStagingConfigResponse() (response *DeleteDcdnSpecificStagingConfigResponse) {
	response = &DeleteDcdnSpecificStagingConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
