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

// DescribeContactList invokes the cms.DescribeContactList API synchronously
func (client *Client) DescribeContactList(request *DescribeContactListRequest) (response *DescribeContactListResponse, err error) {
	response = CreateDescribeContactListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeContactListWithChan invokes the cms.DescribeContactList API asynchronously
func (client *Client) DescribeContactListWithChan(request *DescribeContactListRequest) (<-chan *DescribeContactListResponse, <-chan error) {
	responseChan := make(chan *DescribeContactListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeContactList(request)
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

// DescribeContactListWithCallback invokes the cms.DescribeContactList API asynchronously
func (client *Client) DescribeContactListWithCallback(request *DescribeContactListRequest, callback func(response *DescribeContactListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeContactListResponse
		var err error
		defer close(result)
		response, err = client.DescribeContactList(request)
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

// DescribeContactListRequest is the request struct for api DescribeContactList
type DescribeContactListRequest struct {
	*requests.RpcRequest
	ChanelType  string           `position:"Query" name:"ChanelType"`
	PageNumber  requests.Integer `position:"Query" name:"PageNumber"`
	ContactName string           `position:"Query" name:"ContactName"`
	PageSize    requests.Integer `position:"Query" name:"PageSize"`
	ChanelValue string           `position:"Query" name:"ChanelValue"`
}

// DescribeContactListResponse is the response struct for api DescribeContactList
type DescribeContactListResponse struct {
	*responses.BaseResponse
	Success   bool                          `json:"Success" xml:"Success"`
	Code      string                        `json:"Code" xml:"Code"`
	Message   string                        `json:"Message" xml:"Message"`
	Total     int                           `json:"Total" xml:"Total"`
	RequestId string                        `json:"RequestId" xml:"RequestId"`
	Contacts  ContactsInDescribeContactList `json:"Contacts" xml:"Contacts"`
}

// CreateDescribeContactListRequest creates a request to invoke DescribeContactList API
func CreateDescribeContactListRequest() (request *DescribeContactListRequest) {
	request = &DescribeContactListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "DescribeContactList", "cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeContactListResponse creates a response to parse from DescribeContactList response
func CreateDescribeContactListResponse() (response *DescribeContactListResponse) {
	response = &DescribeContactListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
