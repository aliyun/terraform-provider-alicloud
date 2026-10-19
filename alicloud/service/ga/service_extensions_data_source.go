package ga

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	serviceExtensionsTypeName = "alicloud_ga_service_extensions"
	pageSizeLarge             = 50
)

var (
	_ datasource.DataSource              = (*serviceExtensionsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*serviceExtensionsDataSource)(nil)
)

func NewServiceExtensionsDataSource() datasource.DataSource {
	return &serviceExtensionsDataSource{}
}

type serviceExtensionsDataSource struct {
	fwadapt.DataSourceBase
}

type serviceExtensionsDataSourceModel struct {
	Ids                  types.List   `tfsdk:"ids"`
	NameRegex            types.String `tfsdk:"name_regex"`
	Names                types.List   `tfsdk:"names"`
	ServiceExtensionName types.String `tfsdk:"service_extension_name"`
	Status               types.String `tfsdk:"status"`
	OutputFile           types.String `tfsdk:"output_file"`
	EnableDetails        types.Bool   `tfsdk:"enable_details"`
	ServiceExtensions    types.List   `tfsdk:"service_extensions"`
}

// DS nested object attribute types. The data source surfaces components and
// resources as lists (matching the SDKv2 TypeList), distinct from the
// resource's set semantics.
var (
	dsComponentObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service_component_id": types.StringType,
			"priority":             types.Int64Type,
			"timeout":              types.Int64Type,
			"config":               types.StringType,
			"fail_policy":          types.StringType,
			"component_name":       types.StringType,
		},
	}
	dsResourceObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"accelerator_id": types.StringType,
			"resource_type":  types.StringType,
			"resource_id":    types.StringType,
			"associate_id":   types.StringType,
		},
	}
	dsServiceExtensionObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                   types.StringType,
			"service_extension_id": types.StringType,
			"name":                 types.StringType,
			"description":          types.StringType,
			"type":                 types.StringType,
			"state":                types.StringType,
			"resource_group_id":    types.StringType,
			"tags":                 types.MapType{ElemType: types.StringType},
			"components":           types.ListType{ElemType: dsComponentObjectType},
			"resources":            types.ListType{ElemType: dsResourceObjectType},
			"create_time":          types.StringType,
			"update_time":          types.StringType,
		},
	}
)

func (d *serviceExtensionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = serviceExtensionsTypeName
}

func (d *serviceExtensionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the list of Global Accelerator (GA) Service Extensions and their basic information.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of service extension IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by service extension name.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of service extension names.",
			},
			"service_extension_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the service extension.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "The status of the service extension.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results.",
			},
			"enable_details": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Default to false. Set it to true to get the details of each service extension, including components, resources and tags.",
			},
			"service_extensions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of service extensions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the service extension, same as `service_extension_id`.",
						},
						"service_extension_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the service extension.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the service extension.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the service extension.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the service extension.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The state of the service extension.",
						},
						"resource_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the resource group.",
						},
						"tags": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The tags of the service extension.",
						},
						"components": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"service_component_id": schema.StringAttribute{Computed: true},
									"priority":             schema.Int64Attribute{Computed: true},
									"timeout":              schema.Int64Attribute{Computed: true},
									"config":               schema.StringAttribute{Computed: true},
									"fail_policy":          schema.StringAttribute{Computed: true},
									"component_name":       schema.StringAttribute{Computed: true},
								},
							},
						},
						"resources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"accelerator_id": schema.StringAttribute{Computed: true},
									"resource_type":  schema.StringAttribute{Computed: true},
									"resource_id":    schema.StringAttribute{Computed: true},
									"associate_id":   schema.StringAttribute{Computed: true},
								},
							},
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time when the service extension was created.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time when the service extension was last updated.",
						},
					},
				},
			},
		},
	}
}

