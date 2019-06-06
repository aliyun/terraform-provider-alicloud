package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyDBInstanceNetworkType invokes the gpdb.ModifyDBInstanceNetworkType API synchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancenetworktype.html
func (client *Client) ModifyDBInstanceNetworkType(request *ModifyDBInstanceNetworkTypeRequest) (response *ModifyDBInstanceNetworkTypeResponse, err error) {
	response = CreateModifyDBInstanceNetworkTypeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyDBInstanceNetworkTypeWithChan invokes the gpdb.ModifyDBInstanceNetworkType API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancenetworktype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyDBInstanceNetworkTypeWithChan(request *ModifyDBInstanceNetworkTypeRequest) (<-chan *ModifyDBInstanceNetworkTypeResponse, <-chan error) {
	responseChan := make(chan *ModifyDBInstanceNetworkTypeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyDBInstanceNetworkType(request)
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

// ModifyDBInstanceNetworkTypeWithCallback invokes the gpdb.ModifyDBInstanceNetworkType API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancenetworktype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyDBInstanceNetworkTypeWithCallback(request *ModifyDBInstanceNetworkTypeRequest, callback func(response *ModifyDBInstanceNetworkTypeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyDBInstanceNetworkTypeResponse
		var err error
		defer close(result)
		response, err = client.ModifyDBInstanceNetworkType(request)
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

// ModifyDBInstanceNetworkTypeRequest is the request struct for api ModifyDBInstanceNetworkType
type ModifyDBInstanceNetworkTypeRequest struct {
	*requests.RpcRequest
	VSwitchId           string `position:"Query" name:"VSwitchId"`
	PrivateIpAddress    string `position:"Query" name:"PrivateIpAddress"`
	VPCId               string `position:"Query" name:"VPCId"`
	DBInstanceId        string `position:"Query" name:"DBInstanceId"`
	InstanceNetworkType string `position:"Query" name:"InstanceNetworkType"`
}

// ModifyDBInstanceNetworkTypeResponse is the response struct for api ModifyDBInstanceNetworkType
type ModifyDBInstanceNetworkTypeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyDBInstanceNetworkTypeRequest creates a request to invoke ModifyDBInstanceNetworkType API
func CreateModifyDBInstanceNetworkTypeRequest() (request *ModifyDBInstanceNetworkTypeRequest) {
	request = &ModifyDBInstanceNetworkTypeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "ModifyDBInstanceNetworkType", "gpdb", "openAPI")
	return
}

// CreateModifyDBInstanceNetworkTypeResponse creates a response to parse from ModifyDBInstanceNetworkType response
func CreateModifyDBInstanceNetworkTypeResponse() (response *ModifyDBInstanceNetworkTypeResponse) {
	response = &ModifyDBInstanceNetworkTypeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
