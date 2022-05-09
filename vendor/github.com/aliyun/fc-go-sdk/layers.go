package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	layersPath             = "/layers"
	singleLayerPath        = layersPath + "/%s"
	singleLayerVersionPath = singleLayerPath + "/versions/%d"
	layerArnPath           = "/layerarn"
	adminLayersPath        = "/adminlayers/%s/%s/versions/%d"
	AnyRunTime             = "Any"
)

// OutputCodeLocation represents the output code location.
type OutputCodeLocation struct {
	RepositoryType *string `json:"repositoryType"`
	Location       *string `json:"location"`
}

type Layer struct {
	LayerName         string             `json:"layerName"`
	Version           int32              `json:"version"`
	Description       string             `json:"description"`
	Code              OutputCodeLocation `json:"code"`
	CodeSize          int64              `json:"codeSize"`
	CodeChecksum      string             `json:"codeChecksum"`
	CreateTime        string             `json:"createTime"`
	ACL               int32              `json:"acl"`
	CompatibleRuntime []string           `json:"compatibleRuntime"`
	Arn               string             `json:"arn"`
}

// PublishLayerVersionInput defines input to create layer version
type PublishLayerVersionInput struct {
	LayerName         string   `json:"layerName"`
	Description       string   `json:"description"`
	Code              *Code    `json:"code"`
	CompatibleRuntime []string `json:"compatibleRuntime"`
}

func NewPublishLayerVersionInput() *PublishLayerVersionInput {
	return &PublishLayerVersionInput{}
}

func (s *PublishLayerVersionInput) WithLayerName(layerName string) *PublishLayerVersionInput {
	s.LayerName = layerName
	return s
}

func (s *PublishLayerVersionInput) WithDescription(description string) *PublishLayerVersionInput {
	s.Description = description
	return s
}

func (s *PublishLayerVersionInput) WithCompatibleRuntime(compatibleRuntime []string) *PublishLayerVersionInput {
	s.CompatibleRuntime = compatibleRuntime
	return s
}

func (s *PublishLayerVersionInput) WithCode(code *Code) *PublishLayerVersionInput {
	s.Code = code
	return s
}

func (i *PublishLayerVersionInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *PublishLayerVersionInput) GetPath() string {
	return layersPath + "/" + i.LayerName + "/versions"
}

func (i *PublishLayerVersionInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *PublishLayerVersionInput) GetPayload() interface{} {
	return i
}

func (i *PublishLayerVersionInput) Validate() error {
	return nil
}

type PublishPublicLayerVersionOutput struct {
	Header http.Header
}

// PublishLayerVersionOutput define publish layer version response
type PublishLayerVersionOutput struct {
	Header http.Header
	Layer
}

func (o PublishLayerVersionInput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o PublishLayerVersionOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o PublishLayerVersionOutput) GetEtag() string {
	return GetEtag(o.Header)
}

// GetLayerVersionOutput define get layer version response
type GetLayerVersionOutput struct {
	Header http.Header
	Layer
}

