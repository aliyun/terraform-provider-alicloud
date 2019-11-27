package elasticsearch

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

// ConvertPayType invokes the elasticsearch.ConvertPayType API synchronously
// api document: https://help.aliyun.com/api/elasticsearch/convertpaytype.html
func (client *Client) ConvertPayType(request *ConvertPayTypeRequest) (response *ConvertPayTypeResponse, err error) {
	response = CreateConvertPayTypeResponse()
	err = client.DoAction(request, response)
	return
}

// ConvertPayTypeWithChan invokes the elasticsearch.ConvertPayType API asynchronously
// api document: https://help.aliyun.com/api/elasticsearch/convertpaytype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ConvertPayTypeWithChan(request *ConvertPayTypeRequest) (<-chan *ConvertPayTypeResponse, <-chan error) {
	responseChan := make(chan *ConvertPayTypeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ConvertPayType(request)
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

// ConvertPayTypeWithCallback invokes the elasticsearch.ConvertPayType API asynchronously
// api document: https://help.aliyun.com/api/elasticsearch/convertpaytype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ConvertPayTypeWithCallback(request *ConvertPayTypeRequest, callback func(response *ConvertPayTypeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ConvertPayTypeResponse
		var err error
		defer close(result)
		response, err = client.ConvertPayType(request)
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

// ConvertPayTypeRequest is the request struct for api ConvertPayType
type ConvertPayTypeRequest struct {
	*requests.RoaRequest
	InstanceId  string `position:"Path" name:"InstanceId"`
	ClientToken string `position:"Query" name:"clientToken"`
}

// ConvertPayTypeResponse is the response struct for api ConvertPayType
type ConvertPayTypeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Result    bool   `json:"Result" xml:"Result"`
}

// CreateConvertPayTypeRequest creates a request to invoke ConvertPayType API
func CreateConvertPayTypeRequest() (request *ConvertPayTypeRequest) {
	request = &ConvertPayTypeRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "ConvertPayType", "/openapi/instances/[InstanceId]/convert-pay-type", "elasticsearch", "openAPI")
	request.Method = requests.POST
	return
}

// CreateConvertPayTypeResponse creates a response to parse from ConvertPayType response
func CreateConvertPayTypeResponse() (response *ConvertPayTypeResponse) {
	response = &ConvertPayTypeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
