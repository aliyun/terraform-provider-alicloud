package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DeleteDBInstance invokes the gpdb.DeleteDBInstance API synchronously
// api document: https://help.aliyun.com/api/gpdb/deletedbinstance.html
func (client *Client) DeleteDBInstance(request *DeleteDBInstanceRequest) (response *DeleteDBInstanceResponse, err error) {
	response = CreateDeleteDBInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteDBInstanceWithChan invokes the gpdb.DeleteDBInstance API asynchronously
// api document: https://help.aliyun.com/api/gpdb/deletedbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDBInstanceWithChan(request *DeleteDBInstanceRequest) (<-chan *DeleteDBInstanceResponse, <-chan error) {
	responseChan := make(chan *DeleteDBInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteDBInstance(request)
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

// DeleteDBInstanceWithCallback invokes the gpdb.DeleteDBInstance API asynchronously
// api document: https://help.aliyun.com/api/gpdb/deletedbinstance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDBInstanceWithCallback(request *DeleteDBInstanceRequest, callback func(response *DeleteDBInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteDBInstanceResponse
		var err error
		defer close(result)
		response, err = client.DeleteDBInstance(request)
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

// DeleteDBInstanceRequest is the request struct for api DeleteDBInstance
type DeleteDBInstanceRequest struct {
	*requests.RpcRequest
	ClientToken  string           `position:"Query" name:"ClientToken"`
	DBInstanceId string           `position:"Query" name:"DBInstanceId"`
	OwnerId      requests.Integer `position:"Query" name:"OwnerId"`
}

// DeleteDBInstanceResponse is the response struct for api DeleteDBInstance
type DeleteDBInstanceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteDBInstanceRequest creates a request to invoke DeleteDBInstance API
func CreateDeleteDBInstanceRequest() (request *DeleteDBInstanceRequest) {
	request = &DeleteDBInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "DeleteDBInstance", "gpdb", "openAPI")
	return
}

// CreateDeleteDBInstanceResponse creates a response to parse from DeleteDBInstance response
func CreateDeleteDBInstanceResponse() (response *DeleteDBInstanceResponse) {
	response = &DeleteDBInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
