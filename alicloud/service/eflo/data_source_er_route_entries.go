package eflo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const erRouteEntriesTypeName = "alicloud_eflo_er_route_entries"

var (
	_ datasource.DataSource              = &erRouteEntriesDataSource{}
	_ datasource.DataSourceWithConfigure = &erRouteEntriesDataSource{}
)

// NewErRouteEntriesDataSource is the factory registered against the Eflo
// ServicePackage so the framework provider server can construct the data source
// on first use.
func NewErRouteEntriesDataSource() datasource.DataSource {
	return &erRouteEntriesDataSource{}
}

type erRouteEntriesDataSource struct {
	fwadapt.DataSourceBase
}

// erRouteEntryModel is the per-entry shape of the `entries` list. The `id`
// field is the Terraform-facing composite "<er_id>:<er_route_entry_id>" used
// by the `ids` filter, not the raw API ErRouteEntryId.
type erRouteEntryModel struct {
	Id                   types.String `tfsdk:"id"`
	ErRouteEntryId       types.String `tfsdk:"er_route_entry_id"`
	ErId                 types.String `tfsdk:"er_id"`
	RouteType            types.String `tfsdk:"route_type"`
	NextHopId            types.String `tfsdk:"next_hop_id"`
	NextHopType          types.String `tfsdk:"next_hop_type"`
	Status               types.String `tfsdk:"status"`
	DestinationCidrBlock types.String `tfsdk:"destination_cidr_block"`
	GmtModified          types.String `tfsdk:"gmt_modified"`
	RegionId             types.String `tfsdk:"region_id"`
	TenantId             types.String `tfsdk:"tenant_id"`
	ResourceGroupId      types.String `tfsdk:"resource_group_id"`
}

type erRouteEntriesModel struct {
	Id                   types.String        `tfsdk:"id"`
	ErId                 types.String        `tfsdk:"er_id"`
	RouteType            types.String        `tfsdk:"route_type"`
	NextHopId            types.String        `tfsdk:"next_hop_id"`
	NextHopType          types.String        `tfsdk:"next_hop_type"`
	Status               types.String        `tfsdk:"status"`
	DestinationCidrBlock types.String        `tfsdk:"destination_cidr_block"`
	Ids                  types.List          `tfsdk:"ids"`
	PageNumber           types.Int64         `tfsdk:"page_number"`
	PageSize             types.Int64         `tfsdk:"page_size"`
	OutputFile           types.String        `tfsdk:"output_file"`
	Entries              []erRouteEntryModel `tfsdk:"entries"`
}

func (d *erRouteEntriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = erRouteEntriesTypeName
}

func (d *erRouteEntriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Eflo Er Route Entries available to the user. `er_id` is required; " +
			"optional filters narrow the list. Each entry id is `<er_id>:<er_route_entry_id>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, computed as the hash of the Er Route Entry ids.",
			},
			"er_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Eflo Er (灵骏 HUB). The route entries of this Er are listed.",
			},
			"route_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the type of the route, e.g. `VCC`.",
			},
			"next_hop_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the ID of the next hop instance.",
			},
			"next_hop_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the type of the next hop instance, e.g. `VCC`.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the status of the route entry. Valid values: `Synchronizing`, `Available`.",
			},
			"destination_cidr_block": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the destination network segment.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of Er Route Entry ids. When set, only entries whose `<er_id>:<er_route_entry_id>` is in this list are returned.",
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The page number. When set, only that single page is returned (single-page mode); otherwise all pages are walked.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The page size. Defaults to 50 when unset.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Writes the entries list as JSON to this path. No file is written when unset.",
			},
		},
		// ListNestedBlock (not ListNestedAttribute): the SDKv2+Framework
		// tf5mux server is protocol v5, which cannot express nested
		// *Attributes* (only nested *Blocks*, which map to SDKv2 TypeList +
		// Elem: Resource). The provider's framework/schema convertBlock helper
		// confirms the sanctioned v5 form. A Block implements schema.Block
		// (not schema.Attribute), so it lives in schema.Schema.Blocks.
		Blocks: map[string]schema.Block{
			"entries": schema.ListNestedBlock{
				Description: "The list of Er Route Entries.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The composite id `<er_id>:<er_route_entry_id>`.",
						},
						"er_route_entry_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the Er Route Entry.",
						},
						"er_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the Er.",
						},
						"route_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the route.",
						},
						"next_hop_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the next hop instance.",
						},
						"next_hop_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the next hop instance.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the route entry.",
						},
						"destination_cidr_block": schema.StringAttribute{
							Computed:    true,
							Description: "The destination network segment.",
						},
						"gmt_modified": schema.StringAttribute{
							Computed:    true,
							Description: "The time the route entry was last modified.",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "The region ID.",
						},
						"tenant_id": schema.StringAttribute{
							Computed:    true,
							Description: "The tenant ID.",
						},
						"resource_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the resource group.",
						},
					},
				},
			},
		},
	}
}

