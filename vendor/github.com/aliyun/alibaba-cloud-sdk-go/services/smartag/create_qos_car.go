package smartag

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

// CreateQosCar invokes the smartag.CreateQosCar API synchronously
// api document: https://help.aliyun.com/api/smartag/createqoscar.html
func (client *Client) CreateQosCar(request *CreateQosCarRequest) (response *CreateQosCarResponse, err error) {
	response = CreateCreateQosCarResponse()
	err = client.DoAction(request, response)
	return
}

// CreateQosCarWithChan invokes the smartag.CreateQosCar API asynchronously
// api document: https://help.aliyun.com/api/smartag/createqoscar.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateQosCarWithChan(request *CreateQosCarRequest) (<-chan *CreateQosCarResponse, <-chan error) {
	responseChan := make(chan *CreateQosCarResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateQosCar(request)
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

// CreateQosCarWithCallback invokes the smartag.CreateQosCar API asynchronously
// api document: https://help.aliyun.com/api/smartag/createqoscar.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateQosCarWithCallback(request *CreateQosCarRequest, callback func(response *CreateQosCarResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateQosCarResponse
		var err error
		defer close(result)
		response, err = client.CreateQosCar(request)
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

// CreateQosCarRequest is the request struct for api CreateQosCar
type CreateQosCarRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	MinBandwidthAbs      requests.Integer `position:"Query" name:"MinBandwidthAbs"`
	Description          string           `position:"Query" name:"Description"`
	PercentSourceType    string           `position:"Query" name:"PercentSourceType"`
	QosId                string           `position:"Query" name:"QosId"`
	MaxBandwidthAbs      requests.Integer `position:"Query" name:"MaxBandwidthAbs"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	MaxBandwidthPercent  requests.Integer `position:"Query" name:"MaxBandwidthPercent"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Priority             requests.Integer `position:"Query" name:"Priority"`
	MinBandwidthPercent  requests.Integer `position:"Query" name:"MinBandwidthPercent"`
	LimitType            string           `position:"Query" name:"LimitType"`
	Name                 string           `position:"Query" name:"Name"`
}

// CreateQosCarResponse is the response struct for api CreateQosCar
type CreateQosCarResponse struct {
	*responses.BaseResponse
	RequestId           string `json:"RequestId" xml:"RequestId"`
	QosId               string `json:"QosId" xml:"QosId"`
	QosCarId            string `json:"QosCarId" xml:"QosCarId"`
	Description         string `json:"Description" xml:"Description"`
	Priority            int    `json:"Priority" xml:"Priority"`
	LimitType           string `json:"LimitType" xml:"LimitType"`
	MinBandwidthAbs     int    `json:"MinBandwidthAbs" xml:"MinBandwidthAbs"`
	MaxBandwidthAbs     int    `json:"MaxBandwidthAbs" xml:"MaxBandwidthAbs"`
	MinBandwidthPercent int    `json:"MinBandwidthPercent" xml:"MinBandwidthPercent"`
	MaxBandwidthPercent int    `json:"MaxBandwidthPercent" xml:"MaxBandwidthPercent"`
	PercentSourceType   string `json:"PercentSourceType" xml:"PercentSourceType"`
}

// CreateCreateQosCarRequest creates a request to invoke CreateQosCar API
func CreateCreateQosCarRequest() (request *CreateQosCarRequest) {
	request = &CreateQosCarRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "CreateQosCar", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateQosCarResponse creates a response to parse from CreateQosCar response
func CreateCreateQosCarResponse() (response *CreateQosCarResponse) {
	response = &CreateQosCarResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
