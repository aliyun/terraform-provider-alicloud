package dataworks

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const nodeOnBaselineTypeName = "alicloud_data_works_node_on_baseline"

var (
	_ datasource.DataSource              = &nodeOnBaselineDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeOnBaselineDataSource{}
)

func NewNodeOnBaselineDataSource() datasource.DataSource {
	return &nodeOnBaselineDataSource{}
}

type nodeOnBaselineDataSource struct {
	fwadapt.DataSourceBase
}

// nodeOnBaselineModel is the singular data source's tfsdk model. The filter
// arguments (owner, project_id) are Optional + Computed: a user may supply them
// to narrow the lookup, or leave them unset and read the matched node's value.
type nodeOnBaselineModel struct {
	BaselineId types.String `tfsdk:"baseline_id"`
	NodeId     types.String `tfsdk:"node_id"`
	Owner      types.String `tfsdk:"owner"`
	ProjectId  types.String `tfsdk:"project_id"`
	NameRegex  types.String `tfsdk:"name_regex"`
	OutputFile types.String `tfsdk:"output_file"`
	Id         types.String `tfsdk:"id"`
	NodeName   types.String `tfsdk:"node_name"`
	RegionId   types.String `tfsdk:"region_id"`
}

func (d *nodeOnBaselineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = nodeOnBaselineTypeName
}

func (d *nodeOnBaselineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the first Data Works Node On Baseline of the current Alibaba Cloud user that matches the supplied filters.",
		Attributes: map[string]schema.Attribute{
			"baseline_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DataWorks Baseline that the node belongs to.",
			},
			"node_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the result by the specified node ID. The first matching node is returned.",
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Filter the result by the specified node owner. The first matching node is returned.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Filter the result by the specified DataWorks project ID. The first matching node is returned.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string used to filter the result by node name. The first matching node is returned.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in the form of `<baseline_id>:<node_id>`.",
			},
			"node_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the node mounted on the baseline.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the node.",
			},
		},
	}
}

func (d *nodeOnBaselineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model nodeOnBaselineModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	action := "GetNodeOnBaseline"
	request := map[string]interface{}{
		"BaselineId": model.BaselineId.ValueString(),
	}

	var nameRegex *regexp.Regexp
	if !model.NameRegex.IsNull() && model.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(model.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", nodeOnBaselineTypeName),
				errs.WrapError(err).Error(),
			)
			return
		}
		nameRegex = r
	}

	nodeIDFilter := model.NodeId.ValueString()
	ownerFilter := model.Owner.ValueString()
	projectIDFilter := model.ProjectId.ValueString()

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", nodeOnBaselineTypeName, action),
			errs.WrapErrorf(err, errs.DataDefaultErrorMsg, nodeOnBaselineTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	data, gerr := jsonpath.Get("$.Data", response)
	if gerr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: extracting $.Data from %s", nodeOnBaselineTypeName, action),
			errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, action, "$.Data", response).Error(),
		)
		return
	}
	result, _ := data.([]interface{})

	var matched map[string]interface{}
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if nodeIDFilter != "" && fmt.Sprint(item["NodeId"]) != nodeIDFilter {
			continue
		}
		if ownerFilter != "" && fmt.Sprint(item["Owner"]) != ownerFilter {
			continue
		}
		if projectIDFilter != "" && fmt.Sprint(item["ProjectId"]) != projectIDFilter {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["NodeName"])) {
			continue
		}
		matched = item
		break
	}

	if matched == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: no matching node on baseline %s", nodeOnBaselineTypeName, model.BaselineId.ValueString()),
			errs.WrapErrorf(errs.NotFoundErr("DataWorks NodeOnBaseline", model.BaselineId.ValueString()), errs.NotFoundWithResponse, response).Error(),
		)
		return
	}

	nodeID := fmt.Sprint(matched["NodeId"])
	model.Id = types.StringValue(fmt.Sprintf("%s:%s", model.BaselineId.ValueString(), nodeID))
	model.NodeId = types.StringValue(nodeID)
	model.NodeName = types.StringValue(fmt.Sprint(matched["NodeName"]))
	model.Owner = types.StringValue(fmt.Sprint(matched["Owner"]))
	model.ProjectId = types.StringValue(fmt.Sprint(matched["ProjectId"]))
	model.RegionId = types.StringValue(client.RegionId)

	if !model.OutputFile.IsNull() && model.OutputFile.ValueString() != "" {
		out := map[string]interface{}{
			"id":          model.Id.ValueString(),
			"baseline_id": model.BaselineId.ValueString(),
			"node_id":     nodeID,
			"node_name":   matched["NodeName"],
			"owner":       matched["Owner"],
			"project_id":  matched["ProjectId"],
			"region_id":   client.RegionId,
		}
		if werr := helper.WriteToFile(model.OutputFile.ValueString(), out); werr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", nodeOnBaselineTypeName),
				errs.WrapError(werr).Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
