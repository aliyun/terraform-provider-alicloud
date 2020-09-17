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

// DescribeDcdnDomainOriginTrafficData invokes the dcdn.DescribeDcdnDomainOriginTrafficData API synchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomainorigintrafficdata.html
func (client *Client) DescribeDcdnDomainOriginTrafficData(request *DescribeDcdnDomainOriginTrafficDataRequest) (response *DescribeDcdnDomainOriginTrafficDataResponse, err error) {
	response = CreateDescribeDcdnDomainOriginTrafficDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainOriginTrafficDataWithChan invokes the dcdn.DescribeDcdnDomainOriginTrafficData API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomainorigintrafficdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnDomainOriginTrafficDataWithChan(request *DescribeDcdnDomainOriginTrafficDataRequest) (<-chan *DescribeDcdnDomainOriginTrafficDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainOriginTrafficDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainOriginTrafficData(request)
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

// DescribeDcdnDomainOriginTrafficDataWithCallback invokes the dcdn.DescribeDcdnDomainOriginTrafficData API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomainorigintrafficdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnDomainOriginTrafficDataWithCallback(request *DescribeDcdnDomainOriginTrafficDataRequest, callback func(response *DescribeDcdnDomainOriginTrafficDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainOriginTrafficDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainOriginTrafficData(request)
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

// DescribeDcdnDomainOriginTrafficDataRequest is the request struct for api DescribeDcdnDomainOriginTrafficData
type DescribeDcdnDomainOriginTrafficDataRequest struct {
	*requests.RpcRequest
	StartTime  string           `position:"Query" name:"StartTime"`
	DomainName string           `position:"Query" name:"DomainName"`
	EndTime    string           `position:"Query" name:"EndTime"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
	Interval   string           `position:"Query" name:"Interval"`
}

// DescribeDcdnDomainOriginTrafficDataResponse is the response struct for api DescribeDcdnDomainOriginTrafficData
type DescribeDcdnDomainOriginTrafficDataResponse struct {
	*responses.BaseResponse
	RequestId                    string                       `json:"RequestId" xml:"RequestId"`
	DomainName                   string                       `json:"DomainName" xml:"DomainName"`
	StartTime                    string                       `json:"StartTime" xml:"StartTime"`
	EndTime                      string                       `json:"EndTime" xml:"EndTime"`
	DataInterval                 string                       `json:"DataInterval" xml:"DataInterval"`
	OriginTrafficDataPerInterval OriginTrafficDataPerInterval `json:"OriginTrafficDataPerInterval" xml:"OriginTrafficDataPerInterval"`
}

// CreateDescribeDcdnDomainOriginTrafficDataRequest creates a request to invoke DescribeDcdnDomainOriginTrafficData API
func CreateDescribeDcdnDomainOriginTrafficDataRequest() (request *DescribeDcdnDomainOriginTrafficDataRequest) {
	request = &DescribeDcdnDomainOriginTrafficDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainOriginTrafficData", "", "")
	return
}

// CreateDescribeDcdnDomainOriginTrafficDataResponse creates a response to parse from DescribeDcdnDomainOriginTrafficData response
func CreateDescribeDcdnDomainOriginTrafficDataResponse() (response *DescribeDcdnDomainOriginTrafficDataResponse) {
	response = &DescribeDcdnDomainOriginTrafficDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
