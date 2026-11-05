package ecs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
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
					fwTagMapValidator{allowEmptyValue: true},
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
				Description: "The pagination token to start reading from. When it is configured, the first request sends it and the read starts at that position; when it is not configured, the first request omits it and the read starts at the beginning. Each following request sends the token of the previous response, and the read stops when the server returns none. A token that repeats, a blank token, or 1,000 consecutive requests that each still return a token aborts the read with an error instead of returning a partial result. The value stored afterwards is the last token the server returned, and it is empty when every page has been read.",
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

	idsFilter := make(map[string]struct{}, len(model.Ids.Elements()))
	for _, id := range fwElementsAsStrings(model.Ids) {
		idsFilter[id] = struct{}{}
	}

	matched := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		executionId := fwString(object["ExecutionId"]).ValueString()
		if len(idsFilter) > 0 {
			if _, ok := idsFilter[executionId]; !ok {
				continue
			}
		}
		matched = append(matched, object)
	}
	// Sort the objects themselves so ids, executions and every other
	// index-aligned view describe the same element.
	sort.Slice(matched, func(i, j int) bool {
		return fwString(matched[i]["ExecutionId"]).ValueString() < fwString(matched[j]["ExecutionId"]).ValueString()
	})

	ids := make([]string, 0, len(matched))
	elements := make([]attr.Value, 0, len(matched))
	for _, object := range matched {
		ids = append(ids, fwString(object["ExecutionId"]).ValueString())
		element, err := imagePipelineExecutionObject(object)
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
	model.Ids = fwStringListValue(ids)
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
	fwSetRequestString(request, "ImagePipelineId", model.ImagePipelineId)
	fwSetRequestString(request, "ExecutionId", model.ExecutionId)
	fwSetRequestString(request, "Status", imagePipelineExecutionsNormalizedStatus(model.Status))
	if !model.MaxResults.IsNull() && !model.MaxResults.IsUnknown() {
		request["MaxResults"] = int(model.MaxResults.ValueInt64())
	}
	fwApplyTags(request, model.Tags)
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
	return fwCollectTokenPages(ctx, imagePipelineExecutionsMaxPages, initialToken, imagePipelineExecutionsAction, "ExecutionId", func(token string) ([]map[string]interface{}, string, error) {
		if token == "" {
			delete(request, "NextToken")
		} else {
			request["NextToken"] = token
		}
		response, err := fwRpcPost(ctx, client, imagePipelineExecutionsTypeName, imagePipelineExecutionsAction, request)
		if err != nil {
			return nil, "", err
		}
		items, err := fwDecodeRecordList(response, "ImagePipelineExecution", "ImagePipelineExecutionSet", imagePipelineExecutionsAction)
		if err != nil {
			return nil, "", err
		}
		nextToken, err := fwDecodeToken(response, "NextToken", imagePipelineExecutionsAction)
		if err != nil {
			return nil, "", err
		}
		return items, nextToken, nil
	})
}

func imagePipelineExecutionObject(object map[string]interface{}) (types.Object, error) {
	identity, err := fwRecordIdentity(object, "ExecutionId", imagePipelineExecutionsAction)
	if err != nil {
		return types.ObjectNull(imagePipelineExecutionAttrTypes), err
	}
	tags, err := fwDecodeTags(object["Tags"])
	if err != nil {
		return types.ObjectNull(imagePipelineExecutionAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	return types.ObjectValueMust(imagePipelineExecutionAttrTypes, map[string]attr.Value{
		"execution_id":      fwString(object["ExecutionId"]),
		"image_pipeline_id": fwString(object["ImagePipelineId"]),
		"status":            fwString(object["Status"]),
		"creation_time":     fwString(object["CreationTime"]),
		"modified_time":     fwString(object["ModifiedTime"]),
		"resource_group_id": fwString(object["ResourceGroupId"]),
		"message":           fwString(object["Message"]),
		"image_id":          fwString(object["ImageId"]),
		"error_code":        fwString(object["ErrorCode"]),
		"error_info":        fwString(object["ErrorInfo"]),
		"tags":              tags,
	}), nil
}
