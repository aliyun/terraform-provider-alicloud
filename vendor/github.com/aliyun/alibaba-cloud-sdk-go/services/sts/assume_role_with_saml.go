package sts

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

// AssumeRoleWithSAML invokes the sts.AssumeRoleWithSAML API synchronously
func (client *Client) AssumeRoleWithSAML(request *AssumeRoleWithSAMLRequest) (response *AssumeRoleWithSAMLResponse, err error) {
	response = CreateAssumeRoleWithSAMLResponse()
	err = client.DoAction(request, response)
	return
}

// AssumeRoleWithSAMLWithChan invokes the sts.AssumeRoleWithSAML API asynchronously
func (client *Client) AssumeRoleWithSAMLWithChan(request *AssumeRoleWithSAMLRequest) (<-chan *AssumeRoleWithSAMLResponse, <-chan error) {
	responseChan := make(chan *AssumeRoleWithSAMLResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AssumeRoleWithSAML(request)
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

// AssumeRoleWithSAMLWithCallback invokes the sts.AssumeRoleWithSAML API asynchronously
func (client *Client) AssumeRoleWithSAMLWithCallback(request *AssumeRoleWithSAMLRequest, callback func(response *AssumeRoleWithSAMLResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AssumeRoleWithSAMLResponse
		var err error
		defer close(result)
		response, err = client.AssumeRoleWithSAML(request)
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

// AssumeRoleWithSAMLRequest is the request struct for api AssumeRoleWithSAML
type AssumeRoleWithSAMLRequest struct {
	*requests.RpcRequest
	SAMLAssertion   string           `position:"Query" name:"SAMLAssertion"`
	RoleArn         string           `position:"Query" name:"RoleArn"`
	SAMLProviderArn string           `position:"Query" name:"SAMLProviderArn"`
	DurationSeconds requests.Integer `position:"Query" name:"DurationSeconds"`
	Policy          string           `position:"Query" name:"Policy"`
}

// AssumeRoleWithSAMLResponse is the response struct for api AssumeRoleWithSAML
type AssumeRoleWithSAMLResponse struct {
	*responses.BaseResponse
	RequestId         string            `json:"RequestId" xml:"RequestId"`
	SAMLAssertionInfo SAMLAssertionInfo `json:"SAMLAssertionInfo" xml:"SAMLAssertionInfo"`
	AssumedRoleUser   AssumedRoleUser   `json:"AssumedRoleUser" xml:"AssumedRoleUser"`
	Credentials       Credentials       `json:"Credentials" xml:"Credentials"`
}

// CreateAssumeRoleWithSAMLRequest creates a request to invoke AssumeRoleWithSAML API
func CreateAssumeRoleWithSAMLRequest() (request *AssumeRoleWithSAMLRequest) {
	request = &AssumeRoleWithSAMLRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Sts", "2015-04-01", "AssumeRoleWithSAML", "sts", "openAPI")
	request.Method = requests.POST
	return
}

// CreateAssumeRoleWithSAMLResponse creates a response to parse from AssumeRoleWithSAML response
func CreateAssumeRoleWithSAMLResponse() (response *AssumeRoleWithSAMLResponse) {
	response = &AssumeRoleWithSAMLResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
