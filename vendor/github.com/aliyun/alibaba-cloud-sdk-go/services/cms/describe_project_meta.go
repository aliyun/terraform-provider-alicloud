package cms

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

// DescribeProjectMeta invokes the cms.DescribeProjectMeta API synchronously
func (client *Client) DescribeProjectMeta(request *DescribeProjectMetaRequest) (response *DescribeProjectMetaResponse, err error) {
	response = CreateDescribeProjectMetaResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeProjectMetaWithChan invokes the cms.DescribeProjectMeta API asynchronously
func (client *Client) DescribeProjectMetaWithChan(request *DescribeProjectMetaRequest) (<-chan *DescribeProjectMetaResponse, <-chan error) {
	responseChan := make(chan *DescribeProjectMetaResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeProjectMeta(request)
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

// DescribeProjectMetaWithCallback invokes the cms.DescribeProjectMeta API asynchronously
func (client *Client) DescribeProjectMetaWithCallback(request *DescribeProjectMetaRequest, callback func(response *DescribeProjectMetaResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeProjectMetaResponse
		var err error
		defer close(result)
		response, err = client.DescribeProjectMeta(request)
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

// DescribeProjectMetaRequest is the request struct for api DescribeProjectMeta
type DescribeProjectMetaRequest struct {
	*requests.RpcRequest
	PageSize   requests.Integer `position:"Query" name:"PageSize"`
	PageNumber requests.Integer `position:"Query" name:"PageNumber"`
	Labels     string           `position:"Query" name:"Labels"`
}

// DescribeProjectMetaResponse is the response struct for api DescribeProjectMeta
type DescribeProjectMetaResponse struct {
	*responses.BaseResponse
	RequestId  string                         `json:"RequestId" xml:"RequestId"`
	Success    bool                           `json:"Success" xml:"Success"`
	Code       string                         `json:"Code" xml:"Code"`
	Message    string                         `json:"Message" xml:"Message"`
	PageSize   string                         `json:"PageSize" xml:"PageSize"`
	PageNumber string                         `json:"PageNumber" xml:"PageNumber"`
	Total      string                         `json:"Total" xml:"Total"`
	Resources  ResourcesInDescribeProjectMeta `json:"Resources" xml:"Resources"`
}

// CreateDescribeProjectMetaRequest creates a request to invoke DescribeProjectMeta API
func CreateDescribeProjectMetaRequest() (request *DescribeProjectMetaRequest) {
	request = &DescribeProjectMetaRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "DescribeProjectMeta", "Cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeProjectMetaResponse creates a response to parse from DescribeProjectMeta response
func CreateDescribeProjectMetaResponse() (response *DescribeProjectMetaResponse) {
	response = &DescribeProjectMetaResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
