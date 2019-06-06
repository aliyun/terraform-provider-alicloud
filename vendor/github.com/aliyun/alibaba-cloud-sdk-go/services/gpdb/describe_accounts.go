package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeAccounts invokes the gpdb.DescribeAccounts API synchronously
// api document: https://help.aliyun.com/api/gpdb/describeaccounts.html
func (client *Client) DescribeAccounts(request *DescribeAccountsRequest) (response *DescribeAccountsResponse, err error) {
	response = CreateDescribeAccountsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAccountsWithChan invokes the gpdb.DescribeAccounts API asynchronously
// api document: https://help.aliyun.com/api/gpdb/describeaccounts.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAccountsWithChan(request *DescribeAccountsRequest) (<-chan *DescribeAccountsResponse, <-chan error) {
	responseChan := make(chan *DescribeAccountsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAccounts(request)
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

// DescribeAccountsWithCallback invokes the gpdb.DescribeAccounts API asynchronously
// api document: https://help.aliyun.com/api/gpdb/describeaccounts.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAccountsWithCallback(request *DescribeAccountsRequest, callback func(response *DescribeAccountsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAccountsResponse
		var err error
		defer close(result)
		response, err = client.DescribeAccounts(request)
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

// DescribeAccountsRequest is the request struct for api DescribeAccounts
type DescribeAccountsRequest struct {
	*requests.RpcRequest
	AccountName  string `position:"Query" name:"AccountName"`
	DBInstanceId string `position:"Query" name:"DBInstanceId"`
}

// DescribeAccountsResponse is the response struct for api DescribeAccounts
type DescribeAccountsResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Accounts  Accounts `json:"Accounts" xml:"Accounts"`
}

// CreateDescribeAccountsRequest creates a request to invoke DescribeAccounts API
func CreateDescribeAccountsRequest() (request *DescribeAccountsRequest) {
	request = &DescribeAccountsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "DescribeAccounts", "gpdb", "openAPI")
	return
}

// CreateDescribeAccountsResponse creates a response to parse from DescribeAccounts response
func CreateDescribeAccountsResponse() (response *DescribeAccountsResponse) {
	response = &DescribeAccountsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
