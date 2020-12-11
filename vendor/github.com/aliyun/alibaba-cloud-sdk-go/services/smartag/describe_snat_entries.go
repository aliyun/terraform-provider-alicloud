package smartag

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

// DescribeSnatEntries invokes the smartag.DescribeSnatEntries API synchronously
func (client *Client) DescribeSnatEntries(request *DescribeSnatEntriesRequest) (response *DescribeSnatEntriesResponse, err error) {
	response = CreateDescribeSnatEntriesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSnatEntriesWithChan invokes the smartag.DescribeSnatEntries API asynchronously
func (client *Client) DescribeSnatEntriesWithChan(request *DescribeSnatEntriesRequest) (<-chan *DescribeSnatEntriesResponse, <-chan error) {
	responseChan := make(chan *DescribeSnatEntriesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSnatEntries(request)
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

// DescribeSnatEntriesWithCallback invokes the smartag.DescribeSnatEntries API asynchronously
func (client *Client) DescribeSnatEntriesWithCallback(request *DescribeSnatEntriesRequest, callback func(response *DescribeSnatEntriesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSnatEntriesResponse
		var err error
		defer close(result)
		response, err = client.DescribeSnatEntries(request)
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

// DescribeSnatEntriesRequest is the request struct for api DescribeSnatEntries
type DescribeSnatEntriesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	SmartAGId            string           `position:"Query" name:"SmartAGId"`
}

// DescribeSnatEntriesResponse is the response struct for api DescribeSnatEntries
type DescribeSnatEntriesResponse struct {
	*responses.BaseResponse
	RequestId   string      `json:"RequestId" xml:"RequestId"`
	TotalCount  int         `json:"TotalCount" xml:"TotalCount"`
	PageNumber  int         `json:"PageNumber" xml:"PageNumber"`
	PageSize    int         `json:"PageSize" xml:"PageSize"`
	SnatEntries SnatEntries `json:"SnatEntries" xml:"SnatEntries"`
}

// CreateDescribeSnatEntriesRequest creates a request to invoke DescribeSnatEntries API
func CreateDescribeSnatEntriesRequest() (request *DescribeSnatEntriesRequest) {
	request = &DescribeSnatEntriesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DescribeSnatEntries", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSnatEntriesResponse creates a response to parse from DescribeSnatEntries response
func CreateDescribeSnatEntriesResponse() (response *DescribeSnatEntriesResponse) {
	response = &DescribeSnatEntriesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
