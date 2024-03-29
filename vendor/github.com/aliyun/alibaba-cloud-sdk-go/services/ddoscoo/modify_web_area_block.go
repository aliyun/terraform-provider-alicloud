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

// ModifyWebAreaBlock invokes the ddoscoo.ModifyWebAreaBlock API synchronously
func (client *Client) ModifyWebAreaBlock(request *ModifyWebAreaBlockRequest) (response *ModifyWebAreaBlockResponse, err error) {
	response = CreateModifyWebAreaBlockResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyWebAreaBlockWithChan invokes the ddoscoo.ModifyWebAreaBlock API asynchronously
func (client *Client) ModifyWebAreaBlockWithChan(request *ModifyWebAreaBlockRequest) (<-chan *ModifyWebAreaBlockResponse, <-chan error) {
	responseChan := make(chan *ModifyWebAreaBlockResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyWebAreaBlock(request)
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

// ModifyWebAreaBlockWithCallback invokes the ddoscoo.ModifyWebAreaBlock API asynchronously
func (client *Client) ModifyWebAreaBlockWithCallback(request *ModifyWebAreaBlockRequest, callback func(response *ModifyWebAreaBlockResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyWebAreaBlockResponse
		var err error
		defer close(result)
		response, err = client.ModifyWebAreaBlock(request)
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

// ModifyWebAreaBlockRequest is the request struct for api ModifyWebAreaBlock
type ModifyWebAreaBlockRequest struct {
	*requests.RpcRequest
	Regions         *[]string `position:"Query" name:"Regions"  type:"Repeated"`
	ResourceGroupId string    `position:"Query" name:"ResourceGroupId"`
	SourceIp        string    `position:"Query" name:"SourceIp"`
	Domain          string    `position:"Query" name:"Domain"`
}

// ModifyWebAreaBlockResponse is the response struct for api ModifyWebAreaBlock
type ModifyWebAreaBlockResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyWebAreaBlockRequest creates a request to invoke ModifyWebAreaBlock API
func CreateModifyWebAreaBlockRequest() (request *ModifyWebAreaBlockRequest) {
	request = &ModifyWebAreaBlockRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "ModifyWebAreaBlock", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyWebAreaBlockResponse creates a response to parse from ModifyWebAreaBlock response
func CreateModifyWebAreaBlockResponse() (response *ModifyWebAreaBlockResponse) {
	response = &ModifyWebAreaBlockResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
