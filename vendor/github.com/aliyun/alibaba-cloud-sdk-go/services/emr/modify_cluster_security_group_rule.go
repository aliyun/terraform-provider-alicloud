package emr

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

// ModifyClusterSecurityGroupRule invokes the emr.ModifyClusterSecurityGroupRule API synchronously
func (client *Client) ModifyClusterSecurityGroupRule(request *ModifyClusterSecurityGroupRuleRequest) (response *ModifyClusterSecurityGroupRuleResponse, err error) {
	response = CreateModifyClusterSecurityGroupRuleResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyClusterSecurityGroupRuleWithChan invokes the emr.ModifyClusterSecurityGroupRule API asynchronously
func (client *Client) ModifyClusterSecurityGroupRuleWithChan(request *ModifyClusterSecurityGroupRuleRequest) (<-chan *ModifyClusterSecurityGroupRuleResponse, <-chan error) {
	responseChan := make(chan *ModifyClusterSecurityGroupRuleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyClusterSecurityGroupRule(request)
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

// ModifyClusterSecurityGroupRuleWithCallback invokes the emr.ModifyClusterSecurityGroupRule API asynchronously
func (client *Client) ModifyClusterSecurityGroupRuleWithCallback(request *ModifyClusterSecurityGroupRuleRequest, callback func(response *ModifyClusterSecurityGroupRuleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyClusterSecurityGroupRuleResponse
		var err error
		defer close(result)
		response, err = client.ModifyClusterSecurityGroupRule(request)
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

// ModifyClusterSecurityGroupRuleRequest is the request struct for api ModifyClusterSecurityGroupRule
type ModifyClusterSecurityGroupRuleRequest struct {
	*requests.RpcRequest
	NicType         string           `position:"Query" name:"NicType"`
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PortRange       string           `position:"Query" name:"PortRange"`
	IpProtocol      string           `position:"Query" name:"IpProtocol"`
	ClusterId       string           `position:"Query" name:"ClusterId"`
	WhiteIp         string           `position:"Query" name:"WhiteIp"`
	ModifyType      string           `position:"Query" name:"ModifyType"`
}

// ModifyClusterSecurityGroupRuleResponse is the response struct for api ModifyClusterSecurityGroupRule
type ModifyClusterSecurityGroupRuleResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyClusterSecurityGroupRuleRequest creates a request to invoke ModifyClusterSecurityGroupRule API
func CreateModifyClusterSecurityGroupRuleRequest() (request *ModifyClusterSecurityGroupRuleRequest) {
	request = &ModifyClusterSecurityGroupRuleRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "ModifyClusterSecurityGroupRule", "emr", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyClusterSecurityGroupRuleResponse creates a response to parse from ModifyClusterSecurityGroupRule response
func CreateModifyClusterSecurityGroupRuleResponse() (response *ModifyClusterSecurityGroupRuleResponse) {
	response = &ModifyClusterSecurityGroupRuleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
