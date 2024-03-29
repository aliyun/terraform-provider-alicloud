package ddoscoo

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

// EnableWebCC invokes the ddoscoo.EnableWebCC API synchronously
func (client *Client) EnableWebCC(request *EnableWebCCRequest) (response *EnableWebCCResponse, err error) {
	response = CreateEnableWebCCResponse()
	err = client.DoAction(request, response)
	return
}

// EnableWebCCWithChan invokes the ddoscoo.EnableWebCC API asynchronously
func (client *Client) EnableWebCCWithChan(request *EnableWebCCRequest) (<-chan *EnableWebCCResponse, <-chan error) {
	responseChan := make(chan *EnableWebCCResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.EnableWebCC(request)
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

// EnableWebCCWithCallback invokes the ddoscoo.EnableWebCC API asynchronously
func (client *Client) EnableWebCCWithCallback(request *EnableWebCCRequest, callback func(response *EnableWebCCResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *EnableWebCCResponse
		var err error
		defer close(result)
		response, err = client.EnableWebCC(request)
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

// EnableWebCCRequest is the request struct for api EnableWebCC
type EnableWebCCRequest struct {
	*requests.RpcRequest
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
	Domain          string `position:"Query" name:"Domain"`
}

// EnableWebCCResponse is the response struct for api EnableWebCC
type EnableWebCCResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateEnableWebCCRequest creates a request to invoke EnableWebCC API
func CreateEnableWebCCRequest() (request *EnableWebCCRequest) {
	request = &EnableWebCCRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "EnableWebCC", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateEnableWebCCResponse creates a response to parse from EnableWebCC response
func CreateEnableWebCCResponse() (response *EnableWebCCResponse) {
	response = &EnableWebCCResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
