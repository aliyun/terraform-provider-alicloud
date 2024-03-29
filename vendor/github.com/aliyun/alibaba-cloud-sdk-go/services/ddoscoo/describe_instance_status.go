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

// DescribeInstanceStatus invokes the ddoscoo.DescribeInstanceStatus API synchronously
func (client *Client) DescribeInstanceStatus(request *DescribeInstanceStatusRequest) (response *DescribeInstanceStatusResponse, err error) {
	response = CreateDescribeInstanceStatusResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeInstanceStatusWithChan invokes the ddoscoo.DescribeInstanceStatus API asynchronously
func (client *Client) DescribeInstanceStatusWithChan(request *DescribeInstanceStatusRequest) (<-chan *DescribeInstanceStatusResponse, <-chan error) {
	responseChan := make(chan *DescribeInstanceStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeInstanceStatus(request)
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

// DescribeInstanceStatusWithCallback invokes the ddoscoo.DescribeInstanceStatus API asynchronously
func (client *Client) DescribeInstanceStatusWithCallback(request *DescribeInstanceStatusRequest, callback func(response *DescribeInstanceStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeInstanceStatusResponse
		var err error
		defer close(result)
		response, err = client.DescribeInstanceStatus(request)
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

// DescribeInstanceStatusRequest is the request struct for api DescribeInstanceStatus
type DescribeInstanceStatusRequest struct {
	*requests.RpcRequest
	ProductType requests.Integer `position:"Query" name:"ProductType"`
	InstanceId  string           `position:"Query" name:"InstanceId"`
	SourceIp    string           `position:"Query" name:"SourceIp"`
}

// DescribeInstanceStatusResponse is the response struct for api DescribeInstanceStatus
type DescribeInstanceStatusResponse struct {
	*responses.BaseResponse
	InstanceStatus int    `json:"InstanceStatus" xml:"InstanceStatus"`
	RequestId      string `json:"RequestId" xml:"RequestId"`
	InstanceId     string `json:"InstanceId" xml:"InstanceId"`
}

// CreateDescribeInstanceStatusRequest creates a request to invoke DescribeInstanceStatus API
func CreateDescribeInstanceStatusRequest() (request *DescribeInstanceStatusRequest) {
	request = &DescribeInstanceStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeInstanceStatus", "ddoscoo", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeInstanceStatusResponse creates a response to parse from DescribeInstanceStatus response
func CreateDescribeInstanceStatusResponse() (response *DescribeInstanceStatusResponse) {
	response = &DescribeInstanceStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
