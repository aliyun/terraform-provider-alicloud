package dds

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

// CheckCloudResourceAuthorized invokes the dds.CheckCloudResourceAuthorized API synchronously
// api document: https://help.aliyun.com/api/dds/checkcloudresourceauthorized.html
func (client *Client) CheckCloudResourceAuthorized(request *CheckCloudResourceAuthorizedRequest) (response *CheckCloudResourceAuthorizedResponse, err error) {
	response = CreateCheckCloudResourceAuthorizedResponse()
	err = client.DoAction(request, response)
	return
}

// CheckCloudResourceAuthorizedWithChan invokes the dds.CheckCloudResourceAuthorized API asynchronously
// api document: https://help.aliyun.com/api/dds/checkcloudresourceauthorized.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CheckCloudResourceAuthorizedWithChan(request *CheckCloudResourceAuthorizedRequest) (<-chan *CheckCloudResourceAuthorizedResponse, <-chan error) {
	responseChan := make(chan *CheckCloudResourceAuthorizedResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CheckCloudResourceAuthorized(request)
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

// CheckCloudResourceAuthorizedWithCallback invokes the dds.CheckCloudResourceAuthorized API asynchronously
// api document: https://help.aliyun.com/api/dds/checkcloudresourceauthorized.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CheckCloudResourceAuthorizedWithCallback(request *CheckCloudResourceAuthorizedRequest, callback func(response *CheckCloudResourceAuthorizedResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CheckCloudResourceAuthorizedResponse
		var err error
		defer close(result)
		response, err = client.CheckCloudResourceAuthorized(request)
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

// CheckCloudResourceAuthorizedRequest is the request struct for api CheckCloudResourceAuthorized
type CheckCloudResourceAuthorizedRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	TargetRegionId       string           `position:"Query" name:"TargetRegionId"`
}

// CheckCloudResourceAuthorizedResponse is the response struct for api CheckCloudResourceAuthorized
type CheckCloudResourceAuthorizedResponse struct {
	*responses.BaseResponse
	RequestId          string `json:"RequestId" xml:"RequestId"`
	AuthorizationState int    `json:"AuthorizationState" xml:"AuthorizationState"`
	RoleArn            string `json:"RoleArn" xml:"RoleArn"`
}

// CreateCheckCloudResourceAuthorizedRequest creates a request to invoke CheckCloudResourceAuthorized API
func CreateCheckCloudResourceAuthorizedRequest() (request *CheckCloudResourceAuthorizedRequest) {
	request = &CheckCloudResourceAuthorizedRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "CheckCloudResourceAuthorized", "Dds", "openAPI")
	return
}

// CreateCheckCloudResourceAuthorizedResponse creates a response to parse from CheckCloudResourceAuthorized response
func CreateCheckCloudResourceAuthorizedResponse() (response *CheckCloudResourceAuthorizedResponse) {
	response = &CheckCloudResourceAuthorizedResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
