package ens

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const networkInterfacesTypeName = "alicloud_ens_network_interfaces"

var (
	_ datasource.DataSource              = &networkInterfacesDataSource{}
	_ datasource.DataSourceWithConfigure = &networkInterfacesDataSource{}
)

func NewNetworkInterfacesDataSource() datasource.DataSource {
	return &networkInterfacesDataSource{}
}

type networkInterfacesDataSource struct {
	fwadapt.DataSourceBase
}

// networkInterfacesModel holds the full data source state. Computed list
// outputs use Go slices so the tfsdk reflection maps them to the
// ListAttribute / ListNestedAttribute in the schema directly.
type networkInterfacesModel struct {
	Ids                  types.List                       `tfsdk:"ids"`
	NameRegex            types.String                     `tfsdk:"name_regex"`
	EnsRegionId          types.String                     `tfsdk:"ens_region_id"`
	InstanceId           types.String                     `tfsdk:"instance_id"`
	NetworkId            types.String                     `tfsdk:"network_id"`
	NetworkInterfaceId   types.String                     `tfsdk:"network_interface_id"`
	NetworkInterfaceName types.String                     `tfsdk:"network_interface_name"`
	Status               types.String                     `tfsdk:"status"`
	VswitchId            types.String                     `tfsdk:"vswitch_id"`
	OutputFile           types.String                     `tfsdk:"output_file"`
	EnableDetails        types.Bool                       `tfsdk:"enable_details"`
	Names                []string                         `tfsdk:"names"`
	Interfaces           []networkInterfaceInterfaceModel `tfsdk:"interfaces"`
	Id                   types.String                     `tfsdk:"id"`
}

// networkInterfaceInterfaceModel is one element of the computed `interfaces`
// list. Field names mirror the SDKv2 data source so existing docs and configs
// keep working.
type networkInterfaceInterfaceModel struct {
	Id                   types.String `tfsdk:"id"`
	CreateTime           types.String `tfsdk:"create_time"`
	Description          types.String `tfsdk:"description"`
	EnsRegionId          types.String `tfsdk:"ens_region_id"`
	InstanceId           types.String `tfsdk:"instance_id"`
	MacAddress           types.String `tfsdk:"mac_address"`
	NetworkId            types.String `tfsdk:"network_id"`
	NetworkInterfaceId   types.String `tfsdk:"network_interface_id"`
	NetworkInterfaceName types.String `tfsdk:"network_interface_name"`
	PrimaryIp            types.String `tfsdk:"primary_ip"`
	PrimaryIpType        types.String `tfsdk:"primary_ip_type"`
	SecurityGroupIds     types.List   `tfsdk:"security_group_ids"`
	Status               types.String `tfsdk:"status"`
	VswitchId            types.String `tfsdk:"vswitch_id"`
	VmncLearn            types.Bool   `tfsdk:"vmnc_learn"`
}

func (d *networkInterfacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = networkInterfacesTypeName
}

func (d *networkInterfacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of ENS Network Interfaces owned by an Alibaba Cloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is a hash of the returned NetworkInterface IDs.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of NetworkInterface IDs. Also used as a filter: only matching IDs are returned.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of NetworkInterface names.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by NetworkInterface name.",
			},
			"ens_region_id": schema.StringAttribute{
				Optional:    true,
				Description: "The node ID of ENS.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the instance bound to the ENI.",
			},
			"network_id": schema.StringAttribute{
				Optional:    true,
				Description: "The network ID.",
			},
			"network_interface_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the ENI.",
			},
			"network_interface_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the ENI.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Status of the network card.",
			},
			"vswitch_id": schema.StringAttribute{
				Optional:    true,
				Description: "The vSwitch ID.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results (after running terraform plan).",
			},
			"enable_details": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set to true to populate the per-interface security_group_ids. Default false.",
			},
			"interfaces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the NetworkInterface.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the ENI.",
						},
						"ens_region_id": schema.StringAttribute{
							Computed:    true,
							Description: "The node ID of ENS.",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the instance bound to the ENI.",
						},
						"mac_address": schema.StringAttribute{
							Computed:    true,
							Description: "The MAC address of the ENI.",
						},
						"network_id": schema.StringAttribute{
							Computed:    true,
							Description: "The network ID.",
						},
						"network_interface_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the ENI.",
						},
						"network_interface_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the ENI.",
						},
						"primary_ip": schema.StringAttribute{
							Computed:    true,
							Description: "The primary private network IP address.",
						},
						"primary_ip_type": schema.StringAttribute{
							Computed:    true,
							Description: "The primary IP address type.",
						},
						"security_group_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The IDs of the security groups. Populated only when enable_details is true.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the network card.",
						},
						"vswitch_id": schema.StringAttribute{
							Computed:    true,
							Description: "The vSwitch ID.",
						},
						"vmnc_learn": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to enable NIC route learning.",
						},
					},
				},
			},
		},
	}
}

