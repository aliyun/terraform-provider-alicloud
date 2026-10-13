package cas

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// pcaCertSyncAction is a terraform-plugin-framework Action
// (alicloud_ssl_certificates_service_pca_cert_sync) that synchronizes a batch of
// PCA client certificates to the SSL Certificates service. Actions are imperative
// operations invoked from configuration (Terraform >= 1.14); unlike resources they
// do not manage state.

var _ action.Action = &pcaCertSyncAction{}

// NewPcaCertSyncAction returns the alicloud_ssl_certificates_service_pca_cert_sync action.
func NewPcaCertSyncAction() action.Action {
	return &pcaCertSyncAction{}
}

type pcaCertSyncAction struct {
	fwadapt.ActionBase
}

type pcaCertSyncConfigModel struct {
	Ids types.List `tfsdk:"ids"`
}

func (a *pcaCertSyncAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "alicloud_ssl_certificates_service_pca_cert_sync"
}

func (a *pcaCertSyncAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Synchronizes a batch of PCA client certificates to the SSL Certificates service.",
		Attributes: map[string]actionschema.Attribute{
			"ids": actionschema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The identifiers of the PCA client certificates to synchronize.",
			},
		},
	}
}

func (a *pcaCertSyncAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config pcaCertSyncConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := make([]string, 0, len(config.Ids.Elements()))
	resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := a.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			"Provider Client Not Injected",
			"The alicloud_ssl_certificates_service_pca_cert_sync action ended Configure without a provider client. Please report this issue to the provider developers.",
		)
		return
	}

	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Synchronizing %d PCA client certificate(s) to the SSL Certificates service", len(ids)),
		})
	}

	const actionName = "UploadPcaCertToCas"

	response, err := client.RpcPost("cas", "2020-06-30", actionName, map[string]interface{}{}, map[string]interface{}{
		"Ids": strings.Join(ids, ","),
	}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error invoking %s: calling %s", "alicloud_ssl_certificates_service_pca_cert_sync", actionName),
			err.Error(),
		)
		return
	}

	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Synchronization requested with RequestId %v", response["RequestId"]),
		})
	}
}
