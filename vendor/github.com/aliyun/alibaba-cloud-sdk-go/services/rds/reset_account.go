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

// ResetAccount invokes the rds.ResetAccount API synchronously
func (client *Client) ResetAccount(request *ResetAccountRequest) (response *ResetAccountResponse, err error) {
	response = CreateResetAccountResponse()
	err = client.DoAction(request, response)
	return
}

// ResetAccountWithChan invokes the rds.ResetAccount API asynchronously
func (client *Client) ResetAccountWithChan(request *ResetAccountRequest) (<-chan *ResetAccountResponse, <-chan error) {
	responseChan := make(chan *ResetAccountResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ResetAccount(request)
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

// ResetAccountWithCallback invokes the rds.ResetAccount API asynchronously
func (client *Client) ResetAccountWithCallback(request *ResetAccountRequest, callback func(response *ResetAccountResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ResetAccountResponse
		var err error
		defer close(result)
		response, err = client.ResetAccount(request)
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

// ResetAccountRequest is the request struct for api ResetAccount
type ResetAccountRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	AccountName          string           `position:"Query" name:"AccountName"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	AccountPassword      string           `position:"Query" name:"AccountPassword"`
}

// ResetAccountResponse is the response struct for api ResetAccount
type ResetAccountResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateResetAccountRequest creates a request to invoke ResetAccount API
func CreateResetAccountRequest() (request *ResetAccountRequest) {
	request = &ResetAccountRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "ResetAccount", "rds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateResetAccountResponse creates a response to parse from ResetAccount response
func CreateResetAccountResponse() (response *ResetAccountResponse) {
	response = &ResetAccountResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
