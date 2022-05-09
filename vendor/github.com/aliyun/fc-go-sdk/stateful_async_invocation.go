package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// StatefulAsyncInvocation represents the stateful async invocation records.
type StatefulAsyncInvocation struct {
	FunctionName *string `json:"functionName"`
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	InvocationID *string `json:"invocationId"`
	Status       *string `json:"status"`
	// in ms
	StartedTime            *int64  `json:"startedTime"`
	EndTime                *int64  `json:"endTime"`
	DestinationStatus      *string `json:"destinationStatus"`
	InvocationErrorMessage *string `json:"invocationErrorMessage"`
	InvocationPayload      *string `json:"invocationPayload"`
	RequestID              *string `json:"requestId"`
	AlreadyRetriedTimes    *int64  `json:"alreadyRetriedTimes"`
}

type StatefulAsyncInvocationResponse struct {
	StatefulAsyncInvocation
}

// GetFunctionAsyncInvokeConfigInput defines function creation input
type GetStatefulAsyncInvocationInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	InvocationID *string `json:"invocationId"`
}

func NewGetStatefulAsyncInvocationInput(serviceName, functionName, invocationID string) *GetStatefulAsyncInvocationInput {
	return &GetStatefulAsyncInvocationInput{ServiceName: &serviceName, FunctionName: &functionName, InvocationID: &invocationID}
}

func (i *GetStatefulAsyncInvocationInput) WithQualifier(qualifier string) *GetStatefulAsyncInvocationInput {
	i.Qualifier = &qualifier
	return i
}

func (i *GetStatefulAsyncInvocationInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *GetStatefulAsyncInvocationInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(statefulAsyncInvocationWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName), pathEscape(*i.InvocationID))
	}
	return fmt.Sprintf(statefulAsyncInvocationPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName), pathEscape(*i.InvocationID))
}

func (i *GetStatefulAsyncInvocationInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *GetStatefulAsyncInvocationInput) GetPayload() interface{} {
	return nil
}

func (i *GetStatefulAsyncInvocationInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	if IsBlank(i.InvocationID) {
		return fmt.Errorf("InvocationID is required but not provided")
	}
	return nil
}

// GetFunctionAsyncInvokeConfigOutput define get data response
type GetStatefulAsyncInvocationOutput struct {
	Header http.Header
	StatefulAsyncInvocationResponse
}

func (o GetStatefulAsyncInvocationOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetStatefulAsyncInvocationOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

// ListFunctionAsyncInvokeConfigsOutput defines ListFunctionAsyncInvocationsOutput result
type ListStatefulAsyncInvocationsOutput struct {
	Header      http.Header
	Invocations []*StatefulAsyncInvocation `json:"invocations"`
	NextToken   *string                    `json:"nextToken"`
}

func (o ListStatefulAsyncInvocationsOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListStatefulAsyncInvocationsOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type ListStatefulAsyncInvocationsInput struct {
	ServiceName        *string `json:"serviceName"`
	FunctionName       *string `json:"functionName"`
	Qualifier          *string `json:"qualifier"`
	NextToken          *string `json:"nextToken"`
	Limit              *int32  `json:"limit"`
	Status             *string `json:"status"`
	StartedTimeBegin   *int64  `json:"startedTimeBegin"`
	StartedTimeEnd     *int64  `json:"startedTimeEnd"`
	SortOrderByTime    *string `json:"sortOrderByTime"`
	InvocationIDPrefix *string `json:"invocationIdPrefix"`
	IncludePayload     *bool   `json:"includePayload"`
}

func NewListStatefulAsyncInvocationsInput(serviceName, functionName string) *ListStatefulAsyncInvocationsInput {
	return &ListStatefulAsyncInvocationsInput{
		ServiceName:  &serviceName,
		FunctionName: &functionName,
	}
}

func (i *ListStatefulAsyncInvocationsInput) WithQualifier(qualifier string) *ListStatefulAsyncInvocationsInput {
	i.Qualifier = &qualifier
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithNextToken(nextToken string) *ListStatefulAsyncInvocationsInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithLimit(limit int32) *ListStatefulAsyncInvocationsInput {
	i.Limit = &limit
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithStatus(status string) *ListStatefulAsyncInvocationsInput {
	i.Status = &status
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithStartedTimeBegin(startedTimeBegin int64) *ListStatefulAsyncInvocationsInput {
	i.StartedTimeBegin = &startedTimeBegin
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithStartedTimeEnd(startedTimeEnd int64) *ListStatefulAsyncInvocationsInput {
	i.StartedTimeEnd = &startedTimeEnd
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithSortOrderByTime(sortOrder string) *ListStatefulAsyncInvocationsInput {
	i.SortOrderByTime = &sortOrder
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithInvocationIDPrefix(prefix string) *ListStatefulAsyncInvocationsInput {
	i.InvocationIDPrefix = &prefix
	return i
}

func (i *ListStatefulAsyncInvocationsInput) WithIncludePayload(includePayload bool) *ListStatefulAsyncInvocationsInput {
	i.IncludePayload = &includePayload
	return i
}

func (i *ListStatefulAsyncInvocationsInput) GetQueryParams() url.Values {
	out := url.Values{}

	if i.NextToken != nil {
		out.Set("nextToken", *i.NextToken)
	}

	if i.Limit != nil {
		out.Set("limit", strconv.FormatInt(int64(*i.Limit), 10))
	}

	if i.Status != nil {
		out.Set("status", *i.Status)
	}

	if i.StartedTimeBegin != nil {
		out.Set("startedTimeBegin", strconv.FormatInt(int64(*i.StartedTimeBegin), 10))
	}

	if i.StartedTimeEnd != nil {
		out.Set("startedTimeEnd", strconv.FormatInt(int64(*i.StartedTimeEnd), 10))
	}

	if i.SortOrderByTime != nil {
		out.Set("sortOrderByTime", *i.SortOrderByTime)
	}

	if i.InvocationIDPrefix != nil {
		out.Set("invocationIdPrefix", *i.InvocationIDPrefix)
	}

	if i.IncludePayload != nil {
		out.Set("includePayload", strconv.FormatBool(*i.IncludePayload))
	}

	return out
}

func (i *ListStatefulAsyncInvocationsInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(listStatefulAsyncInvocationsWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
	}
	return fmt.Sprintf(listStatefulAsyncInvocationsPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
}

func (i *ListStatefulAsyncInvocationsInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListStatefulAsyncInvocationsInput) GetPayload() interface{} {
	return nil
}

func (i *ListStatefulAsyncInvocationsInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}

// StopStatefulAsyncInvocationInput defines function creation input
type StopStatefulAsyncInvocationInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	InvocationID *string `json:"invocationId"`
}

func NewStopStatefulAsyncInvocationInput(serviceName, functionName, invocationID string) *StopStatefulAsyncInvocationInput {
	return &StopStatefulAsyncInvocationInput{ServiceName: &serviceName, FunctionName: &functionName, InvocationID: &invocationID}
}

func (i *StopStatefulAsyncInvocationInput) WithQualifier(qualifier string) *StopStatefulAsyncInvocationInput {
	i.Qualifier = &qualifier
	return i
}

func (i *StopStatefulAsyncInvocationInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *StopStatefulAsyncInvocationInput) GetPath() string {
	if !IsBlank(i.Qualifier) {
		return fmt.Sprintf(statefulAsyncInvocationWithQualifierPath,
			pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName), pathEscape(*i.InvocationID))
	}
	return fmt.Sprintf(statefulAsyncInvocationPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName), pathEscape(*i.InvocationID))
}

func (i *StopStatefulAsyncInvocationInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *StopStatefulAsyncInvocationInput) GetPayload() interface{} {
	return nil
}

func (i *StopStatefulAsyncInvocationInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	if IsBlank(i.InvocationID) {
		return fmt.Errorf("InvocationID is required but not provided")
	}
	return nil
}

// StopStatefulAsyncInvocationOutput ...
type StopStatefulAsyncInvocationOutput struct {
	Header http.Header
}

func (o StopStatefulAsyncInvocationOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o StopStatefulAsyncInvocationOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}
