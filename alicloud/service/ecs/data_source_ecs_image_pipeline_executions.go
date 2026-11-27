package ecs

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	imagePipelineExecutionsTypeName          = "alicloud_ecs_image_pipeline_executions"
	imagePipelineExecutionsAction            = "DescribeImagePipelineExecutions"
	imagePipelineExecutionsDefaultMaxResults = 50
	imagePipelineExecutionsMaxPages          = 1000
	imagePipelineExecutionsReadBudget        = 30 * time.Minute
)

var imagePipelineExecutionStatuses = []string{
	"PREPARING",
	"REPAIRING",
	"BUILDING",
	"TESTING",
	"DISTRIBUTING",
	"RELEASING",
	"SUCCESS",
	"PARTITION_SUCCESS",
	"FAILED",
	"TEST_FAILED",
	"CANCELLING",
	"CANCELLED",
}

var (
	_ datasource.DataSource                     = (*imagePipelineExecutionsDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*imagePipelineExecutionsDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*imagePipelineExecutionsDataSource)(nil)
)

func NewImagePipelineExecutionsDataSource() datasource.DataSource {
	return &imagePipelineExecutionsDataSource{}
}

type imagePipelineExecutionsDataSource struct {
	fwadapt.DataSourceBase
}

type imagePipelineExecutionsModel struct {
	Id              types.String `tfsdk:"id"`
	ImagePipelineId types.String `tfsdk:"image_pipeline_id"`
	ExecutionId     types.String `tfsdk:"execution_id"`
	Status          types.String `tfsdk:"status"`
	Tags            types.Map    `tfsdk:"tags"`
	Ids             types.List   `tfsdk:"ids"`
	MaxResults      types.Int64  `tfsdk:"max_results"`
	NextToken       types.String `tfsdk:"next_token"`
	TotalCount      types.Int64  `tfsdk:"total_count"`
	Executions      types.List   `tfsdk:"executions"`
}

var imagePipelineExecutionAttrTypes = map[string]attr.Type{
	"execution_id":      types.StringType,
	"image_pipeline_id": types.StringType,
	"status":            types.StringType,
	"creation_time":     types.StringType,
	"modified_time":     types.StringType,
	"resource_group_id": types.StringType,
	"message":           types.StringType,
	"image_id":          types.StringType,
	"error_code":        types.StringType,
	"error_info":        types.StringType,
	"tags":              types.MapType{ElemType: types.StringType},
}

var imagePipelineExecutionObjectType = types.ObjectType{AttrTypes: imagePipelineExecutionAttrTypes}

// imagePipelineExecutionsStatusValidator accepts the comma separated multi
// value form documented for Status and rejects any token outside the enum. An
// explicitly empty value keeps the documented "all statuses" meaning, but a
// non-empty value must not carry empty tokens: dropping them would silently
// widen the query to every status.
type imagePipelineExecutionsStatusValidator struct{}

func (v imagePipelineExecutionsStatusValidator) Description(context.Context) string {
	return "The value must be one status or a comma separated list of statuses from the DescribeImagePipelineExecutions enum."
}

func (v imagePipelineExecutionsStatusValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v imagePipelineExecutionsStatusValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	raw := req.ConfigValue.ValueString()
	if raw == "" {
		// An omitted status is documented to query every status.
		return
	}
	if strings.TrimSpace(raw) == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Status",
			"The status must not be blank; omit it to query every status.",
		)
		return
	}
	for _, token := range strings.Split(raw, ",") {
		token = strings.TrimSpace(token)
		if token == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Status",
				"Every comma separated status must be non-empty; omit the attribute to query every status.",
			)
			return
		}
		valid := false
		for _, status := range imagePipelineExecutionStatuses {
			if token == status {
				valid = true
				break
			}
		}
		if !valid {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Status",
				fmt.Sprintf("%q is not a valid image pipeline execution status. Valid values: %s.", token, strings.Join(imagePipelineExecutionStatuses, ", ")),
			)
			return
		}
	}
}

func (d *imagePipelineExecutionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = imagePipelineExecutionsTypeName
}

