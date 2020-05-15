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

// DescribeDBClusterAttribute invokes the polardb.DescribeDBClusterAttribute API synchronously
// api document: https://help.aliyun.com/api/polardb/describedbclusterattribute.html
func (client *Client) DescribeDBClusterAttribute(request *DescribeDBClusterAttributeRequest) (response *DescribeDBClusterAttributeResponse, err error) {
	response = CreateDescribeDBClusterAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDBClusterAttributeWithChan invokes the polardb.DescribeDBClusterAttribute API asynchronously
// api document: https://help.aliyun.com/api/polardb/describedbclusterattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBClusterAttributeWithChan(request *DescribeDBClusterAttributeRequest) (<-chan *DescribeDBClusterAttributeResponse, <-chan error) {
	responseChan := make(chan *DescribeDBClusterAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDBClusterAttribute(request)
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

// DescribeDBClusterAttributeWithCallback invokes the polardb.DescribeDBClusterAttribute API asynchronously
// api document: https://help.aliyun.com/api/polardb/describedbclusterattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBClusterAttributeWithCallback(request *DescribeDBClusterAttributeRequest, callback func(response *DescribeDBClusterAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDBClusterAttributeResponse
		var err error
		defer close(result)
		response, err = client.DescribeDBClusterAttribute(request)
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

// DescribeDBClusterAttributeRequest is the request struct for api DescribeDBClusterAttribute
type DescribeDBClusterAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeDBClusterAttributeResponse is the response struct for api DescribeDBClusterAttribute
type DescribeDBClusterAttributeResponse struct {
	*responses.BaseResponse
	RequestId                 string   `json:"RequestId" xml:"RequestId"`
	RegionId                  string   `json:"RegionId" xml:"RegionId"`
	DBClusterNetworkType      string   `json:"DBClusterNetworkType" xml:"DBClusterNetworkType"`
	VPCId                     string   `json:"VPCId" xml:"VPCId"`
	VSwitchId                 string   `json:"VSwitchId" xml:"VSwitchId"`
	PayType                   string   `json:"PayType" xml:"PayType"`
	DBClusterId               string   `json:"DBClusterId" xml:"DBClusterId"`
	DBClusterStatus           string   `json:"DBClusterStatus" xml:"DBClusterStatus"`
	DBClusterDescription      string   `json:"DBClusterDescription" xml:"DBClusterDescription"`
	Engine                    string   `json:"Engine" xml:"Engine"`
	DBType                    string   `json:"DBType" xml:"DBType"`
	DBVersion                 string   `json:"DBVersion" xml:"DBVersion"`
	LockMode                  string   `json:"LockMode" xml:"LockMode"`
	DeletionLock              int      `json:"DeletionLock" xml:"DeletionLock"`
	CreationTime              string   `json:"CreationTime" xml:"CreationTime"`
	ExpireTime                string   `json:"ExpireTime" xml:"ExpireTime"`
	Expired                   string   `json:"Expired" xml:"Expired"`
	MaintainTime              string   `json:"MaintainTime" xml:"MaintainTime"`
	StorageUsed               int64    `json:"StorageUsed" xml:"StorageUsed"`
	StorageMax                int64    `json:"StorageMax" xml:"StorageMax"`
	ZoneIds                   string   `json:"ZoneIds" xml:"ZoneIds"`
	SQLSize                   int64    `json:"SQLSize" xml:"SQLSize"`
	IsLatestVersion           bool     `json:"IsLatestVersion" xml:"IsLatestVersion"`
	ResourceGroupId           string   `json:"ResourceGroupId" xml:"ResourceGroupId"`
	DataLevel1BackupChainSize int64    `json:"DataLevel1BackupChainSize" xml:"DataLevel1BackupChainSize"`
	Tags                      []Tag    `json:"Tags" xml:"Tags"`
	DBNodes                   []DBNode `json:"DBNodes" xml:"DBNodes"`
}

// CreateDescribeDBClusterAttributeRequest creates a request to invoke DescribeDBClusterAttribute API
func CreateDescribeDBClusterAttributeRequest() (request *DescribeDBClusterAttributeRequest) {
	request = &DescribeDBClusterAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("polardb", "2017-08-01", "DescribeDBClusterAttribute", "polardb", "openAPI")
	return
}

// CreateDescribeDBClusterAttributeResponse creates a response to parse from DescribeDBClusterAttribute response
func CreateDescribeDBClusterAttributeResponse() (response *DescribeDBClusterAttributeResponse) {
	response = &DescribeDBClusterAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
