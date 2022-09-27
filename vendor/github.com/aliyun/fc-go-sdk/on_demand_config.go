package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DefaultListLimit = 100
)

type OnDemandConfig struct {
	Resource             *string `json:"resource"`
	MaximumInstanceCount *int64  `json:"maximumInstanceCount"`
}

type OnDemandConfigs struct {
	Configs   []*OnDemandConfig `json:"configs"`
	NextToken *string           `json:"nextToken,omitempty"`
}

type ListOnDemandConfigsInput struct {
	Prefix    *string `json:"prefix"`
	StartKey  *string `json:"startKey"`
	NextToken *string `json:"nextToken,omitempty"`
	Limit     int     `json:"limit,omitempty"`
}

type ListOnDemandConfigsOutput struct {
	Header http.Header
	OnDemandConfigs
}

func NewListOnDemandConfigsInput() *ListOnDemandConfigsInput {
	emptyStr := ""
	return &ListOnDemandConfigsInput{
		Prefix:    &emptyStr,
		StartKey:  &emptyStr,
		NextToken: &emptyStr,
		Limit:     DefaultListLimit,
	}
}

func (i *ListOnDemandConfigsInput) WithPrefix(prefix string) *ListOnDemandConfigsInput {
	i.Prefix = &prefix
	return i
}

func (i *ListOnDemandConfigsInput) WithStartKey(startKey string) *ListOnDemandConfigsInput {
	i.StartKey = &startKey
	return i
}

func (i *ListOnDemandConfigsInput) WithNextToken(nextToken string) *ListOnDemandConfigsInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListOnDemandConfigsInput) WithLimit(limit int) *ListOnDemandConfigsInput {
	i.Limit = limit
	return i
}

func (i *ListOnDemandConfigsInput) GetPath() string {
	return onDemandConfigPath
}

func (i *ListOnDemandConfigsInput) GetHeaders() Header {
	return make(Header)
}

func (i *ListOnDemandConfigsInput) GetPayload() interface{} {
	return nil
}

func (i *ListOnDemandConfigsInput) GetQueryParams() url.Values {
	return url.Values{
		"prefix":    []string{*i.Prefix},
		"startKey":  []string{*i.StartKey},
		"nextToken": []string{*i.NextToken},
		"limit":     []string{fmt.Sprintf("%d", i.Limit)},
	}
}

func (i *ListOnDemandConfigsInput) Validate() error {
	return nil
}

func (o ListOnDemandConfigsOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListOnDemandConfigsOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type PutOnDemandConfigObject struct {
	MaximumInstanceCount *int64 `json:"maximumInstanceCount"`
}

type PutOnDemandConfigInput struct {
	PutOnDemandConfigObject
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	IfMatch      *string
}

func NewPutOnDemandConfigInput(serviceName, qualifier, functionName string) *PutOnDemandConfigInput {
	return &PutOnDemandConfigInput{
		ServiceName:  &serviceName,
		Qualifier:    &qualifier,
		FunctionName: &functionName,
	}
}

func (i *PutOnDemandConfigInput) WithMaximumInstanceCount(maximumInstanceCount int64) *PutOnDemandConfigInput {
	i.MaximumInstanceCount = &maximumInstanceCount
	return i
}

func (i *PutOnDemandConfigInput) WithIfMatch(ifMatch string) *PutOnDemandConfigInput {
	i.IfMatch = &ifMatch
	return i
}

func (i *PutOnDemandConfigInput) GetPath() string {
	return fmt.Sprintf(onDemandConfigWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
}

func (i *PutOnDemandConfigInput) GetHeaders() Header {
	header := make(Header)
	if i.IfMatch != nil {
		header[ifMatch] = *i.IfMatch
	}
	return header
}

func (i *PutOnDemandConfigInput) GetPayload() interface{} {
	return i.PutOnDemandConfigObject
}

func (i *PutOnDemandConfigInput) GetQueryParams() url.Values {
	return url.Values{}
}

func (i *PutOnDemandConfigInput) Validate() error {
	return nil
}

type PutOnDemandConfigOutput struct {
	Header http.Header
	OnDemandConfig
}

func (o PutOnDemandConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o PutOnDemandConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o PutOnDemandConfigOutput) GetEtag() string {
	return GetEtag(o.Header)
}

type GetOnDemandConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
}

func NewGetOnDemandConfigInput(serviceName, qualifier, functionName string) *GetOnDemandConfigInput {
	return &GetOnDemandConfigInput{
		ServiceName:  &serviceName,
		Qualifier:    &qualifier,
		FunctionName: &functionName,
	}
}

func (i *GetOnDemandConfigInput) GetPath() string {
	return fmt.Sprintf(onDemandConfigWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
}

func (i *GetOnDemandConfigInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *GetOnDemandConfigInput) GetPayload() interface{} {
	return nil
}

func (i *GetOnDemandConfigInput) GetQueryParams() url.Values {
	return url.Values{}
}

func (i *GetOnDemandConfigInput) Validate() error {
	return nil
}

type GetOnDemandConfigOutput struct {
	Header http.Header
	OnDemandConfig
}

func (o GetOnDemandConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetOnDemandConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o GetOnDemandConfigOutput) GetEtag() string {
	return GetEtag(o.Header)
}

type DeleteOnDemandConfigInput struct {
	ServiceName  *string `json:"serviceName"`
	Qualifier    *string `json:"qualifier"`
	FunctionName *string `json:"functionName"`
	IfMatch      *string
}

func NewDeleteOnDemandConfigInput(serviceName, qualifier, functionName string) *DeleteOnDemandConfigInput {
	return &DeleteOnDemandConfigInput{
		ServiceName:  &serviceName,
		Qualifier:    &qualifier,
		FunctionName: &functionName,
	}
}

func (s *DeleteOnDemandConfigInput) WithIfMatch(ifMatch string) *DeleteOnDemandConfigInput {
	s.IfMatch = &ifMatch
	return s
}

func (i *DeleteOnDemandConfigInput) GetPath() string {
	return fmt.Sprintf(onDemandConfigWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName))
}

func (i *DeleteOnDemandConfigInput) GetHeaders() Header {
	header := make(Header)
	if i.IfMatch != nil {
		header[ifMatch] = *i.IfMatch
	}
	return header
}

func (i *DeleteOnDemandConfigInput) GetPayload() interface{} {
	return nil
}

func (i *DeleteOnDemandConfigInput) GetQueryParams() url.Values {
	return url.Values{}
}

func (i *DeleteOnDemandConfigInput) Validate() error {
	return nil
}

type DeleteOnDemandConfigOutput struct {
	Header http.Header
}

func (o DeleteOnDemandConfigOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o DeleteOnDemandConfigOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}