func (d *imagePipelineExecutionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// The execution list is a ListAttribute of ObjectType rather than a
	// NestedAttribute: protocol version 5 rejects nested attributes.
	resp.Schema = schema.Schema{
		Description: "Provides a list of image pipeline executions in a region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the matched execution IDs.",
			},
			"image_pipeline_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the image pipeline whose executions are queried. At least one of `image_pipeline_id` and `execution_id` must be set.",
			},
			"execution_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of a specific image pipeline execution. At least one of `image_pipeline_id` and `execution_id` must be set.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "The execution status to filter by. Multiple statuses can be separated by commas. Valid values: `PREPARING`, `REPAIRING`, `BUILDING`, `TESTING`, `DISTRIBUTING`, `RELEASING`, `SUCCESS`, `PARTITION_SUCCESS`, `FAILED`, `TEST_FAILED`, `CANCELLING` and `CANCELLED`.",
				Validators: []validator.String{
					imagePipelineExecutionsStatusValidator{},
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed. A key must be 1 to 128 characters, must not start with `aliyun` or `acs:`, and must not contain a URL. A value may be empty, is limited to 128 characters and must not contain a URL. A null value is rejected: omit the entry to leave the tag out of the request, or use an empty string to send an empty value.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
					apgTagMapValidator{},
				},
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of execution IDs to filter by. The filter is applied on the provider side after all pages are read, so a server that ignores it cannot widen the result. Blank IDs are dropped and repeated IDs are applied once.",
			},
			"max_results": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 500. Defaults to 50. Only the server decides when the result ends: a page shorter than this value is still followed when it returns a continuation token, and no short-page completeness check applies.",
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The pagination token to start reading from. A non-empty configured token is sent in the first request and starts the read at that position. An omitted or empty token is left out of the first request and starts the read at the beginning. Each following request sends the token of the previous response, and the read stops when the server omits the token or returns an empty string. A token that repeats, a whitespace-only token, or 1,000 consecutive requests that each still return a token aborts the read with an error instead of returning a partial result. The value stored afterwards is the last token the server returned, and it is empty when every page has been read.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of execution records this read returned after all configured filters, including the provider-side `ids` filter. `DescribeImagePipelineExecutions` also returns a `TotalCount` for the query, but that value does not account for the `ids` filter, so the attribute reports the rows actually returned.",
			},
			"executions": schema.ListAttribute{
				Computed:    true,
				ElementType: imagePipelineExecutionObjectType,
				Description: "A list of image pipeline executions that match the configured filters.",
			},
		},
	}
}

func (d *imagePipelineExecutionsDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("image_pipeline_id"),
			path.MatchRoot("execution_id"),
		),
	}
}

func (d *imagePipelineExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !req.Config.Raw.IsFullyKnown() {
		if req.ClientCapabilities.DeferralAllowed {
			resp.Deferred = &datasource.Deferred{Reason: datasource.DeferredReasonDataSourceConfigUnknown}
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: the configuration is not known yet", imagePipelineExecutionsTypeName),
			"The configuration still contains unknown values, so the filters they define cannot be sent. Apply the configuration that defines them and read the data source again.",
		)
		return
	}

	var model imagePipelineExecutionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasImagePipelineId := !model.ImagePipelineId.IsNull() && !model.ImagePipelineId.IsUnknown() && model.ImagePipelineId.ValueString() != ""
	hasExecutionId := !model.ExecutionId.IsNull() && !model.ExecutionId.IsUnknown() && model.ExecutionId.ValueString() != ""
	// The configuration validator enforces this before Read; the check keeps a
	// direct Read invocation from sending a request the API always rejects.
	if !hasImagePipelineId && !hasExecutionId {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: invalid configuration", imagePipelineExecutionsTypeName),
			"At least one of `image_pipeline_id` and `execution_id` must be set; DescribeImagePipelineExecutions cannot list executions without one of them.",
		)
		return
	}

	// One deadline covers every page, every retry and the final state write.
	ctx, cancel := context.WithTimeout(ctx, imagePipelineExecutionsReadBudget)
	defer cancel()

	model, err := d.load(ctx, d.Meta(), d.Meta().RegionId, model)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", imagePipelineExecutionsTypeName, imagePipelineExecutionsAction),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// load performs the whole read: it builds the request, follows every page,
