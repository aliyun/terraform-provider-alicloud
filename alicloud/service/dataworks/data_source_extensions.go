package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	extensionsDataSourceTypeName = "alicloud_data_works_extensions"
	dataworksProduct             = "dataworks-public"
	dataworksAPIVersion          = "2020-05-18"
	// PageSizeLarge mirrors the SDK v2 suite's page size for list queries.
	pageSizeLarge = 50
)

var (
	_ datasource.DataSource              = &extensionsDataSource{}
	_ datasource.DataSourceWithConfigure = &extensionsDataSource{}
)

// NewExtensionsDataSource is the factory registered under
// alicloud_data_works_extensions; it returns a Framework-native data source.
func NewExtensionsDataSource() datasource.DataSource {
	return &extensionsDataSource{}
}

type extensionsDataSource struct {
	fwadapt.DataSourceBase
}

// extensionsDataSourceModel is the tfsdk state/config model.
type extensionsDataSourceModel struct {
	Ids           types.List       `tfsdk:"ids"`
	ExtensionCode types.String     `tfsdk:"extension_code"`
	NameRegex     types.String     `tfsdk:"name_regex"`
	OutputFile    types.String     `tfsdk:"output_file"`
	Extensions    []extensionModel `tfsdk:"extensions"`
}

// extensionModel is the nested element type of the extensions list.
type extensionModel struct {
	Id                types.String `tfsdk:"id"`
	ExtensionCode     types.String `tfsdk:"extension_code"`
	ExtensionName     types.String `tfsdk:"extension_name"`
	ExtensionDesc     types.String `tfsdk:"extension_desc"`
	Status            types.String `tfsdk:"status"`
	HelpDocUrl        types.String `tfsdk:"help_doc_url"`
	ProjectTesting    types.String `tfsdk:"project_testing"`
	DetailUrl         types.String `tfsdk:"detail_url"`
	ParameterSetting  types.String `tfsdk:"parameter_setting"`
	OptionSetting     types.String `tfsdk:"option_setting"`
	BindEventList     types.String `tfsdk:"bind_event_list"`
	EventCategoryList types.String `tfsdk:"event_category_list"`
	RegionId          types.String `tfsdk:"region_id"`
}

func (d *extensionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = extensionsDataSourceTypeName
}

func (d *extensionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Data Works Extensions of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of Extension IDs. The ID is the same as extension_code.",
			},
			"extension_code": schema.StringAttribute{
				Optional:    true,
				Description: "The unique code of the extension. If specified, the data source queries the details of this single extension via GetExtension; otherwise, it lists all extensions via ListExtensions.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter extensions by extension_name.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running terraform plan).",
			},
			"extensions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of Data Works Extensions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                  schema.StringAttribute{Computed: true, Description: "The ID of the extension, same as extension_code."},
						"extension_code":      schema.StringAttribute{Computed: true, Description: "The unique code of the extension."},
						"extension_name":      schema.StringAttribute{Computed: true, Description: "The name of the extension."},
						"extension_desc":      schema.StringAttribute{Computed: true, Description: "The description of the extension."},
						"status":              schema.StringAttribute{Computed: true, Description: "The status of the extension."},
						"help_doc_url":        schema.StringAttribute{Computed: true, Description: "The help document URL of the extension."},
						"project_testing":     schema.StringAttribute{Computed: true, Description: "Whether the extension is in project testing."},
						"detail_url":          schema.StringAttribute{Computed: true, Description: "The detail URL of the extension."},
						"parameter_setting":   schema.StringAttribute{Computed: true, Description: "The parameter setting of the extension, serialized as a JSON string."},
						"option_setting":      schema.StringAttribute{Computed: true, Description: "The option setting of the extension, serialized as a JSON string."},
						"bind_event_list":     schema.StringAttribute{Computed: true, Description: "The list of bound events, serialized as a JSON string."},
						"event_category_list": schema.StringAttribute{Computed: true, Description: "The list of event categories, serialized as a JSON string."},
						"region_id":           schema.StringAttribute{Computed: true, Description: "The region ID of the extension."},
					},
				},
			},
		},
	}
}

