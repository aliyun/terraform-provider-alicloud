package edas

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

// UpdateAccountInfo invokes the edas.UpdateAccountInfo API synchronously
func (client *Client) UpdateAccountInfo(request *UpdateAccountInfoRequest) (response *UpdateAccountInfoResponse, err error) {
	response = CreateUpdateAccountInfoResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateAccountInfoWithChan invokes the edas.UpdateAccountInfo API asynchronously
func (client *Client) UpdateAccountInfoWithChan(request *UpdateAccountInfoRequest) (<-chan *UpdateAccountInfoResponse, <-chan error) {
	responseChan := make(chan *UpdateAccountInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateAccountInfo(request)
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

// UpdateAccountInfoWithCallback invokes the edas.UpdateAccountInfo API asynchronously
func (client *Client) UpdateAccountInfoWithCallback(request *UpdateAccountInfoRequest, callback func(response *UpdateAccountInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateAccountInfoResponse
		var err error
		defer close(result)
		response, err = client.UpdateAccountInfo(request)
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

// UpdateAccountInfoRequest is the request struct for api UpdateAccountInfo
type UpdateAccountInfoRequest struct {
	*requests.RoaRequest
	Name      string `position:"Query" name:"Name"`
	Telephone string `position:"Query" name:"Telephone"`
	Email     string `position:"Query" name:"Email"`
}

// UpdateAccountInfoResponse is the response struct for api UpdateAccountInfo
type UpdateAccountInfoResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateAccountInfoRequest creates a request to invoke UpdateAccountInfo API
func CreateUpdateAccountInfoRequest() (request *UpdateAccountInfoRequest) {
	request = &UpdateAccountInfoRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "UpdateAccountInfo", "/pop/v5/account/edit_account_info", "Edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUpdateAccountInfoResponse creates a response to parse from UpdateAccountInfo response
func CreateUpdateAccountInfoResponse() (response *UpdateAccountInfoResponse) {
	response = &UpdateAccountInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
