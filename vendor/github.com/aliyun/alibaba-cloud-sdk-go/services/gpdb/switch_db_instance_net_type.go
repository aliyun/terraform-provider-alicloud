package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// SwitchDBInstanceNetType invokes the gpdb.SwitchDBInstanceNetType API synchronously
// api document: https://help.aliyun.com/api/gpdb/switchdbinstancenettype.html
func (client *Client) SwitchDBInstanceNetType(request *SwitchDBInstanceNetTypeRequest) (response *SwitchDBInstanceNetTypeResponse, err error) {
	response = CreateSwitchDBInstanceNetTypeResponse()
	err = client.DoAction(request, response)
	return
}

// SwitchDBInstanceNetTypeWithChan invokes the gpdb.SwitchDBInstanceNetType API asynchronously
// api document: https://help.aliyun.com/api/gpdb/switchdbinstancenettype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SwitchDBInstanceNetTypeWithChan(request *SwitchDBInstanceNetTypeRequest) (<-chan *SwitchDBInstanceNetTypeResponse, <-chan error) {
	responseChan := make(chan *SwitchDBInstanceNetTypeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SwitchDBInstanceNetType(request)
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

// SwitchDBInstanceNetTypeWithCallback invokes the gpdb.SwitchDBInstanceNetType API asynchronously
// api document: https://help.aliyun.com/api/gpdb/switchdbinstancenettype.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SwitchDBInstanceNetTypeWithCallback(request *SwitchDBInstanceNetTypeRequest, callback func(response *SwitchDBInstanceNetTypeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SwitchDBInstanceNetTypeResponse
		var err error
		defer close(result)
		response, err = client.SwitchDBInstanceNetType(request)
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

// SwitchDBInstanceNetTypeRequest is the request struct for api SwitchDBInstanceNetType
type SwitchDBInstanceNetTypeRequest struct {
	*requests.RpcRequest
	ConnectionStringPrefix string `position:"Query" name:"ConnectionStringPrefix"`
	Port                   string `position:"Query" name:"Port"`
	DBInstanceId           string `position:"Query" name:"DBInstanceId"`
}

// SwitchDBInstanceNetTypeResponse is the response struct for api SwitchDBInstanceNetType
type SwitchDBInstanceNetTypeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSwitchDBInstanceNetTypeRequest creates a request to invoke SwitchDBInstanceNetType API
func CreateSwitchDBInstanceNetTypeRequest() (request *SwitchDBInstanceNetTypeRequest) {
	request = &SwitchDBInstanceNetTypeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "SwitchDBInstanceNetType", "gpdb", "openAPI")
	return
}

// CreateSwitchDBInstanceNetTypeResponse creates a response to parse from SwitchDBInstanceNetType response
func CreateSwitchDBInstanceNetTypeResponse() (response *SwitchDBInstanceNetTypeResponse) {
	response = &SwitchDBInstanceNetTypeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
