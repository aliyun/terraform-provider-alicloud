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

// DescribeSagHa invokes the smartag.DescribeSagHa API synchronously
func (client *Client) DescribeSagHa(request *DescribeSagHaRequest) (response *DescribeSagHaResponse, err error) {
	response = CreateDescribeSagHaResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSagHaWithChan invokes the smartag.DescribeSagHa API asynchronously
func (client *Client) DescribeSagHaWithChan(request *DescribeSagHaRequest) (<-chan *DescribeSagHaResponse, <-chan error) {
	responseChan := make(chan *DescribeSagHaResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSagHa(request)
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

// DescribeSagHaWithCallback invokes the smartag.DescribeSagHa API asynchronously
func (client *Client) DescribeSagHaWithCallback(request *DescribeSagHaRequest, callback func(response *DescribeSagHaResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSagHaResponse
		var err error
		defer close(result)
		response, err = client.DescribeSagHa(request)
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

// DescribeSagHaRequest is the request struct for api DescribeSagHa
type DescribeSagHaRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	SmartAGId            string           `position:"Query" name:"SmartAGId"`
	SmartAGSn            string           `position:"Query" name:"SmartAGSn"`
}

// DescribeSagHaResponse is the response struct for api DescribeSagHa
type DescribeSagHaResponse struct {
	*responses.BaseResponse
	RequestId  string      `json:"RequestId" xml:"RequestId"`
	Mode       string      `json:"Mode" xml:"Mode"`
	Ports      []Port      `json:"Ports" xml:"Ports"`
	TaskStates []TaskState `json:"TaskStates" xml:"TaskStates"`
}

// CreateDescribeSagHaRequest creates a request to invoke DescribeSagHa API
func CreateDescribeSagHaRequest() (request *DescribeSagHaRequest) {
	request = &DescribeSagHaRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "DescribeSagHa", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeSagHaResponse creates a response to parse from DescribeSagHa response
func CreateDescribeSagHaResponse() (response *DescribeSagHaResponse) {
	response = &DescribeSagHaResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
