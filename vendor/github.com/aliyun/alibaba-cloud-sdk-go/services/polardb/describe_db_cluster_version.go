package polardb

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

// DescribeDBClusterVersion invokes the polardb.DescribeDBClusterVersion API synchronously
func (client *Client) DescribeDBClusterVersion(request *DescribeDBClusterVersionRequest) (response *DescribeDBClusterVersionResponse, err error) {
	response = CreateDescribeDBClusterVersionResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDBClusterVersionWithChan invokes the polardb.DescribeDBClusterVersion API asynchronously
func (client *Client) DescribeDBClusterVersionWithChan(request *DescribeDBClusterVersionRequest) (<-chan *DescribeDBClusterVersionResponse, <-chan error) {
	responseChan := make(chan *DescribeDBClusterVersionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDBClusterVersion(request)
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

// DescribeDBClusterVersionWithCallback invokes the polardb.DescribeDBClusterVersion API asynchronously
func (client *Client) DescribeDBClusterVersionWithCallback(request *DescribeDBClusterVersionRequest, callback func(response *DescribeDBClusterVersionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDBClusterVersionResponse
		var err error
		defer close(result)
		response, err = client.DescribeDBClusterVersion(request)
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

// DescribeDBClusterVersionRequest is the request struct for api DescribeDBClusterVersion
type DescribeDBClusterVersionRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDBClusterVersionResponse is the response struct for api DescribeDBClusterVersion
type DescribeDBClusterVersionResponse struct {
	*responses.BaseResponse
	RequestId             string `json:"RequestId" xml:"RequestId"`
	DBClusterId           string `json:"DBClusterId" xml:"DBClusterId"`
	DBVersion             string `json:"DBVersion" xml:"DBVersion"`
	DBMinorVersion        string `json:"DBMinorVersion" xml:"DBMinorVersion"`
	DBRevisionVersion     string `json:"DBRevisionVersion" xml:"DBRevisionVersion"`
	DBVersionStatus       string `json:"DBVersionStatus" xml:"DBVersionStatus"`
	IsLatestVersion       string `json:"IsLatestVersion" xml:"IsLatestVersion"`
	LatestRevisionVersion string `json:"LatestRevisionVersion" xml:"LatestRevisionVersion"`
}

// CreateDescribeDBClusterVersionRequest creates a request to invoke DescribeDBClusterVersion API
func CreateDescribeDBClusterVersionRequest() (request *DescribeDBClusterVersionRequest) {
	request = &DescribeDBClusterVersionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("polardb", "2017-08-01", "DescribeDBClusterVersion", "polardb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDBClusterVersionResponse creates a response to parse from DescribeDBClusterVersion response
func CreateDescribeDBClusterVersionResponse() (response *DescribeDBClusterVersionResponse) {
	response = &DescribeDBClusterVersionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
