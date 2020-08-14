package cdn

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

// SetVideoSeekConfig invokes the cdn.SetVideoSeekConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/setvideoseekconfig.html
func (client *Client) SetVideoSeekConfig(request *SetVideoSeekConfigRequest) (response *SetVideoSeekConfigResponse, err error) {
	response = CreateSetVideoSeekConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetVideoSeekConfigWithChan invokes the cdn.SetVideoSeekConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setvideoseekconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetVideoSeekConfigWithChan(request *SetVideoSeekConfigRequest) (<-chan *SetVideoSeekConfigResponse, <-chan error) {
	responseChan := make(chan *SetVideoSeekConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetVideoSeekConfig(request)
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

// SetVideoSeekConfigWithCallback invokes the cdn.SetVideoSeekConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setvideoseekconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetVideoSeekConfigWithCallback(request *SetVideoSeekConfigRequest, callback func(response *SetVideoSeekConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetVideoSeekConfigResponse
		var err error
		defer close(result)
		response, err = client.SetVideoSeekConfig(request)
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

// SetVideoSeekConfigRequest is the request struct for api SetVideoSeekConfig
type SetVideoSeekConfigRequest struct {
	*requests.RpcRequest
	Enable     string           `position:"Query" name:"Enable"`
	DomainName string           `position:"Query" name:"DomainName"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
	ConfigId   requests.Integer `position:"Query" name:"ConfigId"`
}

// SetVideoSeekConfigResponse is the response struct for api SetVideoSeekConfig
type SetVideoSeekConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetVideoSeekConfigRequest creates a request to invoke SetVideoSeekConfig API
func CreateSetVideoSeekConfigRequest() (request *SetVideoSeekConfigRequest) {
	request = &SetVideoSeekConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "SetVideoSeekConfig", "", "")
	request.Method = requests.POST
	return
}

// CreateSetVideoSeekConfigResponse creates a response to parse from SetVideoSeekConfig response
func CreateSetVideoSeekConfigResponse() (response *SetVideoSeekConfigResponse) {
	response = &SetVideoSeekConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
