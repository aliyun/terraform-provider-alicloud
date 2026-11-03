package polardb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const networkChannelsTypeName = "alicloud_polardb_network_channels"

var (
	_ datasource.DataSource              = &networkChannelsDataSource{}
	_ datasource.DataSourceWithConfigure = &networkChannelsDataSource{}
)

func NewNetworkChannelsDataSource() datasource.DataSource {
	return &networkChannelsDataSource{}
}

type networkChannelsDataSource struct {
	fwadapt.DataSourceBase
}

type networkChannelsModel struct {
	Id           types.String               `tfsdk:"id"`
	DBClusterId  types.String               `tfsdk:"db_cluster_id"`
	ChannelName  types.String               `tfsdk:"channel_name"`
	VpcId        types.String               `tfsdk:"vpc_id"`
	NameRegex    types.String               `tfsdk:"name_regex"`
	Ids          types.List                 `tfsdk:"ids"`
	Names        types.List                 `tfsdk:"names"`
	ChannelInfos []networkChannelsInfoModel `tfsdk:"channel_infos"`
}

type networkChannelsInfoModel struct {
	ChannelName       types.String `tfsdk:"channel_name"`
	DBClusterId       types.String `tfsdk:"db_cluster_id"`
	TargetDBClusterId types.String `tfsdk:"target_db_cluster_id"`
	TargetIp          types.String `tfsdk:"target_ip"`
	TargetPort        types.String `tfsdk:"target_port"`
	TargetType        types.String `tfsdk:"target_type"`
	VpcId             types.String `tfsdk:"vpc_id"`
	Notes             types.String `tfsdk:"notes"`
	RegionId          types.String `tfsdk:"region_id"`
}

var networkChannelsInfoObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"channel_name":         types.StringType,
	"db_cluster_id":        types.StringType,
	"target_db_cluster_id": types.StringType,
	"target_ip":            types.StringType,
	"target_port":          types.StringType,
	"target_type":          types.StringType,
	"vpc_id":               types.StringType,
	"notes":                types.StringType,
	"region_id":            types.StringType,
}}

func (d *networkChannelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = networkChannelsTypeName
}

func (d *networkChannelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the PolarDB cluster.",
			},
			"db_cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the cluster.",
			},
			"channel_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the network channel.",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the VPC to which the access endpoint resides.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by channel name.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of network channel IDs. The value is formulated as <db_cluster_id>:<channel_name>.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of network channel names.",
			},
			"channel_infos": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"channel_name":         schema.StringAttribute{Computed: true, Description: "Channel Name."},
						"db_cluster_id":        schema.StringAttribute{Computed: true, Description: "Source cluster ID."},
						"target_db_cluster_id": schema.StringAttribute{Computed: true, Description: "Target cluster ID."},
						"target_ip":            schema.StringAttribute{Computed: true, Description: "Target cluster IP."},
						"target_port":          schema.StringAttribute{Computed: true, Description: "Target cluster port."},
						"target_type":          schema.StringAttribute{Computed: true, Description: "Target cluster type."},
						"vpc_id":               schema.StringAttribute{Computed: true, Description: "VPC ID."},
						"notes":                schema.StringAttribute{Computed: true, Description: "Notes."},
						"region_id":            schema.StringAttribute{Computed: true, Description: "Region ID."},
					},
				},
				Computed:    true,
				Description: "The list of network channels.",
			},
		},
	}
}

func (d *networkChannelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config networkChannelsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId":    d.Meta().RegionId,
		"DBClusterId": config.DBClusterId.ValueString(),
	}
	if !config.ChannelName.IsNull() && !config.ChannelName.IsUnknown() {
		request["ChannelName"] = config.ChannelName.ValueString()
	}
	if !config.VpcId.IsNull() && !config.VpcId.IsUnknown() {
		request["VpcId"] = config.VpcId.ValueString()
	}

	response, err := retryRpcPost(d.Meta(), "polardb", "2017-08-01", "DescribeNetworkChannel", request, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading alicloud_polardb_network_channels",
			errs.WrapErrorf(err, "Reading alicloud_polardb_network_channels: calling DescribeNetworkChannel").Error(),
		)
		return
	}

	// Navigate to ChannelInfos (direct array per OpenAPI definition).
	var rawChannels []map[string]interface{}
	if channelInfos, ok := response["ChannelInfos"]; ok {
		switch v := channelInfos.(type) {
		case []interface{}:
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					rawChannels = append(rawChannels, m)
				}
			}
		case map[string]interface{}:
			// Some RPC responses wrap arrays; handle both.
			for _, val := range v {
				if arr, ok := val.([]interface{}); ok {
					for _, item := range arr {
						if m, ok := item.(map[string]interface{}); ok {
							rawChannels = append(rawChannels, m)
						}
					}
				}
			}
			if len(rawChannels) == 0 {
				rawChannels = append(rawChannels, v)
			}
		}
	}

	// Apply name_regex filter.
	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && !config.NameRegex.IsUnknown() {
		nameRegex, err = regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid name_regex", err.Error())
			return
		}
	}

	var channelInfos []networkChannelsInfoModel
	var ids []attr.Value
	var names []attr.Value
	for _, ch := range rawChannels {
		channelName := fwStringFromMap(ch, "ChannelName")
		if nameRegex != nil && !nameRegex.MatchString(channelName) {
			continue
		}
		dbClusterId := fwStringFromMap(ch, "DBClusterId")
		if dbClusterId == "" {
			dbClusterId = config.DBClusterId.ValueString()
		}
		channelInfos = append(channelInfos, networkChannelsInfoModel{
			ChannelName:       types.StringValue(channelName),
			DBClusterId:       types.StringValue(dbClusterId),
			TargetDBClusterId: types.StringValue(fwStringFromMap(ch, "TargetDBClusterId")),
			TargetIp:          types.StringValue(fwStringFromMap(ch, "TargetIp")),
			TargetPort:        types.StringValue(fwStringFromMap(ch, "TargetPort")),
			TargetType:        types.StringValue(fwStringFromMap(ch, "TargetType")),
			VpcId:             types.StringValue(fwStringFromMap(ch, "VpcId")),
			Notes:             types.StringValue(fwStringFromMap(ch, "Notes")),
			RegionId:          types.StringValue(fwStringFromMap(ch, "RegionId")),
		})
		ids = append(ids, types.StringValue(fmt.Sprintf("%s:%s", dbClusterId, channelName)))
		names = append(names, types.StringValue(channelName))
	}

	config.Id = types.StringValue(config.DBClusterId.ValueString())
	if len(ids) == 0 {
		config.Ids = types.ListNull(types.StringType)
		config.Names = types.ListNull(types.StringType)
	} else {
		config.Ids, resp.Diagnostics = types.ListValue(types.StringType, ids)
		config.Names, resp.Diagnostics = types.ListValue(types.StringType, names)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	config.ChannelInfos = channelInfos
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fwStringFromMap safely extracts a string from a map[string]interface{}.
func fwStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}
