package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func ProviderContext(ctx context.Context, traceID, region string) context.Context {
	if traceID != "" {
		ctx = tflog.SetField(ctx, "alicloud_trace_id", traceID)
	}
	if region != "" {
		ctx = tflog.SetField(ctx, "alicloud_region", region)
	}
	return ctx
}

func (client *AliyunClient) ProviderContext(ctx context.Context) context.Context {
	if client == nil || client.config == nil {
		return ctx
	}
	return ProviderContext(ctx, client.config.TerraformTraceId, client.config.RegionId)
}

func RequestContext(ctx context.Context, service, operation string) context.Context {
	if service != "" {
		ctx = tflog.SetField(ctx, "alicloud_service", service)
	}
	return tflog.SetField(ctx, "alicloud_operation", operation)
}
