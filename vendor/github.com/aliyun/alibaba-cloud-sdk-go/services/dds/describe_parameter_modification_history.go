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

// DescribeParameterModificationHistory invokes the dds.DescribeParameterModificationHistory API synchronously
// api document: https://help.aliyun.com/api/dds/describeparametermodificationhistory.html
func (client *Client) DescribeParameterModificationHistory(request *DescribeParameterModificationHistoryRequest) (response *DescribeParameterModificationHistoryResponse, err error) {
	response = CreateDescribeParameterModificationHistoryResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeParameterModificationHistoryWithChan invokes the dds.DescribeParameterModificationHistory API asynchronously
// api document: https://help.aliyun.com/api/dds/describeparametermodificationhistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeParameterModificationHistoryWithChan(request *DescribeParameterModificationHistoryRequest) (<-chan *DescribeParameterModificationHistoryResponse, <-chan error) {
	responseChan := make(chan *DescribeParameterModificationHistoryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeParameterModificationHistory(request)
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

// DescribeParameterModificationHistoryWithCallback invokes the dds.DescribeParameterModificationHistory API asynchronously
// api document: https://help.aliyun.com/api/dds/describeparametermodificationhistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeParameterModificationHistoryWithCallback(request *DescribeParameterModificationHistoryRequest, callback func(response *DescribeParameterModificationHistoryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeParameterModificationHistoryResponse
		var err error
		defer close(result)
		response, err = client.DescribeParameterModificationHistory(request)
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

// DescribeParameterModificationHistoryRequest is the request struct for api DescribeParameterModificationHistory
type DescribeParameterModificationHistoryRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	EndTime              string           `position:"Query" name:"EndTime"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	StartTime            string           `position:"Query" name:"StartTime"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	NodeId               string           `position:"Query" name:"NodeId"`
}

// DescribeParameterModificationHistoryResponse is the response struct for api DescribeParameterModificationHistory
type DescribeParameterModificationHistoryResponse struct {
	*responses.BaseResponse
	RequestId            string               `json:"RequestId" xml:"RequestId"`
	HistoricalParameters HistoricalParameters `json:"HistoricalParameters" xml:"HistoricalParameters"`
}

// CreateDescribeParameterModificationHistoryRequest creates a request to invoke DescribeParameterModificationHistory API
func CreateDescribeParameterModificationHistoryRequest() (request *DescribeParameterModificationHistoryRequest) {
	request = &DescribeParameterModificationHistoryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "DescribeParameterModificationHistory", "dds", "openAPI")
	return
}

// CreateDescribeParameterModificationHistoryResponse creates a response to parse from DescribeParameterModificationHistory response
func CreateDescribeParameterModificationHistoryResponse() (response *DescribeParameterModificationHistoryResponse) {
	response = &DescribeParameterModificationHistoryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
