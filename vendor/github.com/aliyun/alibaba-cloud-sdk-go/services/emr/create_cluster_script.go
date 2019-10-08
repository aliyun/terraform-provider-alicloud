package emr

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

// CreateClusterScript invokes the emr.CreateClusterScript API synchronously
// api document: https://help.aliyun.com/api/emr/createclusterscript.html
func (client *Client) CreateClusterScript(request *CreateClusterScriptRequest) (response *CreateClusterScriptResponse, err error) {
	response = CreateCreateClusterScriptResponse()
	err = client.DoAction(request, response)
	return
}

// CreateClusterScriptWithChan invokes the emr.CreateClusterScript API asynchronously
// api document: https://help.aliyun.com/api/emr/createclusterscript.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateClusterScriptWithChan(request *CreateClusterScriptRequest) (<-chan *CreateClusterScriptResponse, <-chan error) {
	responseChan := make(chan *CreateClusterScriptResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateClusterScript(request)
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

// CreateClusterScriptWithCallback invokes the emr.CreateClusterScript API asynchronously
// api document: https://help.aliyun.com/api/emr/createclusterscript.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateClusterScriptWithCallback(request *CreateClusterScriptRequest, callback func(response *CreateClusterScriptResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateClusterScriptResponse
		var err error
		defer close(result)
		response, err = client.CreateClusterScript(request)
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

// CreateClusterScriptRequest is the request struct for api CreateClusterScript
type CreateClusterScriptRequest struct {
	*requests.RpcRequest
	Args            string           `position:"Query" name:"Args"`
	Path            string           `position:"Query" name:"Path"`
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Name            string           `position:"Query" name:"Name"`
	ClusterId       string           `position:"Query" name:"ClusterId"`
	NodeIdList      string           `position:"Query" name:"NodeIdList"`
}

// CreateClusterScriptResponse is the response struct for api CreateClusterScript
type CreateClusterScriptResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Id        string `json:"Id" xml:"Id"`
}

// CreateCreateClusterScriptRequest creates a request to invoke CreateClusterScript API
func CreateCreateClusterScriptRequest() (request *CreateClusterScriptRequest) {
	request = &CreateClusterScriptRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "CreateClusterScript", "emr", "openAPI")
	return
}

// CreateCreateClusterScriptResponse creates a response to parse from CreateClusterScript response
func CreateCreateClusterScriptResponse() (response *CreateClusterScriptResponse) {
	response = &CreateClusterScriptResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
