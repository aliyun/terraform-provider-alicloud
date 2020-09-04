package alidns

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

// AddDomainRecord invokes the alidns.AddDomainRecord API synchronously
// api document: https://help.aliyun.com/api/alidns/adddomainrecord.html
func (client *Client) AddDomainRecord(request *AddDomainRecordRequest) (response *AddDomainRecordResponse, err error) {
	response = CreateAddDomainRecordResponse()
	err = client.DoAction(request, response)
	return
}

// AddDomainRecordWithChan invokes the alidns.AddDomainRecord API asynchronously
// api document: https://help.aliyun.com/api/alidns/adddomainrecord.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddDomainRecordWithChan(request *AddDomainRecordRequest) (<-chan *AddDomainRecordResponse, <-chan error) {
	responseChan := make(chan *AddDomainRecordResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddDomainRecord(request)
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

// AddDomainRecordWithCallback invokes the alidns.AddDomainRecord API asynchronously
// api document: https://help.aliyun.com/api/alidns/adddomainrecord.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddDomainRecordWithCallback(request *AddDomainRecordRequest, callback func(response *AddDomainRecordResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddDomainRecordResponse
		var err error
		defer close(result)
		response, err = client.AddDomainRecord(request)
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

// AddDomainRecordRequest is the request struct for api AddDomainRecord
type AddDomainRecordRequest struct {
	*requests.RpcRequest
	RR           string           `position:"Query" name:"RR"`
	Line         string           `position:"Query" name:"Line"`
	Type         string           `position:"Query" name:"Type"`
	Lang         string           `position:"Query" name:"Lang"`
	Value        string           `position:"Query" name:"Value"`
	DomainName   string           `position:"Query" name:"DomainName"`
	Priority     requests.Integer `position:"Query" name:"Priority"`
	TTL          requests.Integer `position:"Query" name:"TTL"`
	UserClientIp string           `position:"Query" name:"UserClientIp"`
}

// AddDomainRecordResponse is the response struct for api AddDomainRecord
type AddDomainRecordResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	RecordId  string `json:"RecordId" xml:"RecordId"`
}

// CreateAddDomainRecordRequest creates a request to invoke AddDomainRecord API
func CreateAddDomainRecordRequest() (request *AddDomainRecordRequest) {
	request = &AddDomainRecordRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "AddDomainRecord", "alidns", "openAPI")
	request.Method = requests.POST
	return
}

// CreateAddDomainRecordResponse creates a response to parse from AddDomainRecord response
func CreateAddDomainRecordResponse() (response *AddDomainRecordResponse) {
	response = &AddDomainRecordResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
