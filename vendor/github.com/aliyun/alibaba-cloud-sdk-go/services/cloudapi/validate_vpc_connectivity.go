package cloudapi

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

// ValidateVpcConnectivity invokes the cloudapi.ValidateVpcConnectivity API synchronously
func (client *Client) ValidateVpcConnectivity(request *ValidateVpcConnectivityRequest) (response *ValidateVpcConnectivityResponse, err error) {
	response = CreateValidateVpcConnectivityResponse()
	err = client.DoAction(request, response)
	return
}

// ValidateVpcConnectivityWithChan invokes the cloudapi.ValidateVpcConnectivity API asynchronously
func (client *Client) ValidateVpcConnectivityWithChan(request *ValidateVpcConnectivityRequest) (<-chan *ValidateVpcConnectivityResponse, <-chan error) {
	responseChan := make(chan *ValidateVpcConnectivityResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ValidateVpcConnectivity(request)
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

// ValidateVpcConnectivityWithCallback invokes the cloudapi.ValidateVpcConnectivity API asynchronously
func (client *Client) ValidateVpcConnectivityWithCallback(request *ValidateVpcConnectivityRequest, callback func(response *ValidateVpcConnectivityResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ValidateVpcConnectivityResponse
		var err error
		defer close(result)
		response, err = client.ValidateVpcConnectivity(request)
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

// ValidateVpcConnectivityRequest is the request struct for api ValidateVpcConnectivity
type ValidateVpcConnectivityRequest struct {
	*requests.RpcRequest
	InstanceId    string `position:"Query" name:"InstanceId"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
	VpcAccessId   string `position:"Query" name:"VpcAccessId"`
}

// ValidateVpcConnectivityResponse is the response struct for api ValidateVpcConnectivity
type ValidateVpcConnectivityResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Connected bool   `json:"Connected" xml:"Connected"`
	IpType    string `json:"IpType" xml:"IpType"`
}

// CreateValidateVpcConnectivityRequest creates a request to invoke ValidateVpcConnectivity API
func CreateValidateVpcConnectivityRequest() (request *ValidateVpcConnectivityRequest) {
	request = &ValidateVpcConnectivityRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "ValidateVpcConnectivity", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateValidateVpcConnectivityResponse creates a response to parse from ValidateVpcConnectivity response
func CreateValidateVpcConnectivityResponse() (response *ValidateVpcConnectivityResponse) {
	response = &ValidateVpcConnectivityResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
