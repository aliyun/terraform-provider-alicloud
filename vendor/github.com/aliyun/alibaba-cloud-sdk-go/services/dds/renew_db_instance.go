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

// RenewDBInstance invokes the dds.RenewDBInstance API synchronously
// api document: https://help.aliyun.com/api/dds/renewdbinstance.html
func (client *Client) RenewDBInstance(request *RenewDBInstanceRequest) (response *RenewDBInstanceResponse, err error) {
	response = CreateRenewDBInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// RenewDBInstanceWithChan invokes the dds.RenewDBInstance API asynchronously
// api document: https://help.aliyun.com/api/dds/renewdbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RenewDBInstanceWithChan(request *RenewDBInstanceRequest) (<-chan *RenewDBInstanceResponse, <-chan error) {
	responseChan := make(chan *RenewDBInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RenewDBInstance(request)
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

// RenewDBInstanceWithCallback invokes the dds.RenewDBInstance API asynchronously
// api document: https://help.aliyun.com/api/dds/renewdbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RenewDBInstanceWithCallback(request *RenewDBInstanceRequest, callback func(response *RenewDBInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RenewDBInstanceResponse
		var err error
		defer close(result)
		response, err = client.RenewDBInstance(request)
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

// RenewDBInstanceRequest is the request struct for api RenewDBInstance
type RenewDBInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	CouponNo             string           `position:"Query" name:"CouponNo"`
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	BusinessInfo         string           `position:"Query" name:"BusinessInfo"`
	Period               requests.Integer `position:"Query" name:"Period"`
	AutoPay              requests.Boolean `position:"Query" name:"AutoPay"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// RenewDBInstanceResponse is the response struct for api RenewDBInstance
type RenewDBInstanceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	OrderId   string `json:"OrderId" xml:"OrderId"`
}

// CreateRenewDBInstanceRequest creates a request to invoke RenewDBInstance API
func CreateRenewDBInstanceRequest() (request *RenewDBInstanceRequest) {
	request = &RenewDBInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "RenewDBInstance", "Dds", "openAPI")
	return
}

// CreateRenewDBInstanceResponse creates a response to parse from RenewDBInstance response
func CreateRenewDBInstanceResponse() (response *RenewDBInstanceResponse) {
	response = &RenewDBInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
