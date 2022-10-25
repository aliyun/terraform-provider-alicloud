package fc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/resty.v1"
	"github.com/gorilla/websocket"
)

// Client defines fc client
type Client struct {
	Config  *Config
	Connect *Connection
}

// NewClient new fc client
func NewClient(endpoint, apiVersion, accessKeyID, accessKeySecret string, opts ...ClientOption) (*Client, error) {
	config := NewConfig()
	config.APIVersion = apiVersion
	config.AccessKeyID = accessKeyID
	config.AccessKeySecret = accessKeySecret
	config.Endpoint, config.host = GetAccessPoint(endpoint)
	connect := NewConnection()
	client := &Client{config, connect}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// GetAccountSettings returns account settings from fc
func (c *Client) GetAccountSettings(input *GetAccountSettingsInput) (*GetAccountSettingsOutput, error) {
	if input == nil {
		input = new(GetAccountSettingsInput)
	}

	var output = new(GetAccountSettingsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetService returns service metadata from fc
func (c *Client) GetService(input *GetServiceInput) (*GetServiceOutput, error) {
	if input == nil {
		input = new(GetServiceInput)
	}

	var output = new(GetServiceOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListServices returns list of services from fc
func (c *Client) ListServices(input *ListServicesInput) (*ListServicesOutput, error) {
	if input == nil {
		input = new(ListServicesInput)
	}

	var output = new(ListServicesOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UpdateService updates service
func (c *Client) UpdateService(input *UpdateServiceInput) (*UpdateServiceOutput, error) {
	if input == nil {
		input = new(UpdateServiceInput)
	}

	var output = new(UpdateServiceOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// CreateService creates service
func (c *Client) CreateService(input *CreateServiceInput) (*CreateServiceOutput, error) {
	if input == nil {
		input = new(CreateServiceInput)
	}

	var output = new(CreateServiceOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteService deletes service
func (c *Client) DeleteService(input *DeleteServiceInput) (*DeleteServiceOutput, error) {
	if input == nil {
		input = new(DeleteServiceInput)
	}
	var output = new(DeleteServiceOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// PublishServiceVersion publishes service version
func (c *Client) PublishServiceVersion(input *PublishServiceVersionInput) (*PublishServiceVersionOutput, error) {
	if input == nil {
		input = new(PublishServiceVersionInput)
	}
	var output = new(PublishServiceVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListServiceVersions returns list of service versions
func (c *Client) ListServiceVersions(input *ListServiceVersionsInput) (*ListServiceVersionsOutput, error) {
	if input == nil {
		input = new(ListServiceVersionsInput)
	}

	var output = new(ListServiceVersionsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteServiceVersion marks service version as deleted
func (c *Client) DeleteServiceVersion(input *DeleteServiceVersionInput) (*DeleteServiceVersionOutput, error) {
	if input == nil {
		input = new(DeleteServiceVersionInput)
	}
	var output = new(DeleteServiceVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// CreateAlias creates alias
func (c *Client) CreateAlias(input *CreateAliasInput) (*CreateAliasOutput, error) {
	if input == nil {
		input = new(CreateAliasInput)
	}

	var output = new(CreateAliasOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UpdateAlias updates alias
func (c *Client) UpdateAlias(input *UpdateAliasInput) (*UpdateAliasOutput, error) {
	if input == nil {
		input = new(UpdateAliasInput)
	}

	var output = new(UpdateAliasOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetAlias returns alias metadata from fc
func (c *Client) GetAlias(input *GetAliasInput) (*GetAliasOutput, error) {
	if input == nil {
		input = new(GetAliasInput)
	}

	var output = new(GetAliasOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListAliases returns list of aliases from fc
func (c *Client) ListAliases(input *ListAliasesInput) (*ListAliasesOutput, error) {
	if input == nil {
		input = new(ListAliasesInput)
	}

	var output = new(ListAliasesOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteAlias deletes service
func (c *Client) DeleteAlias(input *DeleteAliasInput) (*DeleteAliasOutput, error) {
	if input == nil {
		input = new(DeleteAliasInput)
	}
	var output = new(DeleteAliasOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// CreateFunction creates function
func (c *Client) CreateFunction(input *CreateFunctionInput) (*CreateFunctionOutput, error) {
	if input == nil {
		input = new(CreateFunctionInput)
	}
	var output = new(CreateFunctionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteFunction deletes function from service
func (c *Client) DeleteFunction(input *DeleteFunctionInput) (*DeleteFunctionOutput, error) {
	if input == nil {
		input = new(DeleteFunctionInput)
	}

	var output = new(DeleteFunctionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	return output, nil
}

// GetFunction returns function metadata from service
func (c *Client) GetFunction(input *GetFunctionInput) (*GetFunctionOutput, error) {
	if input == nil {
		input = new(GetFunctionInput)
	}

	var output = new(GetFunctionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetFunctionCode returns function code
func (c *Client) GetFunctionCode(input *GetFunctionCodeInput) (*GetFunctionCodeOutput, error) {
	if input == nil {
		input = new(GetFunctionCodeInput)
	}

	var output = new(GetFunctionCodeOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListFunctions returns list of functions
func (c *Client) ListFunctions(input *ListFunctionsInput) (*ListFunctionsOutput, error) {
	if input == nil {
		input = new(ListFunctionsInput)
	}

	var output = new(ListFunctionsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UpdateFunction updates function
func (c *Client) UpdateFunction(input *UpdateFunctionInput) (*UpdateFunctionOutput, error) {
	if input == nil {
		input = new(UpdateFunctionInput)
	}

	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	var output = new(UpdateFunctionOutput)
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// CreateTrigger creates trigger
func (c *Client) CreateTrigger(input *CreateTriggerInput) (*CreateTriggerOutput, error) {
	if input == nil {
		input = new(CreateTriggerInput)
	}

	var output = new(CreateTriggerOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetTrigger returns trigger metadata
func (c *Client) GetTrigger(input *GetTriggerInput) (*GetTriggerOutput, error) {
	if input == nil {
		input = new(GetTriggerInput)
	}

	var output = new(GetTriggerOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UpdateTrigger updates trigger
func (c *Client) UpdateTrigger(input *UpdateTriggerInput) (*UpdateTriggerOutput, error) {
	if input == nil {
		input = new(UpdateTriggerInput)
	}

	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	var output = new(UpdateTriggerOutput)
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteTrigger deletes trigger
func (c *Client) DeleteTrigger(input *DeleteTriggerInput) (*DeleteTriggerOutput, error) {
	if input == nil {
		input = new(DeleteTriggerInput)
	}

	var output = new(DeleteTriggerOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	return output, nil
}

// ListTriggers returns list of triggers
func (c *Client) ListTriggers(input *ListTriggersInput) (*ListTriggersOutput, error) {
	if input == nil {
		input = new(ListTriggersInput)
	}

	var output = new(ListTriggersOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// TagResource make a resource with tags
func (c *Client) TagResource(input *TagResourceInput) (*TagResourceOut, error) {
	if input == nil {
		input = new(TagResourceInput)
	}

	var output = new(TagResourceOut)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetResourceTags ...
func (c *Client) GetResourceTags(input *GetResourceTagsInput) (*GetResourceTagsOut, error) {
	if input == nil {
		input = new(GetResourceTagsInput)
	}

	var output = new(GetResourceTagsOut)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UnTagResource ...
func (c *Client) UnTagResource(input *UnTagResourceInput) (*UnTagResourceOut, error) {
	if input == nil {
		input = new(UnTagResourceInput)
	}

	var output = new(UnTagResourceOut)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// PutProvisionConfig put provision config
func (c *Client) PutProvisionConfig(input *PutProvisionConfigInput) (*PutProvisionConfigOutput, error) {
	if input == nil {
		input = new(PutProvisionConfigInput)
	}

	var output = new(PutProvisionConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetProvisionConfig return provision config from fc
func (c *Client) GetProvisionConfig(input *GetProvisionConfigInput) (*GetProvisionConfigOutput, error) {
	if input == nil {
		input = new(GetProvisionConfigInput)
	}

	var output = new(GetProvisionConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListProvisionConfigs return list of provision configs from fc
func (c *Client) ListProvisionConfigs(input *ListProvisionConfigsInput) (*ListProvisionConfigsOutput, error) {
	if input == nil {
		input = new(ListProvisionConfigsInput)
	}

	var output = new(ListProvisionConfigsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// InvokeFunction : invoke function in fc
func (c *Client) InvokeFunction(input *InvokeFunctionInput) (*InvokeFunctionOutput, error) {
	if input == nil {
		input = new(InvokeFunctionInput)
	}

	var output = new(InvokeFunctionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	output.Payload = httpResponse.Body()

	return output, nil
}

// ListReservedCapacities returns list of reserved capacity from fc
func (c *Client) ListReservedCapacities(input *ListReservedCapacitiesInput) (*ListReservedCapacitiesOutput, error) {
	if input == nil {
		input = new(ListReservedCapacitiesInput)
	}

	var output = new(ListReservedCapacitiesOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// CreateCustomDomain creates custom domain
func (c *Client) CreateCustomDomain(input *CreateCustomDomainInput) (*CreateCustomDomainOutput, error) {
	if input == nil {
		input = new(CreateCustomDomainInput)
	}

	var output = new(CreateCustomDomainOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// UpdateCustomDomain updates custom domain
func (c *Client) UpdateCustomDomain(input *UpdateCustomDomainInput) (*UpdateCustomDomainOutput, error) {
	if input == nil {
		input = new(UpdateCustomDomainInput)
	}

	var output = new(UpdateCustomDomainOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetCustomDomain returns custom domain metadata from fc
func (c *Client) GetCustomDomain(input *GetCustomDomainInput) (*GetCustomDomainOutput, error) {
	if input == nil {
		input = new(GetCustomDomainInput)
	}

	var output = new(GetCustomDomainOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteCustomDomain deletes custom domain
func (c *Client) DeleteCustomDomain(input *DeleteCustomDomainInput) (*DeleteCustomDomainOutput, error) {
	if input == nil {
		input = new(DeleteCustomDomainInput)
	}
	var output = new(DeleteCustomDomainOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// ListCustomDomains returns list of custom domains from fc
func (c *Client) ListCustomDomains(input *ListCustomDomainsInput) (*ListCustomDomainsOutput, error) {
	if input == nil {
		input = new(ListCustomDomainsInput)
	}

	var output = new(ListCustomDomainsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

func (c *Client) sendRequest(input ServiceInput, httpMethod string) (*resty.Response, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	var serviceError = new(ServiceError)
	path := "/" + c.Config.APIVersion + input.GetPath()

	headerParams := make(map[string]string)
	for k, v := range input.GetHeaders() {
		headerParams[k] = v
	}
	headerParams["Host"] = c.Config.host
	headerParams[HTTPHeaderAccountID] = c.Config.AccountID
	headerParams[HTTPHeaderUserAgent] = c.Config.UserAgent
	headerParams["Accept"] = "application/json"
	// Caution: should not declare this as byte[] whose zero value is an empty byte array
	// if input has no payload, the http body should not be populated at all.
	var rawBody interface{}
	if input.GetPayload() != nil {
		switch input.GetPayload().(type) {
		case *[]byte:
			headerParams["Content-Type"] = "application/octet-stream"
			b := input.GetPayload().(*[]byte)
			headerParams["Content-MD5"] = MD5(*b)
			rawBody = *b
		default:
			headerParams["Content-Type"] = "application/json"
			b, err := json.Marshal(input.GetPayload())
			if err != nil {
				// TODO: return client side error
				return nil, nil
			}
			headerParams["Content-MD5"] = MD5(b)
			rawBody = b
		}
	}
	headerParams["Date"] = time.Now().UTC().Format(http.TimeFormat)
	if c.Config.SecurityToken != "" {
		headerParams[HTTPHeaderSecurityToken] = c.Config.SecurityToken
	}
	headerParams["Authorization"] = GetAuthStr(c.Config.AccessKeyID, c.Config.AccessKeySecret, httpMethod, headerParams, path)
	resp, err := c.Connect.SendRequest(c.Config.Endpoint+path, httpMethod, rawBody, headerParams, input.GetQueryParams())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 300 {
		serviceError.RequestID = resp.Header().Get(HTTPHeaderRequestID)
		serviceError.HTTPStatus = resp.StatusCode()
		json.Unmarshal(resp.Body(), &serviceError)
		return nil, serviceError
	}
	return resp, nil
}

// GetFunctionAsyncInvokeConfig returns async config from fc
func (c *Client) GetFunctionAsyncInvokeConfig(input *GetFunctionAsyncInvokeConfigInput) (*GetFunctionAsyncInvokeConfigOutput, error) {
	if input == nil {
		input = new(GetFunctionAsyncInvokeConfigInput)
	}

	var output = new(GetFunctionAsyncInvokeConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListFunctionAsyncInvokeConfigs returns list of async configs from fc
func (c *Client) ListFunctionAsyncInvokeConfigs(input *ListFunctionAsyncInvokeConfigsInput) (*ListFunctionAsyncInvokeConfigsOutput, error) {
	if input == nil {
		input = new(ListFunctionAsyncInvokeConfigsInput)
	}

	var output = new(ListFunctionAsyncInvokeConfigsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// PutFunctionAsyncInvokeConfig creates or updates an async config
func (c *Client) PutFunctionAsyncInvokeConfig(input *PutFunctionAsyncInvokeConfigInput) (*PutFunctionAsyncInvokeConfigOutput, error) {
	if input == nil {
		input = new(PutFunctionAsyncInvokeConfigInput)
	}

	var output = new(PutFunctionAsyncInvokeConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DeleteFunctionAsyncInvokeConfig deletes async config
func (c *Client) DeleteFunctionAsyncInvokeConfig(input *DeleteFunctionAsyncInvokeConfigInput) (*DeleteFunctionAsyncInvokeConfigOutput, error) {
	if input == nil {
		input = new(DeleteFunctionAsyncInvokeConfigInput)
	}
	var output = new(DeleteFunctionAsyncInvokeConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// ListLayers returns list of layers from fc
func (c *Client) ListLayers(input *ListLayersInput) (*ListLayersOutput, error) {
	if input == nil {
		input = new(ListLayersInput)
	}

	var output = new(ListLayersOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListLayerVersions returns list of layer versions of a specific layer from fc
func (c *Client) ListLayerVersions(input *ListLayerVersionsInput) (*ListLayerVersionsOutput, error) {
	if input == nil {
		input = new(ListLayerVersionsInput)
	}

	var output = new(ListLayerVersionsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetLayerVersion returns layer version information from fc
func (c *Client) GetLayerVersion(input *GetLayerVersionInput) (*GetLayerVersionOutput, error) {
	if input == nil {
		input = new(GetLayerVersionInput)
	}

	var output = new(GetLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// GetLayerVersionByArn returns layer version information from fc
func (c *Client) GetLayerVersionByArn(input *GetLayerVersionByArnInput) (*GetLayerVersionOutput, error) {
	if input == nil {
		input = new(GetLayerVersionByArnInput)
	}

	var output = new(GetLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// PublishLayerVersion creates a new layer version
func (c *Client) PublishLayerVersion(input *PublishLayerVersionInput) (*PublishLayerVersionOutput, error) {
	if input == nil {
		input = new(PublishLayerVersionInput)
	}

	var output = new(PublishLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// PublishPublicLayerVersion publish a new exiting layer version as public
func (c *Client) PublishPublicLayerVersion(input *GetLayerVersionInput) (*PublishPublicLayerVersionOutput, error) {
	if input == nil {
		input = new(GetLayerVersionInput)
	}

	var output = new(PublishPublicLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// PermanentDeleteVersion delete a soft deleted layer version permanently
func (c *Client) PermanentDeleteLayerVersion(input *PermanentDeleteLayerVersionInput) (*PublishPublicLayerVersionOutput, error) {
	if input == nil {
		return nil, nil
	}

	var output = new(PublishPublicLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPost)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	return output, nil
}

// DeleteLayerVersion deletes a layer version
func (c *Client) DeleteLayerVersion(input *DeleteLayerVersionInput) (*DeleteLayerVersionOutput, error) {
	if input == nil {
		input = new(DeleteLayerVersionInput)
	}

	var output = new(DeleteLayerVersionOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	return output, nil
}

// GetStatefulAsyncInvocation returns stateful async invocation record
func (c *Client) GetStatefulAsyncInvocation(input *GetStatefulAsyncInvocationInput) (*GetStatefulAsyncInvocationOutput, error) {
	if input == nil {
		input = new(GetStatefulAsyncInvocationInput)
	}

	var output = new(GetStatefulAsyncInvocationOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// ListStatefulAsyncInvocations returns list of stateful async invocation records
func (c *Client) ListStatefulAsyncInvocations(input *ListStatefulAsyncInvocationsInput) (*ListStatefulAsyncInvocationsOutput, error) {
	if input == nil {
		input = new(ListStatefulAsyncInvocationsInput)
	}

	var output = new(ListStatefulAsyncInvocationsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// StopStatefulAsyncInvocation ...
func (c *Client) StopStatefulAsyncInvocation(input *StopStatefulAsyncInvocationInput) (*StopStatefulAsyncInvocationOutput, error) {
	if input == nil {
		input = new(StopStatefulAsyncInvocationInput)
	}

	var output = new(StopStatefulAsyncInvocationOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}
	output.Header = httpResponse.Header()
	json.Unmarshal(httpResponse.Body(), output)
	return output, nil
}

// DoHttpRequest returns function http invocation response
func (c *Client) DoHttpRequest(req *http.Request) (*http.Response, error) {
	headerParams := make(map[string]string)
	if req.Header != nil {
		for k, _ := range req.Header {
			headerParams[k] = req.Header.Get(k)
		}
	}
	// CONTENT-MD5
	if req.Body != nil {
		buf, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewReader(buf))
		b, err := json.Marshal(buf)
		if err != nil {
			return nil, err
		}
		headerParams[HTTPHeaderContentMD5] = MD5(b)
	}
	// CONTENT-TYPE
	headerParams[HTTPHeaderContentType] = req.Header.Get(HTTPHeaderContentType)
	// DATE
	headerParams[HTTPHeaderDate] = time.Now().UTC().Format(http.TimeFormat)
	// Canonicalized
	canonicalizedResource := req.URL.Path
	params := req.URL.Query()
	canonicalizedResource = GetSignResourceWithQueries(req.URL.Path, params)
	// Build Authorization header
	headerParams["Authorization"] = GetAuthStr(c.Config.AccessKeyID, c.Config.AccessKeySecret, req.Method, headerParams, canonicalizedResource)
	// Prepare and send request.
	preparedRequest := c.Connect.PrepareRequest(req.Body, headerParams, params).SetDoNotParseResponse(true)
	resp, err := preparedRequest.Execute(req.Method, c.Config.Endpoint+req.URL.Path)
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// SignURL : sign an URL with signature in queries for HTTP function
func (c *Client) SignURL(signURLInput *SignURLInput) (string, error) {
	conf := c.Config
	return signURLInput.signURL(conf.APIVersion, conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret, conf.SecurityToken)
}

// ListOnDemandConfigs return list of provision configs from fc
func (c *Client) ListOnDemandConfigs(input *ListOnDemandConfigsInput) (*ListOnDemandConfigsOutput, error) {
	if input == nil {
		input = NewListOnDemandConfigsInput()
	}

	var output = new(ListOnDemandConfigsOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	if err := json.Unmarshal(httpResponse.Body(), output); err != nil {
		return nil, err
	}
	return output, nil
}

// PutOnDemandConfig put on-demand config
func (c *Client) PutOnDemandConfig(input *PutOnDemandConfigInput) (*PutOnDemandConfigOutput, error) {
	if input == nil {
		input = new(PutOnDemandConfigInput)
	}

	var output = new(PutOnDemandConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodPut)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	if err := json.Unmarshal(httpResponse.Body(), output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetOnDemandConfig return on-demand config from fc
func (c *Client) GetOnDemandConfig(input *GetOnDemandConfigInput) (*GetOnDemandConfigOutput, error) {
	if input == nil {
		input = new(GetOnDemandConfigInput)
	}

	var output = new(GetOnDemandConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	if err := json.Unmarshal(httpResponse.Body(), output); err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteOnDemandConfig delete on-demand config
func (c *Client) DeleteOnDemandConfig(input *DeleteOnDemandConfigInput) (*DeleteOnDemandConfigOutput, error) {
	if input == nil {
		input = new(DeleteOnDemandConfigInput)
	}

	var output = new(DeleteOnDemandConfigOutput)
	httpResponse, err := c.sendRequest(input, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	output.Header = httpResponse.Header()
	return output, nil
}

// ListInstances ...
func (c *Client) ListInstances(input *ListInstancesInput) (*ListInstancesOutput, error) {
	if input == nil {
		input = new(ListInstancesInput)
	}

	var output = new(ListInstancesOutput)
	httpResponse, err := c.sendRequest(input, http.MethodGet)
	if err != nil {
		return nil, err
	}
	data := httpResponse.Body()
	fmt.Printf("%s\n", data)
	json.Unmarshal(data, output)
	output.Header = httpResponse.Header()
	return output, nil
}

// InstanceExec ...
func (c *Client) InstanceExec(input *InstanceExecInput) (*InstanceExecOutput, error) {
	if input == nil {
		input = new(InstanceExecInput)
	}

	var output = new(InstanceExecOutput)
	ws, err := c.openWebSocketConn(input)
	if err != nil {
		return nil, err
	}
	output.WebsocketConnection = ws
	output.start(input)
	return output, nil
}

// buildWebSocket ...
func (c *Client) openWebSocketConn(input ServiceInput) (*websocket.Conn, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	path := "/" + c.Config.APIVersion + input.GetPath()

	headerParams := make(map[string]string)
	for k, v := range input.GetHeaders() {
		headerParams[k] = v
	}

	headerParams["Host"] = c.Config.host
	headerParams[HTTPHeaderAccountID] = c.Config.AccountID
	headerParams[HTTPHeaderUserAgent] = c.Config.UserAgent
	headerParams["Accept"] = "application/json"
	// Caution: should not declare this as byte[] whose zero value is an empty byte array
	// if input has no payload, the http body should not be populated at all.
	if input.GetPayload() != nil {
		switch input.GetPayload().(type) {
		case *[]byte:
			headerParams["Content-Type"] = "application/octet-stream"
			b := input.GetPayload().(*[]byte)
			headerParams["Content-MD5"] = MD5(*b)
		default:
			headerParams["Content-Type"] = "application/json"
			b, err := json.Marshal(input.GetPayload())
			if err != nil {
				// TODO: return client side error
				return nil, nil
			}
			headerParams["Content-MD5"] = MD5(b)
		}
	}
	headerParams["Date"] = time.Now().UTC().Format(http.TimeFormat)
	if c.Config.SecurityToken != "" {
		headerParams[HTTPHeaderSecurityToken] = c.Config.SecurityToken
	}
	switch c.Config.APIVersion {
	case APIVersionV1:
		headerParams["Authorization"] = GetAuthStr(
			c.Config.AccessKeyID, c.Config.AccessKeySecret, http.MethodGet, headerParams, path)
	default:
		return nil, fmt.Errorf("unsupported api version: '%s'", c.Config.APIVersion)
	}

	u := &url.URL{Scheme: "ws", Host: c.Config.host, Path: path, RawQuery: input.GetQueryParams().Encode()}
	if strings.HasPrefix(c.Config.host, "https://") {
		u.Scheme = "wss"
	}
	header := make(http.Header)
	for headerKey, headerValue := range headerParams {
		if headerKey == "Connection" || headerKey == "Upgrade" || headerKey == "Sec-Websocket-Version" {
			continue
		}
		header.Set(headerKey, headerValue)
	}

	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		if resp != nil {
			content, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("%v: %s", err, content)
		}
		return nil, err
	}

	return ws, nil
}

