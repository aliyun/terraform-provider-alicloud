package r_kvstore

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

// CreateTairInstance invokes the r_kvstore.CreateTairInstance API synchronously
func (client *Client) CreateTairInstance(request *CreateTairInstanceRequest) (response *CreateTairInstanceResponse, err error) {
	response = CreateCreateTairInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// CreateTairInstanceWithChan invokes the r_kvstore.CreateTairInstance API asynchronously
func (client *Client) CreateTairInstanceWithChan(request *CreateTairInstanceRequest) (<-chan *CreateTairInstanceResponse, <-chan error) {
	responseChan := make(chan *CreateTairInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateTairInstance(request)
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

// CreateTairInstanceWithCallback invokes the r_kvstore.CreateTairInstance API asynchronously
func (client *Client) CreateTairInstanceWithCallback(request *CreateTairInstanceRequest, callback func(response *CreateTairInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateTairInstanceResponse
		var err error
		defer close(result)
		response, err = client.CreateTairInstance(request)
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

// CreateTairInstanceRequest is the request struct for api CreateTairInstance
type CreateTairInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId        requests.Integer         `position:"Query" name:"ResourceOwnerId"`
	SecondaryZoneId        string                   `position:"Query" name:"SecondaryZoneId"`
	CouponNo               string                   `position:"Query" name:"CouponNo"`
	EngineVersion          string                   `position:"Query" name:"EngineVersion"`
	StorageType            string                   `position:"Query" name:"StorageType"`
	ResourceGroupId        string                   `position:"Query" name:"ResourceGroupId"`
	Password               string                   `position:"Query" name:"Password"`
	SecurityToken          string                   `position:"Query" name:"SecurityToken"`
	Tag                    *[]CreateTairInstanceTag `position:"Query" name:"Tag"  type:"Repeated"`
	GlobalSecurityGroupIds string                   `position:"Query" name:"GlobalSecurityGroupIds"`
	BusinessInfo           string                   `position:"Query" name:"BusinessInfo"`
	ShardCount             requests.Integer         `position:"Query" name:"ShardCount"`
	AutoRenewPeriod        string                   `position:"Query" name:"AutoRenewPeriod"`
	Period                 requests.Integer         `position:"Query" name:"Period"`
	DryRun                 requests.Boolean         `position:"Query" name:"DryRun"`
	BackupId               string                   `position:"Query" name:"BackupId"`
	OwnerId                requests.Integer         `position:"Query" name:"OwnerId"`
	ShardType              string                   `position:"Query" name:"ShardType"`
	VSwitchId              string                   `position:"Query" name:"VSwitchId"`
	PrivateIpAddress       string                   `position:"Query" name:"PrivateIpAddress"`
	InstanceName           string                   `position:"Query" name:"InstanceName"`
	AutoRenew              string                   `position:"Query" name:"AutoRenew"`
	Port                   requests.Integer         `position:"Query" name:"Port"`
	ZoneId                 string                   `position:"Query" name:"ZoneId"`
	ClientToken            string                   `position:"Query" name:"ClientToken"`
	AutoUseCoupon          string                   `position:"Query" name:"AutoUseCoupon"`
	Storage                requests.Integer         `position:"Query" name:"Storage"`
	InstanceClass          string                   `position:"Query" name:"InstanceClass"`
	InstanceType           string                   `position:"Query" name:"InstanceType"`
	AutoPay                requests.Boolean         `position:"Query" name:"AutoPay"`
	ResourceOwnerAccount   string                   `position:"Query" name:"ResourceOwnerAccount"`
	SrcDBInstanceId        string                   `position:"Query" name:"SrcDBInstanceId"`
	OwnerAccount           string                   `position:"Query" name:"OwnerAccount"`
	GlobalInstanceId       string                   `position:"Query" name:"GlobalInstanceId"`
	ParamGroupId           string                   `position:"Query" name:"ParamGroupId"`
	VpcId                  string                   `position:"Query" name:"VpcId"`
	ReadOnlyCount          requests.Integer         `position:"Query" name:"ReadOnlyCount"`
	ChargeType             string                   `position:"Query" name:"ChargeType"`
}

// CreateTairInstanceTag is a repeated param struct in CreateTairInstanceRequest
type CreateTairInstanceTag struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// CreateTairInstanceResponse is the response struct for api CreateTairInstance
type CreateTairInstanceResponse struct {
	*responses.BaseResponse
	QPS              int64  `json:"QPS" xml:"QPS"`
	ConnectionDomain string `json:"ConnectionDomain" xml:"ConnectionDomain"`
	ChargeType       string `json:"ChargeType" xml:"ChargeType"`
	InstanceId       string `json:"InstanceId" xml:"InstanceId"`
	Port             int    `json:"Port" xml:"Port"`
	Config           string `json:"Config" xml:"Config"`
	RegionId         string `json:"RegionId" xml:"RegionId"`
	RequestId        string `json:"RequestId" xml:"RequestId"`
	Bandwidth        int64  `json:"Bandwidth" xml:"Bandwidth"`
	Connections      int64  `json:"Connections" xml:"Connections"`
	InstanceName     string `json:"InstanceName" xml:"InstanceName"`
	ZoneId           string `json:"ZoneId" xml:"ZoneId"`
	InstanceStatus   string `json:"InstanceStatus" xml:"InstanceStatus"`
	TaskId           string `json:"TaskId" xml:"TaskId"`
	OrderId          int64  `json:"OrderId" xml:"OrderId"`
}

// CreateCreateTairInstanceRequest creates a request to invoke CreateTairInstance API
func CreateCreateTairInstanceRequest() (request *CreateTairInstanceRequest) {
	request = &CreateTairInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("R-kvstore", "2015-01-01", "CreateTairInstance", "redisa", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateTairInstanceResponse creates a response to parse from CreateTairInstance response
func CreateCreateTairInstanceResponse() (response *CreateTairInstanceResponse) {
	response = &CreateTairInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
