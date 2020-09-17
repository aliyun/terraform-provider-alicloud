package smartag

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

// UnlockSmartAccessGateway invokes the smartag.UnlockSmartAccessGateway API synchronously
// api document: https://help.aliyun.com/api/smartag/unlocksmartaccessgateway.html
func (client *Client) UnlockSmartAccessGateway(request *UnlockSmartAccessGatewayRequest) (response *UnlockSmartAccessGatewayResponse, err error) {
	response = CreateUnlockSmartAccessGatewayResponse()
	err = client.DoAction(request, response)
	return
}

// UnlockSmartAccessGatewayWithChan invokes the smartag.UnlockSmartAccessGateway API asynchronously
// api document: https://help.aliyun.com/api/smartag/unlocksmartaccessgateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnlockSmartAccessGatewayWithChan(request *UnlockSmartAccessGatewayRequest) (<-chan *UnlockSmartAccessGatewayResponse, <-chan error) {
	responseChan := make(chan *UnlockSmartAccessGatewayResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UnlockSmartAccessGateway(request)
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

// UnlockSmartAccessGatewayWithCallback invokes the smartag.UnlockSmartAccessGateway API asynchronously
// api document: https://help.aliyun.com/api/smartag/unlocksmartaccessgateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnlockSmartAccessGatewayWithCallback(request *UnlockSmartAccessGatewayRequest, callback func(response *UnlockSmartAccessGatewayResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UnlockSmartAccessGatewayResponse
		var err error
		defer close(result)
		response, err = client.UnlockSmartAccessGateway(request)
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

// UnlockSmartAccessGatewayRequest is the request struct for api UnlockSmartAccessGateway
type UnlockSmartAccessGatewayRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	SmartAGId            string           `position:"Query" name:"SmartAGId"`
}

// UnlockSmartAccessGatewayResponse is the response struct for api UnlockSmartAccessGateway
type UnlockSmartAccessGatewayResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUnlockSmartAccessGatewayRequest creates a request to invoke UnlockSmartAccessGateway API
func CreateUnlockSmartAccessGatewayRequest() (request *UnlockSmartAccessGatewayRequest) {
	request = &UnlockSmartAccessGatewayRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "UnlockSmartAccessGateway", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUnlockSmartAccessGatewayResponse creates a response to parse from UnlockSmartAccessGateway response
func CreateUnlockSmartAccessGatewayResponse() (response *UnlockSmartAccessGatewayResponse) {
	response = &UnlockSmartAccessGatewayResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
