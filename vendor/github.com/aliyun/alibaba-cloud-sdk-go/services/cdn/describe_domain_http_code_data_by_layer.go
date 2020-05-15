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

// DescribeDomainHttpCodeDataByLayer invokes the cdn.DescribeDomainHttpCodeDataByLayer API synchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedatabylayer.html
func (client *Client) DescribeDomainHttpCodeDataByLayer(request *DescribeDomainHttpCodeDataByLayerRequest) (response *DescribeDomainHttpCodeDataByLayerResponse, err error) {
	response = CreateDescribeDomainHttpCodeDataByLayerResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainHttpCodeDataByLayerWithChan invokes the cdn.DescribeDomainHttpCodeDataByLayer API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedatabylayer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainHttpCodeDataByLayerWithChan(request *DescribeDomainHttpCodeDataByLayerRequest) (<-chan *DescribeDomainHttpCodeDataByLayerResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainHttpCodeDataByLayerResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainHttpCodeDataByLayer(request)
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

// DescribeDomainHttpCodeDataByLayerWithCallback invokes the cdn.DescribeDomainHttpCodeDataByLayer API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainhttpcodedatabylayer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainHttpCodeDataByLayerWithCallback(request *DescribeDomainHttpCodeDataByLayerRequest, callback func(response *DescribeDomainHttpCodeDataByLayerResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainHttpCodeDataByLayerResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainHttpCodeDataByLayer(request)
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

// DescribeDomainHttpCodeDataByLayerRequest is the request struct for api DescribeDomainHttpCodeDataByLayer
type DescribeDomainHttpCodeDataByLayerRequest struct {
	*requests.RpcRequest
	LocationNameEn string           `position:"Query" name:"LocationNameEn"`
	StartTime      string           `position:"Query" name:"StartTime"`
	IspNameEn      string           `position:"Query" name:"IspNameEn"`
	Layer          string           `position:"Query" name:"Layer"`
	DomainName     string           `position:"Query" name:"DomainName"`
	EndTime        string           `position:"Query" name:"EndTime"`
	OwnerId        requests.Integer `position:"Query" name:"OwnerId"`
	Interval       string           `position:"Query" name:"Interval"`
}

// DescribeDomainHttpCodeDataByLayerResponse is the response struct for api DescribeDomainHttpCodeDataByLayer
type DescribeDomainHttpCodeDataByLayerResponse struct {
	*responses.BaseResponse
	RequestId            string               `json:"RequestId" xml:"RequestId"`
	DataInterval         string               `json:"DataInterval" xml:"DataInterval"`
	HttpCodeDataInterval HttpCodeDataInterval `json:"HttpCodeDataInterval" xml:"HttpCodeDataInterval"`
}

// CreateDescribeDomainHttpCodeDataByLayerRequest creates a request to invoke DescribeDomainHttpCodeDataByLayer API
func CreateDescribeDomainHttpCodeDataByLayerRequest() (request *DescribeDomainHttpCodeDataByLayerRequest) {
	request = &DescribeDomainHttpCodeDataByLayerRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "DescribeDomainHttpCodeDataByLayer", "", "")
	return
}

// CreateDescribeDomainHttpCodeDataByLayerResponse creates a response to parse from DescribeDomainHttpCodeDataByLayer response
func CreateDescribeDomainHttpCodeDataByLayerResponse() (response *DescribeDomainHttpCodeDataByLayerResponse) {
	response = &DescribeDomainHttpCodeDataByLayerResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
