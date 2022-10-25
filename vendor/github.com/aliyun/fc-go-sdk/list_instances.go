package fc

import (
	"fmt"
	"net/http"
	"net/url"
)

type Instance struct {
	InstanceID string `json:"instanceId"`
	VersionID  int    `json:"versionId"`
}

type ListInstancesOutput struct {
	Header    http.Header
	Instances []*Instance
}

// ListInstancesInput define publish layer version response
type ListInstancesInput struct {
	Header       http.Header
	ServiceName  *string
	FunctionName *string
	Qualifier    *string
	Limit        *int
	NextToken    *string
}

func NewListInstancesInput(serviceName, functionName string) *ListInstancesInput {
	return &ListInstancesInput{
		ServiceName:  &serviceName,
		FunctionName: &functionName,
	}
}

func (i *ListInstancesInput) WithServiceName(serviceName string) *ListInstancesInput {
	i.ServiceName = &serviceName
	return i
}

func (i *ListInstancesInput) WithFunctionName(functionName string) *ListInstancesInput {
	i.FunctionName = &functionName
	return i
}

func (i *ListInstancesInput) WithQualifier(qualifier string) *ListInstancesInput {
	i.Qualifier = &qualifier
	return i
}

func (i *ListInstancesInput) WithLimit(limit int) *ListInstancesInput {
	i.Limit = &limit
	return i
}

func (i *ListInstancesInput) WithNextToken(nextToken string) *ListInstancesInput {
	i.NextToken = &nextToken
	return i
}

func (i ListInstancesInput) GetQueryParams() url.Values {
	queries := make(url.Values)
	if i.Limit != nil {
		queries.Add("limit", fmt.Sprint(*i.Limit))
	}
	if i.NextToken != nil {
		queries.Add("nextToken", *i.NextToken)
	}
	return queries
}

func (i ListInstancesInput) GetPath() string {
	if i.Qualifier != nil {
		return fmt.Sprintf(listInstancesWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
	} else {
		return fmt.Sprintf(listInstancesSourcesPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName))
	}
}

func (i ListInstancesInput) GetHeaders() Header {
	return make(Header)
}

func (i ListInstancesInput) GetPayload() interface{} {
	return nil
}

func (i ListInstancesInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	return nil
}