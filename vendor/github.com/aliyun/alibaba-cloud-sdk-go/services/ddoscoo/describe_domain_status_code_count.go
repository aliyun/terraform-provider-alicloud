package ddoscoo

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

// DescribeDomainStatusCodeCount invokes the ddoscoo.DescribeDomainStatusCodeCount API synchronously
// api document: https://help.aliyun.com/api/ddoscoo/describedomainstatuscodecount.html
func (client *Client) DescribeDomainStatusCodeCount(request *DescribeDomainStatusCodeCountRequest) (response *DescribeDomainStatusCodeCountResponse, err error) {
	response = CreateDescribeDomainStatusCodeCountResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainStatusCodeCountWithChan invokes the ddoscoo.DescribeDomainStatusCodeCount API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describedomainstatuscodecount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainStatusCodeCountWithChan(request *DescribeDomainStatusCodeCountRequest) (<-chan *DescribeDomainStatusCodeCountResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainStatusCodeCountResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainStatusCodeCount(request)
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

// DescribeDomainStatusCodeCountWithCallback invokes the ddoscoo.DescribeDomainStatusCodeCount API asynchronously
// api document: https://help.aliyun.com/api/ddoscoo/describedomainstatuscodecount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainStatusCodeCountWithCallback(request *DescribeDomainStatusCodeCountRequest, callback func(response *DescribeDomainStatusCodeCountResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainStatusCodeCountResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainStatusCodeCount(request)
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

// DescribeDomainStatusCodeCountRequest is the request struct for api DescribeDomainStatusCodeCount
type DescribeDomainStatusCodeCountRequest struct {
	*requests.RpcRequest
	EndTime         requests.Integer `position:"Query" name:"EndTime"`
	StartTime       requests.Integer `position:"Query" name:"StartTime"`
	ResourceGroupId string           `position:"Query" name:"ResourceGroupId"`
	SourceIp        string           `position:"Query" name:"SourceIp"`
	Domain          string           `position:"Query" name:"Domain"`
}

// DescribeDomainStatusCodeCountResponse is the response struct for api DescribeDomainStatusCodeCount
type DescribeDomainStatusCodeCountResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Status2XX int64  `json:"Status2XX" xml:"Status2XX"`
	Status501 int64  `json:"Status501" xml:"Status501"`
	Status502 int64  `json:"Status502" xml:"Status502"`
	Status503 int64  `json:"Status503" xml:"Status503"`
	Status504 int64  `json:"Status504" xml:"Status504"`
	Status200 int64  `json:"Status200" xml:"Status200"`
	Status405 int64  `json:"Status405" xml:"Status405"`
	Status5XX int64  `json:"Status5XX" xml:"Status5XX"`
	Status4XX int64  `json:"Status4XX" xml:"Status4XX"`
	Status403 int64  `json:"Status403" xml:"Status403"`
	Status404 int64  `json:"Status404" xml:"Status404"`
	Status3XX int64  `json:"Status3XX" xml:"Status3XX"`
}

// CreateDescribeDomainStatusCodeCountRequest creates a request to invoke DescribeDomainStatusCodeCount API
func CreateDescribeDomainStatusCodeCountRequest() (request *DescribeDomainStatusCodeCountRequest) {
	request = &DescribeDomainStatusCodeCountRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ddoscoo", "2020-01-01", "DescribeDomainStatusCodeCount", "ddoscoo", "openAPI")
	return
}

// CreateDescribeDomainStatusCodeCountResponse creates a response to parse from DescribeDomainStatusCodeCount response
func CreateDescribeDomainStatusCodeCountResponse() (response *DescribeDomainStatusCodeCountResponse) {
	response = &DescribeDomainStatusCodeCountResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
