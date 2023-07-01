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

// DescribeDiagnosisSQLInfo invokes the gpdb.DescribeDiagnosisSQLInfo API synchronously
func (client *Client) DescribeDiagnosisSQLInfo(request *DescribeDiagnosisSQLInfoRequest) (response *DescribeDiagnosisSQLInfoResponse, err error) {
	response = CreateDescribeDiagnosisSQLInfoResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDiagnosisSQLInfoWithChan invokes the gpdb.DescribeDiagnosisSQLInfo API asynchronously
func (client *Client) DescribeDiagnosisSQLInfoWithChan(request *DescribeDiagnosisSQLInfoRequest) (<-chan *DescribeDiagnosisSQLInfoResponse, <-chan error) {
	responseChan := make(chan *DescribeDiagnosisSQLInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDiagnosisSQLInfo(request)
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

// DescribeDiagnosisSQLInfoWithCallback invokes the gpdb.DescribeDiagnosisSQLInfo API asynchronously
func (client *Client) DescribeDiagnosisSQLInfoWithCallback(request *DescribeDiagnosisSQLInfoRequest, callback func(response *DescribeDiagnosisSQLInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDiagnosisSQLInfoResponse
		var err error
		defer close(result)
		response, err = client.DescribeDiagnosisSQLInfo(request)
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

// DescribeDiagnosisSQLInfoRequest is the request struct for api DescribeDiagnosisSQLInfo
type DescribeDiagnosisSQLInfoRequest struct {
	*requests.RpcRequest
	Database     string `position:"Query" name:"Database"`
	DBInstanceId string `position:"Query" name:"DBInstanceId"`
	QueryID      string `position:"Query" name:"QueryID"`
}

// DescribeDiagnosisSQLInfoResponse is the response struct for api DescribeDiagnosisSQLInfo
type DescribeDiagnosisSQLInfoResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	QueryID       string `json:"QueryID" xml:"QueryID"`
	SessionID     string `json:"SessionID" xml:"SessionID"`
	StartTime     int64  `json:"StartTime" xml:"StartTime"`
	Duration      int    `json:"Duration" xml:"Duration"`
	SQLStmt       string `json:"SQLStmt" xml:"SQLStmt"`
	User          string `json:"User" xml:"User"`
	Database      string `json:"Database" xml:"Database"`
	Status        string `json:"Status" xml:"Status"`
	QueryPlan     string `json:"QueryPlan" xml:"QueryPlan"`
	TextPlan      string `json:"TextPlan" xml:"TextPlan"`
	SortedMetrics string `json:"SortedMetrics" xml:"SortedMetrics"`
	MaxOutputRows string `json:"MaxOutputRows" xml:"MaxOutputRows"`
}

// CreateDescribeDiagnosisSQLInfoRequest creates a request to invoke DescribeDiagnosisSQLInfo API
func CreateDescribeDiagnosisSQLInfoRequest() (request *DescribeDiagnosisSQLInfoRequest) {
	request = &DescribeDiagnosisSQLInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "DescribeDiagnosisSQLInfo", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDiagnosisSQLInfoResponse creates a response to parse from DescribeDiagnosisSQLInfo response
func CreateDescribeDiagnosisSQLInfoResponse() (response *DescribeDiagnosisSQLInfoResponse) {
	response = &DescribeDiagnosisSQLInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
