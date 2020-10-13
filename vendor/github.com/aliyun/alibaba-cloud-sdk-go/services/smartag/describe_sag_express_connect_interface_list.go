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

// DescribeSagExpressConnectInterfaceList invokes the smartag.DescribeSagExpressConnectInterfaceList API synchronously
func (client *Client) DescribeSagExpressConnectInterfaceList(request *DescribeSagExpressConnectInterfaceListRequest) (response *DescribeSagExpressConnectInterfaceListResponse, err error) {
	response = CreateDescribeSagExpressConnectInterfaceListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSagExpressConnectInterfaceListWithChan invokes the smartag.DescribeSagExpressConnectInterfaceList API asynchronously
func (client *Client) DescribeSagExpressConnectInterfaceListWithChan(request *DescribeSagExpressConnectInterfaceListRequest) (<-chan *DescribeSagExpressConnectInterfaceListResponse, <-chan error) {
	responseChan := make(chan *DescribeSagExpressConnectInterfaceListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSagExpressConnectInterfaceList(request)
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

// DescribeSagExpressConnectInterfaceListWithCallback invokes the smartag.DescribeSagExpressConnectInterfaceList API asynchronously
func (client *Client) DescribeSagExpressConnectInterfaceListWithCallback(request *DescribeSagExpressConnectInterfaceListRequest, callback func(response *DescribeSagExpressConnectInterfaceListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSagExpressConnectInterfaceListResponse
		var err error
		defer close(result)
		response, err = client.DescribeSagExpressConnectInterfaceList(request)
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

// DescribeSagExpressConnectInterfaceListRequest is the request struct for api DescribeSagExpressConnectInterfaceList
type DescribeSagExpressConnectInterfaceListRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	SmartAGId            string           `position:"Query" name:"SmartAGId"`
	SmartAGSn            string           `position:"Query" name:"SmartAGSn"`
	PortName             string           `position:"Query" name:"PortName"`
}

// DescribeSagExpressConnectInterfaceListResponse is the response struct for api DescribeSagExpressConnectInterfaceList
type DescribeSagExpressConnectInterfaceListResponse struct {
	*responses.BaseResponse
	RequestId  string      `json:"RequestId" xml:"RequestId"`
	Interfaces []Interface `json:"Interfaces" xml:"Interfaces"`
	TaskStates []TaskState `json:"TaskStates" xml:"TaskStates"`
}

// CreateDescribeSagExpressConnectInterfaceListRequest creates a request to invoke DescribeSagExpressConnectInterfaceList API
func CreateDescribeSagExpressConnectInterfaceListRequest() (request *DescribeSagExpressConnectInterfaceListRequest) {
	request = &DescribeSagExpressConnectInterfaceListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DescribeSagExpressConnectInterfaceList", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSagExpressConnectInterfaceListResponse creates a response to parse from DescribeSagExpressConnectInterfaceList response
func CreateDescribeSagExpressConnectInterfaceListResponse() (response *DescribeSagExpressConnectInterfaceListResponse) {
	response = &DescribeSagExpressConnectInterfaceListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
