package vpc

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

// DeleteCommonBandwidthPackage invokes the vpc.DeleteCommonBandwidthPackage API synchronously
// api document: https://help.aliyun.com/api/vpc/deletecommonbandwidthpackage.html
func (client *Client) DeleteCommonBandwidthPackage(request *DeleteCommonBandwidthPackageRequest) (response *DeleteCommonBandwidthPackageResponse, err error) {
	response = CreateDeleteCommonBandwidthPackageResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteCommonBandwidthPackageWithChan invokes the vpc.DeleteCommonBandwidthPackage API asynchronously
// api document: https://help.aliyun.com/api/vpc/deletecommonbandwidthpackage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteCommonBandwidthPackageWithChan(request *DeleteCommonBandwidthPackageRequest) (<-chan *DeleteCommonBandwidthPackageResponse, <-chan error) {
	responseChan := make(chan *DeleteCommonBandwidthPackageResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteCommonBandwidthPackage(request)
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

// DeleteCommonBandwidthPackageWithCallback invokes the vpc.DeleteCommonBandwidthPackage API asynchronously
// api document: https://help.aliyun.com/api/vpc/deletecommonbandwidthpackage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteCommonBandwidthPackageWithCallback(request *DeleteCommonBandwidthPackageRequest, callback func(response *DeleteCommonBandwidthPackageResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteCommonBandwidthPackageResponse
		var err error
		defer close(result)
		response, err = client.DeleteCommonBandwidthPackage(request)
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

// DeleteCommonBandwidthPackageRequest is the request struct for api DeleteCommonBandwidthPackage
type DeleteCommonBandwidthPackageRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	BandwidthPackageId   string           `position:"Query" name:"BandwidthPackageId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Force                string           `position:"Query" name:"Force"`
}

// DeleteCommonBandwidthPackageResponse is the response struct for api DeleteCommonBandwidthPackage
type DeleteCommonBandwidthPackageResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteCommonBandwidthPackageRequest creates a request to invoke DeleteCommonBandwidthPackage API
func CreateDeleteCommonBandwidthPackageRequest() (request *DeleteCommonBandwidthPackageRequest) {
	request = &DeleteCommonBandwidthPackageRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DeleteCommonBandwidthPackage", "Vpc", "openAPI")
	return
}

// CreateDeleteCommonBandwidthPackageResponse creates a response to parse from DeleteCommonBandwidthPackage response
func CreateDeleteCommonBandwidthPackageResponse() (response *DeleteCommonBandwidthPackageResponse) {
	response = &DeleteCommonBandwidthPackageResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