func (d *extensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config extensionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			"Provider Client Not Injected",
			fmt.Sprintf("%s ended Configure without a provider client. Please report this issue to the provider developers.", extensionsDataSourceTypeName),
		)
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: invalid name_regex", extensionsDataSourceTypeName), err.Error())
			return
		}
		nameRegex = r
	}

	idsFilter := map[string]struct{}{}
	if !config.Ids.IsNull() {
		var ids []string
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			idsFilter[id] = struct{}{}
		}
	}

	raw, err := queryExtensions(client, config.ExtensionCode, idsFilter, nameRegex)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", extensionsDataSourceTypeName), err.Error())
		return
	}

	extensions := make([]extensionModel, 0, len(raw))
	ids := make([]string, 0, len(raw))
	output := make([]map[string]interface{}, 0, len(raw))
	for _, obj := range raw {
		mapped := extensionMap(obj)
		output = append(output, mapped)
		extensions = append(extensions, extensionModel{
			Id:                types.StringValue(asString(mapped["id"])),
			ExtensionCode:     types.StringValue(asString(mapped["extension_code"])),
			ExtensionName:     types.StringValue(asString(mapped["extension_name"])),
			ExtensionDesc:     types.StringValue(asString(mapped["extension_desc"])),
			Status:            types.StringValue(asString(mapped["status"])),
			HelpDocUrl:        types.StringValue(asString(mapped["help_doc_url"])),
			ProjectTesting:    types.StringValue(asString(mapped["project_testing"])),
			DetailUrl:         types.StringValue(asString(mapped["detail_url"])),
			ParameterSetting:  types.StringValue(asString(mapped["parameter_setting"])),
			OptionSetting:     types.StringValue(asString(mapped["option_setting"])),
			BindEventList:     types.StringValue(asString(mapped["bind_event_list"])),
			EventCategoryList: types.StringValue(asString(mapped["event_category_list"])),
			RegionId:          types.StringValue(asString(mapped["region_id"])),
		})
		ids = append(ids, asString(mapped["extension_code"]))
	}

	idList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := extensionsDataSourceModel{
		Ids:           idList,
		ExtensionCode: config.ExtensionCode,
		NameRegex:     config.NameRegex,
		OutputFile:    config.OutputFile,
		Extensions:    extensions,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := writeToFile(config.OutputFile.ValueString(), output); err != nil {
			resp.Diagnostics.AddWarning(fmt.Sprintf("Reading %s: writing output_file", extensionsDataSourceTypeName), err.Error())
		}
	}
}

// queryExtensions resolves the API objects to list. When extension_code is set,
// GetExtension returns the single detail; otherwise ListExtensions paginates.
func queryExtensions(client *connectivity.AliyunClient, extensionCode types.String, idsFilter map[string]struct{}, nameRegex *regexp.Regexp) ([]map[string]interface{}, error) {
	if !extensionCode.IsNull() && extensionCode.ValueString() != "" {
		const action = "GetExtension"
		request := map[string]interface{}{"ExtensionCode": extensionCode.ValueString()}
		response, err := client.RpcPost(dataworksProduct, dataworksAPIVersion, action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("calling %s: %w", action, err)
		}
		data, _ := response["Data"].(map[string]interface{})
		if len(data) > 0 {
			return []map[string]interface{}{data}, nil
		}
		return nil, nil
	}

	const action = "ListExtensions"
	pageNumber := 1
	var objects []map[string]interface{}
	for {
		request := map[string]interface{}{
			"PageSize":   pageSizeLarge,
			"PageNumber": pageNumber,
		}
		response, err := client.RpcPost(dataworksProduct, dataworksAPIVersion, action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("calling %s: %w", action, err)
		}
		data, _ := response["Data"].(map[string]interface{})
		list, _ := data["Extensions"].([]interface{})
		for _, item := range list {
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if len(idsFilter) > 0 {
				if _, ok := idsFilter[asString(obj["ExtensionCode"])]; !ok {
					continue
				}
			}
			if nameRegex != nil {
				if !nameRegex.MatchString(asString(obj["ExtensionName"])) {
					continue
				}
			}
			objects = append(objects, obj)
		}
		if len(list) < pageSizeLarge {
			break
		}
		pageNumber++
	}
	return objects, nil
}

// extensionMap normalizes a raw API object into the snake_case shape exported to
// state and output_file; complex fields are JSON-serialized to mirror the SDK v2
// data source.
func extensionMap(obj map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                  asString(obj["ExtensionCode"]),
		"extension_code":      asString(obj["ExtensionCode"]),
		"extension_name":      asString(obj["ExtensionName"]),
		"extension_desc":      asString(obj["ExtensionDesc"]),
		"status":              asString(obj["Status"]),
		"help_doc_url":        asString(obj["HelpDocUrl"]),
		"project_testing":     asString(obj["ProjectTesting"]),
		"detail_url":          asString(obj["DetailUrl"]),
		"parameter_setting":   toJSONString(obj["ParameterSetting"]),
		"option_setting":      toJSONString(obj["OptionSetting"]),
		"bind_event_list":     toJSONString(obj["BindEventList"]),
		"event_category_list": toJSONString(obj["EventCategoryList"]),
		"region_id":           asString(obj["RegionId"]),
	}
}

// asString coerces a value to a string, mapping nil to "".
func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// toJSONString serializes complex API fields (objects/arrays) to a JSON string;
// a string value is returned verbatim, mirroring the SDK v2 convertObjectToJsonString.
func toJSONString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

// writeToFile persists the resolved results as JSON, matching the SDK v2 suite.
func writeToFile(filePath string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, b, 0o644)
}
