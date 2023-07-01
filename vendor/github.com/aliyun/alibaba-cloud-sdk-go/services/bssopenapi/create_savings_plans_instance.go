package bssopenapi

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

// CreateSavingsPlansInstance invokes the bssopenapi.CreateSavingsPlansInstance API synchronously
func (client *Client) CreateSavingsPlansInstance(request *CreateSavingsPlansInstanceRequest) (response *CreateSavingsPlansInstanceResponse, err error) {
	response = CreateCreateSavingsPlansInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// CreateSavingsPlansInstanceWithChan invokes the bssopenapi.CreateSavingsPlansInstance API asynchronously
func (client *Client) CreateSavingsPlansInstanceWithChan(request *CreateSavingsPlansInstanceRequest) (<-chan *CreateSavingsPlansInstanceResponse, <-chan error) {
	responseChan := make(chan *CreateSavingsPlansInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateSavingsPlansInstance(request)
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

// CreateSavingsPlansInstanceWithCallback invokes the bssopenapi.CreateSavingsPlansInstance API asynchronously
func (client *Client) CreateSavingsPlansInstanceWithCallback(request *CreateSavingsPlansInstanceRequest, callback func(response *CreateSavingsPlansInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateSavingsPlansInstanceResponse
		var err error
		defer close(result)
		response, err = client.CreateSavingsPlansInstance(request)
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

// CreateSavingsPlansInstanceRequest is the request struct for api CreateSavingsPlansInstance
type CreateSavingsPlansInstanceRequest struct {
	*requests.RpcRequest
	Specification string            `position:"Query" name:"Specification"`
	PoolValue     string            `position:"Query" name:"PoolValue"`
	CommodityCode string            `position:"Query" name:"CommodityCode"`
	Type          string            `position:"Query" name:"Type"`
	EffectiveDate string            `position:"Query" name:"EffectiveDate"`
	Duration      string            `position:"Query" name:"Duration"`
	SpecType      string            `position:"Query" name:"SpecType"`
	ExtendMap     map[string]string `position:"Query" name:"ExtendMap"  type:"Map"`
	PayMode       string            `position:"Query" name:"PayMode"`
	Region        string            `position:"Query" name:"Region"`
	PricingCycle  string            `position:"Query" name:"PricingCycle"`
}

// CreateSavingsPlansInstanceResponse is the response struct for api CreateSavingsPlansInstance
type CreateSavingsPlansInstanceResponse struct {
	*responses.BaseResponse
	Message   string                           `json:"Message" xml:"Message"`
	RequestId string                           `json:"RequestId" xml:"RequestId"`
	Code      string                           `json:"Code" xml:"Code"`
	Success   bool                             `json:"Success" xml:"Success"`
	Data      DataInCreateSavingsPlansInstance `json:"Data" xml:"Data"`
}

// CreateCreateSavingsPlansInstanceRequest creates a request to invoke CreateSavingsPlansInstance API
func CreateCreateSavingsPlansInstanceRequest() (request *CreateSavingsPlansInstanceRequest) {
	request = &CreateSavingsPlansInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "CreateSavingsPlansInstance", "bssopenapi", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateSavingsPlansInstanceResponse creates a response to parse from CreateSavingsPlansInstance response
func CreateCreateSavingsPlansInstanceResponse() (response *CreateSavingsPlansInstanceResponse) {
	response = &CreateSavingsPlansInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
