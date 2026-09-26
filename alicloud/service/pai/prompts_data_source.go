package pai

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const promptsDataSourceTypeName = "alicloud_pai_workspace_prompts"

// Mirrors the deleted root package's PageSizeLarge.
const promptsPageSize = 50

// The legacy SDKv2 data source ran each page fetch under the SDK default
// read timeout; the page loop keeps that budget.
const promptsReadTimeout = 20 * time.Minute

var (
	_ datasource.DataSource              = &promptsDataSource{}
	_ datasource.DataSourceWithConfigure = &promptsDataSource{}
)

func NewPromptsDataSource() datasource.DataSource {
	return &promptsDataSource{}
}

type promptsDataSource struct {
	fwadapt.DataSourceBase
}

type promptsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Ids           types.List   `tfsdk:"ids"`
	Names         types.List   `tfsdk:"names"`
	FrameworkType types.String `tfsdk:"framework_type"`
	NameRegex     types.String `tfsdk:"name_regex"`
	OutputFile    types.String `tfsdk:"output_file"`
	Prompts       types.List   `tfsdk:"prompts"`
	WorkspaceId   types.String `tfsdk:"workspace_id"`
}

type promptsItemModel struct {
	Id               types.String `tfsdk:"id"`
	Accessibility    types.String `tfsdk:"accessibility"`
	CreateTime       types.String `tfsdk:"create_time"`
	Description      types.String `tfsdk:"description"`
	FrameworkContent types.String `tfsdk:"framework_content"`
	FrameworkType    types.String `tfsdk:"framework_type"`
	ModifyTime       types.String `tfsdk:"modify_time"`
	PromptId         types.String `tfsdk:"prompt_id"`
	PromptName       types.String `tfsdk:"prompt_name"`
}

func (d *promptsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = promptsDataSourceTypeName
}

func (d *promptsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of PAI Workspace Prompts available to the user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Prompt IDs. The value is formulated as `<workspace_id>:<prompt_id>`.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by Prompt name.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Prompt names.",
			},
			"framework_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cue word frame type.",
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "Workspace id.",
			},
			// The muxed provider serves proto5, which cannot express nested
			// attributes: ListNestedAttribute fails provider schema loading with
			// "protocol version 5 cannot have Attributes set". The flat
			// list-of-object form addresses state identically (prompts.0.<attr>).
			"prompts": schema.ListAttribute{
				Computed:    true,
				Description: "A list of Prompts. Each element contains the following attributes:",
				ElementType: promptItemObjectType,
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File path where data source results are saved.",
			},
		},
	}
}

func (d *promptsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config promptsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", promptsDataSourceTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if !config.Ids.IsNull() {
		var idsConfig []types.String
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &idsConfig, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range idsConfig {
			if v.IsNull() {
				continue
			}
			idsMap[v.ValueString()] = v.ValueString()
		}
	}

	workspaceId := config.WorkspaceId.ValueString()
	query := map[string]*string{
		"WorkspaceId": tea.String(workspaceId),
		"PageSize":    tea.String(strconv.Itoa(promptsPageSize)),
		"PageNumber":  tea.String("1"),
	}
	if !config.FrameworkType.IsNull() {
		query["FrameworkType"] = tea.String(config.FrameworkType.ValueString())
	}

	const action = "/api/v1/prompts"
	var objects []map[string]interface{}
	for {
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(promptsReadTimeout, func() *retry.RetryError {
			response, err = d.Client().RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListPrompts", promptsDataSourceTypeName),
				errs.WrapErrorf(err, errs.DataDefaultErrorMsg, promptsDataSourceTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}

		raw, _ := jsonpath.Get("$.Prompts[*]", response)
		result, _ := raw.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["PromptName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(workspaceId, ":", item["PromptId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < promptsPageSize {
			break
		}
		pageNum, _ := strconv.Atoi(tea.StringValue(query["PageNumber"]))
		query["PageNumber"] = tea.String(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	items := make([]promptsItemModel, 0)
	for _, object := range objects {
		id := fmt.Sprint(workspaceId, ":", object["PromptId"])
		items = append(items, promptsItemModel{
			Id:               types.StringValue(id),
			Accessibility:    objectStringValue(object["Accessibility"]),
			CreateTime:       objectStringValue(object["CreateTime"]),
			Description:      objectStringValue(object["Description"]),
			FrameworkContent: objectStringValue(object["FrameworkContent"]),
			FrameworkType:    objectStringValue(object["FrameworkType"]),
			ModifyTime:       objectStringValue(object["ModifyTime"]),
			PromptId:         objectStringValue(object["PromptId"]),
			PromptName:       objectStringValue(object["PromptName"]),
		})
		ids = append(ids, id)
		if name, ok := object["PromptName"].(string); ok {
			names = append(names, name)
		}
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	promptsList, diags := types.ListValueFrom(ctx, promptItemObjectType, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := promptsDataSourceModel{
		Id:            types.StringValue(helper.DataResourceIdHash(ids)),
		Ids:           idsList,
		Names:         namesList,
		FrameworkType: config.FrameworkType,
		NameRegex:     config.NameRegex,
		OutputFile:    config.OutputFile,
		Prompts:       promptsList,
		WorkspaceId:   types.StringValue(workspaceId),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		output := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			output = append(output, map[string]interface{}{
				"id":                item.Id.ValueString(),
				"accessibility":     item.Accessibility.ValueStringPointer(),
				"create_time":       item.CreateTime.ValueStringPointer(),
				"description":       item.Description.ValueStringPointer(),
				"framework_content": item.FrameworkContent.ValueStringPointer(),
				"framework_type":    item.FrameworkType.ValueStringPointer(),
				"modify_time":       item.ModifyTime.ValueStringPointer(),
				"prompt_id":         item.PromptId.ValueStringPointer(),
				"prompt_name":       item.PromptName.ValueStringPointer(),
			})
		}
		if err := helper.WriteToFile(config.OutputFile.ValueString(), output); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output file", promptsDataSourceTypeName),
				err.Error(),
			)
		}
	}
}

var promptItemObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                types.StringType,
		"accessibility":     types.StringType,
		"create_time":       types.StringType,
		"description":       types.StringType,
		"framework_content": types.StringType,
		"framework_type":    types.StringType,
		"modify_time":       types.StringType,
		"prompt_id":         types.StringType,
		"prompt_name":       types.StringType,
	},
}
