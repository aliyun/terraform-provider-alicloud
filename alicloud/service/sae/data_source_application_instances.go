package sae

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const applicationInstancesTypeName = "alicloud_sae_application_instances"

var (
	_ datasource.DataSource              = &applicationInstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationInstancesDataSource{}
)

func NewApplicationInstancesDataSource() datasource.DataSource {
	return &applicationInstancesDataSource{}
}

type applicationInstancesDataSource struct {
	fwadapt.DataSourceBase
}

type applicationInstancesDataSourceModel struct {
	ApplicationID types.String `tfsdk:"application_id"`
	GroupID       types.String `tfsdk:"group_id"`
	IDs           types.List   `tfsdk:"ids"`
	OutputFile    types.String `tfsdk:"output_file"`
	ID            types.String `tfsdk:"id"`
	Instances     types.List   `tfsdk:"instances"`
}

type instancesModel struct {
	InstanceID              types.String
	GroupID                 types.String
	InstanceContainerIP     types.String
	InstanceContainerStatus types.String
	InstanceHealthStatus    types.String
	MainContainerStatus     types.String
	ImageURL                types.String
	PackageVersion          types.String
	VSwitchID               types.String
}

var instanceObjectAttrTypes = map[string]attr.Type{
	"instance_id":               types.StringType,
	"group_id":                  types.StringType,
	"instance_container_ip":     types.StringType,
	"instance_container_status": types.StringType,
	"instance_health_status":    types.StringType,
	"main_container_status":     types.StringType,
	"image_url":                 types.StringType,
	"package_version":           types.StringType,
	"vswitch_id":                types.StringType,
}

func (d *applicationInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = applicationInstancesTypeName
}

func (d *applicationInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of application instances of the specified Serverless App Engine (SAE) application, including the ENI IP addresses of the instances.",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the SAE application.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the application group to which the instances belong. If omitted, the data source queries the instances of every group of the application and aggregates the results.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of instance IDs, used as a filter. It is also exported after the query.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is the same as `application_id`.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				Description: "A list of SAE application instances.",
				ElementType: types.ObjectType{
					AttrTypes: instanceObjectAttrTypes,
				},
			},
		},
	}
}

func (d *applicationInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model applicationInstancesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(d.Client())
	applicationId := model.ApplicationID.ValueString()

	var objects []map[string]interface{}
	groupId := model.GroupID.ValueString()
	if groupId == "" {
		groups, err := svc.DescribeApplicationGroups(applicationId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: querying application groups", applicationInstancesTypeName), err.Error())
			return
		}
		for _, group := range groups {
			groupId, _ := group["GroupId"].(string)
			if groupId == "" {
				continue
			}
			result, err := svc.DescribeApplicationInstances(applicationId, groupId)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: querying application instances", applicationInstancesTypeName), err.Error())
				return
			}
			objects = append(objects, result...)
		}
	} else {
		result, err := svc.DescribeApplicationInstances(applicationId, groupId)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: querying application instances", applicationInstancesTypeName), err.Error())
			return
		}
		objects = result
	}

	idFilter := make(map[string]bool)
	if !model.IDs.IsNull() && !model.IDs.IsUnknown() {
		var filterIds []string
		resp.Diagnostics.Append(model.IDs.ElementsAs(ctx, &filterIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range filterIds {
			idFilter[id] = true
		}
	}

	ids := make([]string, 0)
	instanceRows := make([]instancesModel, 0)
	results := make([]map[string]interface{}, 0)
	for _, object := range objects {
		instanceId := fmt.Sprint(object["InstanceId"])
		if len(idFilter) > 0 && !idFilter[instanceId] {
			continue
		}
		instanceRows = append(instanceRows, instancesModel{
			InstanceID:              types.StringValue(instanceId),
			GroupID:                 saeStringValue(object["GroupId"]),
			InstanceContainerIP:     saeStringValue(object["InstanceContainerIp"]),
			InstanceContainerStatus: saeStringValue(object["InstanceContainerStatus"]),
			InstanceHealthStatus:    saeStringValue(object["InstanceHealthStatus"]),
			MainContainerStatus:     saeStringValue(object["MainContainerStatus"]),
			ImageURL:                saeStringValue(object["ImageUrl"]),
			PackageVersion:          saeStringValue(object["PackageVersion"]),
			VSwitchID:               saeStringValue(object["VSwitchId"]),
		})
		ids = append(ids, instanceId)
		results = append(results, map[string]interface{}{
			"instance_id":               instanceId,
			"group_id":                  object["GroupId"],
			"instance_container_ip":     object["InstanceContainerIp"],
			"instance_container_status": object["InstanceContainerStatus"],
			"instance_health_status":    object["InstanceHealthStatus"],
			"main_container_status":     object["MainContainerStatus"],
			"image_url":                 object["ImageUrl"],
			"package_version":           object["PackageVersion"],
			"vswitch_id":                object["VSwitchId"],
		})
	}

	instanceValues := make([]attr.Value, 0, len(instanceRows))
	for _, row := range instanceRows {
		instanceValues = append(instanceValues, types.ObjectValueMust(instanceObjectAttrTypes, map[string]attr.Value{
			"instance_id":               row.InstanceID,
			"group_id":                  row.GroupID,
			"instance_container_ip":     row.InstanceContainerIP,
			"instance_container_status": row.InstanceContainerStatus,
			"instance_health_status":    row.InstanceHealthStatus,
			"main_container_status":     row.MainContainerStatus,
			"image_url":                 row.ImageURL,
			"package_version":           row.PackageVersion,
			"vswitch_id":                row.VSwitchID,
		}))
	}

	idValues := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		idValues = append(idValues, types.StringValue(id))
	}
	model.IDs = types.ListValueMust(types.StringType, idValues)
	model.ID = types.StringValue(applicationId)
	model.Instances = types.ListValueMust(types.ObjectType{AttrTypes: instanceObjectAttrTypes}, instanceValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if outputFile := model.OutputFile.ValueString(); outputFile != "" {
		if err := helper.WriteToFile(outputFile, results); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: writing output file", applicationInstancesTypeName), err.Error())
			return
		}
	}
}

func saeStringValue(v interface{}) types.String {
	if s, ok := v.(string); ok {
		return types.StringValue(s)
	}
	return types.StringNull()
}
