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

// DescribeAuditFiles invokes the dds.DescribeAuditFiles API synchronously
// api document: https://help.aliyun.com/api/dds/describeauditfiles.html
func (client *Client) DescribeAuditFiles(request *DescribeAuditFilesRequest) (response *DescribeAuditFilesResponse, err error) {
	response = CreateDescribeAuditFilesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAuditFilesWithChan invokes the dds.DescribeAuditFiles API asynchronously
// api document: https://help.aliyun.com/api/dds/describeauditfiles.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAuditFilesWithChan(request *DescribeAuditFilesRequest) (<-chan *DescribeAuditFilesResponse, <-chan error) {
	responseChan := make(chan *DescribeAuditFilesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAuditFiles(request)
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

// DescribeAuditFilesWithCallback invokes the dds.DescribeAuditFiles API asynchronously
// api document: https://help.aliyun.com/api/dds/describeauditfiles.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAuditFilesWithCallback(request *DescribeAuditFilesRequest, callback func(response *DescribeAuditFilesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAuditFilesResponse
		var err error
		defer close(result)
		response, err = client.DescribeAuditFiles(request)
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

// DescribeAuditFilesRequest is the request struct for api DescribeAuditFiles
type DescribeAuditFilesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	NodeId               string           `position:"Query" name:"NodeId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeAuditFilesResponse is the response struct for api DescribeAuditFiles
type DescribeAuditFilesResponse struct {
	*responses.BaseResponse
	RequestId        string                    `json:"RequestId" xml:"RequestId"`
	TotalRecordCount int                       `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageNumber       int                       `json:"PageNumber" xml:"PageNumber"`
	PageRecordCount  int                       `json:"PageRecordCount" xml:"PageRecordCount"`
	DBInstanceId     string                    `json:"DBInstanceId" xml:"DBInstanceId"`
	Items            ItemsInDescribeAuditFiles `json:"Items" xml:"Items"`
}

// CreateDescribeAuditFilesRequest creates a request to invoke DescribeAuditFiles API
func CreateDescribeAuditFilesRequest() (request *DescribeAuditFilesRequest) {
	request = &DescribeAuditFilesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "DescribeAuditFiles", "Dds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAuditFilesResponse creates a response to parse from DescribeAuditFiles response
func CreateDescribeAuditFilesResponse() (response *DescribeAuditFilesResponse) {
	response = &DescribeAuditFilesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
