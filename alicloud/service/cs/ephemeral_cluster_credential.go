package cs

import (
	"context"
	"fmt"
	"time"

	roaCS "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/retry"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const clusterCredentialEphemeralTypeName = "alicloud_cs_cluster_credential"

var (
	_ ephemeral.EphemeralResource              = &clusterCredentialEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &clusterCredentialEphemeral{}
)

func NewClusterCredentialEphemeralResource() ephemeral.EphemeralResource {
	return &clusterCredentialEphemeral{}
}

type clusterCredentialEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type clusterCredentialEphemeralModel struct {
	ClusterId                types.String `tfsdk:"cluster_id"`
	TemporaryDurationMinutes types.Int64  `tfsdk:"temporary_duration_minutes"`
	KubeConfig               types.String `tfsdk:"kube_config"`
	Expiration               types.String `tfsdk:"expiration"`
}

func (e *clusterCredentialEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = clusterCredentialEphemeralTypeName
}

func (e *clusterCredentialEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"temporary_duration_minutes": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(15, 4320),
				},
			},
			"kube_config": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"expiration": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *clusterCredentialEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "DescribeClusterUserKubeconfig"

	var config clusterCredentialEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roaClient, err := e.Client().NewRoaCsClient()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: building the CS client", clusterCredentialEphemeralTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}

	request := &roaCS.DescribeClusterUserKubeconfigRequest{
		PrivateIpAddress: tea.Bool(false),
	}
	if !config.TemporaryDurationMinutes.IsNull() {
		if minutes := config.TemporaryDurationMinutes.ValueInt64(); minutes > 0 {
			request.TemporaryDurationMinutes = tea.Int64(minutes)
		}
	}

	clusterId := config.ClusterId.ValueString()
	wait := retry.IncrementalWait(3*time.Second, 3*time.Second)
	var credential *roaCS.DescribeClusterUserKubeconfigResponseBody
	err = sdkretry.Retry(5*time.Minute, func() *sdkretry.RetryError {
		response, callErr := roaClient.DescribeClusterUserKubeconfig(tea.String(clusterId), request)
		if callErr != nil {
			if retry.NeedRetry(callErr) {
				wait()
				return sdkretry.RetryableError(callErr)
			}
			return sdkretry.NonRetryableError(callErr)
		}
		credential = response.Body
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", clusterCredentialEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", clusterCredentialEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	kubeConfig := tea.StringValue(credential.Config)
	if kubeConfig == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", clusterCredentialEphemeralTypeName, action),
			"The response carried no kubeconfig",
		)
		return
	}
	state.KubeConfig = types.StringValue(kubeConfig)

	if credential.Expiration != nil {
		state.Expiration = types.StringValue(*credential.Expiration)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
