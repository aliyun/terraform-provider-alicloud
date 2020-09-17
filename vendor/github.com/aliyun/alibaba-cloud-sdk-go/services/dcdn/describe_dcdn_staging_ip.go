package dcdn

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

// DescribeDcdnStagingIp invokes the dcdn.DescribeDcdnStagingIp API synchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdnstagingip.html
func (client *Client) DescribeDcdnStagingIp(request *DescribeDcdnStagingIpRequest) (response *DescribeDcdnStagingIpResponse, err error) {
	response = CreateDescribeDcdnStagingIpResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnStagingIpWithChan invokes the dcdn.DescribeDcdnStagingIp API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdnstagingip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnStagingIpWithChan(request *DescribeDcdnStagingIpRequest) (<-chan *DescribeDcdnStagingIpResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnStagingIpResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnStagingIp(request)
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

// DescribeDcdnStagingIpWithCallback invokes the dcdn.DescribeDcdnStagingIp API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdnstagingip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnStagingIpWithCallback(request *DescribeDcdnStagingIpRequest, callback func(response *DescribeDcdnStagingIpResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnStagingIpResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnStagingIp(request)
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

// DescribeDcdnStagingIpRequest is the request struct for api DescribeDcdnStagingIp
type DescribeDcdnStagingIpRequest struct {
	*requests.RpcRequest
	OwnerId requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDcdnStagingIpResponse is the response struct for api DescribeDcdnStagingIp
type DescribeDcdnStagingIpResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	IPV4s     IPV4s  `json:"IPV4s" xml:"IPV4s"`
}

// CreateDescribeDcdnStagingIpRequest creates a request to invoke DescribeDcdnStagingIp API
func CreateDescribeDcdnStagingIpRequest() (request *DescribeDcdnStagingIpRequest) {
	request = &DescribeDcdnStagingIpRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnStagingIp", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnStagingIpResponse creates a response to parse from DescribeDcdnStagingIp response
func CreateDescribeDcdnStagingIpResponse() (response *DescribeDcdnStagingIpResponse) {
	response = &DescribeDcdnStagingIpResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
