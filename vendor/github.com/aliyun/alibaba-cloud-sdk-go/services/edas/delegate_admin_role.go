package edas

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

// DelegateAdminRole invokes the edas.DelegateAdminRole API synchronously
// api document: https://help.aliyun.com/api/edas/delegateadminrole.html
func (client *Client) DelegateAdminRole(request *DelegateAdminRoleRequest) (response *DelegateAdminRoleResponse, err error) {
	response = CreateDelegateAdminRoleResponse()
	err = client.DoAction(request, response)
	return
}

// DelegateAdminRoleWithChan invokes the edas.DelegateAdminRole API asynchronously
// api document: https://help.aliyun.com/api/edas/delegateadminrole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DelegateAdminRoleWithChan(request *DelegateAdminRoleRequest) (<-chan *DelegateAdminRoleResponse, <-chan error) {
	responseChan := make(chan *DelegateAdminRoleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DelegateAdminRole(request)
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

// DelegateAdminRoleWithCallback invokes the edas.DelegateAdminRole API asynchronously
// api document: https://help.aliyun.com/api/edas/delegateadminrole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DelegateAdminRoleWithCallback(request *DelegateAdminRoleRequest, callback func(response *DelegateAdminRoleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DelegateAdminRoleResponse
		var err error
		defer close(result)
		response, err = client.DelegateAdminRole(request)
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

// DelegateAdminRoleRequest is the request struct for api DelegateAdminRole
type DelegateAdminRoleRequest struct {
	*requests.RoaRequest
	TargetUserId string `position:"Query" name:"TargetUserId"`
}

// DelegateAdminRoleResponse is the response struct for api DelegateAdminRole
type DelegateAdminRoleResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDelegateAdminRoleRequest creates a request to invoke DelegateAdminRole API
func CreateDelegateAdminRoleRequest() (request *DelegateAdminRoleRequest) {
	request = &DelegateAdminRoleRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "DelegateAdminRole", "/pop/v5/account/delegate_admin_role", "Edas", "openAPI")
	request.Method = requests.PUT
	return
}

// CreateDelegateAdminRoleResponse creates a response to parse from DelegateAdminRole response
func CreateDelegateAdminRoleResponse() (response *DelegateAdminRoleResponse) {
	response = &DelegateAdminRoleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
