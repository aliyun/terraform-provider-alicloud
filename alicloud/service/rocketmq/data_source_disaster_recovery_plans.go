package rocketmq

import (
	"context"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const disasterRecoveryPlansDataSourceTypeName = "alicloud_rocketmq_disaster_recovery_plans"

var _ datasource.DataSource = (*disasterRecoveryPlansDataSource)(nil)

func NewDisasterRecoveryPlansDataSource() datasource.DataSource {
	return &disasterRecoveryPlansDataSource{}
}

type disasterRecoveryPlansDataSource struct {
	fwadapt.DataSourceBase
}

type disasterRecoveryPlansDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Ids        types.List   `tfsdk:"ids"`
	InstanceId types.String `tfsdk:"instance_id"`
	NameRegex  types.String `tfsdk:"name_regex"`
	OutputFile types.String `tfsdk:"output_file"`
	Plans      types.List   `tfsdk:"plans"`
}

// planItemModel is a single element of the plans list. It mirrors the resource
// attributes (all computed) so callers can interpolate the same field names.
type planItemModel struct {
	AutoSyncCheckpoint    types.Bool   `tfsdk:"auto_sync_checkpoint"`
	CreateTime            types.String `tfsdk:"create_time"`
	Instances             types.List   `tfsdk:"instances"`
	PlanDesc              types.String `tfsdk:"plan_desc"`
	PlanId                types.String `tfsdk:"plan_id"`
	PlanName              types.String `tfsdk:"plan_name"`
	PlanType              types.String `tfsdk:"plan_type"`
	RegionId              types.String `tfsdk:"region_id"`
	Status                types.String `tfsdk:"status"`
	SyncCheckpointEnabled types.Bool   `tfsdk:"sync_checkpoint_enabled"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

func (d *disasterRecoveryPlansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = disasterRecoveryPlansDataSourceTypeName
}

func (d *disasterRecoveryPlansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of RocketMQ Disaster Recovery Plan items in the current region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"instance_id": schema.StringAttribute{Optional: true},
			"name_regex":  schema.StringAttribute{Optional: true},
			"output_file": schema.StringAttribute{Optional: true},
			"plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auto_sync_checkpoint": schema.BoolAttribute{Computed: true},
						"create_time":          schema.StringAttribute{Computed: true},
						"instances": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: dataSourceInstanceAttributes(),
							},
						},
						"plan_desc":               schema.StringAttribute{Computed: true},
						"plan_id":                 schema.StringAttribute{Computed: true},
						"plan_name":               schema.StringAttribute{Computed: true},
						"plan_type":               schema.StringAttribute{Computed: true},
						"region_id":               schema.StringAttribute{Computed: true},
						"status":                  schema.StringAttribute{Computed: true},
						"sync_checkpoint_enabled": schema.BoolAttribute{Computed: true},
						"update_time":             schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// dataSourceInstanceAttributes returns the instance block attributes for the
// data source plans. Every field is computed: the data source never writes.
func dataSourceInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"auth_type":         schema.StringAttribute{Computed: true},
		"consumer_group_id": schema.StringAttribute{Computed: true},
		"endpoint_url":      schema.StringAttribute{Computed: true},
		"instance_id":       schema.StringAttribute{Computed: true},
		"instance_role":     schema.StringAttribute{Computed: true},
		"instance_type":     schema.StringAttribute{Computed: true},
		"message_property": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"property_key":   schema.StringAttribute{Computed: true},
					"property_value": schema.StringAttribute{Computed: true},
				},
			},
		},
		"network_type":      schema.StringAttribute{Computed: true},
		"password":          schema.StringAttribute{Computed: true},
		"region_id":         schema.StringAttribute{Computed: true},
		"security_group_id": schema.StringAttribute{Computed: true},
		"username":          schema.StringAttribute{Computed: true},
		"vpc_id":            schema.StringAttribute{Computed: true},
		"vswitch_id":        schema.StringAttribute{Computed: true},
	}
}

func (d *disasterRecoveryPlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config disasterRecoveryPlansDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	objects, err := listDisasterRecoveryPlans(client, config.InstanceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListDisasterRecoveryPlans", disasterRecoveryPlansDataSourceTypeName),
			err.Error(),
		)
		return
	}

	var nameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); v != "" {
		r, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError("Compiling name_regex", err.Error())
			return
		}
		nameRegex = r
	}

	idsSet := make(map[string]struct{})
	for _, id := range config.Ids.Elements() {
		idsSet[id.(types.String).ValueString()] = struct{}{}
	}

	planItems, ids := filterPlans(ctx, resp, client.RegionId, objects, nameRegex, idsSet)

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	plansList, diags := types.ListValueFrom(ctx, plansElementType(), planItems)
	resp.Diagnostics.Append(diags...)

	config.Id = types.StringValue(dataResourceIdHash(ids))
	config.Ids = idsList
	config.Plans = plansList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// filterPlans applies the client-side name_regex and ids filters to the raw
// list response and projects the survivors onto plan models. The data source
// filters client-side because ListDisasterRecoveryPlans exposes no server
// filter beyond instanceId.
func filterPlans(ctx context.Context, resp *datasource.ReadResponse, regionId string, objects []map[string]interface{}, nameRegex *regexp.Regexp, idsSet map[string]struct{}) ([]planItemModel, []string) {
	planItems := make([]planItemModel, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, obj := range objects {
		planId := fmt.Sprint(obj["planId"])
		planName := fmt.Sprint(obj["planName"])
		if nameRegex != nil && !nameRegex.MatchString(planName) {
			continue
		}
		if len(idsSet) > 0 {
			if _, ok := idsSet[planId]; !ok {
				continue
			}
		}
		instances := flattenInstances(obj["instances"], nil)
		instancesList, diags := types.ListValueFrom(ctx, instancesElementType(), instances)
		resp.Diagnostics.Append(diags...)
		planItems = append(planItems, planItemModel{
			AutoSyncCheckpoint:    boolFrom(obj["autoSyncCheckpoint"]),
			CreateTime:            strFrom(obj["createTime"]),
			Instances:             instancesList,
			PlanDesc:              strFrom(obj["planDesc"]),
			PlanId:                types.StringValue(planId),
			PlanName:              strFrom(obj["planName"]),
			PlanType:              strFrom(obj["planType"]),
			RegionId:              strFrom(regionId),
			Status:                strFrom(obj["planStatus"]),
			SyncCheckpointEnabled: boolFrom(obj["syncCheckpointEnabled"]),
			UpdateTime:            strFrom(obj["updateTime"]),
		})
		ids = append(ids, planId)
	}
	return planItems, ids
}

// plansElementType is the Framework object type of a plans list element,
// matching the planItemModel schema.
func plansElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"auto_sync_checkpoint":    types.BoolType,
			"create_time":             types.StringType,
			"instances":               types.ListType{ElemType: instancesElementType()},
			"plan_desc":               types.StringType,
			"plan_id":                 types.StringType,
			"plan_name":               types.StringType,
			"plan_type":               types.StringType,
			"region_id":               types.StringType,
			"status":                  types.StringType,
			"sync_checkpoint_enabled": types.BoolType,
			"update_time":             types.StringType,
		},
	}
}

// dataResourceIdHash mirrors the SDKv2 helper.Hashcode over a joined id list so
// the data source id matches the legacy Plugin SDK data source output.
func dataResourceIdHash(ids []string) string {
	var buf []byte
	sep := []byte("-")
	for _, id := range ids {
		buf = append(buf, id...)
		buf = append(buf, sep...)
	}
	h := fnv.New32a()
	_, _ = h.Write(buf)
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
