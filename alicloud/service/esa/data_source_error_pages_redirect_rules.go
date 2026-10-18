package esa

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &errorPagesRedirectRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &errorPagesRedirectRulesDataSource{}
)

// NewErrorPagesRedirectRulesDataSource is the factory the framework provider
// registry calls to instantiate the data source.
func NewErrorPagesRedirectRulesDataSource() datasource.DataSource {
	return &errorPagesRedirectRulesDataSource{}
}

type errorPagesRedirectRulesDataSource struct {
	fwadapt.DataSourceBase
}

// errorPagesRedirectRuleDSItemType is the framework object type of a single
// element of the rules output list. error_pages_redirect is a nested list of
// objects reusing the resource's element type.
var errorPagesRedirectRuleDSItemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                   types.StringType,
		"config_id":            types.StringType,
		"config_type":          types.StringType,
		"site_version":         types.Int64Type,
		"rule_enable":          types.StringType,
		"rule_name":            types.StringType,
		"rule":                 types.StringType,
		"sequence":             types.Int64Type,
		"error_pages_redirect": types.ListType{ElemType: errorPagesRedirectItemType},
	},
}

type errorPagesRedirectRuleDSItem struct {
	Id                 types.String `tfsdk:"id"`
	ConfigId           types.String `tfsdk:"config_id"`
	ConfigType         types.String `tfsdk:"config_type"`
	SiteVersion        types.Int64  `tfsdk:"site_version"`
	RuleEnable         types.String `tfsdk:"rule_enable"`
	RuleName           types.String `tfsdk:"rule_name"`
	Rule               types.String `tfsdk:"rule"`
	Sequence           types.Int64  `tfsdk:"sequence"`
	ErrorPagesRedirect types.List   `tfsdk:"error_pages_redirect"`
}

type errorPagesRedirectRulesModel struct {
	Id          types.String `tfsdk:"id"`
	Ids         types.List   `tfsdk:"ids"`
	Names       types.List   `tfsdk:"names"`
	Rules       types.List   `tfsdk:"rules"`
	SiteId      types.String `tfsdk:"site_id"`
	NameRegex   types.String `tfsdk:"name_regex"`
	ConfigId    types.String `tfsdk:"config_id"`
	ConfigType  types.String `tfsdk:"config_type"`
	RuleName    types.String `tfsdk:"rule_name"`
	SiteVersion types.Int64  `tfsdk:"site_version"`
	OutputFile  types.String `tfsdk:"output_file"`
}

func (d *errorPagesRedirectRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = errorPagesRedirectRulesTypeName
}