func (o GetLayerVersionOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetLayerVersionOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

func (o GetLayerVersionOutput) GetEtag() string {
	return GetEtag(o.Header)
}

// ListLayersOutput defines List layers result
type ListLayersOutput struct {
	Header    http.Header
	Layers    []*Layer `json:"layers"`
	NextToken *string  `json:"nextToken,omitempty"`
}

func (o ListLayersOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListLayersOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type ListLayersInput struct {
	Query
	public bool
}

func NewListLayersInput() *ListLayersInput {
	return &ListLayersInput{}
}

func (i *ListLayersInput) WithPrefix(prefix string) *ListLayersInput {
	i.Prefix = &prefix
	return i
}

func (i *ListLayersInput) WithStartKey(startKey string) *ListLayersInput {
	i.StartKey = &startKey
	return i
}

func (i *ListLayersInput) WithNextToken(nextToken string) *ListLayersInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListLayersInput) WithLimit(limit int32) *ListLayersInput {
	i.Limit = &limit
	return i
}

func (i *ListLayersInput) WithPublic(public bool) *ListLayersInput {
	i.public = public
	return i
}

func (i *ListLayersInput) GetQueryParams() url.Values {
	out := url.Values{}
	if i.Prefix != nil {
		out.Set("prefix", *i.Prefix)
	}

	if i.StartKey != nil {
		out.Set("startKey", *i.StartKey)
	}

	if i.NextToken != nil {
		out.Set("nextToken", *i.NextToken)
	}

	if i.Limit != nil {
		out.Set("limit", strconv.FormatInt(int64(*i.Limit), 10))
	}

	if i.public {
		out.Set("public", "true")
	}

	if i.Tags != nil {
		for k, v := range i.Tags {
			out.Set(tagQueryPrefix+k, v)
		}
	}

	return out
}

func (i *ListLayersInput) GetPath() string {
	return layersPath
}

func (i *ListLayersInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListLayersInput) GetPayload() interface{} {
	return nil
}

func (i *ListLayersInput) Validate() error {
	return nil
}

type GetLayerVersionInput struct {
	LayerName string
	Version   int32
}

func NewGetLayerVersionInput(layerName string, version int32) *GetLayerVersionInput {
	return &GetLayerVersionInput{LayerName: layerName, Version: version}
}

func (i *GetLayerVersionInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *GetLayerVersionInput) GetPath() string {
	return fmt.Sprintf(singleLayerVersionPath, pathEscape(i.LayerName), i.Version)
}

func (i *GetLayerVersionInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *GetLayerVersionInput) GetPayload() interface{} {
	return nil
}

func (i *GetLayerVersionInput) Validate() error {
	if i.Version <= 0 {
		return fmt.Errorf("Version must be a positive number.")
	}

	return nil
}

type GetLayerVersionByArnInput struct {
	Arn string
}

func NewGetLayerVersionByArnInput(arn string) *GetLayerVersionByArnInput {
	return &GetLayerVersionByArnInput{Arn: arn}
}

func (i *GetLayerVersionByArnInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *GetLayerVersionByArnInput) GetPath() string {
	return layerArnPath + "/" + pathEscape(i.Arn)
}

func (i *GetLayerVersionByArnInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *GetLayerVersionByArnInput) GetPayload() interface{} {
	return nil
}

func (i *GetLayerVersionByArnInput) Validate() error {
	if len(i.Arn) == 0 {
		return fmt.Errorf("Arm should not be empty.")
	}

	return nil
}

type DeleteLayerVersionInput struct {
	LayerName string
	Version   int32
}

func NewDeleteLayerVersionInput(layerName string, version int32) *DeleteLayerVersionInput {
	return &DeleteLayerVersionInput{LayerName: layerName, Version: version}
}

func (i *DeleteLayerVersionInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *DeleteLayerVersionInput) GetPath() string {
	return fmt.Sprintf(singleLayerVersionPath, pathEscape(i.LayerName), i.Version)
}

func (i *DeleteLayerVersionInput) GetHeaders() Header {
	header := make(Header)
	return header
}

func (i *DeleteLayerVersionInput) GetPayload() interface{} {
	return nil
}

func (i *DeleteLayerVersionInput) Validate() error {
	if i.Version <= 0 {
		return fmt.Errorf("Layer version must be a positive number.")
	}
	return nil
}

type DeleteLayerVersionOutput struct {
	Header http.Header
}

func (o DeleteLayerVersionOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o DeleteLayerVersionOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type ListLayerVersionsInput struct {
	LayerName    string
	StartVersion int32
	Limit        int32
}

func NewListLayerVersionsInput(layerName string, version int32) *ListLayerVersionsInput {
	return &ListLayerVersionsInput{LayerName: layerName, StartVersion: version, Limit: 20}
}

func (i *ListLayerVersionsInput) WithLimit(limit int32) *ListLayerVersionsInput {
	i.Limit = limit
	return i
}

func (i *ListLayerVersionsInput) GetQueryParams() url.Values {
	out := url.Values{}
	out.Set("startVersion", strconv.FormatInt(int64(i.StartVersion), 10))

	if i.Limit != 0 {
		out.Set("limit", strconv.FormatInt(int64(i.Limit), 10))
	}

	return out
}

func (i *ListLayerVersionsInput) GetPath() string {
	return fmt.Sprintf(singleLayerPath, pathEscape(i.LayerName)) + "/versions"
}

func (i *ListLayerVersionsInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListLayerVersionsInput) GetPayload() interface{} {
	return nil
}

func (i *ListLayerVersionsInput) Validate() error {
	if i.StartVersion <= 0 {
		return fmt.Errorf("Version must be a positive number.")
	}

	return nil
}

type ListLayerVersionsOutput struct {
	Header      http.Header
	Layers      []*Layer `json:"layers"`
	NextVersion *int32   `json:"nextVersion,omitempty"`
}

func (o ListLayerVersionsOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListLayerVersionsOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

type PermanentDeleteLayerVersionInput struct {
	LayerName        string
	Version          int32
	LayerOwnerUserID string
}

func NewPermanentDeleteLayerVersionInput(layerOwnerUserID, layerName string, version int32) *PermanentDeleteLayerVersionInput {
	return &PermanentDeleteLayerVersionInput{LayerOwnerUserID: layerOwnerUserID, LayerName: layerName, Version: version}
}

func (i *PermanentDeleteLayerVersionInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (i *PermanentDeleteLayerVersionInput) GetPath() string {
	return fmt.Sprintf(adminLayersPath, i.LayerOwnerUserID, pathEscape(i.LayerName), i.Version)
}

func (i *PermanentDeleteLayerVersionInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *PermanentDeleteLayerVersionInput) GetPayload() interface{} {
	return nil
}

func (i *PermanentDeleteLayerVersionInput) Validate() error {
	if i.Version <= 0 {
		return fmt.Errorf("Version must be a positive number")
	}

	if len(i.LayerOwnerUserID) <= 0 {
		return fmt.Errorf("LayerOwnerUserID must be specified")
	}

	return nil
}