func (d *erRouteEntriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "ListErRouteEntries"

	var config erRouteEntriesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", erRouteEntriesTypeName, action),
			"The data source ended Configure without a provider client. Please report this issue to the provider developers.",
		)
		return
	}

	erId := config.ErId.ValueString()
	if erId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", erRouteEntriesTypeName, action),
			"The required argument `er_id` is empty.",
		)
		return
	}

	// EnablePage defaults to enabled on the server side, but the vpds data
	// source (the closest eflo analog) sets it explicitly, so we do too.
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ErId":       erId,
		"EnablePage": true,
		"PageNumber": 1,
		"PageSize":   defaultErRouteEntriesPageSize,
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
	if v := config.Status.ValueString(); v != "" {
		request["Status"] = v
	}
	if v := config.DestinationCidrBlock.ValueString(); v != "" {
		request["DestinationCidrBlock"] = v
	}

	// page_number set (>0) selects single-page mode, mirroring
	// isPagingRequest/setPagingRequest in the SDKv2 eflo data sources.
	singlePage := false
	if !config.PageNumber.IsNull() {
		if pn := config.PageNumber.ValueInt64(); pn > 0 {
			request["PageNumber"] = int(pn)
			singlePage = true
		}
	}
	if !config.PageSize.IsNull() {
		if ps := config.PageSize.ValueInt64(); ps > 0 {
			request["PageSize"] = int(ps)
		}
	}

	// ids filter: composite "<er_id>:<er_route_entry_id>" values the user
	// passes to narrow the result. Null/empty means "no filter".
	idsFilter := map[string]bool{}
	if !config.Ids.IsNull() {
		var ids []string
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			idsFilter[id] = true
		}
	}

	svc := NewService(client)
	objects, _, err := svc.ListErRouteEntries(request, singlePage)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", erRouteEntriesTypeName, action),
			err.Error(),
		)
		return
	}

	entries := make([]erRouteEntryModel, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, obj := range objects {
		entryId := strField(obj, "ErRouteEntryId")
		id := fmt.Sprintf("%s:%s", erId, entryId)
		if len(idsFilter) > 0 && !idsFilter[id] {
			continue
		}
		entries = append(entries, erRouteEntryModel{
			Id:                   types.StringValue(id),
			ErRouteEntryId:       toStr(entryId),
			ErId:                 toStr(strField(obj, "ErId")),
			RouteType:            toStr(strField(obj, "RouteType")),
			NextHopId:            toStr(strField(obj, "NextHopId")),
			NextHopType:          toStr(strField(obj, "NextHopType")),
			Status:               toStr(strField(obj, "Status")),
			DestinationCidrBlock: toStr(strField(obj, "DestinationCidrBlock")),
			GmtModified:          toStr(strField(obj, "GmtModified")),
			RegionId:             toStr(strField(obj, "RegionId")),
			TenantId:             toStr(strField(obj, "TenantId")),
			ResourceGroupId:      toStr(strField(obj, "ResourceGroupId")),
		})
		ids = append(ids, id)
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := erRouteEntriesModel{
		Id:                   types.StringValue(erRouteEntriesIdHash(ids)),
		ErId:                 config.ErId,
		RouteType:            config.RouteType,
		NextHopId:            config.NextHopId,
		NextHopType:          config.NextHopType,
		Status:               config.Status,
		DestinationCidrBlock: config.DestinationCidrBlock,
		Ids:                  idsList,
		PageNumber:           config.PageNumber,
		PageSize:             config.PageSize,
		OutputFile:           config.OutputFile,
		Entries:              entries,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if output := config.OutputFile.ValueString(); output != "" {
		if err := erRouteEntriesWriteToFile(output, entries); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output_file", erRouteEntriesTypeName),
				err.Error(),
			)
		}
	}
}

// strField returns the string value of obj[key], or "" when missing/not a
// string. The Eflo API returns every entry field as a JSON string, so the only
// non-string shapes tolerated here are nil/missing.
func strField(obj map[string]interface{}, key string) string {
	if v, ok := obj[key].(string); ok {
		return v
	}
	return ""
}

// toStr mirrors the ims data source helper: an empty string becomes null
// rather than "", so omitted API fields read back as absent.
func toStr(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// erRouteEntriesIdHash mirrors the root-package dataResourceIdHash: the data
// source id is helper.Hashcode of "<id1>-<id2>-...", stable for a given entry
// set and changing when the set changes. Reproduced locally because the root
// helper lives in package alicloud, which would form an import cycle.
func erRouteEntriesIdHash(ids []string) string {
	var buf bytes.Buffer
	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}
	return fmt.Sprintf("%d", helper.Hashcode(buf.String()))
}

// erRouteEntriesWriteToFile mirrors the root-package writeToFile so the
// output_file argument works without importing package alicloud.
func erRouteEntriesWriteToFile(filePath string, data interface{}) error {
	var out string
	switch v := data.(type) {
	case string:
		out = v
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}
	if strings.HasPrefix(filePath, "~") {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			filePath = filepath.Join(home, strings.TrimPrefix(filePath, "~"))
		}
	}
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, []byte(out), 0o644)
}