func (d *errorPagesRedirectRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the ESA ErrorPagesRedirectRule list.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of rule IDs. Only rules whose ID is in this list are returned.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by rule_name. Read fails if the value is not a valid regular expression.",
			},
			"site_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ESA site.",
			},
			"config_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the error pages redirect rule.",
			},
			"config_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the configuration.",
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the error pages redirect rule.",
			},
			"site_version": schema.Int64Attribute{
				Optional:    true,
				Description: "The version of the ESA site.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Save the result to a file in JSON.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of rule names.",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of error pages redirect rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The rule ID in the format `<site_id>:<config_id>`.",
						},
						"config_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the error pages redirect rule.",
						},
						"config_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the configuration.",
						},
						"site_version": schema.Int64Attribute{
							Computed:    true,
							Description: "The version of the ESA site.",
						},
						"rule_enable": schema.StringAttribute{
							Computed:    true,
							Description: "Whether the rule is enabled.",
						},
						"rule_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the error pages redirect rule.",
						},
						"rule": schema.StringAttribute{
							Computed:    true,
							Description: "The rule expression.",
						},
						"sequence": schema.Int64Attribute{
							Computed:    true,
							Description: "The sequence of the rule.",
						},
						"error_pages_redirect": schema.ListNestedAttribute{
							Computed:    true,
							Description: "The error pages redirect targets.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"target_url": schema.StringAttribute{
										Computed:    true,
										Description: "The URL to which the error page is redirected.",
									},
									"status_code": schema.StringAttribute{
										Computed:    true,
										Description: "The HTTP status code that triggers the redirect.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *errorPagesRedirectRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config errorPagesRedirectRulesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteId := config.SiteId.ValueString()
	request := map[string]interface{}{
		"RegionId":   d.Client().RegionId,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
		"SiteId":     siteId,
	}
	if !config.ConfigId.IsNull() {
		request["ConfigId"] = config.ConfigId.ValueString()
	}
	if !config.ConfigType.IsNull() {
		request["ConfigType"] = config.ConfigType.ValueString()
	}
	if !config.RuleName.IsNull() {
		request["RuleName"] = config.RuleName.ValueString()
	}
	if !config.SiteVersion.IsNull() {
		request["SiteVersion"] = config.SiteVersion.ValueInt64()
	}

	// ids filter: when set, only rules whose composite id is in this map are kept.
	idsFilter := map[string]bool{}
	if !config.Ids.IsNull() && !config.Ids.IsUnknown() {
		var ids []string
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			idsFilter[id] = true
		}
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s", errorPagesRedirectRulesTypeName),
				fmt.Sprintf("compiling name_regex: %v", err),
			)
			return
		}
		nameRegex = r
	}

	var objects []map[string]interface{}
	for {
		response, err := d.Client().RpcGet(errorPagesRedirectRuleProduct, errorPagesRedirectRuleVersion, "ListErrorPagesRedirectRules", request, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListErrorPagesRedirectRules", errorPagesRedirectRulesTypeName),
				errs.WrapErrorf(err, "Reading %s: calling ListErrorPagesRedirectRules", errorPagesRedirectRulesTypeName).Error(),
			)
			return
		}

		page, _ := response["Configs"].([]interface{})
		for _, v := range page {
			item, _ := v.(map[string]interface{})
			id := fmt.Sprintf("%s:%s", siteId, toString(item["ConfigId"]))
			if len(idsFilter) > 0 && !idsFilter[id] {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
				continue
			}
			objects = append(objects, item)
		}

		if len(page) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0, len(objects))
	names := make([]string, 0, len(objects))
	rules := make([]errorPagesRedirectRuleDSItem, 0, len(objects))
	for _, obj := range objects {
		rule := errorPagesRedirectRuleDSItem{
			Id:                 types.StringValue(fmt.Sprintf("%s:%s", siteId, toString(obj["ConfigId"]))),
			ConfigId:           stringVal(obj["ConfigId"]),
			ConfigType:         stringVal(obj["ConfigType"]),
			SiteVersion:        int64Val(obj["SiteVersion"]),
			RuleEnable:         stringVal(obj["RuleEnable"]),
			RuleName:           stringVal(obj["RuleName"]),
			Rule:               stringVal(obj["Rule"]),
			Sequence:           int64Val(obj["Sequence"]),
			ErrorPagesRedirect: errorPagesRedirectList(ctx, obj["ErrorPagesRedirect"], &resp.Diagnostics),
		}
		rules = append(rules, rule)
		ids = append(ids, rule.Id.ValueString())
		names = append(names, stringOrEmpty(obj["RuleName"]))
	}

	state := errorPagesRedirectRulesModel{
		SiteId:      config.SiteId,
		NameRegex:   config.NameRegex,
		ConfigId:    config.ConfigId,
		ConfigType:  config.ConfigType,
		RuleName:    config.RuleName,
		SiteVersion: config.SiteVersion,
		OutputFile:  config.OutputFile,
	}
	var ok bool
	if state.Ids, ok = listValueFromStrings(ctx, ids, &resp.Diagnostics); !ok {
		return
	}
	if state.Names, ok = listValueFromStrings(ctx, names, &resp.Diagnostics); !ok {
		return
	}
	if state.Rules, ok = listValueFromObjects(ctx, errorPagesRedirectRuleDSItemType, rules, &resp.Diagnostics); !ok {
		return
	}
	state.Id = types.StringValue(helper.DataResourceIdHash(ids))

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := helper.WriteToFile(config.OutputFile.ValueString(), rules); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s", errorPagesRedirectRulesTypeName),
				fmt.Sprintf("writing output_file: %v", err),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// listValueFromStrings builds a types.List of strings from a Go slice. Errors
// are appended to diags; ok=false means the caller should return immediately.
func listValueFromStrings(ctx context.Context, in []string, diags *diag.Diagnostics) (types.List, bool) {
	if in == nil {
		in = []string{}
	}
	l, d := types.ListValueFrom(ctx, types.StringType, in)
	diags.Append(d...)
	return l, !d.HasError()
}

// listValueFromObjects builds a types.List of object values from a Go slice of
// structs with tfsdk tags matching elementType. Errors are appended to diags.
func listValueFromObjects(ctx context.Context, elementType attr.Type, in interface{}, diags *diag.Diagnostics) (types.List, bool) {
	l, d := types.ListValueFrom(ctx, elementType, in)
	diags.Append(d...)
	return l, !d.HasError()
}

// stringOrEmpty returns the string form of v, "" when nil. Used for the names
// list where a null is not appropriate (names is a Computed list of strings).
func stringOrEmpty(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
