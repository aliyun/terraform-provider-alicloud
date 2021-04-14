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

// DescribeDefenseCountStatistics invokes the ddoscoo.DescribeDefenseCountStatistics API synchronously
func (client *Client) DescribeDefenseCountStatistics(request *DescribeDefenseCountStatisticsRequest) (response *DescribeDefenseCountStatisticsResponse, err error) {
	response = CreateDescribeDefenseCountStatisticsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDefenseCountStatisticsWithChan invokes the ddoscoo.DescribeDefenseCountStatistics API asynchronously
func (client *Client) DescribeDefenseCountStatisticsWithChan(request *DescribeDefenseCountStatisticsRequest) (<-chan *DescribeDefenseCountStatisticsResponse, <-chan error) {
	responseChan := make(chan *DescribeDefenseCountStatisticsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDefenseCountStatistics(request)
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

// DescribeDefenseCountStatisticsWithCallback invokes the ddoscoo.DescribeDefenseCountStatistics API asynchronously
func (client *Client) DescribeDefenseCountStatisticsWithCallback(request *DescribeDefenseCountStatisticsRequest, callback func(response *DescribeDefenseCountStatisticsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDefenseCountStatisticsResponse
		var err error
		defer close(result)
		response, err = client.DescribeDefenseCountStatistics(request)
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

// DescribeDefenseCountStatisticsRequest is the request struct for api DescribeDefenseCountStatistics
type DescribeDefenseCountStatisticsRequest struct {
	*requests.RpcRequest
	ResourceGroupId string `position:"Query" name:"ResourceGroupId"`
	SourceIp        string `position:"Query" name:"SourceIp"`
}

// DescribeDefenseCountStatisticsResponse is the response struct for api DescribeDefenseCountStatistics
type DescribeDefenseCountStatisticsResponse struct {
	*responses.BaseResponse
	RequestId              string                 `json:"RequestId" xml:"RequestId"`
	DefenseCountStatistics DefenseCountStatistics `json:"DefenseCountStatistics" xml:"DefenseCountStatistics"`
}

// CreateDescribeDefenseCountStatisticsRequest creates a request to invoke DescribeDefenseCountStatistics API
func CreateDescribeDefenseCountStatisticsRequest() (request *DescribeDefenseCountStatisticsRequest) {
	request = &DescribeDefenseCountStatisticsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeDefenseCountStatistics", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDefenseCountStatisticsResponse creates a response to parse from DescribeDefenseCountStatistics response
func CreateDescribeDefenseCountStatisticsResponse() (response *DescribeDefenseCountStatisticsResponse) {
	response = &DescribeDefenseCountStatisticsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