func (d *serviceExtensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serviceExtensionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := d.Client()

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", serviceExtensionsTypeName),
				errs.WrapError(err).Error(),
			)
			return
		}
		nameRegex = r
	}

	idsFilter := make(map[string]struct{})
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

	status := config.Status.ValueString()
	statusSet := !config.Status.IsNull() && status != ""
	enableDetails := !config.EnableDetails.IsNull() && config.EnableDetails.ValueBool()

	action := "ListServiceExtensions"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": pageSizeLarge,
	}
	if !config.ServiceExtensionName.IsNull() {
		request["Name"] = config.ServiceExtensionName.ValueString()
	}

	var objects []map[string]interface{}
	for {
		response, err := dsListRetry(client, action, request)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", serviceExtensionsTypeName, action),
				errs.WrapErrorf(err, errs.DataDefaultErrorMsg, serviceExtensionsTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		raw, err := jsonpath.Get("$.ServiceExtensions", response)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: parsing %s response", serviceExtensionsTypeName, action),
				errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ServiceExtensions", response).Error(),
			)
			return
		}
		list, _ := raw.([]interface{})
		for _, v := range list {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["Name"])
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}
			if len(idsFilter) > 0 {
				if _, ok := idsFilter[fmt.Sprint(item["ServiceExtensionId"])]; !ok {
					continue
				}
			}
			if statusSet && status != fmt.Sprint(item["State"]) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
			continue
		}
		break
	}

	var ids []string
	var names []string
	var items []attr.Value
	for _, object := range objects {
		id := fmt.Sprint(object["ServiceExtensionId"])
		ids = append(ids, id)
		names = append(names, fmt.Sprint(object["Name"]))

		obj, d := buildDSExtensionObject(ctx, client, object, id, enableDetails)
		resp.Diagnostics.Append(d...)
		items = append(items, obj)
	}

	config.Ids = stringListValue(ctx, ids)
	config.Names = stringListValue(ctx, names)
	seList, diags := types.ListValue(dsServiceExtensionObjectType, items)
	resp.Diagnostics.Append(diags...)
	config.ServiceExtensions = seList

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

	if !config.OutputFile.IsNull() {
		if err := writeToFile(config.OutputFile.ValueString(), serializeDSItems(ctx, seList)); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output_file", serviceExtensionsTypeName),
				err.Error(),
			)
		}
	}
}

// dsListRetry wraps ListServiceExtensions with the SDKv2 data-source retry
// policy (5m, NeedRetry-driven), matching the original resource.
func dsListRetry(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	return rpcPostRetry(client, action, request, 5*time.Minute, true)
}

// buildDSExtensionObject constructs the nested object for one service
// extension. When enable_details is set it fetches components/resources and
// tags; otherwise it surfaces only the list-level fields.
func buildDSExtensionObject(ctx context.Context, client *connectivity.AliyunClient, object map[string]interface{}, id string, enableDetails bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	componentsList, d := types.ListValue(dsComponentObjectType, []attr.Value{})
	diags.Append(d...)
	resourcesList, d := types.ListValue(dsResourceObjectType, []attr.Value{})
	diags.Append(d...)
	tagsMap, d := types.MapValue(types.StringType, map[string]attr.Value{})
	diags.Append(d...)

	if enableDetails {
		detail, err := describeServiceExtension(client, id)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Reading %s: calling GetServiceExtension", serviceExtensionsTypeName),
				errs.WrapError(err).Error(),
			)
			return types.Object{}, diags
		}
		if v, ok := detail["Components"].([]interface{}); ok {
			l, dd := dsComponentList(ctx, v)
			componentsList = l
			diags.Append(dd...)
		}
		if v, ok := detail["Resources"].([]interface{}); ok {
			l, dd := dsResourceList(ctx, v)
			resourcesList = l
			diags.Append(dd...)
		}
		tags, terr := listServiceExtensionTags(client, id)
		if terr != nil {
			if isExpectedErrors(terr, []string{"Invalid.ResourceType"}) {
				log.Printf("[WARN] Data source %s tagging is not supported by the backend (Invalid.ResourceType), skipping tags for %s: %s", serviceExtensionsTypeName, id, terr)
			} else {
				diags.AddError(
					fmt.Sprintf("Reading %s: calling ListTagResources", serviceExtensionsTypeName),
					errs.WrapError(terr).Error(),
				)
				return types.Object{}, diags
			}
		}
		tagsMap = tagsMapToValue(ctx, tags, &diags)
	}

	obj, d := types.ObjectValue(dsServiceExtensionObjectType.AttrTypes, map[string]attr.Value{
		"id":                   types.StringValue(id),
		"service_extension_id": stringValue(object["ServiceExtensionId"]),
		"name":                 stringValue(object["Name"]),
		"description":          stringValue(object["Description"]),
		"type":                 stringValue(object["Type"]),
		"state":                stringValue(object["State"]),
		"resource_group_id":    stringValue(object["ResourceGroupId"]),
		"tags":                 tagsMap,
		"components":           componentsList,
		"resources":            resourcesList,
		"create_time":          stringValue(object["CreateTime"]),
		"update_time":          stringValue(object["UpdateTime"]),
	})
	diags.Append(d...)
	return obj, diags
}

