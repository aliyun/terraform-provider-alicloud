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

// TransformInstanceChargeType invokes the dds.TransformInstanceChargeType API synchronously
func (client *Client) TransformInstanceChargeType(request *TransformInstanceChargeTypeRequest) (response *TransformInstanceChargeTypeResponse, err error) {
	response = CreateTransformInstanceChargeTypeResponse()
	err = client.DoAction(request, response)
	return
}

// TransformInstanceChargeTypeWithChan invokes the dds.TransformInstanceChargeType API asynchronously
func (client *Client) TransformInstanceChargeTypeWithChan(request *TransformInstanceChargeTypeRequest) (<-chan *TransformInstanceChargeTypeResponse, <-chan error) {
	responseChan := make(chan *TransformInstanceChargeTypeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.TransformInstanceChargeType(request)
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

// TransformInstanceChargeTypeWithCallback invokes the dds.TransformInstanceChargeType API asynchronously
func (client *Client) TransformInstanceChargeTypeWithCallback(request *TransformInstanceChargeTypeRequest, callback func(response *TransformInstanceChargeTypeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *TransformInstanceChargeTypeResponse
		var err error
		defer close(result)
		response, err = client.TransformInstanceChargeType(request)
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

// TransformInstanceChargeTypeRequest is the request struct for api TransformInstanceChargeType
type TransformInstanceChargeTypeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	CouponNo             string           `position:"Query" name:"CouponNo"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	BusinessInfo         string           `position:"Query" name:"BusinessInfo"`
	Period               requests.Integer `position:"Query" name:"Period"`
	AutoPay              requests.Boolean `position:"Query" name:"AutoPay"`
	FromApp              string           `position:"Query" name:"FromApp"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	InstanceId           string           `position:"Query" name:"InstanceId"`
	AutoRenew            string           `position:"Query" name:"AutoRenew"`
	ChargeType           string           `position:"Query" name:"ChargeType"`
	PricingCycle         string           `position:"Query" name:"PricingCycle"`
}

// TransformInstanceChargeTypeResponse is the response struct for api TransformInstanceChargeType
type TransformInstanceChargeTypeResponse struct {
	*responses.BaseResponse
	EndTime   string `json:"EndTime" xml:"EndTime"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	OrderId   string `json:"OrderId" xml:"OrderId"`
}

// CreateTransformInstanceChargeTypeRequest creates a request to invoke TransformInstanceChargeType API
func CreateTransformInstanceChargeTypeRequest() (request *TransformInstanceChargeTypeRequest) {
	request = &TransformInstanceChargeTypeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "TransformInstanceChargeType", "dds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateTransformInstanceChargeTypeResponse creates a response to parse from TransformInstanceChargeType response
func CreateTransformInstanceChargeTypeResponse() (response *TransformInstanceChargeTypeResponse) {
	response = &TransformInstanceChargeTypeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
