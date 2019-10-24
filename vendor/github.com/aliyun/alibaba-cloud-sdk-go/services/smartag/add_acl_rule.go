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

// AddACLRule invokes the smartag.AddACLRule API synchronously
// api document: https://help.aliyun.com/api/smartag/addaclrule.html
func (client *Client) AddACLRule(request *AddACLRuleRequest) (response *AddACLRuleResponse, err error) {
	response = CreateAddACLRuleResponse()
	err = client.DoAction(request, response)
	return
}

// AddACLRuleWithChan invokes the smartag.AddACLRule API asynchronously
// api document: https://help.aliyun.com/api/smartag/addaclrule.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddACLRuleWithChan(request *AddACLRuleRequest) (<-chan *AddACLRuleResponse, <-chan error) {
	responseChan := make(chan *AddACLRuleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddACLRule(request)
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

// AddACLRuleWithCallback invokes the smartag.AddACLRule API asynchronously
// api document: https://help.aliyun.com/api/smartag/addaclrule.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddACLRuleWithCallback(request *AddACLRuleRequest, callback func(response *AddACLRuleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddACLRuleResponse
		var err error
		defer close(result)
		response, err = client.AddACLRule(request)
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

// AddACLRuleRequest is the request struct for api AddACLRule
type AddACLRuleRequest struct {
	*requests.RpcRequest
	AclId                string           `position:"Query" name:"AclId"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SourcePortRange      string           `position:"Query" name:"SourcePortRange"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	IpProtocol           string           `position:"Query" name:"IpProtocol"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	SourceCidr           string           `position:"Query" name:"SourceCidr"`
	Description          string           `position:"Query" name:"Description"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Priority             requests.Integer `position:"Query" name:"Priority"`
	DestCidr             string           `position:"Query" name:"DestCidr"`
	DestPortRange        string           `position:"Query" name:"DestPortRange"`
	Direction            string           `position:"Query" name:"Direction"`
	Policy               string           `position:"Query" name:"Policy"`
}

// AddACLRuleResponse is the response struct for api AddACLRule
type AddACLRuleResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	AcrId           string `json:"AcrId" xml:"AcrId"`
	AclId           string `json:"AclId" xml:"AclId"`
	Description     string `json:"Description" xml:"Description"`
	Direction       string `json:"Direction" xml:"Direction"`
	SourceCidr      string `json:"SourceCidr" xml:"SourceCidr"`
	DestCidr        string `json:"DestCidr" xml:"DestCidr"`
	IpProtocol      string `json:"IpProtocol" xml:"IpProtocol"`
	SourcePortRange string `json:"SourcePortRange" xml:"SourcePortRange"`
	DestPortRange   string `json:"DestPortRange" xml:"DestPortRange"`
	Policy          string `json:"Policy" xml:"Policy"`
	Priority        int    `json:"Priority" xml:"Priority"`
	GmtCreate       int    `json:"GmtCreate" xml:"GmtCreate"`
}

// CreateAddACLRuleRequest creates a request to invoke AddACLRule API
func CreateAddACLRuleRequest() (request *AddACLRuleRequest) {
	request = &AddACLRuleRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "AddACLRule", "smartag", "openAPI")
	return
}

// CreateAddACLRuleResponse creates a response to parse from AddACLRule response
func CreateAddACLRuleResponse() (response *AddACLRuleResponse) {
	response = &AddACLRuleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