func dsComponentList(ctx context.Context, items []interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		obj, d := types.ObjectValue(dsComponentObjectType.AttrTypes, map[string]attr.Value{
			"service_component_id": stringValue(item["ServiceComponentId"]),
			"priority":             int64Value(item["Priority"]),
			"timeout":              int64Value(item["Timeout"]),
			"config":               stringValue(item["Config"]),
			"fail_policy":          stringValue(item["FailPolicy"]),
			"component_name":       stringValue(item["ComponentName"]),
		})
		diags.Append(d...)
		elements = append(elements, obj)
	}
	list, d := types.ListValue(dsComponentObjectType, elements)
	diags.Append(d...)
	return list, diags
}

func dsResourceList(ctx context.Context, items []interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		obj, d := types.ObjectValue(dsResourceObjectType.AttrTypes, map[string]attr.Value{
			"accelerator_id": stringValue(item["AcceleratorId"]),
			"resource_type":  stringValue(item["ResourceType"]),
			"resource_id":    stringValue(item["ResourceId"]),
			"associate_id":   stringValue(item["AssociateId"]),
		})
		diags.Append(d...)
		elements = append(elements, obj)
	}
	list, d := types.ListValue(dsResourceObjectType, elements)
	diags.Append(d...)
	return list, diags
}

func stringListValue(ctx context.Context, values []string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	list, _ := types.ListValue(types.StringType, elements)
	return list
}

// serializeDSItems pulls the freshly-built list back into a plain slice so
// writeToFile emits JSON the SDKv2 data source did, instead of a tfsdk value.
func serializeDSItems(ctx context.Context, list types.List) []map[string]interface{} {
	var items []serviceExtensionItemModel
	_ = list.ElementsAs(ctx, &items, false)
	out := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		out = append(out, map[string]interface{}{
			"id":                   it.Id.ValueString(),
			"service_extension_id": it.ServiceExtensionId.ValueString(),
			"name":                 it.Name.ValueString(),
			"description":          it.Description.ValueString(),
			"type":                 it.Type.ValueString(),
			"state":                it.State.ValueString(),
			"resource_group_id":    it.ResourceGroupId.ValueString(),
			"create_time":          it.CreateTime.ValueString(),
			"update_time":          it.UpdateTime.ValueString(),
		})
	}
	return out
}

// serviceExtensionItemModel mirrors the data source's nested object for
// ElementsAs when serializing the output file.
type serviceExtensionItemModel struct {
	Id                 types.String `tfsdk:"id"`
	ServiceExtensionId types.String `tfsdk:"service_extension_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
	ResourceGroupId    types.String `tfsdk:"resource_group_id"`
	Tags               types.Map    `tfsdk:"tags"`
	Components         types.List   `tfsdk:"components"`
	Resources          types.List   `tfsdk:"resources"`
	CreateTime         types.String `tfsdk:"create_time"`
	UpdateTime         types.String `tfsdk:"update_time"`
}
