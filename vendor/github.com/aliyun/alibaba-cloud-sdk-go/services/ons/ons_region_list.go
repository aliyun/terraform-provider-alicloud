package ons

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

// OnsRegionList invokes the ons.OnsRegionList API synchronously
// api document: https://help.aliyun.com/api/ons/onsregionlist.html
func (client *Client) OnsRegionList(request *OnsRegionListRequest) (response *OnsRegionListResponse, err error) {
	response = CreateOnsRegionListResponse()
	err = client.DoAction(request, response)
	return
}

// OnsRegionListWithChan invokes the ons.OnsRegionList API asynchronously
// api document: https://help.aliyun.com/api/ons/onsregionlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) OnsRegionListWithChan(request *OnsRegionListRequest) (<-chan *OnsRegionListResponse, <-chan error) {
	responseChan := make(chan *OnsRegionListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.OnsRegionList(request)
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

// OnsRegionListWithCallback invokes the ons.OnsRegionList API asynchronously
// api document: https://help.aliyun.com/api/ons/onsregionlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) OnsRegionListWithCallback(request *OnsRegionListRequest, callback func(response *OnsRegionListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *OnsRegionListResponse
		var err error
		defer close(result)
		response, err = client.OnsRegionList(request)
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

// OnsRegionListRequest is the request struct for api OnsRegionList
type OnsRegionListRequest struct {
	*requests.RpcRequest
}

// OnsRegionListResponse is the response struct for api OnsRegionList
type OnsRegionListResponse struct {
	*responses.BaseResponse
	RequestId string              `json:"RequestId" xml:"RequestId"`
	HelpUrl   string              `json:"HelpUrl" xml:"HelpUrl"`
	Data      DataInOnsRegionList `json:"Data" xml:"Data"`
}

// CreateOnsRegionListRequest creates a request to invoke OnsRegionList API
func CreateOnsRegionListRequest() (request *OnsRegionListRequest) {
	request = &OnsRegionListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ons", "2019-02-14", "OnsRegionList", "ons", "openAPI")
	request.Method = requests.POST
	return
}

// CreateOnsRegionListResponse creates a response to parse from OnsRegionList response
func CreateOnsRegionListResponse() (response *OnsRegionListResponse) {
	response = &OnsRegionListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
