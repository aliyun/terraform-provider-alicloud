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

// DescribeSecurityIps invokes the dds.DescribeSecurityIps API synchronously
func (client *Client) DescribeSecurityIps(request *DescribeSecurityIpsRequest) (response *DescribeSecurityIpsResponse, err error) {
	response = CreateDescribeSecurityIpsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSecurityIpsWithChan invokes the dds.DescribeSecurityIps API asynchronously
func (client *Client) DescribeSecurityIpsWithChan(request *DescribeSecurityIpsRequest) (<-chan *DescribeSecurityIpsResponse, <-chan error) {
	responseChan := make(chan *DescribeSecurityIpsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSecurityIps(request)
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

// DescribeSecurityIpsWithCallback invokes the dds.DescribeSecurityIps API asynchronously
func (client *Client) DescribeSecurityIpsWithCallback(request *DescribeSecurityIpsRequest, callback func(response *DescribeSecurityIpsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSecurityIpsResponse
		var err error
		defer close(result)
		response, err = client.DescribeSecurityIps(request)
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

// DescribeSecurityIpsRequest is the request struct for api DescribeSecurityIps
type DescribeSecurityIpsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeSecurityIpsResponse is the response struct for api DescribeSecurityIps
type DescribeSecurityIpsResponse struct {
	*responses.BaseResponse
	RequestId        string           `json:"RequestId" xml:"RequestId"`
	SecurityIps      string           `json:"SecurityIps" xml:"SecurityIps"`
	SecurityIpGroups SecurityIpGroups `json:"SecurityIpGroups" xml:"SecurityIpGroups"`
}

// CreateDescribeSecurityIpsRequest creates a request to invoke DescribeSecurityIps API
func CreateDescribeSecurityIpsRequest() (request *DescribeSecurityIpsRequest) {
	request = &DescribeSecurityIpsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "DescribeSecurityIps", "dds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSecurityIpsResponse creates a response to parse from DescribeSecurityIps response
func CreateDescribeSecurityIpsResponse() (response *DescribeSecurityIpsResponse) {
	response = &DescribeSecurityIpsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
