package ddosbgp

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

// ReleaseInstance invokes the ddosbgp.ReleaseInstance API synchronously
// api document: https://help.aliyun.com/api/ddosbgp/releaseinstance.html
func (client *Client) ReleaseInstance(request *ReleaseInstanceRequest) (response *ReleaseInstanceResponse, err error) {
	response = CreateReleaseInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// ReleaseInstanceWithChan invokes the ddosbgp.ReleaseInstance API asynchronously
// api document: https://help.aliyun.com/api/ddosbgp/releaseinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ReleaseInstanceWithChan(request *ReleaseInstanceRequest) (<-chan *ReleaseInstanceResponse, <-chan error) {
	responseChan := make(chan *ReleaseInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ReleaseInstance(request)
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

// ReleaseInstanceWithCallback invokes the ddosbgp.ReleaseInstance API asynchronously
// api document: https://help.aliyun.com/api/ddosbgp/releaseinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ReleaseInstanceWithCallback(request *ReleaseInstanceRequest, callback func(response *ReleaseInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ReleaseInstanceResponse
		var err error
		defer close(result)
		response, err = client.ReleaseInstance(request)
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

// ReleaseInstanceRequest is the request struct for api ReleaseInstance
type ReleaseInstanceRequest struct {
	*requests.RpcRequest
	SourceIp   string `position:"Query" name:"SourceIp"`
	InstanceId string `position:"Query" name:"InstanceId"`
}

// ReleaseInstanceResponse is the response struct for api ReleaseInstance
type ReleaseInstanceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Code      string `json:"Code" xml:"Code"`
}

// CreateReleaseInstanceRequest creates a request to invoke ReleaseInstance API
func CreateReleaseInstanceRequest() (request *ReleaseInstanceRequest) {
	request = &ReleaseInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddosbgp", "2018-02-01", "ReleaseInstance", "ddosbgp", "openAPI")
	return
}

// CreateReleaseInstanceResponse creates a response to parse from ReleaseInstance response
func CreateReleaseInstanceResponse() (response *ReleaseInstanceResponse) {
	response = &ReleaseInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
