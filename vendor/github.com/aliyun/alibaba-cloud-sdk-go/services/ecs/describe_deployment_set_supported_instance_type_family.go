package ecs

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

// DescribeDeploymentSetSupportedInstanceTypeFamily invokes the ecs.DescribeDeploymentSetSupportedInstanceTypeFamily API synchronously
// api document: https://help.aliyun.com/api/ecs/describedeploymentsetsupportedinstancetypefamily.html
func (client *Client) DescribeDeploymentSetSupportedInstanceTypeFamily(request *DescribeDeploymentSetSupportedInstanceTypeFamilyRequest) (response *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse, err error) {
	response = CreateDescribeDeploymentSetSupportedInstanceTypeFamilyResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDeploymentSetSupportedInstanceTypeFamilyWithChan invokes the ecs.DescribeDeploymentSetSupportedInstanceTypeFamily API asynchronously
// api document: https://help.aliyun.com/api/ecs/describedeploymentsetsupportedinstancetypefamily.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDeploymentSetSupportedInstanceTypeFamilyWithChan(request *DescribeDeploymentSetSupportedInstanceTypeFamilyRequest) (<-chan *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse, <-chan error) {
	responseChan := make(chan *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDeploymentSetSupportedInstanceTypeFamily(request)
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

// DescribeDeploymentSetSupportedInstanceTypeFamilyWithCallback invokes the ecs.DescribeDeploymentSetSupportedInstanceTypeFamily API asynchronously
// api document: https://help.aliyun.com/api/ecs/describedeploymentsetsupportedinstancetypefamily.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDeploymentSetSupportedInstanceTypeFamilyWithCallback(request *DescribeDeploymentSetSupportedInstanceTypeFamilyRequest, callback func(response *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse
		var err error
		defer close(result)
		response, err = client.DescribeDeploymentSetSupportedInstanceTypeFamily(request)
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

// DescribeDeploymentSetSupportedInstanceTypeFamilyRequest is the request struct for api DescribeDeploymentSetSupportedInstanceTypeFamily
type DescribeDeploymentSetSupportedInstanceTypeFamilyRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDeploymentSetSupportedInstanceTypeFamilyResponse is the response struct for api DescribeDeploymentSetSupportedInstanceTypeFamily
type DescribeDeploymentSetSupportedInstanceTypeFamilyResponse struct {
	*responses.BaseResponse
	RequestId            string `json:"RequestId" xml:"RequestId"`
	InstanceTypeFamilies string `json:"InstanceTypeFamilies" xml:"InstanceTypeFamilies"`
}

// CreateDescribeDeploymentSetSupportedInstanceTypeFamilyRequest creates a request to invoke DescribeDeploymentSetSupportedInstanceTypeFamily API
func CreateDescribeDeploymentSetSupportedInstanceTypeFamilyRequest() (request *DescribeDeploymentSetSupportedInstanceTypeFamilyRequest) {
	request = &DescribeDeploymentSetSupportedInstanceTypeFamilyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeDeploymentSetSupportedInstanceTypeFamily", "ecs", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDeploymentSetSupportedInstanceTypeFamilyResponse creates a response to parse from DescribeDeploymentSetSupportedInstanceTypeFamily response
func CreateDescribeDeploymentSetSupportedInstanceTypeFamilyResponse() (response *DescribeDeploymentSetSupportedInstanceTypeFamilyResponse) {
	response = &DescribeDeploymentSetSupportedInstanceTypeFamilyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
