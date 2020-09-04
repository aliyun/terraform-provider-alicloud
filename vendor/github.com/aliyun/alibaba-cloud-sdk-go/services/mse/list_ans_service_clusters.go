package mse

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

// ListAnsServiceClusters invokes the mse.ListAnsServiceClusters API synchronously
// api document: https://help.aliyun.com/api/mse/listansserviceclusters.html
func (client *Client) ListAnsServiceClusters(request *ListAnsServiceClustersRequest) (response *ListAnsServiceClustersResponse, err error) {
	response = CreateListAnsServiceClustersResponse()
	err = client.DoAction(request, response)
	return
}

// ListAnsServiceClustersWithChan invokes the mse.ListAnsServiceClusters API asynchronously
// api document: https://help.aliyun.com/api/mse/listansserviceclusters.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListAnsServiceClustersWithChan(request *ListAnsServiceClustersRequest) (<-chan *ListAnsServiceClustersResponse, <-chan error) {
	responseChan := make(chan *ListAnsServiceClustersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListAnsServiceClusters(request)
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

// ListAnsServiceClustersWithCallback invokes the mse.ListAnsServiceClusters API asynchronously
// api document: https://help.aliyun.com/api/mse/listansserviceclusters.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListAnsServiceClustersWithCallback(request *ListAnsServiceClustersRequest, callback func(response *ListAnsServiceClustersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListAnsServiceClustersResponse
		var err error
		defer close(result)
		response, err = client.ListAnsServiceClusters(request)
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

// ListAnsServiceClustersRequest is the request struct for api ListAnsServiceClusters
type ListAnsServiceClustersRequest struct {
	*requests.RpcRequest
	ClusterName string           `position:"Query" name:"ClusterName"`
	ClusterId   string           `position:"Query" name:"ClusterId"`
	PageNum     requests.Integer `position:"Query" name:"PageNum"`
	GroupName   string           `position:"Query" name:"GroupName"`
	NamespaceId string           `position:"Query" name:"NamespaceId"`
	RequestPars string           `position:"Query" name:"RequestPars"`
	PageSize    requests.Integer `position:"Query" name:"PageSize"`
	ServiceName string           `position:"Query" name:"ServiceName"`
}

// ListAnsServiceClustersResponse is the response struct for api ListAnsServiceClusters
type ListAnsServiceClustersResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Message   string `json:"Message" xml:"Message"`
	ErrorCode string `json:"ErrorCode" xml:"ErrorCode"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateListAnsServiceClustersRequest creates a request to invoke ListAnsServiceClusters API
func CreateListAnsServiceClustersRequest() (request *ListAnsServiceClustersRequest) {
	request = &ListAnsServiceClustersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("mse", "2019-05-31", "ListAnsServiceClusters", "mse", "openAPI")
	request.Method = requests.GET
	return
}

// CreateListAnsServiceClustersResponse creates a response to parse from ListAnsServiceClusters response
func CreateListAnsServiceClustersResponse() (response *ListAnsServiceClustersResponse) {
	response = &ListAnsServiceClustersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
