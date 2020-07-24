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

// ClearSagRouteableAddress invokes the smartag.ClearSagRouteableAddress API synchronously
// api document: https://help.aliyun.com/api/smartag/clearsagrouteableaddress.html
func (client *Client) ClearSagRouteableAddress(request *ClearSagRouteableAddressRequest) (response *ClearSagRouteableAddressResponse, err error) {
	response = CreateClearSagRouteableAddressResponse()
	err = client.DoAction(request, response)
	return
}

// ClearSagRouteableAddressWithChan invokes the smartag.ClearSagRouteableAddress API asynchronously
// api document: https://help.aliyun.com/api/smartag/clearsagrouteableaddress.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearSagRouteableAddressWithChan(request *ClearSagRouteableAddressRequest) (<-chan *ClearSagRouteableAddressResponse, <-chan error) {
	responseChan := make(chan *ClearSagRouteableAddressResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ClearSagRouteableAddress(request)
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

// ClearSagRouteableAddressWithCallback invokes the smartag.ClearSagRouteableAddress API asynchronously
// api document: https://help.aliyun.com/api/smartag/clearsagrouteableaddress.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearSagRouteableAddressWithCallback(request *ClearSagRouteableAddressRequest, callback func(response *ClearSagRouteableAddressResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ClearSagRouteableAddressResponse
		var err error
		defer close(result)
		response, err = client.ClearSagRouteableAddress(request)
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

// ClearSagRouteableAddressRequest is the request struct for api ClearSagRouteableAddress
type ClearSagRouteableAddressRequest struct {
	*requests.RpcRequest
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	SagId                string           `position:"Query" name:"SagId"`
}

// ClearSagRouteableAddressResponse is the response struct for api ClearSagRouteableAddress
type ClearSagRouteableAddressResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateClearSagRouteableAddressRequest creates a request to invoke ClearSagRouteableAddress API
func CreateClearSagRouteableAddressRequest() (request *ClearSagRouteableAddressRequest) {
	request = &ClearSagRouteableAddressRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "ClearSagRouteableAddress", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateClearSagRouteableAddressResponse creates a response to parse from ClearSagRouteableAddress response
func CreateClearSagRouteableAddressResponse() (response *ClearSagRouteableAddressResponse) {
	response = &ClearSagRouteableAddressResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
