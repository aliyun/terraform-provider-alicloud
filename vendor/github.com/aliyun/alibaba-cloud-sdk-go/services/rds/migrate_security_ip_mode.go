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

// MigrateSecurityIPMode invokes the rds.MigrateSecurityIPMode API synchronously
// api document: https://help.aliyun.com/api/rds/migratesecurityipmode.html
func (client *Client) MigrateSecurityIPMode(request *MigrateSecurityIPModeRequest) (response *MigrateSecurityIPModeResponse, err error) {
	response = CreateMigrateSecurityIPModeResponse()
	err = client.DoAction(request, response)
	return
}

// MigrateSecurityIPModeWithChan invokes the rds.MigrateSecurityIPMode API asynchronously
// api document: https://help.aliyun.com/api/rds/migratesecurityipmode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MigrateSecurityIPModeWithChan(request *MigrateSecurityIPModeRequest) (<-chan *MigrateSecurityIPModeResponse, <-chan error) {
	responseChan := make(chan *MigrateSecurityIPModeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.MigrateSecurityIPMode(request)
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

// MigrateSecurityIPModeWithCallback invokes the rds.MigrateSecurityIPMode API asynchronously
// api document: https://help.aliyun.com/api/rds/migratesecurityipmode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MigrateSecurityIPModeWithCallback(request *MigrateSecurityIPModeRequest, callback func(response *MigrateSecurityIPModeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *MigrateSecurityIPModeResponse
		var err error
		defer close(result)
		response, err = client.MigrateSecurityIPMode(request)
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

// MigrateSecurityIPModeRequest is the request struct for api MigrateSecurityIPMode
type MigrateSecurityIPModeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
}

// MigrateSecurityIPModeResponse is the response struct for api MigrateSecurityIPMode
type MigrateSecurityIPModeResponse struct {
	*responses.BaseResponse
	RequestId      string `json:"RequestId" xml:"RequestId"`
	DBInstanceId   string `json:"DBInstanceId" xml:"DBInstanceId"`
	SecurityIPMode string `json:"SecurityIPMode" xml:"SecurityIPMode"`
}

// CreateMigrateSecurityIPModeRequest creates a request to invoke MigrateSecurityIPMode API
func CreateMigrateSecurityIPModeRequest() (request *MigrateSecurityIPModeRequest) {
	request = &MigrateSecurityIPModeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "MigrateSecurityIPMode", "Rds", "openAPI")
	return
}

// CreateMigrateSecurityIPModeResponse creates a response to parse from MigrateSecurityIPMode response
func CreateMigrateSecurityIPModeResponse() (response *MigrateSecurityIPModeResponse) {
	response = &MigrateSecurityIPModeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
