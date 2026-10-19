package eflo

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const vccRouteEntriesTypeName = "alicloud_eflo_vcc_route_entries"

var _ datasource.DataSource = (*vccRouteEntriesDataSource)(nil)

func NewVccRouteEntriesDataSource() datasource.DataSource {
	return &vccRouteEntriesDataSource{}
}

type vccRouteEntriesDataSource struct {
	fwadapt.DataSourceBase
}

// vccRouteEntriesModel is the data source state shape. The `ids` list carries the
// composite identity "<vcc_id>:<vcc_route_entry_id>" of every matched route entry,
// matching the resource `id` format so an entry can be fed straight into `terraform import`.
// The filters mirror the ListVccRouteEntries request parameters (VccId required).
type vccRouteEntriesModel struct {
	Id                   types.String `tfsdk:"id"`
	VccId                types.String `tfsdk:"vcc_id"`
	DestinationCidrBlock types.String `tfsdk:"destination_cidr_block"`
	Status               types.String `tfsdk:"status"`
	RouteType            types.String `tfsdk:"route_type"`
	NextHopId            types.String `tfsdk:"next_hop_id"`
	NextHopType          types.String `tfsdk:"next_hop_type"`
	RegionId             types.String `tfsdk:"region_id"`
	Ids                  types.List   `tfsdk:"ids"`
}

func (d *vccRouteEntriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = vccRouteEntriesTypeName
}

func (d *vccRouteEntriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Eflo Vcc Route Entry list indexed by a VCC instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A stable identifier for this data source result, derived from the matched ids.",
			},
			"vcc_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Eflo VCC instance. Required by ListVccRouteEntries.",
			},
			"destination_cidr_block": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the destination CIDR block of the route entry.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the status of the route entry, e.g. `Available`.",
			},
			"route_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the route type, e.g. `System`, `BGP`.",
			},
			"next_hop_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the next hop instance ID.",
			},
			"next_hop_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the next hop type, e.g. `VCC`.",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region ID. Defaults to the provider region.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The IDs of the matched route entries, each formatted as `<vcc_id>:<vcc_route_entry_id>`.",
			},
		},
	}
}

func (d *vccRouteEntriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config vccRouteEntriesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			"Provider Client Not Injected",
			"The alicloud_eflo_vcc_route_entries data source ended Configure without a provider client. Please report this issue to the provider developers.",
		)
		return
	}

	regionId := regionOrDefault(client, config.RegionId)
	request := map[string]interface{}{
		"RegionId":   regionId,
		"VccId":      config.VccId.ValueString(),
		"PageSize":   50,
		"PageNumber": 1,
	}
	if v := config.DestinationCidrBlock.ValueString(); v != "" {
		request["DestinationCidrBlock"] = v
	}
	if v := config.Status.ValueString(); v != "" {
		request["Status"] = v
	}
	if v := config.RouteType.ValueString(); v != "" {
		request["RouteType"] = v
	}
	if v := config.NextHopId.ValueString(); v != "" {
		request["NextHopId"] = v
	}
	if v := config.NextHopType.ValueString(); v != "" {
		request["NextHopType"] = v
	}

	const action = "ListVccRouteEntries"
	var ids []string
	pageSize := 50
	pageNumber := 1
	for {
		response, err := client.RpcPost(efloProduct, efloVersion, action, nil, request, true)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", vccRouteEntriesTypeName, action),
				errs.WrapErrorf(err, "Reading %s: calling %s", vccRouteEntriesTypeName, action).Error(),
			)
			return
		}
		content, _ := response["Content"].(map[string]interface{})
		if content == nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", vccRouteEntriesTypeName, action),
				fmt.Sprintf("The response carried no Content object: %v", response),
			)
			return
		}
		data, _ := content["Data"].([]interface{})
		for _, raw := range data {
			item, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			vccId, _ := item["VccId"].(string)
			vccRouteEntryId, _ := item["VccRouteEntryId"].(string)
			if vccId == "" || vccRouteEntryId == "" {
				continue
			}
			ids = append(ids, vccId+idSeparator+vccRouteEntryId)
		}
		if len(data) < pageSize {
			break
		}
		pageNumber++
		request["PageNumber"] = pageNumber
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := vccRouteEntriesModel{
		Id:                   types.StringValue(dataSourceId(config.VccId.ValueString(), ids)),
		VccId:                config.VccId,
		DestinationCidrBlock: config.DestinationCidrBlock,
		Status:               config.Status,
		RouteType:            config.RouteType,
		NextHopId:            config.NextHopId,
		NextHopType:          config.NextHopType,
		RegionId:             types.StringValue(regionId),
		Ids:                  idsList,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// dataSourceId returns a stable identifier for the data source result. It hashes the
// VCC id and the matched entry ids so the value changes only when the result set does.
func dataSourceId(vccId string, ids []string) string {
	h := sha256.Sum256([]byte(vccId + "|" + strings.Join(ids, ",")))
	return fmt.Sprintf("%x", h[:16])
}
