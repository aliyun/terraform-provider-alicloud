package ddoscoo

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

// DescribePortAutoCcStatus invokes the ddoscoo.DescribePortAutoCcStatus API synchronously
func (client *Client) DescribePortAutoCcStatus(request *DescribePortAutoCcStatusRequest) (response *DescribePortAutoCcStatusResponse, err error) {
	response = CreateDescribePortAutoCcStatusResponse()
	err = client.DoAction(request, response)
	return
}

// DescribePortAutoCcStatusWithChan invokes the ddoscoo.DescribePortAutoCcStatus API asynchronously
func (client *Client) DescribePortAutoCcStatusWithChan(request *DescribePortAutoCcStatusRequest) (<-chan *DescribePortAutoCcStatusResponse, <-chan error) {
	responseChan := make(chan *DescribePortAutoCcStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribePortAutoCcStatus(request)
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

// DescribePortAutoCcStatusWithCallback invokes the ddoscoo.DescribePortAutoCcStatus API asynchronously
func (client *Client) DescribePortAutoCcStatusWithCallback(request *DescribePortAutoCcStatusRequest, callback func(response *DescribePortAutoCcStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribePortAutoCcStatusResponse
		var err error
		defer close(result)
		response, err = client.DescribePortAutoCcStatus(request)
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

// DescribePortAutoCcStatusRequest is the request struct for api DescribePortAutoCcStatus
type DescribePortAutoCcStatusRequest struct {
	*requests.RpcRequest
	SourceIp    string    `position:"Query" name:"SourceIp"`
	InstanceIds *[]string `position:"Query" name:"InstanceIds"  type:"Repeated"`
}

// DescribePortAutoCcStatusResponse is the response struct for api DescribePortAutoCcStatus
type DescribePortAutoCcStatusResponse struct {
	*responses.BaseResponse
	RequestId        string   `json:"RequestId" xml:"RequestId"`
	PortAutoCcStatus []Status `json:"PortAutoCcStatus" xml:"PortAutoCcStatus"`
}

// CreateDescribePortAutoCcStatusRequest creates a request to invoke DescribePortAutoCcStatus API
func CreateDescribePortAutoCcStatusRequest() (request *DescribePortAutoCcStatusRequest) {
	request = &DescribePortAutoCcStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribePortAutoCcStatus", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribePortAutoCcStatusResponse creates a response to parse from DescribePortAutoCcStatus response
func CreateDescribePortAutoCcStatusResponse() (response *DescribePortAutoCcStatusResponse) {
	response = &DescribePortAutoCcStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
