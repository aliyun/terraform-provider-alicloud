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

const nodeOnBaselinesTypeName = "alicloud_data_works_node_on_baselines"

var (
	_ datasource.DataSource              = &nodeOnBaselinesDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeOnBaselinesDataSource{}
)

func NewNodeOnBaselinesDataSource() datasource.DataSource {
	return &nodeOnBaselinesDataSource{}
}

type nodeOnBaselinesDataSource struct {
	fwadapt.DataSourceBase
}

// nodeOnBaselinesModel is the plural data source's tfsdk model. ids is Optional +
// Computed: a user may narrow the lookup by a set of node IDs, and the field also
// reports the IDs of the nodes that matched. nodes is the exported list.
type nodeOnBaselinesModel struct {
	BaselineId types.String              `tfsdk:"baseline_id"`
	Ids        []types.String            `tfsdk:"ids"`
	NodeId     types.String              `tfsdk:"node_id"`
	Owner      types.String              `tfsdk:"owner"`
	ProjectId  types.String              `tfsdk:"project_id"`
	NameRegex  types.String              `tfsdk:"name_regex"`
	OutputFile types.String              `tfsdk:"output_file"`
	Id         types.String              `tfsdk:"id"`
	Nodes      []nodeOnBaselineNodeModel `tfsdk:"nodes"`
}

// nodeOnBaselineNodeModel is one element of the plural data source's nodes list.
type nodeOnBaselineNodeModel struct {
	Id         types.String `tfsdk:"id"`
	BaselineId types.String `tfsdk:"baseline_id"`
	NodeId     types.String `tfsdk:"node_id"`
	NodeName   types.String `tfsdk:"node_name"`
	Owner      types.String `tfsdk:"owner"`
	ProjectId  types.String `tfsdk:"project_id"`
	RegionId   types.String `tfsdk:"region_id"`
}

func (d *nodeOnBaselinesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = nodeOnBaselinesTypeName
}

func (d *nodeOnBaselinesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of Data Works Node On Baseline resources of the current Alibaba Cloud user that match the supplied filters.",
		Attributes: map[string]schema.Attribute{
			"baseline_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DataWorks Baseline that the nodes belong to.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of node IDs. Only nodes whose IDs match this list are returned.",
			},
			"node_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the result by the specified node ID.",
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the result by the specified node owner.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the result by the specified DataWorks project ID.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string used to filter the result by node name.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID of the data source, computed from the matched node IDs.",
			},
			"nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The resource ID in the form of `<baseline_id>:<node_id>`.",
						},
						"baseline_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the DataWorks Baseline that the node belongs to.",
						},
						"node_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the node mounted on the baseline.",
						},
						"node_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the node mounted on the baseline.",
						},
						"owner": schema.StringAttribute{
							Computed:    true,
							Description: "The owner of the node mounted on the baseline.",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the DataWorks workspace that the node belongs to.",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "The region ID of the node.",
						},
					},
				},
			},
		},
	}
}

func (d *nodeOnBaselinesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model nodeOnBaselinesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	action := "ListNodesByBaseline"
	request := map[string]interface{}{
		"BaselineId": model.BaselineId.ValueString(),
	}

	idsFilter := make(map[string]struct{})
	for _, id := range model.Ids {
		if id.IsNull() {
			continue
		}
		idsFilter[id.ValueString()] = struct{}{}
	}

	var nameRegex *regexp.Regexp
	if !model.NameRegex.IsNull() && model.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(model.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", nodeOnBaselinesTypeName),
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
			fmt.Sprintf("Reading %s: calling %s", nodeOnBaselinesTypeName, action),
			errs.WrapErrorf(err, errs.DataDefaultErrorMsg, nodeOnBaselinesTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	data, gerr := jsonpath.Get("$.Data", response)
	if gerr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: extracting $.Data from %s", nodeOnBaselinesTypeName, action),
			errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, action, "$.Data", response).Error(),
		)
		return
	}
	result, _ := data.([]interface{})

	var nodes []nodeOnBaselineNodeModel
	var ids []types.String
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		nodeID := fmt.Sprint(item["NodeId"])
		if nodeIDFilter != "" && nodeID != nodeIDFilter {
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
		if len(idsFilter) > 0 {
			if _, keep := idsFilter[fmt.Sprintf("%s:%s", model.BaselineId.ValueString(), nodeID)]; !keep {
				continue
			}
		}
		nodes = append(nodes, nodeOnBaselineNodeModel{
			Id:         types.StringValue(fmt.Sprintf("%s:%s", model.BaselineId.ValueString(), nodeID)),
			BaselineId: model.BaselineId,
			NodeId:     types.StringValue(nodeID),
			NodeName:   types.StringValue(fmt.Sprint(item["NodeName"])),
			Owner:      types.StringValue(fmt.Sprint(item["Owner"])),
			ProjectId:  types.StringValue(fmt.Sprint(item["ProjectId"])),
			RegionId:   types.StringValue(client.RegionId),
		})
		ids = append(ids, types.StringValue(fmt.Sprintf("%s:%s", model.BaselineId.ValueString(), nodeID)))
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(toStringSlice(ids)))
	model.Ids = ids
	model.Nodes = nodes

	if !model.OutputFile.IsNull() && model.OutputFile.ValueString() != "" {
		out := make([]map[string]interface{}, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, map[string]interface{}{
				"id":          n.Id.ValueString(),
				"baseline_id": n.BaselineId.ValueString(),
				"node_id":     n.NodeId.ValueString(),
				"node_name":   n.NodeName.ValueString(),
				"owner":       n.Owner.ValueString(),
				"project_id":  n.ProjectId.ValueString(),
				"region_id":   n.RegionId.ValueString(),
			})
		}
		if werr := helper.WriteToFile(model.OutputFile.ValueString(), out); werr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", nodeOnBaselinesTypeName),
				errs.WrapError(werr).Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// toStringSlice flattens a []types.String into a []string for hashing.
func toStringSlice(values []types.String) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v.IsNull() {
			continue
		}
		out = append(out, v.ValueString())
	}
	return out
}
