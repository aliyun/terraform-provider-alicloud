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

// CreateInstantSiteMonitor invokes the cms.CreateInstantSiteMonitor API synchronously
func (client *Client) CreateInstantSiteMonitor(request *CreateInstantSiteMonitorRequest) (response *CreateInstantSiteMonitorResponse, err error) {
	response = CreateCreateInstantSiteMonitorResponse()
	err = client.DoAction(request, response)
	return
}

// CreateInstantSiteMonitorWithChan invokes the cms.CreateInstantSiteMonitor API asynchronously
func (client *Client) CreateInstantSiteMonitorWithChan(request *CreateInstantSiteMonitorRequest) (<-chan *CreateInstantSiteMonitorResponse, <-chan error) {
	responseChan := make(chan *CreateInstantSiteMonitorResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateInstantSiteMonitor(request)
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

// CreateInstantSiteMonitorWithCallback invokes the cms.CreateInstantSiteMonitor API asynchronously
func (client *Client) CreateInstantSiteMonitorWithCallback(request *CreateInstantSiteMonitorRequest, callback func(response *CreateInstantSiteMonitorResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateInstantSiteMonitorResponse
		var err error
		defer close(result)
		response, err = client.CreateInstantSiteMonitor(request)
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

// CreateInstantSiteMonitorRequest is the request struct for api CreateInstantSiteMonitor
type CreateInstantSiteMonitorRequest struct {
	*requests.RpcRequest
	RandomIspCity requests.Integer `position:"Query" name:"RandomIspCity"`
	Address       string           `position:"Query" name:"Address"`
	TaskType      string           `position:"Query" name:"TaskType"`
	TaskName      string           `position:"Query" name:"TaskName"`
	IspCities     string           `position:"Query" name:"IspCities"`
	OptionsJson   string           `position:"Query" name:"OptionsJson"`
}

// CreateInstantSiteMonitorResponse is the response struct for api CreateInstantSiteMonitor
type CreateInstantSiteMonitorResponse struct {
	*responses.BaseResponse
	Code             string                 `json:"Code" xml:"Code"`
	Message          string                 `json:"Message" xml:"Message"`
	RequestId        string                 `json:"RequestId" xml:"RequestId"`
	Success          string                 `json:"Success" xml:"Success"`
	CreateResultList []CreateResultListItem `json:"CreateResultList" xml:"CreateResultList"`
}

// CreateCreateInstantSiteMonitorRequest creates a request to invoke CreateInstantSiteMonitor API
func CreateCreateInstantSiteMonitorRequest() (request *CreateInstantSiteMonitorRequest) {
	request = &CreateInstantSiteMonitorRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "CreateInstantSiteMonitor", "Cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateInstantSiteMonitorResponse creates a response to parse from CreateInstantSiteMonitor response
func CreateCreateInstantSiteMonitorResponse() (response *CreateInstantSiteMonitorResponse) {
	response = &CreateInstantSiteMonitorResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
