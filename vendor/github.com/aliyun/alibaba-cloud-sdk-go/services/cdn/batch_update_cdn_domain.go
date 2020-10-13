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

// BatchUpdateCdnDomain invokes the cdn.BatchUpdateCdnDomain API synchronously
func (client *Client) BatchUpdateCdnDomain(request *BatchUpdateCdnDomainRequest) (response *BatchUpdateCdnDomainResponse, err error) {
	response = CreateBatchUpdateCdnDomainResponse()
	err = client.DoAction(request, response)
	return
}

// BatchUpdateCdnDomainWithChan invokes the cdn.BatchUpdateCdnDomain API asynchronously
func (client *Client) BatchUpdateCdnDomainWithChan(request *BatchUpdateCdnDomainRequest) (<-chan *BatchUpdateCdnDomainResponse, <-chan error) {
	responseChan := make(chan *BatchUpdateCdnDomainResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BatchUpdateCdnDomain(request)
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

// BatchUpdateCdnDomainWithCallback invokes the cdn.BatchUpdateCdnDomain API asynchronously
func (client *Client) BatchUpdateCdnDomainWithCallback(request *BatchUpdateCdnDomainRequest, callback func(response *BatchUpdateCdnDomainResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BatchUpdateCdnDomainResponse
		var err error
		defer close(result)
		response, err = client.BatchUpdateCdnDomain(request)
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

// BatchUpdateCdnDomainRequest is the request struct for api BatchUpdateCdnDomain
type BatchUpdateCdnDomainRequest struct {
	*requests.RpcRequest
	Sources         string           `position:"Query" name:"Sources"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SecurityToken   string           `position:"Query" name:"SecurityToken"`
	TopLevelDomain  string           `position:"Query" name:"TopLevelDomain"`
	DomainName      string           `position:"Query" name:"DomainName"`
	OwnerId         requests.Integer `position:"Query" name:"OwnerId"`
}

// BatchUpdateCdnDomainResponse is the response struct for api BatchUpdateCdnDomain
type BatchUpdateCdnDomainResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateBatchUpdateCdnDomainRequest creates a request to invoke BatchUpdateCdnDomain API
func CreateBatchUpdateCdnDomainRequest() (request *BatchUpdateCdnDomainRequest) {
	request = &BatchUpdateCdnDomainRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "BatchUpdateCdnDomain", "", "")
	request.Method = requests.POST
	return
}

// CreateBatchUpdateCdnDomainResponse creates a response to parse from BatchUpdateCdnDomain response
func CreateBatchUpdateCdnDomainResponse() (response *BatchUpdateCdnDomainResponse) {
	response = &BatchUpdateCdnDomainResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
