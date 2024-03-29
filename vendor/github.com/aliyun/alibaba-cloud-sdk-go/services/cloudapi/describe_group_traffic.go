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

// DescribeGroupTraffic invokes the cloudapi.DescribeGroupTraffic API synchronously
func (client *Client) DescribeGroupTraffic(request *DescribeGroupTrafficRequest) (response *DescribeGroupTrafficResponse, err error) {
	response = CreateDescribeGroupTrafficResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeGroupTrafficWithChan invokes the cloudapi.DescribeGroupTraffic API asynchronously
func (client *Client) DescribeGroupTrafficWithChan(request *DescribeGroupTrafficRequest) (<-chan *DescribeGroupTrafficResponse, <-chan error) {
	responseChan := make(chan *DescribeGroupTrafficResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeGroupTraffic(request)
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

// DescribeGroupTrafficWithCallback invokes the cloudapi.DescribeGroupTraffic API asynchronously
func (client *Client) DescribeGroupTrafficWithCallback(request *DescribeGroupTrafficRequest, callback func(response *DescribeGroupTrafficResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeGroupTrafficResponse
		var err error
		defer close(result)
		response, err = client.DescribeGroupTraffic(request)
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

// DescribeGroupTrafficRequest is the request struct for api DescribeGroupTraffic
type DescribeGroupTrafficRequest struct {
	*requests.RpcRequest
	StageName     string `position:"Query" name:"StageName"`
	GroupId       string `position:"Query" name:"GroupId"`
	EndTime       string `position:"Query" name:"EndTime"`
	StartTime     string `position:"Query" name:"StartTime"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
}

// DescribeGroupTrafficResponse is the response struct for api DescribeGroupTraffic
type DescribeGroupTrafficResponse struct {
	*responses.BaseResponse
	RequestId        string           `json:"RequestId" xml:"RequestId"`
	TrafficPerSecond TrafficPerSecond `json:"TrafficPerSecond" xml:"TrafficPerSecond"`
}

// CreateDescribeGroupTrafficRequest creates a request to invoke DescribeGroupTraffic API
func CreateDescribeGroupTrafficRequest() (request *DescribeGroupTrafficRequest) {
	request = &DescribeGroupTrafficRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeGroupTraffic", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeGroupTrafficResponse creates a response to parse from DescribeGroupTraffic response
func CreateDescribeGroupTrafficResponse() (response *DescribeGroupTrafficResponse) {
	response = &DescribeGroupTrafficResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
