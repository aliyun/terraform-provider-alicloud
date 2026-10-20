// Package paillmtrace contains the Terraform Plugin Framework data sources for
// the Alibaba Cloud PAI LLM Trace (PaiLLMTrace) product.
package paillmtrace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// evaluationTemplatesTypeName is the registered Terraform data source type.
	evaluationTemplatesTypeName = "alicloud_pai_llm_trace_evaluation_templates"

	// evaluationTemplatesPath is the ROA request path of GetEvaluationTemplates.
	// The API is public information (@ram enable:false), takes no business input
	// parameters and returns the EvaluationTemplates array.
	evaluationTemplatesPath = "/api/v1/PAILLMTrace/eval/templates"

	// evaluationTemplatesProduct is the gateway product code of PaiLLMTrace.
	evaluationTemplatesProduct = "PaiLLMTrace"

	// evaluationTemplatesVersion is the API version.
	evaluationTemplatesVersion = "2024-03-11"

	// evaluationTemplatesResponseKey is the response attribute carrying the
	// array of evaluation templates. The element type is free-form (any), so the
	// data source passes each template through as its JSON encoding.
	evaluationTemplatesResponseKey = "EvaluationTemplates"
)

var (
	_ datasource.DataSource              = &evaluationTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &evaluationTemplatesDataSource{}
)

// NewEvaluationTemplatesDataSource is the factory registered with the service
// package; it returns a pointer so the injected client lands on the base.
func NewEvaluationTemplatesDataSource() datasource.DataSource {
	return &evaluationTemplatesDataSource{}
}

type evaluationTemplatesDataSource struct {
	fwadapt.DataSourceBase
}

// evaluationTemplatesModel is the data source state. Templates holds one JSON
// string per returned template, preserving the free-form public payload and the
// collection order reported by the API.
type evaluationTemplatesModel struct {
	Id        types.String   `tfsdk:"id"`
	Templates []types.String `tfsdk:"templates"`
}

func (d *evaluationTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = evaluationTemplatesTypeName
}

func (d *evaluationTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the list of PAI LLM Trace evaluation templates returned by `GetEvaluationTemplates`. " +
			"The API exposes public information and requires no input parameters; every template is passed through as its JSON encoding, in the order returned by the API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier, set to the provider region.",
			},
			"templates": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				MarkdownDescription: "The evaluation templates, each encoded as a JSON string in the order returned by the API. " +
					"The element shape is free-form public information, so the raw payload is preserved verbatim.",
			},
		},
	}
}

func (d *evaluationTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state, diags := d.read(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// read performs the GetEvaluationTemplates call and builds the state. It is kept
// free of Framework state plumbing so the request path, response parsing,
// ordering, empty-set and failure behaviors are unit-testable without a running
// RPC server (which is what attaches the schema to resp.State before Read).
func (d *evaluationTemplatesDataSource) read(ctx context.Context) (evaluationTemplatesModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := d.Client()
	if client == nil {
		diags.AddError(
			fmt.Sprintf("Reading %s", evaluationTemplatesTypeName),
			"The provider client was not injected. Please report this issue to the provider developers.",
		)
		return evaluationTemplatesModel{}, diags
	}

	response, err := client.RoaGet(evaluationTemplatesProduct, evaluationTemplatesVersion, evaluationTemplatesPath, nil, nil, nil)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling GetEvaluationTemplates", evaluationTemplatesTypeName),
			err.Error(),
		)
		return evaluationTemplatesModel{}, diags
	}

	state := evaluationTemplatesModel{
		Id:        types.StringValue(client.RegionId),
		Templates: make([]types.String, 0),
	}

	if raw, ok := response[evaluationTemplatesResponseKey]; ok && raw != nil {
		if arr, ok := raw.([]interface{}); ok {
			for _, item := range arr {
				encoded, encErr := json.Marshal(item)
				if encErr != nil {
					// Keep the element order stable even if one item resists encoding;
					// fall back to its fmt representation rather than dropping it.
					encoded = []byte(fmt.Sprintf("%v", item))
				}
				state.Templates = append(state.Templates, types.StringValue(string(encoded)))
			}
		}
	}

	return state, diags
}
