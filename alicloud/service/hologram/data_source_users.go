package hologram

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const usersTypeName = "alicloud_hologram_users"

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// usersObjectAttrTypes is the element type of the users list. It is a plain
// ObjectType (not a NestedAttributeObject) so the schema stays convertible to
// protocol v5: the provider is muxed with the SDK v2 provider behind a
// tf5muxserver, and protocol v5 cannot express ListNestedAttribute's nested
// Attributes. A ListAttribute of an anonymous Object element type carries the
// object shape inside the cty element type instead, which v5 supports.
var usersObjectAttrTypes = map[string]attr.Type{
	"id":         types.StringType,
	"user_name":  types.StringType,
	"super_user": types.BoolType,
}

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource lists the users of a Hologres instance via ListUsers and
// filters them client-side by name regex and/or id list. ListUsers is the only
// available read API; it returns each user's name as "useName" and nothing
// else, so super_user is exported as null.
type usersDataSource struct {
	fwadapt.DataSourceBase
}

type usersDataSourceModel struct {
	Id         types.String   `tfsdk:"id"`
	InstanceId types.String   `tfsdk:"instance_id"`
	NameRegex  types.String   `tfsdk:"name_regex"`
	Ids        []types.String `tfsdk:"ids"`
	Names      []types.String `tfsdk:"names"`
	Users      []types.Object `tfsdk:"users"`
	OutputFile types.String   `tfsdk:"output_file"`
}

func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = usersTypeName
}

func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of Hologres (Hologram) Users owned by an Alibaba Cloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID, a hash of the listed user ids.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Hologres instance.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by user name.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of User IDs. The value is formulated as `<instance_id>:<user_name>`.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of user names.",
			},
			"users": schema.ListAttribute{
				ElementType: types.ObjectType{AttrTypes: usersObjectAttrTypes},
				Computed:    true,
				Description: "A list of Hologram Users. Each element is an object with `id`, `user_name` and `super_user`.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results (after running `terraform plan`).",
			},
		},
	}
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()

	var nameRegex *regexp.Regexp
	if !state.NameRegex.IsNull() && state.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(state.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", usersTypeName),
				errs.WrapError(err).Error(),
			)
			return
		}
		nameRegex = r
	}

	idsFilter := map[string]bool{}
	for _, id := range state.Ids {
		idsFilter[id.ValueString()] = true
	}

	users, err := listUsers(d.Client(), instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListUsers", usersTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListUsers", usersTypeName).Error(),
		)
		return
	}

	var matched []types.Object
	var ids []types.String
	var names []types.String
	var fileRows []map[string]interface{}
	for _, raw := range users {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		userName := fmt.Sprint(item["useName"])
		if userName == "" || userName == "<nil>" {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(userName) {
			continue
		}
		id := fmt.Sprintf("%s:%s", instanceId, userName)
		if len(idsFilter) > 0 && !idsFilter[id] {
			continue
		}
		obj, d := types.ObjectValue(usersObjectAttrTypes, map[string]attr.Value{
			"id":         types.StringValue(id),
			"user_name":  types.StringValue(userName),
			"super_user": types.BoolNull(),
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		matched = append(matched, obj)
		ids = append(ids, types.StringValue(id))
		names = append(names, types.StringValue(userName))
		fileRows = append(fileRows, map[string]interface{}{
			"id":         id,
			"user_name":  userName,
			"super_user": nil,
		})
	}

	state.Id = types.StringValue(dataSourceIdHash(ids))
	state.Ids = ids
	state.Names = names
	state.Users = matched

	if !state.OutputFile.IsNull() && state.OutputFile.ValueString() != "" {
		if err := writeUsersToFile(state.OutputFile.ValueString(), fileRows); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", usersTypeName),
				errs.WrapError(err).Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func dataSourceIdHash(ids []types.String) string {
	h := sha1.New()
	for _, id := range ids {
		h.Write([]byte(id.ValueString()))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeUsersToFile(path string, rows []map[string]interface{}) error {
	raw, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}