// decodes the records and returns the state model. Read only adds the deadline,
// diagnostics and the state write, so the same request-to-state path can be
// exercised directly in unit tests.
func (d *imagePipelineExecutionsDataSource) load(ctx context.Context, client *connectivity.AliyunClient, regionId string, model imagePipelineExecutionsModel) (imagePipelineExecutionsModel, error) {
	hasImagePipelineId := !model.ImagePipelineId.IsNull() && !model.ImagePipelineId.IsUnknown() && model.ImagePipelineId.ValueString() != ""
	hasExecutionId := !model.ExecutionId.IsNull() && !model.ExecutionId.IsUnknown() && model.ExecutionId.ValueString() != ""
	// The configuration validator enforces this before Read; the check keeps a
	// direct Read invocation from sending a request the API always rejects.
	if !hasImagePipelineId && !hasExecutionId {
		return model, fmt.Errorf("at least one of `image_pipeline_id` and `execution_id` must be set; DescribeImagePipelineExecutions cannot list executions without one of them")
	}

	request := imagePipelineExecutionsRequest(model, regionId)

	objects, lastToken, err := d.imagePipelineExecutionsDescribeAll(ctx, client, request)
	if err != nil {
		return model, err
	}

	// Blank and repeated IDs are dropped before the filter is built, exactly as
	// the argument and its Schema description promise. An all-blank list is
	// therefore an empty filter that selects every execution, not a filter that
	// matches nothing.
	idsFilter := make(map[string]struct{}, len(model.Ids.Elements()))
	for _, id := range ecsUniqueStrings(ecsElementsAsStrings(model.Ids)) {
		idsFilter[id] = struct{}{}
	}

	// The identity is read with the same strict rule the collector applies, so
	// the filter, the sort key and the returned id cannot disagree about which
	// record they describe.
	type matchedExecution struct {
		id     string
		object map[string]interface{}
	}
	matched := make([]matchedExecution, 0, len(objects))
	for _, object := range objects {
		executionId, err := ecsRecordIdentity(object, imagePipelineExecutionsAction, "ExecutionId")
		if err != nil {
			return model, err
		}
		if len(idsFilter) > 0 {
			if _, ok := idsFilter[executionId]; !ok {
				continue
			}
		}
		matched = append(matched, matchedExecution{id: executionId, object: object})
	}
	// Sort the objects themselves so ids, executions and every other
	// index-aligned view describe the same element.
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].id < matched[j].id
	})

	ids := make([]string, 0, len(matched))
	elements := make([]attr.Value, 0, len(matched))
	for _, entry := range matched {
		ids = append(ids, entry.id)
		element, err := imagePipelineExecutionObject(entry.object)
		if err != nil {
			return model, err
		}
		elements = append(elements, element)
	}

	// A cancelled or expired read must not return a partial result to state.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return model, fmt.Errorf("the read was cancelled or ran out of time before the result was written: %v", ctxErr)
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = ecsStringListValue(ids)
	model.TotalCount = types.Int64Value(int64(len(elements)))
	model.NextToken = types.StringValue(lastToken)
	model.Executions = types.ListValueMust(imagePipelineExecutionObjectType, elements)
	return model, nil
}

// imagePipelineExecutionsRequest builds the full DescribeImagePipelineExecutions
// request from the configuration. It is kept separate from Read so the
// request-to-state mapping can be asserted directly.
func imagePipelineExecutionsRequest(model imagePipelineExecutionsModel, regionId string) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":   regionId,
		"MaxResults": int(imagePipelineExecutionsDefaultMaxResults),
	}
	ecsSetRequestString(request, "ImagePipelineId", model.ImagePipelineId)
	ecsSetRequestString(request, "ExecutionId", model.ExecutionId)
	ecsSetRequestString(request, "Status", imagePipelineExecutionsNormalizedStatus(model.Status))
	if !model.MaxResults.IsNull() && !model.MaxResults.IsUnknown() {
		request["MaxResults"] = int(model.MaxResults.ValueInt64())
	}
	ecsApplyTags(request, model.Tags)
	if !model.NextToken.IsNull() && !model.NextToken.IsUnknown() && model.NextToken.ValueString() != "" {
		request["NextToken"] = model.NextToken.ValueString()
	}
	return request
}

// imagePipelineExecutionsNormalizedStatus trims the comma separated tokens so
// the value sent to the API is exactly the form that was validated.
func imagePipelineExecutionsNormalizedStatus(status types.String) types.String {
	if status.IsNull() || status.IsUnknown() {
		return status
	}
	tokens := make([]string, 0)
	for _, token := range strings.Split(status.ValueString(), ",") {
		if token = strings.TrimSpace(token); token != "" {
			tokens = append(tokens, token)
		}
	}
	return types.StringValue(strings.Join(tokens, ","))
}

