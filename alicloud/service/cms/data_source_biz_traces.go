// Package cms. Framework-native data source for alicloud_cms_biz_traces.
// This file is a hand-written addition mirroring the SDKv2 data source (PR #10443)
// while staying self-contained within package cms (no root alicloud import).
package cms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const bizTracesDataSourceTypeName = "alicloud_cms_biz_traces"

var (
	_ datasource.DataSource              = &bizTracesDataSource{}
	_ datasource.DataSourceWithConfigure = &bizTracesDataSource{}
)

func NewBizTracesDataSource() datasource.DataSource {
	return &bizTracesDataSource{}
}

type bizTracesDataSource struct {
	fwadapt.DataSourceBase
}

type bizTracesDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Ids        types.List   `tfsdk:"ids"`
	NameRegex  types.String `tfsdk:"name_regex"`
	OutputFile types.String `tfsdk:"output_file"`
	Workspace  types.String `tfsdk:"workspace"`
	BizTraces  types.List   `tfsdk:"biz_traces"`
}

func (d *bizTracesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = bizTracesDataSourceTypeName
}

func (d *bizTracesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Cms Biz Traces of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID, a hash of the returned Biz Trace IDs.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Biz Trace IDs. When provided, only matching Biz Traces are returned.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by Biz Trace name.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running terraform plan).",
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Description: "The workspace to which the business traces belong.",
			},
			// biz_traces is a list of objects (cty list(object(...))) rather than a
			// nested-attributes attribute: this provider serves its framework half over
			// protocol v5 via tf5muxserver, and protocol v5 cannot represent the
			// framework's nested-attributes (NestedAttributeObject/Attributes-map) form.
			// A plain ListAttribute with an object element type uses the cty type system
			// and is fully representable under protocol v5.
			"biz_traces": schema.ListAttribute{
				Computed:    true,
				Description: "A list of Biz Traces.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"biz_trace_id":    types.StringType,
						"biz_trace_code":  types.StringType,
						"biz_trace_name":  types.StringType,
						"rule_config":     types.StringType,
						"advanced_config": types.StringType,
						"workspace":       types.StringType,
						"create_time":     types.StringType,
						"region_id":       types.StringType,
					},
				},
			},
		},
	}
}