func (d *networkInterfacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config networkInterfacesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); v != "" {
		r, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", networkInterfacesTypeName),
				errs.WrapError(err).Error(),
			)
			return
		}
		nameRegex = r
	}

	idFilter := map[string]bool{}
	for _, el := range config.Ids.Elements() {
		if s, ok := el.(types.String); ok && s.ValueString() != "" {
			idFilter[s.ValueString()] = true
		}
	}

	var selected []map[string]interface{}
	pageNumber := 1
	for {
		request := map[string]interface{}{
			"PageSize":   ensPageSize,
			"PageNumber": pageNumber,
		}
		if v := config.NetworkInterfaceId.ValueString(); v != "" {
			request["NetworkInterfaceId"] = v
		}
		if v := config.EnsRegionId.ValueString(); v != "" {
			request["EnsRegionId"] = v
		}
		if v := config.InstanceId.ValueString(); v != "" {
			request["InstanceId"] = v
		}
		if v := config.NetworkId.ValueString(); v != "" {
			request["NetworkId"] = v
		}
		if v := config.NetworkInterfaceName.ValueString(); v != "" {
			request["NetworkInterfaceName"] = v
		}
		if v := config.Status.ValueString(); v != "" {
			request["Status"] = v
		}
		if v := config.VswitchId.ValueString(); v != "" {
			request["VSwitchId"] = v
		}

		response, err := d.Client().RpcPost(ensProduct, ensVersion, "DescribeNetworkInterfaces", nil, request, true)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", networkInterfacesTypeName),
				errs.WrapErrorf(err, "Reading %s: calling DescribeNetworkInterfaces", networkInterfacesTypeName).Error(),
			)
			return
		}

		list := networkInterfaceList(response)
		for _, raw := range list {
			item, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["NetworkInterfaceName"])
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}
			if len(idFilter) > 0 && !idFilter[fmt.Sprint(item["NetworkInterfaceId"])] {
				continue
			}
			if config.EnableDetails.ValueBool() {
				detail, err := describeNetworkInterface(d.Client(), fmt.Sprint(item["NetworkInterfaceId"]))
				if err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", networkInterfacesTypeName),
						errs.WrapErrorf(err, "Reading %s: detailing %s", networkInterfacesTypeName, fmt.Sprint(item["NetworkInterfaceId"])).Error(),
					)
					return
				}
				if detail != nil {
					item = detail
				}
			}
			selected = append(selected, item)
		}

		if len(list) < ensPageSize {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0, len(selected))
	names := make([]string, 0, len(selected))
	interfaces := make([]networkInterfaceInterfaceModel, 0, len(selected))
	for _, item := range selected {
		id := fmt.Sprint(item["NetworkInterfaceId"])
		ids = append(ids, id)
		names = append(names, fmt.Sprint(item["NetworkInterfaceName"]))
		interfaces = append(interfaces, networkInterfaceInterfaceState(ctx, item))
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := networkInterfacesModel{
		Ids:                  idsList,
		NameRegex:            config.NameRegex,
		EnsRegionId:          config.EnsRegionId,
		InstanceId:           config.InstanceId,
		NetworkId:            config.NetworkId,
		NetworkInterfaceId:   config.NetworkInterfaceId,
		NetworkInterfaceName: config.NetworkInterfaceName,
		Status:               config.Status,
		VswitchId:            config.VswitchId,
		OutputFile:           config.OutputFile,
		EnableDetails:        config.EnableDetails,
		Names:                names,
		Interfaces:           interfaces,
		Id:                   types.StringValue(dataResourceIdHash(ids)),
	}

	if output := config.OutputFile.ValueString(); output != "" {
		_ = os.WriteFile(output, []byte(fmt.Sprintf("%v", selected)), 0o644)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// networkInterfaceInterfaceState maps one DescribeNetworkInterfaces element to
// the nested `interfaces` model. Empty API strings become null.
func networkInterfaceInterfaceState(ctx context.Context, object map[string]interface{}) networkInterfaceInterfaceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	securityGroupIds, diags := types.ListValueFrom(ctx, types.StringType, securityGroupIDs(object))
	if diags.HasError() {
		securityGroupIds = types.ListNull(types.StringType)
	}
	var vmnc types.Bool
	if v, ok := object["VmncLearn"].(bool); ok {
		vmnc = types.BoolValue(v)
	} else {
		vmnc = types.BoolNull()
	}
	return networkInterfaceInterfaceModel{
		Id:                   str("NetworkInterfaceId"),
		CreateTime:           str("CreationTime"),
		Description:          str("Description"),
		EnsRegionId:          str("EnsRegionId"),
		InstanceId:           str("InstanceId"),
		MacAddress:           str("MacAddress"),
		NetworkId:            str("NetworkId"),
		NetworkInterfaceId:   str("NetworkInterfaceId"),
		NetworkInterfaceName: str("NetworkInterfaceName"),
		PrimaryIp:            str("PrimaryIp"),
		PrimaryIpType:        str("PrimaryIpType"),
		SecurityGroupIds:     securityGroupIds,
		Status:               str("Status"),
		VswitchId:            str("VSwitchId"),
		VmncLearn:            vmnc,
	}
}

// dataResourceIdHash returns a stable hex hash of the joined NetworkInterface
// IDs, mirroring the SDKv2 dataResourceIdHash identifier used by ENS data
// sources.
func dataResourceIdHash(ids []string) string {
	h := sha256.Sum256([]byte(strings.Join(ids, "|")))
	return hex.EncodeToString(h[:])
}
