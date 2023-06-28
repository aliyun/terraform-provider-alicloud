package gpdb

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

// DescribeSampleData invokes the gpdb.DescribeSampleData API synchronously
func (client *Client) DescribeSampleData(request *DescribeSampleDataRequest) (response *DescribeSampleDataResponse, err error) {
	response = CreateDescribeSampleDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSampleDataWithChan invokes the gpdb.DescribeSampleData API asynchronously
func (client *Client) DescribeSampleDataWithChan(request *DescribeSampleDataRequest) (<-chan *DescribeSampleDataResponse, <-chan error) {
	responseChan := make(chan *DescribeSampleDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSampleData(request)
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

// DescribeSampleDataWithCallback invokes the gpdb.DescribeSampleData API asynchronously
func (client *Client) DescribeSampleDataWithCallback(request *DescribeSampleDataRequest, callback func(response *DescribeSampleDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSampleDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeSampleData(request)
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

// DescribeSampleDataRequest is the request struct for api DescribeSampleData
type DescribeSampleDataRequest struct {
	*requests.RpcRequest
	DBInstanceId string           `position:"Query" name:"DBInstanceId"`
	OwnerId      requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeSampleDataResponse is the response struct for api DescribeSampleData
type DescribeSampleDataResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	HasSampleData bool   `json:"HasSampleData" xml:"HasSampleData"`
	ErrorMessage  string `json:"ErrorMessage" xml:"ErrorMessage"`
	DBInstanceId  string `json:"DBInstanceId" xml:"DBInstanceId"`
}

// CreateDescribeSampleDataRequest creates a request to invoke DescribeSampleData API
func CreateDescribeSampleDataRequest() (request *DescribeSampleDataRequest) {
	request = &DescribeSampleDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "DescribeSampleData", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeSampleDataResponse creates a response to parse from DescribeSampleData response
func CreateDescribeSampleDataResponse() (response *DescribeSampleDataResponse) {
	response = &DescribeSampleDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}