package dlfnext

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const catalogsDataSourceTypeName = "alicloud_dlf_next_catalogs"

var (
	_ datasource.DataSource              = &catalogsDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogsDataSource{}
)

func NewCatalogsDataSource() datasource.DataSource {
	return &catalogsDataSource{}
}

type catalogsDataSource struct {
	fwadapt.DataSourceBase
}

type catalogsDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	CatalogNamePattern types.String `tfsdk:"catalog_name_pattern"`
	NameRegex          types.String `tfsdk:"name_regex"`
	Ids                types.List   `tfsdk:"ids"`
	Names              types.List   `tfsdk:"names"`
	OutputFile         types.String `tfsdk:"output_file"`
	Catalogs           types.List   `tfsdk:"catalogs"`
}

type catalogsItemModel struct {
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Options   types.Map    `tfsdk:"options"`
	IsShared  types.Bool   `tfsdk:"is_shared"`
	ShareId   types.String `tfsdk:"share_id"`
	Id        types.String `tfsdk:"id"`
	Status    types.String `tfsdk:"status"`
	Owner     types.String `tfsdk:"owner"`
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	UpdatedBy types.String `tfsdk:"updated_by"`
}

// catalogItem is the plain (non-terraform) projection used by the filtering,
// sorting and output-file helpers; keeping it framework-free makes the
// discrimination logic unit-testable without a provider server.
type catalogItem struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Options   map[string]string `json:"options"`
	IsShared  bool              `json:"is_shared"`
	ShareId   string            `json:"share_id"`
	Id        string            `json:"id"`
	Status    string            `json:"status"`
	Owner     string            `json:"owner"`
	CreatedAt string            `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt string            `json:"updated_at"`
	UpdatedBy string            `json:"updated_by"`
}

func (d *catalogsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = catalogsDataSourceTypeName
}

func (d *catalogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of DLF Next Catalogs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"catalog_name_pattern": schema.StringAttribute{
				Optional:    true,
				Description: "The service-side catalog name pattern passed to ListCatalogs as `catalogNamePattern`.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter the returned catalog names.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of catalog names to filter; the matched names are also returned here.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The matched catalog names.",
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			// Protocol 5 (the muxed provider) cannot serve nested attributes;
			// a list of objects carries the same shape.
			"catalogs": schema.ListAttribute{
				Computed:    true,
				ElementType: catalogItemObjectType(),
				Description: "The matched catalogs.",
			},
		},
	}
}

func (d *catalogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model catalogsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pattern := ""
	if !model.CatalogNamePattern.IsNull() && !model.CatalogNamePattern.IsUnknown() {
		pattern = model.CatalogNamePattern.ValueString()
	}
	service := NewService(d.Meta())
	deadline := time.Now().Add(catalogTimeout)
	objects, err := listAllCatalogs(ctx, func(pattern string, token string, max int) ([]interface{}, string, error) {
		return service.ListCatalogsWithTimeout(ctx, pattern, token, max, remainingUntil(deadline))
	}, pattern)
	if err != nil {
		resp.Diagnostics.AddError("Listing "+catalogsDataSourceTypeName, errs.WrapError(err).Error())
		return
	}

	idSet := map[string]bool{}
	if !model.Ids.IsNull() && !model.Ids.IsUnknown() {
		var ids []string
		resp.Diagnostics.Append(model.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			idSet[id] = true
		}
	}

	var nameRegex *regexp.Regexp
	if !model.NameRegex.IsNull() && !model.NameRegex.IsUnknown() && model.NameRegex.ValueString() != "" {
		nameRegex, err = regexp.Compile(model.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid name_regex", err.Error())
			return
		}
	}

	matchedNames, items := filterCatalogs(objects, idSet, nameRegex)
	sort.Strings(matchedNames)
	items = sortCatalogItems(items)

	idsList, diags := types.ListValueFrom(ctx, types.StringType, matchedNames)
	resp.Diagnostics.Append(diags...)
	namesList, diags := types.ListValueFrom(ctx, types.StringType, matchedNames)
	resp.Diagnostics.Append(diags...)

	itemModels := make([]catalogsItemModel, len(items))
	for i, item := range items {
		optionsMap, mapDiags := types.MapValueFrom(ctx, types.StringType, item.Options)
		resp.Diagnostics.Append(mapDiags...)
		itemModels[i] = catalogsItemModel{
			Name:      types.StringValue(item.Name),
			Type:      types.StringValue(item.Type),
			Options:   optionsMap,
			IsShared:  types.BoolValue(item.IsShared),
			ShareId:   types.StringValue(item.ShareId),
			Id:        types.StringValue(item.Id),
			Status:    types.StringValue(item.Status),
			Owner:     types.StringValue(item.Owner),
			CreatedAt: types.StringValue(item.CreatedAt),
			CreatedBy: types.StringValue(item.CreatedBy),
			UpdatedAt: types.StringValue(item.UpdatedAt),
			UpdatedBy: types.StringValue(item.UpdatedBy),
		}
	}
	catalogsList, diags := types.ListValueFrom(ctx, catalogItemObjectType(), itemModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(catalogsID(matchedNames))
	model.Ids = idsList
	model.Names = namesList
	model.Catalogs = catalogsList

	if !model.OutputFile.IsNull() && !model.OutputFile.IsUnknown() && model.OutputFile.ValueString() != "" {
		if writeErr := writeCatalogOutput(model.OutputFile.ValueString(), items); writeErr != nil {
			resp.Diagnostics.AddError("Writing output_file", writeErr.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// listAllCatalogs pages ListCatalogs until the token is empty. The lister is
// injectable so pagination can be unit-tested without a backend.
func listAllCatalogs(ctx context.Context, lister func(pattern, token string, max int) ([]interface{}, string, error), pattern string) ([]interface{}, error) {
	var all []interface{}
	token := ""
	for page := 0; ; page++ {
		if page >= 1000 {
			return nil, fmt.Errorf("ListCatalogs exceeded 1000 pages; refusing to keep paging")
		}
		objects, next, err := lister(pattern, token, 100)
		if err != nil {
			return nil, err
		}
		all = append(all, objects...)
		if next == "" {
			return all, nil
		}
		if next == token {
			return nil, fmt.Errorf("ListCatalogs returned the same page token twice; aborting to avoid a loop")
		}
		token = next
	}
}

// filterCatalogs applies the ids and name_regex filters and projects the raw
// objects. A nil/empty ids set means no id filter.
func filterCatalogs(objects []interface{}, idSet map[string]bool, nameRegex *regexp.Regexp) ([]string, []catalogItem) {
	names := make([]string, 0, len(objects))
	items := make([]catalogItem, 0, len(objects))
	for _, raw := range objects {
		object, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		name := fmt.Sprint(object["name"])
		if len(idSet) > 0 && !idSet[name] {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		names = append(names, name)
		options := map[string]string{}
		if rawOptions, ok := object["options"].(map[string]interface{}); ok {
			for k, v := range rawOptions {
				options[k] = fmt.Sprint(v)
			}
		}
		isShared, _ := object["isShared"].(bool)
		items = append(items, catalogItem{
			Name:      name,
			Type:      stringValue(object["type"]),
			Options:   options,
			IsShared:  isShared,
			ShareId:   stringValue(object["shareId"]),
			Id:        name,
			Status:    stringValue(object["status"]),
			Owner:     stringValue(object["owner"]),
			CreatedAt: catalogTimestampString(object["createdAt"]),
			CreatedBy: stringValue(object["createdBy"]),
			UpdatedAt: catalogTimestampString(object["updatedAt"]),
			UpdatedBy: stringValue(object["updatedBy"]),
		})
	}
	return names, items
}

// sortCatalogItems returns a copy sorted by name so the list is stable.
func sortCatalogItems(items []catalogItem) []catalogItem {
	sorted := make([]catalogItem, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	return sorted
}

// catalogsID returns a stable id from the sorted matched names.
func catalogsID(names []string) string {
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Strings(sorted)
	return helper.DataResourceIdHash(sorted)
}

// writeCatalogOutput writes the matched catalogs to output_file.
func writeCatalogOutput(path string, items []catalogItem) error {
	return helper.WriteToFile(path, items)
}

func catalogItemObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":       types.StringType,
		"type":       types.StringType,
		"options":    types.MapType{ElemType: types.StringType},
		"is_shared":  types.BoolType,
		"share_id":   types.StringType,
		"id":         types.StringType,
		"status":     types.StringType,
		"owner":      types.StringType,
		"created_at": types.StringType,
		"created_by": types.StringType,
		"updated_at": types.StringType,
		"updated_by": types.StringType,
	}}
}
