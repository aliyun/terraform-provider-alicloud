package esa

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const siteEnvironmentsTypeName = "alicloud_esa_site_environments"

var (
	_ datasource.DataSource              = &siteEnvironmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &siteEnvironmentsDataSource{}
)

func NewSiteEnvironmentsDataSource() datasource.DataSource {
	return &siteEnvironmentsDataSource{}
}

type siteEnvironmentsDataSource struct {
	fwadapt.DataSourceBase
}

type siteEnvironmentsModel struct {
	Ids             types.List   `tfsdk:"ids"`
	NameRegex       types.String `tfsdk:"name_regex"`
	SiteId          types.String `tfsdk:"site_id"`
	OutputFile      types.String `tfsdk:"output_file"`
	Id              types.String `tfsdk:"id"`
	Names           types.List   `tfsdk:"names"`
	EnvironmentList types.List   `tfsdk:"environments"`
}

var siteEnvironmentsAttrTypes = map[string]attr.Type{
	"id":               types.StringType,
	"environment_name": types.StringType,
	"rule":             types.StringType,
	"read_only":        types.BoolType,
	"site_version":     types.Int32Type,
	"pre_site_version": types.Int32Type,
	"priority":         types.Int32Type,
	"is_default":       types.BoolType,
	"create_time":      types.StringType,
	"update_time":      types.StringType,
}

func (d *siteEnvironmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = siteEnvironmentsTypeName
}

func (d *siteEnvironmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Esa Site Environments of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Site Environment IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by Site Environment name.",
			},
			"site_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Site.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Site Environment names.",
			},
			"environments": schema.ListAttribute{
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: siteEnvironmentsAttrTypes},
				Description: "The list of environments.",
			},
		},
	}
}

func (d *siteEnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "ListEnvironments"

	var config siteEnvironmentsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}

	setString := func(key string, value types.String) {
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			request[key] = value.ValueString()
		}
	}

	setString("SiteId", config.SiteId)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if !config.Ids.IsNull() {
		ids := make([]string, 0, len(config.Ids.Elements()))
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, v := range ids {
			idsMap[v] = v
		}
	}

	var siteEnvironmentNameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); !config.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", siteEnvironmentsTypeName),
				err.Error(),
			)
			return
		}

		siteEnvironmentNameRegex = compiled
	}

	var response map[string]interface{}
	var err error

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
		response, err = d.Meta().RpcGet("ESA", "2024-09-10", action, request, nil)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "InternalException"}) || helper.NeedRetry(err) {
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
			fmt.Sprintf("Reading %s: calling %s", siteEnvironmentsTypeName, action),
			err.Error(),
		)
		return
	}

	rawEnvironmentList, _ := response["EnvironmentList"].([]interface{})
	for _, v := range rawEnvironmentList {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%v:%v", request["SiteId"], item["EnvironmentName"])]; !ok {
				continue
			}
		}

		if siteEnvironmentNameRegex != nil && !siteEnvironmentNameRegex.MatchString(fmt.Sprint(item["EnvironmentName"])) {
			continue
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	items := make([]attr.Value, 0, len(objects))
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprintf("%v:%v", request["SiteId"], object["EnvironmentName"]),
			"environment_name": fmt.Sprint(object["EnvironmentName"]),
			"rule":             object["Rule"],
			"read_only":        object["ReadOnly"],
			"site_version":     object["SiteVersion"],
			"pre_site_version": object["PreSiteVersion"],
			"priority":         object["Priority"],
			"is_default":       object["IsDefault"],
			"create_time":      object["CreateTime"],
			"update_time":      object["UpdateTime"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["EnvironmentName"])
		s = append(s, mapping)

		items = append(items, types.ObjectValueMust(siteEnvironmentsAttrTypes, map[string]attr.Value{
			"id":               types.StringValue(fmt.Sprintf("%v:%v", request["SiteId"], object["EnvironmentName"])),
			"environment_name": types.StringValue(fmt.Sprint(object["EnvironmentName"])),
			"rule":             types.StringValue(convert.Interface2String(object["Rule"])),
			"read_only":        types.BoolValue(convert.Interface2Bool(object["ReadOnly"])),
			"site_version":     types.Int32Value(convert.Interface2Int32(object["SiteVersion"])),
			"pre_site_version": types.Int32Value(convert.Interface2Int32(object["PreSiteVersion"])),
			"priority":         types.Int32Value(convert.Interface2Int32(object["Priority"])),
			"is_default":       types.BoolValue(convert.Interface2Bool(object["IsDefault"])),
			"create_time":      types.StringValue(convert.Interface2String(object["CreateTime"])),
			"update_time":      types.StringValue(convert.Interface2String(object["UpdateTime"])),
		}))
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, diags := types.ListValue(types.ObjectType{AttrTypes: siteEnvironmentsAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := config
	state.Id = types.StringValue(helper.DataResourceIdHash(ids))
	state.Ids = idsList
	state.Names = namesList
	state.EnvironmentList = environments

	if err := helper.DataResourceWriteOutputFile(state.OutputFile, s); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: writing output_file", siteEnvironmentsTypeName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
