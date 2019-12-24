package rds

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

// MigrateDBInstance invokes the rds.MigrateDBInstance API synchronously
// api document: https://help.aliyun.com/api/rds/migratedbinstance.html
func (client *Client) MigrateDBInstance(request *MigrateDBInstanceRequest) (response *MigrateDBInstanceResponse, err error) {
	response = CreateMigrateDBInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// MigrateDBInstanceWithChan invokes the rds.MigrateDBInstance API asynchronously
// api document: https://help.aliyun.com/api/rds/migratedbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MigrateDBInstanceWithChan(request *MigrateDBInstanceRequest) (<-chan *MigrateDBInstanceResponse, <-chan error) {
	responseChan := make(chan *MigrateDBInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.MigrateDBInstance(request)
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

// MigrateDBInstanceWithCallback invokes the rds.MigrateDBInstance API asynchronously
// api document: https://help.aliyun.com/api/rds/migratedbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MigrateDBInstanceWithCallback(request *MigrateDBInstanceRequest, callback func(response *MigrateDBInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *MigrateDBInstanceResponse
		var err error
		defer close(result)
		response, err = client.MigrateDBInstance(request)
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

// MigrateDBInstanceRequest is the request struct for api MigrateDBInstance
type MigrateDBInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId                requests.Integer `position:"Query" name:"ResourceOwnerId"`
	SpecifiedTime                  string           `position:"Query" name:"SpecifiedTime"`
	TargetDedicatedHostIdForSlave  string           `position:"Query" name:"TargetDedicatedHostIdForSlave"`
	EngineVersion                  string           `position:"Query" name:"EngineVersion"`
	Storage                        requests.Integer `position:"Query" name:"Storage"`
	EffectiveTime                  string           `position:"Query" name:"EffectiveTime"`
	DBInstanceTransType            requests.Integer `position:"Query" name:"DBInstanceTransType"`
	TargetDedicatedHostIdForMaster string           `position:"Query" name:"TargetDedicatedHostIdForMaster"`
	DBInstanceId                   string           `position:"Query" name:"DBInstanceId"`
	DedicatedHostGroupId           string           `position:"Query" name:"DedicatedHostGroupId"`
	ResourceOwnerAccount           string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId                        requests.Integer `position:"Query" name:"OwnerId"`
	TargetDBInstanceClass          string           `position:"Query" name:"TargetDBInstanceClass"`
	VSwitchId                      string           `position:"Query" name:"VSwitchId"`
	TargetDedicatedHostIdForLog    string           `position:"Query" name:"TargetDedicatedHostIdForLog"`
	ZoneId                         string           `position:"Query" name:"ZoneId"`
}

// MigrateDBInstanceResponse is the response struct for api MigrateDBInstance
type MigrateDBInstanceResponse struct {
	*responses.BaseResponse
	RequestId   string `json:"RequestId" xml:"RequestId"`
	TaskId      int    `json:"TaskId" xml:"TaskId"`
	MigrationId int    `json:"MigrationId" xml:"MigrationId"`
}

// CreateMigrateDBInstanceRequest creates a request to invoke MigrateDBInstance API
func CreateMigrateDBInstanceRequest() (request *MigrateDBInstanceRequest) {
	request = &MigrateDBInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "MigrateDBInstance", "rds", "openAPI")
	return
}

// CreateMigrateDBInstanceResponse creates a response to parse from MigrateDBInstance response
func CreateMigrateDBInstanceResponse() (response *MigrateDBInstanceResponse) {
	response = &MigrateDBInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
