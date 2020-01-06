package rds

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

// UnlockAccount invokes the rds.UnlockAccount API synchronously
// api document: https://help.aliyun.com/api/rds/unlockaccount.html
func (client *Client) UnlockAccount(request *UnlockAccountRequest) (response *UnlockAccountResponse, err error) {
	response = CreateUnlockAccountResponse()
	err = client.DoAction(request, response)
	return
}

// UnlockAccountWithChan invokes the rds.UnlockAccount API asynchronously
// api document: https://help.aliyun.com/api/rds/unlockaccount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnlockAccountWithChan(request *UnlockAccountRequest) (<-chan *UnlockAccountResponse, <-chan error) {
	responseChan := make(chan *UnlockAccountResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UnlockAccount(request)
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

// UnlockAccountWithCallback invokes the rds.UnlockAccount API asynchronously
// api document: https://help.aliyun.com/api/rds/unlockaccount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UnlockAccountWithCallback(request *UnlockAccountRequest, callback func(response *UnlockAccountResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UnlockAccountResponse
		var err error
		defer close(result)
		response, err = client.UnlockAccount(request)
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

// UnlockAccountRequest is the request struct for api UnlockAccount
type UnlockAccountRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	AccountName          string           `position:"Query" name:"AccountName"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
}

// UnlockAccountResponse is the response struct for api UnlockAccount
type UnlockAccountResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUnlockAccountRequest creates a request to invoke UnlockAccount API
func CreateUnlockAccountRequest() (request *UnlockAccountRequest) {
	request = &UnlockAccountRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "UnlockAccount", "Rds", "openAPI")
	return
}

// CreateUnlockAccountResponse creates a response to parse from UnlockAccount response
func CreateUnlockAccountResponse() (response *UnlockAccountResponse) {
	response = &UnlockAccountResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
