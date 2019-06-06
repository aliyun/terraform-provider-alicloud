package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyDBInstanceDescription invokes the gpdb.ModifyDBInstanceDescription API synchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancedescription.html
func (client *Client) ModifyDBInstanceDescription(request *ModifyDBInstanceDescriptionRequest) (response *ModifyDBInstanceDescriptionResponse, err error) {
	response = CreateModifyDBInstanceDescriptionResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyDBInstanceDescriptionWithChan invokes the gpdb.ModifyDBInstanceDescription API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancedescription.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyDBInstanceDescriptionWithChan(request *ModifyDBInstanceDescriptionRequest) (<-chan *ModifyDBInstanceDescriptionResponse, <-chan error) {
	responseChan := make(chan *ModifyDBInstanceDescriptionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyDBInstanceDescription(request)
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

// ModifyDBInstanceDescriptionWithCallback invokes the gpdb.ModifyDBInstanceDescription API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifydbinstancedescription.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyDBInstanceDescriptionWithCallback(request *ModifyDBInstanceDescriptionRequest, callback func(response *ModifyDBInstanceDescriptionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyDBInstanceDescriptionResponse
		var err error
		defer close(result)
		response, err = client.ModifyDBInstanceDescription(request)
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

// ModifyDBInstanceDescriptionRequest is the request struct for api ModifyDBInstanceDescription
type ModifyDBInstanceDescriptionRequest struct {
	*requests.RpcRequest
	DBInstanceId          string `position:"Query" name:"DBInstanceId"`
	DBInstanceDescription string `position:"Query" name:"DBInstanceDescription"`
}

// ModifyDBInstanceDescriptionResponse is the response struct for api ModifyDBInstanceDescription
type ModifyDBInstanceDescriptionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyDBInstanceDescriptionRequest creates a request to invoke ModifyDBInstanceDescription API
func CreateModifyDBInstanceDescriptionRequest() (request *ModifyDBInstanceDescriptionRequest) {
	request = &ModifyDBInstanceDescriptionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "ModifyDBInstanceDescription", "gpdb", "openAPI")
	return
}

// CreateModifyDBInstanceDescriptionResponse creates a response to parse from ModifyDBInstanceDescription response
func CreateModifyDBInstanceDescriptionResponse() (response *ModifyDBInstanceDescriptionResponse) {
	response = &ModifyDBInstanceDescriptionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
