package ecs

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const savingPlansTypeName = "alicloud_ecs_saving_plans"

var _ datasource.DataSource = &savingPlansDataSource{}

func NewSavingPlansDataSource() datasource.DataSource {
	return &savingPlansDataSource{}
}

type savingPlansDataSource struct {
	fwadapt.DataSourceBase
}

type savingPlansModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	Ids        types.List   `tfsdk:"ids"`
	NameRegex  types.String `tfsdk:"name_regex"`
	Names      types.List   `tfsdk:"names"`
	OutputFile types.String `tfsdk:"output_file"`
	Plans      types.List   `tfsdk:"plans"`
}

type savingPlanItemModel struct {
	Id              types.String `tfsdk:"id"`
	SavingPlanId    types.String `tfsdk:"saving_plan_id"`
	CommittedAmount types.String `tfsdk:"committed_amount"`
	CreateTime      types.String `tfsdk:"create_time"`
	InstanceFamily  types.String `tfsdk:"instance_family"`
	OfferingType    types.String `tfsdk:"offering_type"`
	PaymentType     types.String `tfsdk:"payment_type"`
	Period          types.Int64  `tfsdk:"period"`
	PlanType        types.String `tfsdk:"plan_type"`
	RegionId        types.String `tfsdk:"region_id"`
	StartTime       types.String `tfsdk:"start_time"`
}

func (d *savingPlansDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = savingPlansTypeName
}

func (d *savingPlansDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of ECS Saving Plans available to the current user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID, set to the region ID.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the saving plan instance. If specified, only the matching plan will be returned.",
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of saving plan IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter saving plans by name.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of saving plan names.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the result after running `terraform plan`.",
			},
			"plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the saving plan, same as `saving_plan_id`.",
						},
						"saving_plan_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the saving plan.",
						},
						"committed_amount": schema.StringAttribute{
							Computed:    true,
							Description: "The hourly commitment amount of the saving plan.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "The creation time of the saving plan.",
						},
						"instance_family": schema.StringAttribute{
							Computed:    true,
							Description: "The instance family of the saving plan.",
						},
						"offering_type": schema.StringAttribute{
							Computed:    true,
							Description: "The payment method of the saving plan.",
						},
						"payment_type": schema.StringAttribute{
							Computed:    true,
							Description: "The payment type of the saving plan.",
						},
						"period": schema.Int64Attribute{
							Computed:    true,
							Description: "The period of the saving plan.",
						},
						"plan_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the saving plan.",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "The region ID of the saving plan.",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "The effective start time of the saving plan.",
						},
					},
				},
				Description: "A list of saving plans.",
			},
		},
	}
}

func (d *savingPlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config savingPlansModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := ""
	if !config.InstanceId.IsNull() {
		instanceId = config.InstanceId.ValueString()
	}

	items, err := querySavingPlans(d.Client(), instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlansTypeName),
			err.Error(),
		)
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() {
		if r, e := regexp.Compile(config.NameRegex.ValueString()); e == nil {
			nameRegex = r
		}
	}

	var ids, names []string
	var planItems []savingPlanItemModel
	var outputPlans []map[string]interface{}
	for _, item := range items {
		savingPlanId := ""
		if v, ok := item["InstanceId"]; ok {
			savingPlanId = fmt.Sprint(v)
		}
		if savingPlanId == "" {
			continue
		}

		name := ""
		if v, ok := item["Name"]; ok {
			name = fmt.Sprint(v)
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		planItems = append(planItems, savingPlanItemModel{
			Id:              types.StringValue(savingPlanId),
			SavingPlanId:    types.StringValue(savingPlanId),
			CommittedAmount: apiStr(item, "PoolValue"),
			CreateTime:      apiStr(item, "StartTimestamp"),
			InstanceFamily:  apiStr(item, "InstanceFamily"),
			OfferingType:    apiStr(item, "PayMode"),
			PaymentType:     apiStr(item, "PaymentType"),
			Period:          periodFromCycle(item["Cycle"]),
			PlanType:        apiStr(item, "SavingsType"),
			RegionId:        apiStr(item, "Region"),
			StartTime:       apiStr(item, "StartTime"),
		})
		outputPlans = append(outputPlans, savingPlanOutputMap(item, savingPlanId))
		ids = append(ids, savingPlanId)
		names = append(names, name)
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)

	planObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":               types.StringType,
			"saving_plan_id":   types.StringType,
			"committed_amount": types.StringType,
			"create_time":      types.StringType,
			"instance_family":  types.StringType,
			"offering_type":    types.StringType,
			"payment_type":     types.StringType,
			"period":           types.Int64Type,
			"plan_type":        types.StringType,
			"region_id":        types.StringType,
			"start_time":       types.StringType,
		},
	}
	plansList, diags := types.ListValueFrom(ctx, planObjType, planItems)
	resp.Diagnostics.Append(diags...)

	state := savingPlansModel{
		Id:         types.StringValue(d.Client().RegionId),
		InstanceId: config.InstanceId,
		Ids:        idsList,
		NameRegex:  config.NameRegex,
		Names:      namesList,
		OutputFile: config.OutputFile,
		Plans:      plansList,
	}

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := writeToFile(config.OutputFile.ValueString(), outputPlans); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output file", savingPlansTypeName),
				err.Error(),
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