func (d *imagePipelineExecutionsDataSource) imagePipelineExecutionsDescribeAll(ctx context.Context, client *connectivity.AliyunClient, request map[string]interface{}) ([]map[string]interface{}, string, error) {
	initialToken := ""
	if v, ok := request["NextToken"].(string); ok {
		initialToken = v
	}
	return fwCollectTokenPages(ctx, imagePipelineExecutionsMaxPages, initialToken, imagePipelineExecutionsAction, "ExecutionId", func(token string) (fwTokenPage, error) {
		if token == "" {
			delete(request, "NextToken")
		} else {
			request["NextToken"] = token
		}
		response, err := fwRpcPost(ctx, client, imagePipelineExecutionsTypeName, imagePipelineExecutionsAction, request)
		if err != nil {
			return fwTokenPage{}, err
		}
		if err := ecsValidateQueryEnvelope(response, imagePipelineExecutionsAction, "ImagePipelineExecution", "TotalCount", token == ""); err != nil {
			return fwTokenPage{}, err
		}
		items, err := ecsDecodeRecordList(response, imagePipelineExecutionsAction, "ImagePipelineExecution", "ImagePipelineExecutionSet")
		if err != nil {
			return fwTokenPage{}, err
		}
		nextToken, err := fwDecodeToken(response, "NextToken", imagePipelineExecutionsAction)
		if err != nil {
			return fwTokenPage{}, err
		}
		// The documented empty result still reports TotalCount, so the count is
		// what separates a legitimate empty page from a page that dropped its
		// records: the collector rejects a total that no record accounts for.
		totalCount, totalKnown, err := fwDecodeCount(response, "TotalCount", imagePipelineExecutionsAction)
		if err != nil {
			return fwTokenPage{}, err
		}
		return fwTokenPage{items: items, nextToken: nextToken, total: totalCount, totalKnown: totalKnown}, nil
	})
}

