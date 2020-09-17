package oos

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

// ListTemplates invokes the oos.ListTemplates API synchronously
// api document: https://help.aliyun.com/api/oos/listtemplates.html
func (client *Client) ListTemplates(request *ListTemplatesRequest) (response *ListTemplatesResponse, err error) {
	response = CreateListTemplatesResponse()
	err = client.DoAction(request, response)
	return
}

// ListTemplatesWithChan invokes the oos.ListTemplates API asynchronously
// api document: https://help.aliyun.com/api/oos/listtemplates.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListTemplatesWithChan(request *ListTemplatesRequest) (<-chan *ListTemplatesResponse, <-chan error) {
	responseChan := make(chan *ListTemplatesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListTemplates(request)
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

// ListTemplatesWithCallback invokes the oos.ListTemplates API asynchronously
// api document: https://help.aliyun.com/api/oos/listtemplates.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListTemplatesWithCallback(request *ListTemplatesRequest, callback func(response *ListTemplatesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListTemplatesResponse
		var err error
		defer close(result)
		response, err = client.ListTemplates(request)
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

// ListTemplatesRequest is the request struct for api ListTemplates
type ListTemplatesRequest struct {
	*requests.RpcRequest
	CreatedDateBefore string                 `position:"Query" name:"CreatedDateBefore"`
	CreatedBy         string                 `position:"Query" name:"CreatedBy"`
	NextToken         string                 `position:"Query" name:"NextToken"`
	TemplateType      string                 `position:"Query" name:"TemplateType"`
	TemplateName      string                 `position:"Query" name:"TemplateName"`
	SortOrder         string                 `position:"Query" name:"SortOrder"`
	ShareType         string                 `position:"Query" name:"ShareType"`
	HasTrigger        requests.Boolean       `position:"Query" name:"HasTrigger"`
	CreatedDateAfter  string                 `position:"Query" name:"CreatedDateAfter"`
	Tags              map[string]interface{} `position:"Query" name:"Tags"`
	MaxResults        requests.Integer       `position:"Query" name:"MaxResults"`
	TemplateFormat    string                 `position:"Query" name:"TemplateFormat"`
	SortField         string                 `position:"Query" name:"SortField"`
	Category          string                 `position:"Query" name:"Category"`
}

// ListTemplatesResponse is the response struct for api ListTemplates
type ListTemplatesResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	MaxResults int        `json:"MaxResults" xml:"MaxResults"`
	NextToken  string     `json:"NextToken" xml:"NextToken"`
	Templates  []Template `json:"Templates" xml:"Templates"`
}

// CreateListTemplatesRequest creates a request to invoke ListTemplates API
func CreateListTemplatesRequest() (request *ListTemplatesRequest) {
	request = &ListTemplatesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("oos", "2019-06-01", "ListTemplates", "oos", "openAPI")
	request.Method = requests.POST
	return
}

// CreateListTemplatesResponse creates a response to parse from ListTemplates response
func CreateListTemplatesResponse() (response *ListTemplatesResponse) {
	response = &ListTemplatesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
