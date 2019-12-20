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

// ModifySQLCollectorPolicy invokes the rds.ModifySQLCollectorPolicy API synchronously
// api document: https://help.aliyun.com/api/rds/modifysqlcollectorpolicy.html
func (client *Client) ModifySQLCollectorPolicy(request *ModifySQLCollectorPolicyRequest) (response *ModifySQLCollectorPolicyResponse, err error) {
	response = CreateModifySQLCollectorPolicyResponse()
	err = client.DoAction(request, response)
	return
}

// ModifySQLCollectorPolicyWithChan invokes the rds.ModifySQLCollectorPolicy API asynchronously
// api document: https://help.aliyun.com/api/rds/modifysqlcollectorpolicy.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySQLCollectorPolicyWithChan(request *ModifySQLCollectorPolicyRequest) (<-chan *ModifySQLCollectorPolicyResponse, <-chan error) {
	responseChan := make(chan *ModifySQLCollectorPolicyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifySQLCollectorPolicy(request)
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

// ModifySQLCollectorPolicyWithCallback invokes the rds.ModifySQLCollectorPolicy API asynchronously
// api document: https://help.aliyun.com/api/rds/modifysqlcollectorpolicy.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySQLCollectorPolicyWithCallback(request *ModifySQLCollectorPolicyRequest, callback func(response *ModifySQLCollectorPolicyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifySQLCollectorPolicyResponse
		var err error
		defer close(result)
		response, err = client.ModifySQLCollectorPolicy(request)
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

// ModifySQLCollectorPolicyRequest is the request struct for api ModifySQLCollectorPolicy
type ModifySQLCollectorPolicyRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	StoragePeriod        requests.Integer `position:"Query" name:"StoragePeriod"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	SQLCollectorStatus   string           `position:"Query" name:"SQLCollectorStatus"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// ModifySQLCollectorPolicyResponse is the response struct for api ModifySQLCollectorPolicy
type ModifySQLCollectorPolicyResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifySQLCollectorPolicyRequest creates a request to invoke ModifySQLCollectorPolicy API
func CreateModifySQLCollectorPolicyRequest() (request *ModifySQLCollectorPolicyRequest) {
	request = &ModifySQLCollectorPolicyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "ModifySQLCollectorPolicy", "rds", "openAPI")
	return
}

// CreateModifySQLCollectorPolicyResponse creates a response to parse from ModifySQLCollectorPolicy response
func CreateModifySQLCollectorPolicyResponse() (response *ModifySQLCollectorPolicyResponse) {
	response = &ModifySQLCollectorPolicyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
