package cs

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

// DescribeClusterV2UserKubeconfig invokes the cs.DescribeClusterV2UserKubeconfig API synchronously
func (client *Client) DescribeClusterV2UserKubeconfig(request *DescribeClusterV2UserKubeconfigRequest) (response *DescribeClusterV2UserKubeconfigResponse, err error) {
	response = CreateDescribeClusterV2UserKubeconfigResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeClusterV2UserKubeconfigWithChan invokes the cs.DescribeClusterV2UserKubeconfig API asynchronously
func (client *Client) DescribeClusterV2UserKubeconfigWithChan(request *DescribeClusterV2UserKubeconfigRequest) (<-chan *DescribeClusterV2UserKubeconfigResponse, <-chan error) {
	responseChan := make(chan *DescribeClusterV2UserKubeconfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeClusterV2UserKubeconfig(request)
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

// DescribeClusterV2UserKubeconfigWithCallback invokes the cs.DescribeClusterV2UserKubeconfig API asynchronously
func (client *Client) DescribeClusterV2UserKubeconfigWithCallback(request *DescribeClusterV2UserKubeconfigRequest, callback func(response *DescribeClusterV2UserKubeconfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeClusterV2UserKubeconfigResponse
		var err error
		defer close(result)
		response, err = client.DescribeClusterV2UserKubeconfig(request)
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

// DescribeClusterV2UserKubeconfigRequest is the request struct for api DescribeClusterV2UserKubeconfig
type DescribeClusterV2UserKubeconfigRequest struct {
	*requests.RoaRequest
	PrivateIpAddress requests.Boolean `position:"Query" name:"PrivateIpAddress"`
	ClusterId        string           `position:"Path" name:"ClusterId"`
}

// DescribeClusterV2UserKubeconfigResponse is the response struct for api DescribeClusterV2UserKubeconfig
type DescribeClusterV2UserKubeconfigResponse struct {
	*responses.BaseResponse
	Config string `json:"config" xml:"config"`
}

// CreateDescribeClusterV2UserKubeconfigRequest creates a request to invoke DescribeClusterV2UserKubeconfig API
func CreateDescribeClusterV2UserKubeconfigRequest() (request *DescribeClusterV2UserKubeconfigRequest) {
	request = &DescribeClusterV2UserKubeconfigRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("CS", "2015-12-15", "DescribeClusterV2UserKubeconfig", "/api/v2/k8s/[ClusterId]/user_config", "", "")
	request.Method = requests.GET
	return
}

// CreateDescribeClusterV2UserKubeconfigResponse creates a response to parse from DescribeClusterV2UserKubeconfig response
func CreateDescribeClusterV2UserKubeconfigResponse() (response *DescribeClusterV2UserKubeconfigResponse) {
	response = &DescribeClusterV2UserKubeconfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