func (d *bizTracesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state bizTracesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the ids filter set.
	idsFilter := map[string]bool{}
	if !state.Ids.IsNull() && !state.Ids.IsUnknown() {
		for _, v := range state.Ids.Elements() {
			if s, ok := v.(types.String); ok && !s.IsNull() {
				idsFilter[s.ValueString()] = true
			}
		}
	}

	// Compile name_regex if provided.
	var nameRegex *regexp.Regexp
	if !state.NameRegex.IsNull() && !state.NameRegex.IsUnknown() {
		r, err := regexp.Compile(state.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", bizTracesDataSourceTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = r
	}

	locations := map[string]string{"workspace": "query", "maxResults": "query", "nextToken": "query", "__api": "ListBizTraces"}
	nextToken := ""
	var objects []map[string]interface{}

	for {
		request := map[string]interface{}{"maxResults": "100"}
		if !state.Workspace.IsNull() && !state.Workspace.IsUnknown() {
			request["workspace"] = state.Workspace.ValueString()
		}
		if nextToken != "" {
			request["nextToken"] = nextToken
		}

		response, err := retryRoaGet(d.Client(), "Cms", "2024-03-30", "/bizTraces", request, 5*time.Minute, locations)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListBizTraces", bizTracesDataSourceTypeName),
				errs.WrapErrorf(err, "Reading %s: calling ListBizTraces", bizTracesDataSourceTypeName).Error(),
			)
			return
		}

		items, _ := response["items"].([]interface{})
		if items == nil {
			items = []interface{}{}
		}
		for _, v := range items {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if len(idsFilter) > 0 {
				bizTraceId := fmt.Sprint(item["bizTraceId"])
				if !idsFilter[bizTraceId] {
					continue
				}
			}
			objects = append(objects, item)
		}

		// Check for nextToken pagination.
		nextToken = ""
		if rawToken, ok := response["nextToken"]; ok {
			nextToken = fmt.Sprint(rawToken)
			if nextToken == "" || nextToken == "<nil>" {
				nextToken = ""
			}
		}
		if nextToken == "" || len(items) == 0 {
			break
		}
	}

	// Apply name_regex filter and build output.
	ids := make([]string, 0)
	bizTraceItems := make([]attr.Value, 0)
	bizTraceObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"biz_trace_id":    types.StringType,
			"biz_trace_code":  types.StringType,
			"biz_trace_name":  types.StringType,
			"rule_config":     types.StringType,
			"advanced_config": types.StringType,
			"workspace":       types.StringType,
			"create_time":     types.StringType,
			"region_id":       types.StringType,
		},
	}
	for _, object := range objects {
		bizTraceName := fmt.Sprint(object["bizTraceName"])
		if nameRegex != nil && !nameRegex.MatchString(bizTraceName) {
			continue
		}
		bizTraceId := fmt.Sprint(object["bizTraceId"])
		ids = append(ids, bizTraceId)

		itemValues := map[string]attr.Value{
			"biz_trace_id":    stringVal(object, "bizTraceId"),
			"biz_trace_code":  stringVal(object, "bizTraceCode"),
			"biz_trace_name":  stringVal(object, "bizTraceName"),
			"rule_config":     stringVal(object, "ruleConfig"),
			"advanced_config": stringVal(object, "advancedConfig"),
			"workspace":       stringVal(object, "workspace"),
			"create_time":     stringVal(object, "createTime"),
			"region_id":       stringVal(object, "regionId"),
		}
		bizTraceItems = append(bizTraceItems, types.ObjectValueMust(bizTraceObjType.AttributeTypes(), itemValues))
	}

	state.Ids = types.ListValueMust(types.StringType, stringSliceToValues(ids))
	state.BizTraces = types.ListValueMust(bizTraceObjType, bizTraceItems)
	state.Id = types.StringValue(bizTracesDataSourceID(ids))

	// Write output file if requested.
	if !state.OutputFile.IsNull() && !state.OutputFile.IsUnknown() {
		mappings := make([]map[string]interface{}, 0, len(objects))
		for _, object := range objects {
			mappings = append(mappings, map[string]interface{}{
				"biz_trace_id":    fmt.Sprint(object["bizTraceId"]),
				"biz_trace_code":  object["bizTraceCode"],
				"biz_trace_name":  object["bizTraceName"],
				"rule_config":     object["ruleConfig"],
				"advanced_config": object["advancedConfig"],
				"workspace":       object["workspace"],
				"create_time":     object["createTime"],
				"region_id":       object["regionId"],
			})
		}
		data, _ := json.Marshal(mappings)
		if err := os.WriteFile(state.OutputFile.ValueString(), data, 0644); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output file", bizTracesDataSourceTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// stringVal extracts a string field from a response object, returning StringNull if missing.
func stringVal(m map[string]interface{}, key string) types.String {
	v, ok := m[key]
	if !ok || v == nil {
		return types.StringNull()
	}
	return types.StringValue(fmt.Sprint(v))
}

func stringSliceToValues(ss []string) []attr.Value {
	out := make([]attr.Value, 0, len(ss))
	for _, s := range ss {
		out = append(out, types.StringValue(s))
	}
	return out
}

// bizTracesDataSourceID computes a stable hash of the returned IDs for the data
// source ID. A local sha256 is used to avoid importing the root alicloud
// package (which would create a circular dependency).
func bizTracesDataSourceID(ids []string) string {
	h := sha256.Sum256([]byte(strings.Join(ids, ",")))
	return hex.EncodeToString(h[:16])
}
