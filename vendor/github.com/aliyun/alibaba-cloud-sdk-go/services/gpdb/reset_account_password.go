package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ResetAccountPassword invokes the gpdb.ResetAccountPassword API synchronously
// api document: https://help.aliyun.com/api/gpdb/resetaccountpassword.html
func (client *Client) ResetAccountPassword(request *ResetAccountPasswordRequest) (response *ResetAccountPasswordResponse, err error) {
	response = CreateResetAccountPasswordResponse()
	err = client.DoAction(request, response)
	return
}

// ResetAccountPasswordWithChan invokes the gpdb.ResetAccountPassword API asynchronously
// api document: https://help.aliyun.com/api/gpdb/resetaccountpassword.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ResetAccountPasswordWithChan(request *ResetAccountPasswordRequest) (<-chan *ResetAccountPasswordResponse, <-chan error) {
	responseChan := make(chan *ResetAccountPasswordResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ResetAccountPassword(request)
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

// ResetAccountPasswordWithCallback invokes the gpdb.ResetAccountPassword API asynchronously
// api document: https://help.aliyun.com/api/gpdb/resetaccountpassword.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ResetAccountPasswordWithCallback(request *ResetAccountPasswordRequest, callback func(response *ResetAccountPasswordResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ResetAccountPasswordResponse
		var err error
		defer close(result)
		response, err = client.ResetAccountPassword(request)
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

// ResetAccountPasswordRequest is the request struct for api ResetAccountPassword
type ResetAccountPasswordRequest struct {
	*requests.RpcRequest
	AccountPassword string `position:"Query" name:"AccountPassword"`
	AccountName     string `position:"Query" name:"AccountName"`
	DBInstanceId    string `position:"Query" name:"DBInstanceId"`
}

// ResetAccountPasswordResponse is the response struct for api ResetAccountPassword
type ResetAccountPasswordResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateResetAccountPasswordRequest creates a request to invoke ResetAccountPassword API
func CreateResetAccountPasswordRequest() (request *ResetAccountPasswordRequest) {
	request = &ResetAccountPasswordRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "ResetAccountPassword", "gpdb", "openAPI")
	return
}

// CreateResetAccountPasswordResponse creates a response to parse from ResetAccountPassword response
func CreateResetAccountPasswordResponse() (response *ResetAccountPasswordResponse) {
	response = &ResetAccountPasswordResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
