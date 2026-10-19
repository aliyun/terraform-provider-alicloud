package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const contactsDataSourceTypeName = "alicloud_cms_contacts"

var (
	_ datasource.DataSource              = &contactsDataSource{}
	_ datasource.DataSourceWithConfigure = &contactsDataSource{}
)

// NewContactsDataSource is the Framework factory for the alicloud_cms_contacts
// data source, registered by ServicePackage.
func NewContactsDataSource() datasource.DataSource {
	return &contactsDataSource{}
}

// contactsDataSource implements alicloud_cms_contacts on the
// terraform-plugin-framework, embedding fwadapt.DataSourceBase so the muxed
// provider injects the *connectivity.AliyunClient during Configure.
type contactsDataSource struct {
	fwadapt.DataSourceBase
}

// contactsDataSourceModel is the tfsdk state shape of alicloud_cms_contacts.
type contactsDataSourceModel struct {
	ContactName            types.String         `tfsdk:"contact_name"`
	Phone                  types.String         `tfsdk:"phone"`
	Email                  types.String         `tfsdk:"email"`
	QueryUngroupedContacts types.Bool           `tfsdk:"query_ungrouped_contacts"`
	NameRegex              types.String         `tfsdk:"name_regex"`
	OutputFile             types.String         `tfsdk:"output_file"`
	Ids                    []string             `tfsdk:"ids"`
	Contacts               []contactObjectModel `tfsdk:"contacts"`
	Id                     types.String         `tfsdk:"id"`
}

// contactObjectModel is one element of the contacts list exported by the data
// source.
type contactObjectModel struct {
	ContactName types.String `tfsdk:"contact_name"`
	Email       types.String `tfsdk:"email"`
	Phone       types.String `tfsdk:"phone"`
	Lang        types.String `tfsdk:"lang"`
	Workspace   types.String `tfsdk:"workspace"`
	ImUserIds   types.Map    `tfsdk:"im_user_ids"`
	ContactId   types.String `tfsdk:"contact_id"`
}

func (d *contactsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = contactsDataSourceTypeName
}

func (d *contactsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of Cloud Monitor Service (CMS) Contacts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID, a hash of the region and the returned contact IDs.",
			},
			"contact_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the alarm contact used to filter results.",
			},
			"phone": schema.StringAttribute{
				Optional:    true,
				Description: "The phone number of the alarm contact used to filter results.",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "The email address of the alarm contact used to filter results.",
			},
			"query_ungrouped_contacts": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to query ungrouped contacts.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by contact name.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results after running `terraform plan`.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of contact IDs.",
			},
			// contacts is a list of contact objects. It is declared as a
			// ListAttribute with an Object ElementType rather than a
			// ListNestedAttribute: under the provider's protocol v5 mux, a nested
			// attribute carrying Attributes (NestedAttributeObject) fails schema
			// conversion with "protocol version 5 cannot have Attributes set",
			// breaking the whole provider schema load. An Object ElementType uses
			// Type (not Attributes) and is protocol-v5 compatible; elements still
			// decode into contactObjectModel via tfsdk tags.
			"contacts": schema.ListAttribute{
				Computed:    true,
				Description: "A list of CMS Contacts. Each element contains the following attributes.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"contact_id":   types.StringType,
						"contact_name": types.StringType,
						"email":        types.StringType,
						"phone":        types.StringType,
						"lang":         types.StringType,
						"workspace":    types.StringType,
						"im_user_ids":  types.MapType{ElemType: types.StringType},
					},
				},
			},
		},
	}
}

func (d *contactsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config contactsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(d.Client())
	objects, err := svc.ListContacts(ctx,
		config.ContactName.ValueString(),
		config.Phone.ValueString(),
		config.Email.ValueString(),
		config.QueryUngroupedContacts.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", contactsDataSourceTypeName), err.Error())
		return
	}

	var nameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); v != "" {
		if r, err := regexp.Compile(v); err == nil {
			nameRegex = r
		}
	}

	contacts := make([]contactObjectModel, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, obj := range objects {
		contactName := fmt.Sprint(obj["name"])
		if nameRegex != nil && !nameRegex.MatchString(contactName) {
			continue
		}
		contacts = append(contacts, contactObjectState(ctx, obj))
		ids = append(ids, fmt.Sprint(obj["contactId"]))
	}

	state := config
	state.Ids = ids
	state.Contacts = contacts
	state.Id = types.StringValue(dataSourceContactsID(d.Client().RegionId, ids, contacts))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if output := config.OutputFile.ValueString(); output != "" {
		if err := writeContactsToFile(output, contacts); err != nil {
			resp.Diagnostics.AddWarning(fmt.Sprintf("Reading %s", contactsDataSourceTypeName), err.Error())
		}
	}
}

// dataSourceContactsID mirrors the legacy SDKv2 data source ID: a stable hash
// of the region and the returned contact IDs, falling back to a region-only
// hash when no contacts matched.
func dataSourceContactsID(region string, ids []string, contacts []contactObjectModel) string {
	if len(ids) == 0 {
		return cmsHash(region)
	}
	return cmsHash(append([]string{region}, ids...)...)
}

// cmsHash returns a stable FNV-1a 64-bit hash of the joined parts, formatted
// as a base-36 string. It is the local equivalent of the root package
// dataResourceIdHash so the service package does not import the root package.
func cmsHash(parts ...string) string {
	const (
		offset64 uint64 = 1469598103934665603
		prime64  uint64 = 1099511628211
	)
	h := offset64
	for _, part := range parts {
		for i := 0; i < len(part); i++ {
			h ^= uint64(part[i])
			h *= prime64
		}
	}
	return strconv.FormatUint(h, 36)
}

// contactObjectState maps a CMS contact object onto the data source element
// model. Absent string keys become null; im_user_ids absent becomes a null map.
func contactObjectState(ctx context.Context, obj map[string]interface{}) contactObjectModel {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	imUserIds := types.MapNull(types.StringType)
	if raw, ok := obj["imUserIds"].(map[string]interface{}); ok && len(raw) > 0 {
		m := make(map[string]string, len(raw))
		for k, v := range raw {
			m[k] = fmt.Sprint(v)
		}
		if mv, diags := types.MapValueFrom(ctx, types.StringType, m); !diags.HasError() {
			imUserIds = mv
		}
	}
	return contactObjectModel{
		ContactId:   str("contactId"),
		ContactName: str("name"),
		Email:       str("email"),
		Phone:       str("phone"),
		Lang:        str("lang"),
		Workspace:   str("workspace"),
		ImUserIds:   imUserIds,
	}
}

// writeContactsToFile writes the contacts list as JSON to output, matching the
// legacy SDKv2 writeToFile behavior for data sources.
func writeContactsToFile(output string, contacts []contactObjectModel) error {
	b, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(output, b, 0o644)
}
