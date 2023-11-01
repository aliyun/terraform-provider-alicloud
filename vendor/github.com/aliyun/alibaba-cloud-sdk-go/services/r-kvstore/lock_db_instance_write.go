package r_kvstore

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

// LockDBInstanceWrite invokes the r_kvstore.LockDBInstanceWrite API synchronously
func (client *Client) LockDBInstanceWrite(request *LockDBInstanceWriteRequest) (response *LockDBInstanceWriteResponse, err error) {
	response = CreateLockDBInstanceWriteResponse()
	err = client.DoAction(request, response)
	return
}

// LockDBInstanceWriteWithChan invokes the r_kvstore.LockDBInstanceWrite API asynchronously
func (client *Client) LockDBInstanceWriteWithChan(request *LockDBInstanceWriteRequest) (<-chan *LockDBInstanceWriteResponse, <-chan error) {
	responseChan := make(chan *LockDBInstanceWriteResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.LockDBInstanceWrite(request)
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

// LockDBInstanceWriteWithCallback invokes the r_kvstore.LockDBInstanceWrite API asynchronously
func (client *Client) LockDBInstanceWriteWithCallback(request *LockDBInstanceWriteRequest, callback func(response *LockDBInstanceWriteResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *LockDBInstanceWriteResponse
		var err error
		defer close(result)
		response, err = client.LockDBInstanceWrite(request)
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

// LockDBInstanceWriteRequest is the request struct for api LockDBInstanceWrite
type LockDBInstanceWriteRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	LockReason           string           `position:"Query" name:"LockReason"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// LockDBInstanceWriteResponse is the response struct for api LockDBInstanceWrite
type LockDBInstanceWriteResponse struct {
	*responses.BaseResponse
	RequestId      string `json:"RequestId" xml:"RequestId"`
	DBInstanceName string `json:"DBInstanceName" xml:"DBInstanceName"`
	TaskId         int64  `json:"TaskId" xml:"TaskId"`
	LockReason     string `json:"LockReason" xml:"LockReason"`
}

// CreateLockDBInstanceWriteRequest creates a request to invoke LockDBInstanceWrite API
func CreateLockDBInstanceWriteRequest() (request *LockDBInstanceWriteRequest) {
	request = &LockDBInstanceWriteRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "LockDBInstanceWrite", "redisa", "openAPI")
	request.Method = requests.POST
	return
}

// CreateLockDBInstanceWriteResponse creates a response to parse from LockDBInstanceWrite response
func CreateLockDBInstanceWriteResponse() (response *LockDBInstanceWriteResponse) {
	response = &LockDBInstanceWriteResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
