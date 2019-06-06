package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyAccountDescription invokes the gpdb.ModifyAccountDescription API synchronously
// api document: https://help.aliyun.com/api/gpdb/modifyaccountdescription.html
func (client *Client) ModifyAccountDescription(request *ModifyAccountDescriptionRequest) (response *ModifyAccountDescriptionResponse, err error) {
	response = CreateModifyAccountDescriptionResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyAccountDescriptionWithChan invokes the gpdb.ModifyAccountDescription API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifyaccountdescription.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAccountDescriptionWithChan(request *ModifyAccountDescriptionRequest) (<-chan *ModifyAccountDescriptionResponse, <-chan error) {
	responseChan := make(chan *ModifyAccountDescriptionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyAccountDescription(request)
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

// ModifyAccountDescriptionWithCallback invokes the gpdb.ModifyAccountDescription API asynchronously
// api document: https://help.aliyun.com/api/gpdb/modifyaccountdescription.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAccountDescriptionWithCallback(request *ModifyAccountDescriptionRequest, callback func(response *ModifyAccountDescriptionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyAccountDescriptionResponse
		var err error
		defer close(result)
		response, err = client.ModifyAccountDescription(request)
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

// ModifyAccountDescriptionRequest is the request struct for api ModifyAccountDescription
type ModifyAccountDescriptionRequest struct {
	*requests.RpcRequest
	AccountName        string `position:"Query" name:"AccountName"`
	DBInstanceId       string `position:"Query" name:"DBInstanceId"`
	AccountDescription string `position:"Query" name:"AccountDescription"`
}

// ModifyAccountDescriptionResponse is the response struct for api ModifyAccountDescription
type ModifyAccountDescriptionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyAccountDescriptionRequest creates a request to invoke ModifyAccountDescription API
func CreateModifyAccountDescriptionRequest() (request *ModifyAccountDescriptionRequest) {
	request = &ModifyAccountDescriptionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "ModifyAccountDescription", "gpdb", "openAPI")
	return
}

// CreateModifyAccountDescriptionResponse creates a response to parse from ModifyAccountDescription response
func CreateModifyAccountDescriptionResponse() (response *ModifyAccountDescriptionResponse) {
	response = &ModifyAccountDescriptionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
