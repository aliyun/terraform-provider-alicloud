package cms

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

// PutGroupMetricRule invokes the cms.PutGroupMetricRule API synchronously
func (client *Client) PutGroupMetricRule(request *PutGroupMetricRuleRequest) (response *PutGroupMetricRuleResponse, err error) {
	response = CreatePutGroupMetricRuleResponse()
	err = client.DoAction(request, response)
	return
}

// PutGroupMetricRuleWithChan invokes the cms.PutGroupMetricRule API asynchronously
func (client *Client) PutGroupMetricRuleWithChan(request *PutGroupMetricRuleRequest) (<-chan *PutGroupMetricRuleResponse, <-chan error) {
	responseChan := make(chan *PutGroupMetricRuleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PutGroupMetricRule(request)
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

// PutGroupMetricRuleWithCallback invokes the cms.PutGroupMetricRule API asynchronously
func (client *Client) PutGroupMetricRuleWithCallback(request *PutGroupMetricRuleRequest, callback func(response *PutGroupMetricRuleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PutGroupMetricRuleResponse
		var err error
		defer close(result)
		response, err = client.PutGroupMetricRule(request)
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

// PutGroupMetricRuleRequest is the request struct for api PutGroupMetricRule
type PutGroupMetricRuleRequest struct {
	*requests.RpcRequest
	Webhook                               string           `position:"Query" name:"Webhook"`
	EscalationsWarnComparisonOperator     string           `position:"Query" name:"Escalations.Warn.ComparisonOperator"`
	RuleName                              string           `position:"Query" name:"RuleName"`
	EscalationsInfoStatistics             string           `position:"Query" name:"Escalations.Info.Statistics"`
	EffectiveInterval                     string           `position:"Query" name:"EffectiveInterval"`
	EscalationsInfoComparisonOperator     string           `position:"Query" name:"Escalations.Info.ComparisonOperator"`
	NoEffectiveInterval                   string           `position:"Query" name:"NoEffectiveInterval"`
	EmailSubject                          string           `position:"Query" name:"EmailSubject"`
	SilenceTime                           requests.Integer `position:"Query" name:"SilenceTime"`
	MetricName                            string           `position:"Query" name:"MetricName"`
	EscalationsWarnTimes                  requests.Integer `position:"Query" name:"Escalations.Warn.Times"`
	Period                                string           `position:"Query" name:"Period"`
	EscalationsWarnThreshold              string           `position:"Query" name:"Escalations.Warn.Threshold"`
	EscalationsCriticalStatistics         string           `position:"Query" name:"Escalations.Critical.Statistics"`
	GroupId                               string           `position:"Query" name:"GroupId"`
	EscalationsInfoTimes                  requests.Integer `position:"Query" name:"Escalations.Info.Times"`
	EscalationsCriticalTimes              requests.Integer `position:"Query" name:"Escalations.Critical.Times"`
	EscalationsWarnStatistics             string           `position:"Query" name:"Escalations.Warn.Statistics"`
	EscalationsInfoThreshold              string           `position:"Query" name:"Escalations.Info.Threshold"`
	Namespace                             string           `position:"Query" name:"Namespace"`
	Interval                              string           `position:"Query" name:"Interval"`
	RuleId                                string           `position:"Query" name:"RuleId"`
	Category                              string           `position:"Query" name:"Category"`
	EscalationsCriticalComparisonOperator string           `position:"Query" name:"Escalations.Critical.ComparisonOperator"`
	EscalationsCriticalThreshold          string           `position:"Query" name:"Escalations.Critical.Threshold"`
	Dimensions                            string           `position:"Query" name:"Dimensions"`
}

// PutGroupMetricRuleResponse is the response struct for api PutGroupMetricRule
type PutGroupMetricRuleResponse struct {
	*responses.BaseResponse
	Success   bool   `json:"Success" xml:"Success"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreatePutGroupMetricRuleRequest creates a request to invoke PutGroupMetricRule API
func CreatePutGroupMetricRuleRequest() (request *PutGroupMetricRuleRequest) {
	request = &PutGroupMetricRuleRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "PutGroupMetricRule", "cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreatePutGroupMetricRuleResponse creates a response to parse from PutGroupMetricRule response
func CreatePutGroupMetricRuleResponse() (response *PutGroupMetricRuleResponse) {
	response = &PutGroupMetricRuleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
