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

// DescribeLogStoreExistStatus invokes the ddoscoo.DescribeLogStoreExistStatus API synchronously
func (client *Client) DescribeLogStoreExistStatus(request *DescribeLogStoreExistStatusRequest) (response *DescribeLogStoreExistStatusResponse, err error) {
	response = CreateDescribeLogStoreExistStatusResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLogStoreExistStatusWithChan invokes the ddoscoo.DescribeLogStoreExistStatus API asynchronously
func (client *Client) DescribeLogStoreExistStatusWithChan(request *DescribeLogStoreExistStatusRequest) (<-chan *DescribeLogStoreExistStatusResponse, <-chan error) {
	responseChan := make(chan *DescribeLogStoreExistStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLogStoreExistStatus(request)
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

// DescribeLogStoreExistStatusWithCallback invokes the ddoscoo.DescribeLogStoreExistStatus API asynchronously
func (client *Client) DescribeLogStoreExistStatusWithCallback(request *DescribeLogStoreExistStatusRequest, callback func(response *DescribeLogStoreExistStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLogStoreExistStatusResponse
		var err error
		defer close(result)
		response, err = client.DescribeLogStoreExistStatus(request)
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

// DescribeLogStoreExistStatusRequest is the request struct for api DescribeLogStoreExistStatus
type DescribeLogStoreExistStatusRequest struct {
	*requests.RpcRequest
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
	Lang            string `position:"Query" name:"Lang"`
}

// DescribeLogStoreExistStatusResponse is the response struct for api DescribeLogStoreExistStatus
type DescribeLogStoreExistStatusResponse struct {
	*responses.BaseResponse
	ExistStatus bool   `json:"ExistStatus" xml:"ExistStatus"`
	RequestId   string `json:"RequestId" xml:"RequestId"`
}

// CreateDescribeLogStoreExistStatusRequest creates a request to invoke DescribeLogStoreExistStatus API
func CreateDescribeLogStoreExistStatusRequest() (request *DescribeLogStoreExistStatusRequest) {
	request = &DescribeLogStoreExistStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeLogStoreExistStatus", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeLogStoreExistStatusResponse creates a response to parse from DescribeLogStoreExistStatus response
func CreateDescribeLogStoreExistStatusResponse() (response *DescribeLogStoreExistStatusResponse) {
	response = &DescribeLogStoreExistStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
