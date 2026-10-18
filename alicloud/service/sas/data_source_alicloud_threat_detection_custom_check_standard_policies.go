package sas

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const customCheckStandardPoliciesTypeName = "alicloud_threat_detection_custom_check_standard_policies"

var _ datasource.DataSource = (*customCheckStandardPoliciesDataSource)(nil)

func NewCustomCheckStandardPoliciesDataSource() datasource.DataSource {
	return &customCheckStandardPoliciesDataSource{}
}

type customCheckStandardPoliciesDataSource struct {
	fwadapt.DataSourceBase
}

type customCheckStandardPoliciesModel struct {
	Id                types.String  `tfsdk:"id"`
	PolicyType        types.String  `tfsdk:"policy_type"`
	DependentPolicyId types.String  `tfsdk:"dependent_policy_id"`
	Type              types.String  `tfsdk:"type"`
	NameRegex         types.String  `tfsdk:"name_regex"`
	OutputFile        types.String  `tfsdk:"output_file"`
	Ids               types.List    `tfsdk:"ids"`
	Names             types.List    `tfsdk:"names"`
	Policies          []policyModel `tfsdk:"policies"`
}

// policyModel is the element type of the `policies` ListNestedBlock. It is a
// block (not a nested attribute) so the schema converts cleanly under the
// provider's protocol v5 mux, which cannot represent attributes-with-attributes.
type policyModel struct {
	Id                types.String `tfsdk:"id"`
	PolicyId          types.String `tfsdk:"policy_id"`
	PolicyShowName    types.String `tfsdk:"policy_show_name"`
	PolicyType        types.String `tfsdk:"policy_type"`
	DependentPolicyId types.String `tfsdk:"dependent_policy_id"`
	Type              types.String `tfsdk:"type"`
	CheckType         types.String `tfsdk:"check_type"`
}

func (d *customCheckStandardPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = customCheckStandardPoliciesTypeName
}

func (d *customCheckStandardPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Threat Detection Custom Check Standard Policies of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, a hash of the returned policy IDs.",
			},
			"policy_type": schema.StringAttribute{
				Required:    true,
				Description: "The policy type of the custom check item rule. Valid values: `STANDARD`, `REQUIREMENT`, `SECTION`.",
				Validators: []validator.String{
					stringvalidator.OneOf("STANDARD", "REQUIREMENT", "SECTION"),
				},
			},
			"dependent_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the associated parent policy.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the major policy category. Valid values: `AISPM`, `IDENTITY_PERMISSION`, `RISK`, `COMPLIANCE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("AISPM", "IDENTITY_PERMISSION", "RISK", "COMPLIANCE"),
				},
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter the results by the `policy_show_name` of the policy.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results after running `terraform plan`.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of Custom Check Standard Policy IDs. Each ID is formatted as `<policy_id>:<policy_type>`.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of policy names.",
			},
		},
		Blocks: map[string]schema.Block{
			"policies": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The resource ID in Terraform, formatted as `<policy_id>:<policy_type>`.",
						},
						"policy_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the custom policy.",
						},
						"policy_show_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the custom policy.",
						},
						"policy_type": schema.StringAttribute{
							Computed:    true,
							Description: "The policy category type.",
						},
						"dependent_policy_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the associated parent policy.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the major policy category.",
						},
						"check_type": schema.StringAttribute{
							Computed:    true,
							Description: "The check type of the policy.",
						},
					},
				},
			},
		},
	}
}

func (d *customCheckStandardPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config customCheckStandardPoliciesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := listCustomCheckStandardPolicies(d.Client(), config.PolicyType.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckPolicies", customCheckStandardPoliciesTypeName),
			err.Error(),
		)
		return
	}

	dependentFilter := config.DependentPolicyId.ValueString()
	typeFilter := config.Type.ValueString()

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", customCheckStandardPoliciesTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = r
	}

	var idsFilter []string
	if !config.Ids.IsNull() {
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &idsFilter, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	idsSet := make(map[string]bool, len(idsFilter))
	for _, id := range idsFilter {
		idsSet[id] = true
	}

	out := make([]policyModel, 0, len(policies))
	ids := make([]string, 0, len(policies))
	names := make([]string, 0, len(policies))

	for _, p := range policies {
		if dependentFilter != "" && policyFieldString(p, "DependentPolicyId") != dependentFilter {
			continue
		}
		if typeFilter != "" && policyFieldString(p, "Type") != typeFilter {
			continue
		}

		showName := policyFieldString(p, "PolicyShowName")
		if nameRegex != nil && !nameRegex.MatchString(showName) {
			continue
		}

		policyId := policyFieldString(p, "PolicyId")
		composite := fmt.Sprintf("%s:%s", policyId, policyFieldString(p, "PolicyType"))
		if len(idsSet) > 0 && !idsSet[composite] && !idsSet[policyId] {
			continue
		}

		out = append(out, policyItemState(p))
		ids = append(ids, composite)
		names = append(names, showName)
	}

	idsList, idsDiags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(idsDiags...)
	namesList, namesDiags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(namesDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := customCheckStandardPoliciesModel{
		Id:                types.StringValue(dataSourceHash(ids)),
		PolicyType:        config.PolicyType,
		DependentPolicyId: config.DependentPolicyId,
		Type:              config.Type,
		NameRegex:         config.NameRegex,
		OutputFile:        config.OutputFile,
		Ids:               idsList,
		Names:             namesList,
		Policies:          out,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := os.WriteFile(config.OutputFile.ValueString(), []byte(fmt.Sprintf("%v", out)), 0o644); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output_file", customCheckStandardPoliciesTypeName),
				err.Error(),
			)
		}
	}
}

// dataSourceHash returns a stable hex id for the data source from the list of
// composite policy ids, mirroring the legacy dataResourceIdHash behavior.
func dataSourceHash(ids []string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", ids)))
	return hex.EncodeToString(sum[:])
}

// policyItemState maps a ListCheckPolicies element onto a policies block model.
func policyItemState(object map[string]interface{}) policyModel {
	str := func(key string) types.String {
		if s := policyFieldString(object, key); s != "" {
			return types.StringValue(s)
		}
		return types.StringNull()
	}

	m := policyModel{
		PolicyId:       str("PolicyId"),
		PolicyShowName: str("PolicyShowName"),
		PolicyType:     str("PolicyType"),
		CheckType:      str("CheckType"),
		Type:           str("Type"),
	}

	policyId := str("PolicyId")
	policyType := str("PolicyType")
	if !policyId.IsNull() && !policyType.IsNull() {
		m.Id = types.StringValue(fmt.Sprintf("%s:%s", policyId.ValueString(), policyType.ValueString()))
	}

	if dependent := policyFieldString(object, "DependentPolicyId"); dependent != "" && dependent != "0" {
		m.DependentPolicyId = types.StringValue(dependent)
	}

	return m
}
