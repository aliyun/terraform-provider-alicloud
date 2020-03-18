package cdn

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

// DescribeLiveStreamRecordIndexFile invokes the cdn.DescribeLiveStreamRecordIndexFile API synchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamrecordindexfile.html
func (client *Client) DescribeLiveStreamRecordIndexFile(request *DescribeLiveStreamRecordIndexFileRequest) (response *DescribeLiveStreamRecordIndexFileResponse, err error) {
	response = CreateDescribeLiveStreamRecordIndexFileResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLiveStreamRecordIndexFileWithChan invokes the cdn.DescribeLiveStreamRecordIndexFile API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamrecordindexfile.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamRecordIndexFileWithChan(request *DescribeLiveStreamRecordIndexFileRequest) (<-chan *DescribeLiveStreamRecordIndexFileResponse, <-chan error) {
	responseChan := make(chan *DescribeLiveStreamRecordIndexFileResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLiveStreamRecordIndexFile(request)
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

// DescribeLiveStreamRecordIndexFileWithCallback invokes the cdn.DescribeLiveStreamRecordIndexFile API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamrecordindexfile.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamRecordIndexFileWithCallback(request *DescribeLiveStreamRecordIndexFileRequest, callback func(response *DescribeLiveStreamRecordIndexFileResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLiveStreamRecordIndexFileResponse
		var err error
		defer close(result)
		response, err = client.DescribeLiveStreamRecordIndexFile(request)
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

// DescribeLiveStreamRecordIndexFileRequest is the request struct for api DescribeLiveStreamRecordIndexFile
type DescribeLiveStreamRecordIndexFileRequest struct {
	*requests.RpcRequest
	AppName       string           `position:"Query" name:"AppName"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	StreamName    string           `position:"Query" name:"StreamName"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	RecordId      string           `position:"Query" name:"RecordId"`
}

// DescribeLiveStreamRecordIndexFileResponse is the response struct for api DescribeLiveStreamRecordIndexFile
type DescribeLiveStreamRecordIndexFileResponse struct {
	*responses.BaseResponse
	RequestId       string          `json:"RequestId" xml:"RequestId"`
	RecordIndexInfo RecordIndexInfo `json:"RecordIndexInfo" xml:"RecordIndexInfo"`
}

// CreateDescribeLiveStreamRecordIndexFileRequest creates a request to invoke DescribeLiveStreamRecordIndexFile API
func CreateDescribeLiveStreamRecordIndexFileRequest() (request *DescribeLiveStreamRecordIndexFileRequest) {
	request = &DescribeLiveStreamRecordIndexFileRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeLiveStreamRecordIndexFile", "", "")
	return
}

// CreateDescribeLiveStreamRecordIndexFileResponse creates a response to parse from DescribeLiveStreamRecordIndexFile response
func CreateDescribeLiveStreamRecordIndexFileResponse() (response *DescribeLiveStreamRecordIndexFileResponse) {
	response = &DescribeLiveStreamRecordIndexFileResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
