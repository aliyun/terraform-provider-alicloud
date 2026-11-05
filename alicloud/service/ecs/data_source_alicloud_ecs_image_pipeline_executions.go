package ecs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
// value form documented for Status and rejects any token outside the enum.
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
	for _, token := range strings.Split(req.ConfigValue.ValueString(), ",") {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
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
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
				},
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of execution IDs to filter by. The filter is applied on the provider side after all pages are read, so a server that ignores it cannot widen the result.",
			},
			"max_results": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 500. Defaults to 50.",
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The pagination token to start reading from. The data source follows the remaining tokens itself; the value stored afterwards is the last token the server returned, and it is empty when every page has been read.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of execution records that match all configured filters.",
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

	request := map[string]interface{}{
		"RegionId":   d.Meta().RegionId,
		"MaxResults": int(imagePipelineExecutionsDefaultMaxResults),
	}
	ecsFwSetRequestString(request, "ImagePipelineId", model.ImagePipelineId)
	ecsFwSetRequestString(request, "ExecutionId", model.ExecutionId)
	ecsFwSetRequestString(request, "Status", imagePipelineExecutionsNormalizedStatus(model.Status))
	if !model.MaxResults.IsNull() && !model.MaxResults.IsUnknown() {
		request["MaxResults"] = int(model.MaxResults.ValueInt64())
	}
	ecsFwApplyTags(request, model.Tags)
	if !model.NextToken.IsNull() && !model.NextToken.IsUnknown() && model.NextToken.ValueString() != "" {
		request["NextToken"] = model.NextToken.ValueString()
	}

	objects, lastToken, err := d.imagePipelineExecutionsDescribeAll(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", imagePipelineExecutionsTypeName, imagePipelineExecutionsAction),
			err.Error(),
		)
		return
	}

	idsFilter := make(map[string]struct{}, len(model.Ids.Elements()))
	for _, id := range ecsFwElementsAsStrings(model.Ids) {
		idsFilter[id] = struct{}{}
	}

	matched := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		executionId := ecsFwString(object["ExecutionId"]).ValueString()
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
		return ecsFwString(matched[i]["ExecutionId"]).ValueString() < ecsFwString(matched[j]["ExecutionId"]).ValueString()
	})

	ids := make([]string, 0, len(matched))
	elements := make([]attr.Value, 0, len(matched))
	for _, object := range matched {
		ids = append(ids, ecsFwString(object["ExecutionId"]).ValueString())
		elements = append(elements, imagePipelineExecutionObject(object))
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = ecsFwStringListValue(ids)
	model.TotalCount = types.Int64Value(int64(len(elements)))
	model.NextToken = types.StringValue(lastToken)
	model.Executions = types.ListValueMust(imagePipelineExecutionObjectType, elements)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
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

func (d *imagePipelineExecutionsDataSource) imagePipelineExecutionsDescribeAll(ctx context.Context, request map[string]interface{}) ([]map[string]interface{}, string, error) {
	initialToken := ""
	if v, ok := request["NextToken"].(string); ok {
		initialToken = v
	}
	deadline := time.Now().Add(imagePipelineExecutionsReadBudget)
	pagesRead := 0
	return imagePipelineExecutionsCollectPages(imagePipelineExecutionsMaxPages, initialToken, func(token string) ([]map[string]interface{}, string, error) {
		if time.Now().After(deadline) {
			return nil, "", fmt.Errorf("%s read exceeded the %s budget after %d pages; the result was not returned because it would be incomplete", imagePipelineExecutionsAction, imagePipelineExecutionsReadBudget, pagesRead)
		}
		pagesRead++
		if token == "" {
			delete(request, "NextToken")
		} else {
			request["NextToken"] = token
		}
		response, err := ecsFwRpcPost(ctx, d.Meta(), imagePipelineExecutionsTypeName, imagePipelineExecutionsAction, request)
		if err != nil {
			return nil, "", err
		}
		return imagePipelineExecutionsItems(response), imagePipelineExecutionsNextToken(response), nil
	})
}

// imagePipelineExecutionsCollectPages follows the NextToken cursor until the
// server returns no further token. A repeated token, or a full page for the
// whole page budget, means the cursor never terminates; in both cases the
// accumulated rows are discarded instead of being returned as a partial result.
func imagePipelineExecutionsCollectPages(maxPages int, initialToken string, fetch func(token string) ([]map[string]interface{}, string, error)) ([]map[string]interface{}, string, error) {
	objects := make([]map[string]interface{}, 0)
	seenTokens := map[string]struct{}{}
	token := initialToken
	for page := 0; page < maxPages; page++ {
		pageItems, nextToken, err := fetch(token)
		if err != nil {
			return nil, "", err
		}
		objects = append(objects, pageItems...)
		if nextToken == "" {
			return objects, "", nil
		}
		if _, duplicate := seenTokens[nextToken]; duplicate {
			return nil, "", fmt.Errorf("%s returned the same continuation token %q twice, so pagination did not terminate; the result was not returned because it would be incomplete", imagePipelineExecutionsAction, nextToken)
		}
		seenTokens[nextToken] = struct{}{}
		token = nextToken
	}
	return nil, "", fmt.Errorf("%s returned a continuation token for %d consecutive requests, so pagination did not terminate; the result was not returned because it would be incomplete", imagePipelineExecutionsAction, maxPages)
}

// imagePipelineExecutionsItems parses the execution array. A response without
// the member yields an empty slice: the API omits it when nothing matches and
// an empty result is not an error.
func imagePipelineExecutionsItems(response map[string]interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return items
	}
	raw, err := jsonpath.Get("$.ImagePipelineExecution.ImagePipelineExecutionSet", response)
	if err != nil {
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, entry := range list {
		if item, ok := entry.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

func imagePipelineExecutionsNextToken(response map[string]interface{}) string {
	if response == nil {
		return ""
	}
	if token, ok := response["NextToken"].(string); ok {
		return token
	}
	return ""
}

func imagePipelineExecutionObject(object map[string]interface{}) attr.Value {
	return types.ObjectValueMust(imagePipelineExecutionAttrTypes, map[string]attr.Value{
		"execution_id":      ecsFwString(object["ExecutionId"]),
		"image_pipeline_id": ecsFwString(object["ImagePipelineId"]),
		"status":            ecsFwString(object["Status"]),
		"creation_time":     ecsFwString(object["CreationTime"]),
		"modified_time":     ecsFwString(object["ModifiedTime"]),
		"resource_group_id": ecsFwString(object["ResourceGroupId"]),
		"message":           ecsFwString(object["Message"]),
		"image_id":          ecsFwString(object["ImageId"]),
		"error_code":        ecsFwString(object["ErrorCode"]),
		"error_info":        ecsFwString(object["ErrorInfo"]),
		"tags":              ecsFwTagsValue(object["Tags"]),
	})
}
