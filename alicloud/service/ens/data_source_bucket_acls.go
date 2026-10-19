package ens

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const bucketAclsDataSourceTypeName = "alicloud_ens_bucket_acls"

var (
	_ datasource.DataSource              = &bucketAclsDataSource{}
	_ datasource.DataSourceWithConfigure = &bucketAclsDataSource{}
)

func NewBucketAclsDataSource() datasource.DataSource {
	return &bucketAclsDataSource{}
}

type bucketAclsDataSource struct {
	fwadapt.DataSourceBase
}

type bucketAclsDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	BucketName types.String `tfsdk:"bucket_name"`
	Ids        types.List   `tfsdk:"ids"`
	BucketAcls types.List   `tfsdk:"bucket_acls"`
}

type bucketAclEntry struct {
	Id         types.String `tfsdk:"id"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketAcl  types.String `tfsdk:"bucket_acl"`
}

// bucketAclEntryType is the framework object type for a bucketAclEntry, used
// to build typed (including empty) lists of entries.
var bucketAclEntryType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"bucket_name": types.StringType,
		"bucket_acl":  types.StringType,
	},
}

func (d *bucketAclsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = bucketAclsDataSourceTypeName
}

func (d *bucketAclsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the ACL of an ENS (Edge Node Service) bucket.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Bucket Acl. The value is the bucket name.",
			},
			"bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ENS bucket to query the ACL for.",
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Bucket Acl IDs. The ID is the bucket name.",
			},
			// bucket_acls is modeled as a ListAttribute of a composite Object
			// element type rather than a ListNestedAttribute. The provider is
			// assembled behind a protocol v5 mux server (SDKv2 + Framework); a
			// ListNestedAttribute serializes nested attributes into the v5 schema,
			// which the v5 mux cannot represent ("protocol version 5 cannot have
			// Attributes set"). A single ListAttribute whose element type is an
			// opaque Object encodes the same shape as a composite type without
			// nested schema attributes, so the list still exposes
			// bucket_acls.0.{id,bucket_name,bucket_acl} to configurations.
			"bucket_acls": schema.ListAttribute{
				Computed:    true,
				ElementType: bucketAclEntryType,
				Description: "A list of ENS Bucket Acls. Each element contains the following attributes.",
			},
		},
	}
}

func (d *bucketAclsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config bucketAclsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := config.BucketName.ValueString()

	acl, err := readBucketAclStable(d.Client(), bucketName, bucketAclReadTimeout)
	if err != nil {
		if isBucketAclNotFound(err) {
			// Match the legacy semantics: a missing bucket reads back as an
			// empty list with the data source id set to the requested name.
			state := bucketAclsDataSourceModel{
				Id:         types.StringValue(bucketName),
				BucketName: types.StringValue(bucketName),
				Ids:        types.ListNull(types.StringType),
				BucketAcls: types.ListNull(bucketAclEntryType),
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketAcl", bucketAclsDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketAcl for bucket %s", bucketAclsDataSourceTypeName, bucketName).Error(),
		)
		return
	}

	ids, diags := types.ListValueFrom(ctx, types.StringType, []string{bucketName})
	resp.Diagnostics.Append(diags...)
	entries, diags := types.ListValueFrom(ctx, bucketAclEntryType, []bucketAclEntry{{
		Id:         types.StringValue(bucketName),
		BucketName: types.StringValue(bucketName),
		BucketAcl:  types.StringValue(acl),
	}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := bucketAclsDataSourceModel{
		Id:         types.StringValue(bucketName),
		BucketName: types.StringValue(bucketName),
		Ids:        ids,
		BucketAcls: entries,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