func imagePipelineExecutionObject(object map[string]interface{}) (types.Object, error) {
	identity, err := ecsRecordIdentity(object, imagePipelineExecutionsAction, "ExecutionId")
	if err != nil {
		return types.ObjectNull(imagePipelineExecutionAttrTypes), err
	}
	attributes, err := imagePipelineExecutionAttributes(object)
	if err != nil {
		return types.ObjectNull(imagePipelineExecutionAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	attributes["execution_id"] = types.StringValue(identity)
	return types.ObjectValueMust(imagePipelineExecutionAttrTypes, attributes), nil
}

// imagePipelineExecutionAttributes decodes every member of one execution. Each
// member is decoded by name so a malformed scalar fails the record with its
// path instead of being written to state as an empty or formatted string.
func imagePipelineExecutionAttributes(object map[string]interface{}) (map[string]attr.Value, error) {
	fields := []ecsStringField{
		{attribute: "execution_id", member: "ExecutionId"},
		{attribute: "image_pipeline_id", member: "ImagePipelineId"},
		{attribute: "status", member: "Status"},
		{attribute: "creation_time", member: "CreationTime"},
		{attribute: "modified_time", member: "ModifiedTime"},
		{attribute: "resource_group_id", member: "ResourceGroupId"},
		{attribute: "message", member: "Message"},
		{attribute: "image_id", member: "ImageId"},
		{attribute: "error_code", member: "ErrorCode"},
		{attribute: "error_info", member: "ErrorInfo"},
	}
	attributes := make(map[string]attr.Value, len(fields)+1)
	if err := ecsDecodeStringFields(object, fields, attributes); err != nil {
		return nil, err
	}
	tags, err := ecsDecodeTags(object, ecsTagDecodeOptions{reservedKey: ecsProviderReservedTagKey, alternateMembers: true})
	if err != nil {
		return nil, err
	}
	attributes["tags"] = tags
	return attributes, nil
}

// fwRetryTimeout is the retry window of one read: no attempt is started, and no
// attempt may run, past it. It is a variable so tests can drive the loop without
// waiting.
var fwRetryTimeout = 5 * time.Minute

// fwRetryBackoff is the first wait between two attempts of one read; it grows
// by the same step up to 30 seconds. It is a variable so tests can drive the
// retry loop without waiting.
var fwRetryBackoff = 3 * time.Second

// fwRpcClient is the context-bound RPC capability the read loop needs. Keeping
// the boundary explicit lets its error classification and cancellation be
// regression-tested without real credentials.
type fwRpcClient interface {
	RpcPostWithContext(ctx context.Context, apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}) (map[string]interface{}, error)
}

// fwRpcCall performs the single transport call. It is a variable so tests can
// simulate a request that is still in flight when the read is cancelled. The
// client call issues exactly one attempt and binds the request to ctx, so
// cancelling the read aborts the call that is already on the wire; every retry
// is owned by fwRpcPost below.
var fwRpcCall = func(ctx context.Context, client fwRpcClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	return client.RpcPostWithContext(ctx, "Ecs", "2014-05-26", action, nil, request)
}

// fwErrorDocument reports a 2xx response whose body is an error document. The
// service may answer a rejected query with a success status and a Code/Message
// body; treating it as a result would turn a refusal into a silent empty list.
// The SDK error keeps the business code and message so the diagnostic names the
// failure the service reported instead of this provider's own wording.
func fwErrorDocument(response map[string]interface{}) error {
	code, _ := response["Code"].(string)
	code = strings.TrimSpace(code)
	if code == "" {
		return nil
	}
	message, _ := response["Message"].(string)
	requestId, _ := response["RequestId"].(string)
	return tea.NewSDKError(map[string]interface{}{
		"code":       code,
		"statusCode": 200,
		"message":    "code: 200, " + strings.TrimSpace(message) + " request id: " + strings.TrimSpace(requestId),
	})
}

// fwRpcPost sends one read-only ECS query with the shared retry policy:
// transient failures are retried inside the bounded window, everything else
// fails on the first response. The context carries the read budget: it gates
// every attempt and the backoff, and the request itself is built from it, so an
// in-flight call is aborted instead of running on after the read returns. Each
// attempt is additionally bounded by what is left of fwRetryTimeout, so one slow
// response cannot consume the retries the window was meant to allow. The
// loop runs inline, so no retry callback can write a response after the caller
// has stopped waiting.
func fwRpcPost(ctx context.Context, client fwRpcClient, typeName, action string, request map[string]interface{}) (map[string]interface{}, error) {
	backoff := fwRetryBackoff
	retryDeadline := time.Now().Add(fwRetryTimeout)
	for {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, errs.WrapErrorf(fmt.Errorf("%s was not sent: %v", action, ctxErr), errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
		}
		// The retry window also bounds the request that is in flight: without
		// this, a single slow response could outlive fwRetryTimeout and leave no
		// budget for the attempt the deadline was meant to allow.
		remaining := time.Until(retryDeadline)
		if remaining <= 0 {
			return nil, errs.WrapErrorf(fmt.Errorf("%s was not sent: the retry window of %v is exhausted", action, fwRetryTimeout), errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
		}
		attemptCtx, cancelAttempt := context.WithTimeout(ctx, remaining)
		response, rpcErr := fwRpcCall(attemptCtx, client, action, request)
		attemptErr := attemptCtx.Err()
		cancelAttempt()
		if rpcErr == nil {
			logs.AddDebug(action, response, request)
			if err := fwErrorDocument(response); err != nil {
				return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
			}
			return response, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, errs.WrapErrorf(fmt.Errorf("%s was cancelled while the request was in flight: %v", action, ctxErr), errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
		}
		shouldRetry := ecsNeedRetry(rpcErr)
		// http.Client.Timeout can match DeadlineExceeded while both read budgets remain valid.
		if !shouldRetry && attemptErr == nil && errors.Is(rpcErr, context.DeadlineExceeded) && !errors.Is(rpcErr, context.Canceled) {
			shouldRetry = true
		}
		if attemptErr != nil || !shouldRetry || !time.Now().Add(backoff).Before(retryDeadline) {
			return nil, errs.WrapErrorf(rpcErr, errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
		}
		if waitErr := reservedInstancesSleepBackoff(ctx, backoff); waitErr != nil {
			return nil, errs.WrapErrorf(waitErr, errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
		}
		if backoff < 30*time.Second {
			backoff += 3 * time.Second
		}
	}
}

func fwDecodeCount(response map[string]interface{}, key, action string) (int64, bool, error) {
	if raw, ok := response[key]; ok && raw == nil {
		return 0, false, fmt.Errorf("%s returned %s as null instead of a non-negative integer", action, key)
	}
	return ecsDecodeCount(response, action, key)
}

// fwDecodeToken reads a string pagination token. A missing member or an empty
// string means the server has no further page; a member of the wrong type or a
// whitespace-only non-empty value is malformed and must not be treated as the
// end of the result or as a cursor that can be followed.
func fwDecodeToken(response map[string]interface{}, key, action string) (string, error) {
	raw, ok := response[key]
	if !ok {
		return "", nil
	}
	if raw == nil {
		return "", fmt.Errorf("%s returned %s as null instead of a string", action, key)
	}
	token, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s returned %s as %T instead of a string", action, key, raw)
	}
	if token != "" && strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("%s returned %s as a blank string instead of a cursor", action, key)
	}
	return token, nil
}

// fwTokenPage is one page of a token paginated query. totalKnown distinguishes
// an explicit TotalCount of 0 from a response that does not report the total.
type fwTokenPage struct {
	items      []map[string]interface{}
	nextToken  string
	total      int64
	totalKnown bool
}

// fwCollectTokenPages follows a continuation token until the server returns
// none. Every record must carry a unique, non-empty identity, so partially
// overlapping pages are rejected as well. A repeated token or a token for the
// whole page budget means the cursor never terminates; the accumulated rows are
// discarded instead of being returned as a partial result.
//
// The server total is part of the success envelope, so it is checked as well: a
// page that omits the records while the total still expects rows, a total that
// changes between pages, and more rows than the total are all inconsistent
// responses and are rejected instead of being read as an empty or complete
// result. A read that started at the beginning of the result must report the
// total once it ends: without it the prefix already read cannot be proven
// complete, so a missing TotalCount is rejected from a start token while a read
// that starts from a configured token cannot see the rows before that cursor and
// completeness is only asserted when the read started at the beginning of the
// result.
//
// The configured start token is already in use before the first request, so it
// is recorded up front: a server that echoes it back would otherwise make the
// loop send the very same request a second time before noticing the repeat.
func fwCollectTokenPages(ctx context.Context, maxPages int, initialToken, action, identityKey string, fetch func(token string) (fwTokenPage, error)) ([]map[string]interface{}, string, error) {
	objects := make([]map[string]interface{}, 0)
	seenTokens := map[string]struct{}{}
	if initialToken != "" {
		seenTokens[initialToken] = struct{}{}
	}
	seenRecords := map[string]struct{}{}
	totalCount := int64(0)
	totalKnown := false
	token := initialToken
	for page := 0; page < maxPages; page++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, "", fmt.Errorf("%s read stopped before page %d was requested: %v; the result was not returned because it would be incomplete", action, page+1, ctxErr)
		}
		result, err := fetch(token)
		if err != nil {
			return nil, "", err
		}
		pageItems, nextToken := result.items, result.nextToken
		if result.totalKnown {
			if totalKnown && result.total != totalCount {
				return nil, "", fmt.Errorf("%s reported TotalCount %d for the page at token %q after %d for an earlier page; the result was not returned because the server total is inconsistent", action, result.total, token, totalCount)
			}
			totalCount = result.total
			totalKnown = true
		}
		for _, item := range pageItems {
			identity, err := ecsRecordIdentity(item, action, identityKey)
			if err != nil {
				return nil, "", err
			}
			if _, duplicate := seenRecords[identity]; duplicate {
				return nil, "", fmt.Errorf("%s returned %s %q more than once; the result was not returned because it would contain duplicate rows", action, identityKey, identity)
			}
			seenRecords[identity] = struct{}{}
		}
		objects = append(objects, pageItems...)
		if totalKnown && int64(len(objects)) > totalCount {
			return nil, "", fmt.Errorf("%s returned %d rows while TotalCount reports %d; the result was not returned because the response is inconsistent", action, len(objects), totalCount)
		}
		if nextToken == "" {
			if initialToken == "" && !totalKnown {
				return nil, "", fmt.Errorf("%s returned no continuation token and no TotalCount for a read that started at the beginning; the result was not returned because its completeness cannot be proven", action)
			}
			if totalKnown && initialToken == "" && int64(len(objects)) != totalCount {
				return nil, "", fmt.Errorf("%s returned %d rows while TotalCount reports %d and the server returned no continuation token; the result was not returned because it would be incomplete", action, len(objects), totalCount)
			}
			return objects, "", nil
		}
		if _, duplicate := seenTokens[nextToken]; duplicate {
			return nil, "", fmt.Errorf("%s returned the same continuation token %q twice, so pagination did not terminate; the result was not returned because it would be incomplete", action, nextToken)
		}
		seenTokens[nextToken] = struct{}{}
		token = nextToken
	}
	return nil, "", fmt.Errorf("%s returned a continuation token for %d consecutive requests, so pagination did not terminate; the result was not returned because it would be incomplete", action, maxPages)
}
