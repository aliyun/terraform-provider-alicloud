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

// DescribeSagPortRouteProtocolList invokes the smartag.DescribeSagPortRouteProtocolList API synchronously
func (client *Client) DescribeSagPortRouteProtocolList(request *DescribeSagPortRouteProtocolListRequest) (response *DescribeSagPortRouteProtocolListResponse, err error) {
	response = CreateDescribeSagPortRouteProtocolListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSagPortRouteProtocolListWithChan invokes the smartag.DescribeSagPortRouteProtocolList API asynchronously
func (client *Client) DescribeSagPortRouteProtocolListWithChan(request *DescribeSagPortRouteProtocolListRequest) (<-chan *DescribeSagPortRouteProtocolListResponse, <-chan error) {
	responseChan := make(chan *DescribeSagPortRouteProtocolListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSagPortRouteProtocolList(request)
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

// DescribeSagPortRouteProtocolListWithCallback invokes the smartag.DescribeSagPortRouteProtocolList API asynchronously
func (client *Client) DescribeSagPortRouteProtocolListWithCallback(request *DescribeSagPortRouteProtocolListRequest, callback func(response *DescribeSagPortRouteProtocolListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSagPortRouteProtocolListResponse
		var err error
		defer close(result)
		response, err = client.DescribeSagPortRouteProtocolList(request)
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

// DescribeSagPortRouteProtocolListRequest is the request struct for api DescribeSagPortRouteProtocolList
type DescribeSagPortRouteProtocolListRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query"`
	ResourceOwnerAccount string           `position:"Query"`
	OwnerAccount         string           `position:"Query"`
	OwnerId              requests.Integer `position:"Query"`
	SmartAGId            string           `position:"Query"`
	SmartAGSn            string           `position:"Query"`
}

// DescribeSagPortRouteProtocolListResponse is the response struct for api DescribeSagPortRouteProtocolList
type DescribeSagPortRouteProtocolListResponse struct {
	*responses.BaseResponse
	RequestId  string      `json:"RequestId" xml:"RequestId"`
	Ports      []Port      `json:"Ports" xml:"Ports"`
	TaskStates []TaskState `json:"TaskStates" xml:"TaskStates"`
}

// CreateDescribeSagPortRouteProtocolListRequest creates a request to invoke DescribeSagPortRouteProtocolList API
func CreateDescribeSagPortRouteProtocolListRequest() (request *DescribeSagPortRouteProtocolListRequest) {
	request = &DescribeSagPortRouteProtocolListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DescribeSagPortRouteProtocolList", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSagPortRouteProtocolListResponse creates a response to parse from DescribeSagPortRouteProtocolList response
func CreateDescribeSagPortRouteProtocolListResponse() (response *DescribeSagPortRouteProtocolListResponse) {
	response = &DescribeSagPortRouteProtocolListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
