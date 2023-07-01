package cbn

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

// AssociateTransitRouterMulticastDomain invokes the cbn.AssociateTransitRouterMulticastDomain API synchronously
func (client *Client) AssociateTransitRouterMulticastDomain(request *AssociateTransitRouterMulticastDomainRequest) (response *AssociateTransitRouterMulticastDomainResponse, err error) {
	response = CreateAssociateTransitRouterMulticastDomainResponse()
	err = client.DoAction(request, response)
	return
}

// AssociateTransitRouterMulticastDomainWithChan invokes the cbn.AssociateTransitRouterMulticastDomain API asynchronously
func (client *Client) AssociateTransitRouterMulticastDomainWithChan(request *AssociateTransitRouterMulticastDomainRequest) (<-chan *AssociateTransitRouterMulticastDomainResponse, <-chan error) {
	responseChan := make(chan *AssociateTransitRouterMulticastDomainResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AssociateTransitRouterMulticastDomain(request)
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

// AssociateTransitRouterMulticastDomainWithCallback invokes the cbn.AssociateTransitRouterMulticastDomain API asynchronously
func (client *Client) AssociateTransitRouterMulticastDomainWithCallback(request *AssociateTransitRouterMulticastDomainRequest, callback func(response *AssociateTransitRouterMulticastDomainResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AssociateTransitRouterMulticastDomainResponse
		var err error
		defer close(result)
		response, err = client.AssociateTransitRouterMulticastDomain(request)
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

// AssociateTransitRouterMulticastDomainRequest is the request struct for api AssociateTransitRouterMulticastDomain
type AssociateTransitRouterMulticastDomainRequest struct {
	*requests.RpcRequest
	ResourceOwnerId                requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken                    string           `position:"Query" name:"ClientToken"`
	VSwitchIds                     *[]string        `position:"Query" name:"VSwitchIds"  type:"Repeated"`
	TransitRouterMulticastDomainId string           `position:"Query" name:"TransitRouterMulticastDomainId"`
	DryRun                         requests.Boolean `position:"Query" name:"DryRun"`
	ResourceOwnerAccount           string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                   string           `position:"Query" name:"OwnerAccount"`
	OwnerId                        requests.Integer `position:"Query" name:"OwnerId"`
	Version                        string           `position:"Query" name:"Version"`
	TransitRouterAttachmentId      string           `position:"Query" name:"TransitRouterAttachmentId"`
}

// AssociateTransitRouterMulticastDomainResponse is the response struct for api AssociateTransitRouterMulticastDomain
type AssociateTransitRouterMulticastDomainResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateAssociateTransitRouterMulticastDomainRequest creates a request to invoke AssociateTransitRouterMulticastDomain API
func CreateAssociateTransitRouterMulticastDomainRequest() (request *AssociateTransitRouterMulticastDomainRequest) {
	request = &AssociateTransitRouterMulticastDomainRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cbn", "2017-09-12", "AssociateTransitRouterMulticastDomain", "", "")
	request.Method = requests.POST
	return
}

// CreateAssociateTransitRouterMulticastDomainResponse creates a response to parse from AssociateTransitRouterMulticastDomain response
func CreateAssociateTransitRouterMulticastDomainResponse() (response *AssociateTransitRouterMulticastDomainResponse) {
	response = &AssociateTransitRouterMulticastDomainResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
