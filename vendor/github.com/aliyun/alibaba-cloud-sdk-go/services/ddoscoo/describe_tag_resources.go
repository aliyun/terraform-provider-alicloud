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

// DescribeTagResources invokes the ddoscoo.DescribeTagResources API synchronously
// api document: https://help.aliyun.com/api/ddoscoo/describetagresources.html
func (client *Client) DescribeTagResources(request *DescribeTagResourcesRequest) (response *DescribeTagResourcesResponse, err error) {
	response = CreateDescribeTagResourcesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeTagResourcesWithChan invokes the ddoscoo.DescribeTagResources API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describetagresources.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeTagResourcesWithChan(request *DescribeTagResourcesRequest) (<-chan *DescribeTagResourcesResponse, <-chan error) {
	responseChan := make(chan *DescribeTagResourcesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeTagResources(request)
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

// DescribeTagResourcesWithCallback invokes the ddoscoo.DescribeTagResources API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describetagresources.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeTagResourcesWithCallback(request *DescribeTagResourcesRequest, callback func(response *DescribeTagResourcesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeTagResourcesResponse
		var err error
		defer close(result)
		response, err = client.DescribeTagResources(request)
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

// DescribeTagResourcesRequest is the request struct for api DescribeTagResources
type DescribeTagResourcesRequest struct {
	*requests.RpcRequest
	ResourceGroupId string                      `position:"Query" name:"ResourceGroupId"`
	SourceIp        string                      `position:"Query" name:"SourceIp"`
	NextToken       string                      `position:"Query" name:"NextToken"`
	ResourceType    string                      `position:"Query" name:"ResourceType"`
	Tags            *[]DescribeTagResourcesTags `position:"Query" name:"Tags"  type:"Repeated"`
	ResourceIds     *[]string                   `position:"Query" name:"ResourceIds"  type:"Repeated"`
}

// DescribeTagResourcesTags is a repeated param struct in DescribeTagResourcesRequest
type DescribeTagResourcesTags struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// DescribeTagResourcesResponse is the response struct for api DescribeTagResources
type DescribeTagResourcesResponse struct {
	*responses.BaseResponse
	RequestId    string       `json:"RequestId" xml:"RequestId"`
	NextToken    string       `json:"NextToken" xml:"NextToken"`
	TagResources TagResources `json:"TagResources" xml:"TagResources"`
}

// CreateDescribeTagResourcesRequest creates a request to invoke DescribeTagResources API
func CreateDescribeTagResourcesRequest() (request *DescribeTagResourcesRequest) {
	request = &DescribeTagResourcesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeTagResources", "ddoscoo", "openAPI")
	return
}

// CreateDescribeTagResourcesResponse creates a response to parse from DescribeTagResources response
func CreateDescribeTagResourcesResponse() (response *DescribeTagResourcesResponse) {
	response = &DescribeTagResourcesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
