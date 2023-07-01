package cms

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

// DescribeAlertLogHistogram invokes the cms.DescribeAlertLogHistogram API synchronously
func (client *Client) DescribeAlertLogHistogram(request *DescribeAlertLogHistogramRequest) (response *DescribeAlertLogHistogramResponse, err error) {
	response = CreateDescribeAlertLogHistogramResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAlertLogHistogramWithChan invokes the cms.DescribeAlertLogHistogram API asynchronously
func (client *Client) DescribeAlertLogHistogramWithChan(request *DescribeAlertLogHistogramRequest) (<-chan *DescribeAlertLogHistogramResponse, <-chan error) {
	responseChan := make(chan *DescribeAlertLogHistogramResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAlertLogHistogram(request)
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

// DescribeAlertLogHistogramWithCallback invokes the cms.DescribeAlertLogHistogram API asynchronously
func (client *Client) DescribeAlertLogHistogramWithCallback(request *DescribeAlertLogHistogramRequest, callback func(response *DescribeAlertLogHistogramResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAlertLogHistogramResponse
		var err error
		defer close(result)
		response, err = client.DescribeAlertLogHistogram(request)
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

// DescribeAlertLogHistogramRequest is the request struct for api DescribeAlertLogHistogram
type DescribeAlertLogHistogramRequest struct {
	*requests.RpcRequest
	SendStatus   string           `position:"Query" name:"SendStatus"`
	ContactGroup string           `position:"Query" name:"ContactGroup"`
	SearchKey    string           `position:"Query" name:"SearchKey"`
	RuleName     string           `position:"Query" name:"RuleName"`
	StartTime    requests.Integer `position:"Query" name:"StartTime"`
	PageNumber   requests.Integer `position:"Query" name:"PageNumber"`
	LastMin      string           `position:"Query" name:"LastMin"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
	MetricName   string           `position:"Query" name:"MetricName"`
	Product      string           `position:"Query" name:"Product"`
	Level        string           `position:"Query" name:"Level"`
	GroupId      string           `position:"Query" name:"GroupId"`
	EndTime      requests.Integer `position:"Query" name:"EndTime"`
	GroupBy      string           `position:"Query" name:"GroupBy"`
	Namespace    string           `position:"Query" name:"Namespace"`
}

// DescribeAlertLogHistogramResponse is the response struct for api DescribeAlertLogHistogram
type DescribeAlertLogHistogramResponse struct {
	*responses.BaseResponse
	Code                  string                      `json:"Code" xml:"Code"`
	Message               string                      `json:"Message" xml:"Message"`
	RequestId             string                      `json:"RequestId" xml:"RequestId"`
	Success               bool                        `json:"Success" xml:"Success"`
	AlertLogHistogramList []AlertLogHistogramListItem `json:"AlertLogHistogramList" xml:"AlertLogHistogramList"`
}

// CreateDescribeAlertLogHistogramRequest creates a request to invoke DescribeAlertLogHistogram API
func CreateDescribeAlertLogHistogramRequest() (request *DescribeAlertLogHistogramRequest) {
	request = &DescribeAlertLogHistogramRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "DescribeAlertLogHistogram", "cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAlertLogHistogramResponse creates a response to parse from DescribeAlertLogHistogram response
func CreateDescribeAlertLogHistogramResponse() (response *DescribeAlertLogHistogramResponse) {
	response = &DescribeAlertLogHistogramResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
