package cloudapi

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

// DescribeGroupQps invokes the cloudapi.DescribeGroupQps API synchronously
func (client *Client) DescribeGroupQps(request *DescribeGroupQpsRequest) (response *DescribeGroupQpsResponse, err error) {
	response = CreateDescribeGroupQpsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeGroupQpsWithChan invokes the cloudapi.DescribeGroupQps API asynchronously
func (client *Client) DescribeGroupQpsWithChan(request *DescribeGroupQpsRequest) (<-chan *DescribeGroupQpsResponse, <-chan error) {
	responseChan := make(chan *DescribeGroupQpsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeGroupQps(request)
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

// DescribeGroupQpsWithCallback invokes the cloudapi.DescribeGroupQps API asynchronously
func (client *Client) DescribeGroupQpsWithCallback(request *DescribeGroupQpsRequest, callback func(response *DescribeGroupQpsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeGroupQpsResponse
		var err error
		defer close(result)
		response, err = client.DescribeGroupQps(request)
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

// DescribeGroupQpsRequest is the request struct for api DescribeGroupQps
type DescribeGroupQpsRequest struct {
	*requests.RpcRequest
	StageName     string `position:"Query" name:"StageName"`
	GroupId       string `position:"Query" name:"GroupId"`
	EndTime       string `position:"Query" name:"EndTime"`
	StartTime     string `position:"Query" name:"StartTime"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
}

// DescribeGroupQpsResponse is the response struct for api DescribeGroupQps
type DescribeGroupQpsResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	GroupQps  GroupQps `json:"GroupQps" xml:"GroupQps"`
}

// CreateDescribeGroupQpsRequest creates a request to invoke DescribeGroupQps API
func CreateDescribeGroupQpsRequest() (request *DescribeGroupQpsRequest) {
	request = &DescribeGroupQpsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeGroupQps", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeGroupQpsResponse creates a response to parse from DescribeGroupQps response
func CreateDescribeGroupQpsResponse() (response *DescribeGroupQpsResponse) {
	response = &DescribeGroupQpsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
