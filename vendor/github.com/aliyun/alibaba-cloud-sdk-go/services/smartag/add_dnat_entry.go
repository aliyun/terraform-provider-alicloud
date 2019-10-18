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

// AddDnatEntry invokes the smartag.AddDnatEntry API synchronously
// api document: https://help.aliyun.com/api/smartag/adddnatentry.html
func (client *Client) AddDnatEntry(request *AddDnatEntryRequest) (response *AddDnatEntryResponse, err error) {
	response = CreateAddDnatEntryResponse()
	err = client.DoAction(request, response)
	return
}

// AddDnatEntryWithChan invokes the smartag.AddDnatEntry API asynchronously
// api document: https://help.aliyun.com/api/smartag/adddnatentry.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddDnatEntryWithChan(request *AddDnatEntryRequest) (<-chan *AddDnatEntryResponse, <-chan error) {
	responseChan := make(chan *AddDnatEntryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddDnatEntry(request)
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

// AddDnatEntryWithCallback invokes the smartag.AddDnatEntry API asynchronously
// api document: https://help.aliyun.com/api/smartag/adddnatentry.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddDnatEntryWithCallback(request *AddDnatEntryRequest, callback func(response *AddDnatEntryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddDnatEntryResponse
		var err error
		defer close(result)
		response, err = client.AddDnatEntry(request)
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

// AddDnatEntryRequest is the request struct for api AddDnatEntry
type AddDnatEntryRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	SagId		     string           `position:"Query" name:"SagId"`
	Type                 string           `position:"Query" name:"Type"`
	ExternalIp           string           `position:"Query" name:"ExternalIp"`
	ExternalPort         string           `position:"Query" name:"ExternalPort"`
	InternalIp           string           `position:"Query" name:"InternalIp"`
	InternalPort         string           `position:"Query" name:"InternalPort"`
	IpProtocol           string           `position:"Query" name:"IpProtocol"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// AddDnatEntryResponse is the response struct for api AddDnatEntry
type AddDnatEntryResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	DnatEntryId   string `json:"DnatEntryId" xml:"DnatEntryId"`
}

// CreateAddDnatEntryRequest creates a request to invoke AddDnatEntry API
func CreateAddDnatEntryRequest() (request *AddDnatEntryRequest) {
	request = &AddDnatEntryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "AddDnatEntry", "smartag", "openAPI")
	return
}

// CreateAddDnatEntryResponse creates a response to parse from AddDnatEntry response
func CreateAddDnatEntryResponse() (response *AddDnatEntryResponse) {
	response = &AddDnatEntryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
