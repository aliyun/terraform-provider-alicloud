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

// TransformClusterMember invokes the edas.TransformClusterMember API synchronously
// api document: https://help.aliyun.com/api/edas/transformclustermember.html
func (client *Client) TransformClusterMember(request *TransformClusterMemberRequest) (response *TransformClusterMemberResponse, err error) {
	response = CreateTransformClusterMemberResponse()
	err = client.DoAction(request, response)
	return
}

// TransformClusterMemberWithChan invokes the edas.TransformClusterMember API asynchronously
// api document: https://help.aliyun.com/api/edas/transformclustermember.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TransformClusterMemberWithChan(request *TransformClusterMemberRequest) (<-chan *TransformClusterMemberResponse, <-chan error) {
	responseChan := make(chan *TransformClusterMemberResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.TransformClusterMember(request)
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

// TransformClusterMemberWithCallback invokes the edas.TransformClusterMember API asynchronously
// api document: https://help.aliyun.com/api/edas/transformclustermember.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TransformClusterMemberWithCallback(request *TransformClusterMemberRequest, callback func(response *TransformClusterMemberResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *TransformClusterMemberResponse
		var err error
		defer close(result)
		response, err = client.TransformClusterMember(request)
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

// TransformClusterMemberRequest is the request struct for api TransformClusterMember
type TransformClusterMemberRequest struct {
	*requests.RoaRequest
	Password        string `position:"Query" name:"Password"`
	InstanceIds     string `position:"Query" name:"InstanceIds"`
	TargetClusterId string `position:"Query" name:"TargetClusterId"`
}

// TransformClusterMemberResponse is the response struct for api TransformClusterMember
type TransformClusterMemberResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      string `json:"Data" xml:"Data"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateTransformClusterMemberRequest creates a request to invoke TransformClusterMember API
func CreateTransformClusterMemberRequest() (request *TransformClusterMemberRequest) {
	request = &TransformClusterMemberRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Edas", "2017-08-01", "TransformClusterMember", "/pop/v5/resource/transform_cluster_member", "edas", "openAPI")
	request.Method = requests.POST
	return
}

// CreateTransformClusterMemberResponse creates a response to parse from TransformClusterMember response
func CreateTransformClusterMemberResponse() (response *TransformClusterMemberResponse) {
	response = &TransformClusterMemberResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
