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

// PreCheckCreateOrderForTempUpgrade invokes the rds.PreCheckCreateOrderForTempUpgrade API synchronously
// api document: https://help.aliyun.com/api/rds/precheckcreateorderfortempupgrade.html
func (client *Client) PreCheckCreateOrderForTempUpgrade(request *PreCheckCreateOrderForTempUpgradeRequest) (response *PreCheckCreateOrderForTempUpgradeResponse, err error) {
	response = CreatePreCheckCreateOrderForTempUpgradeResponse()
	err = client.DoAction(request, response)
	return
}

// PreCheckCreateOrderForTempUpgradeWithChan invokes the rds.PreCheckCreateOrderForTempUpgrade API asynchronously
// api document: https://help.aliyun.com/api/rds/precheckcreateorderfortempupgrade.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) PreCheckCreateOrderForTempUpgradeWithChan(request *PreCheckCreateOrderForTempUpgradeRequest) (<-chan *PreCheckCreateOrderForTempUpgradeResponse, <-chan error) {
	responseChan := make(chan *PreCheckCreateOrderForTempUpgradeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PreCheckCreateOrderForTempUpgrade(request)
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

// PreCheckCreateOrderForTempUpgradeWithCallback invokes the rds.PreCheckCreateOrderForTempUpgrade API asynchronously
// api document: https://help.aliyun.com/api/rds/precheckcreateorderfortempupgrade.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) PreCheckCreateOrderForTempUpgradeWithCallback(request *PreCheckCreateOrderForTempUpgradeRequest, callback func(response *PreCheckCreateOrderForTempUpgradeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PreCheckCreateOrderForTempUpgradeResponse
		var err error
		defer close(result)
		response, err = client.PreCheckCreateOrderForTempUpgrade(request)
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

// PreCheckCreateOrderForTempUpgradeRequest is the request struct for api PreCheckCreateOrderForTempUpgrade
type PreCheckCreateOrderForTempUpgradeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId       requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceStorage     requests.Integer `position:"Query" name:"DBInstanceStorage"`
	NodeType              string           `position:"Query" name:"NodeType"`
	ClientToken           string           `position:"Query" name:"ClientToken"`
	EffectiveTime         string           `position:"Query" name:"EffectiveTime"`
	DBInstanceId          string           `position:"Query" name:"DBInstanceId"`
	DBInstanceStorageType string           `position:"Query" name:"DBInstanceStorageType"`
	BusinessInfo          string           `position:"Query" name:"BusinessInfo"`
	AutoPay               requests.Boolean `position:"Query" name:"AutoPay"`
	ResourceOwnerAccount  string           `position:"Query" name:"ResourceOwnerAccount"`
	Resource              string           `position:"Query" name:"Resource"`
	CommodityCode         string           `position:"Query" name:"CommodityCode"`
	OwnerId               requests.Integer `position:"Query" name:"OwnerId"`
	UsedTime              string           `position:"Query" name:"UsedTime"`
	DBInstanceClass       string           `position:"Query" name:"DBInstanceClass"`
}

// PreCheckCreateOrderForTempUpgradeResponse is the response struct for api PreCheckCreateOrderForTempUpgrade
type PreCheckCreateOrderForTempUpgradeResponse struct {
	*responses.BaseResponse
	RequestId      string                                      `json:"RequestId" xml:"RequestId"`
	PreCheckResult bool                                        `json:"PreCheckResult" xml:"PreCheckResult"`
	Failures       FailuresInPreCheckCreateOrderForTempUpgrade `json:"Failures" xml:"Failures"`
}

// CreatePreCheckCreateOrderForTempUpgradeRequest creates a request to invoke PreCheckCreateOrderForTempUpgrade API
func CreatePreCheckCreateOrderForTempUpgradeRequest() (request *PreCheckCreateOrderForTempUpgradeRequest) {
	request = &PreCheckCreateOrderForTempUpgradeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "PreCheckCreateOrderForTempUpgrade", "rds", "openAPI")
	return
}

// CreatePreCheckCreateOrderForTempUpgradeResponse creates a response to parse from PreCheckCreateOrderForTempUpgrade response
func CreatePreCheckCreateOrderForTempUpgradeResponse() (response *PreCheckCreateOrderForTempUpgradeResponse) {
	response = &PreCheckCreateOrderForTempUpgradeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
