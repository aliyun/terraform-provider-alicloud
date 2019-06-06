package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeSQLLogRecords invokes the gpdb.DescribeSQLLogRecords API synchronously
// api document: https://help.aliyun.com/api/gpdb/describesqllogrecords.html
func (client *Client) DescribeSQLLogRecords(request *DescribeSQLLogRecordsRequest) (response *DescribeSQLLogRecordsResponse, err error) {
	response = CreateDescribeSQLLogRecordsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSQLLogRecordsWithChan invokes the gpdb.DescribeSQLLogRecords API asynchronously
// api document: https://help.aliyun.com/api/gpdb/describesqllogrecords.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSQLLogRecordsWithChan(request *DescribeSQLLogRecordsRequest) (<-chan *DescribeSQLLogRecordsResponse, <-chan error) {
	responseChan := make(chan *DescribeSQLLogRecordsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSQLLogRecords(request)
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

// DescribeSQLLogRecordsWithCallback invokes the gpdb.DescribeSQLLogRecords API asynchronously
// api document: https://help.aliyun.com/api/gpdb/describesqllogrecords.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSQLLogRecordsWithCallback(request *DescribeSQLLogRecordsRequest, callback func(response *DescribeSQLLogRecordsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSQLLogRecordsResponse
		var err error
		defer close(result)
		response, err = client.DescribeSQLLogRecords(request)
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

// DescribeSQLLogRecordsRequest is the request struct for api DescribeSQLLogRecords
type DescribeSQLLogRecordsRequest struct {
	*requests.RpcRequest
	Database      string           `position:"Query" name:"Database"`
	Form          string           `position:"Query" name:"Form"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	EndTime       string           `position:"Query" name:"EndTime"`
	DBInstanceId  string           `position:"Query" name:"DBInstanceId"`
	StartTime     string           `position:"Query" name:"StartTime"`
	User          string           `position:"Query" name:"User"`
	QueryKeywords string           `position:"Query" name:"QueryKeywords"`
	PageNumber    requests.Integer `position:"Query" name:"PageNumber"`
}

// DescribeSQLLogRecordsResponse is the response struct for api DescribeSQLLogRecords
type DescribeSQLLogRecordsResponse struct {
	*responses.BaseResponse
	RequestId        string                       `json:"RequestId" xml:"RequestId"`
	TotalRecordCount int                          `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageNumber       int                          `json:"PageNumber" xml:"PageNumber"`
	PageRecordCount  int                          `json:"PageRecordCount" xml:"PageRecordCount"`
	Items            ItemsInDescribeSQLLogRecords `json:"Items" xml:"Items"`
}

// CreateDescribeSQLLogRecordsRequest creates a request to invoke DescribeSQLLogRecords API
func CreateDescribeSQLLogRecordsRequest() (request *DescribeSQLLogRecordsRequest) {
	request = &DescribeSQLLogRecordsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("gpdb", "2016-05-03", "DescribeSQLLogRecords", "gpdb", "openAPI")
	return
}

// CreateDescribeSQLLogRecordsResponse creates a response to parse from DescribeSQLLogRecords response
func CreateDescribeSQLLogRecordsResponse() (response *DescribeSQLLogRecordsResponse) {
	response = &DescribeSQLLogRecordsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
