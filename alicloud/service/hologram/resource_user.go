// Package hologram holds the Framework resources and data sources for Hologres
// (Hologram). The Hologram instance itself still lives in package alicloud as a
// SDK v2 resource; new Hologram resources are Framework-native and registered
// here.
package hologram

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const userTypeName = "alicloud_hologram_user"

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource manages a Hologres instance user. Hologram 2022-06-01 exposes no
// UpdateUser, so every attribute is ForceNew; there is no Update path. Read is
// ListUsers (the only available read API) filtered client-side by user name.
type userResource struct {
	fwadapt.ResourceBase
}

type userResourceModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	UserName   types.String `tfsdk:"user_name"`
	SuperUser  types.Bool   `tfsdk:"super_user"`
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = userTypeName
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Hologres (Hologram) User resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID of Hologram User. The value is formulated as `<instance_id>:<user_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Hologres instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the instance user. It must start with a letter and can contain lowercase letters, digits, and underscores (_).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"super_user": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether the user is a super user. Hologram provides no UpdateUser API, so this attribute can only be set at creation.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := plan.InstanceId.ValueString()
	userName := plan.UserName.ValueString()
	body := map[string]interface{}{
		"userName": userName,
	}
	if !plan.SuperUser.IsNull() && !plan.SuperUser.IsUnknown() {
		body["superUser"] = plan.SuperUser.ValueBool()
	}

	action := fmt.Sprintf("/api/v1/instances/%s/createUser", instanceId)
	if _, err := r.Client().RoaPost("Hologram", "2022-06-01", action, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUser", userTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateUser", userTypeName).Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", instanceId, userName))
	// ListUsers does not return SuperUser, so the create-time value is the only
	// source of truth; keep it in state rather than overwriting it from a read.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, userName, ok := parseUserId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", userTypeName),
			fmt.Sprintf("Invalid resource id %q: expected <instance_id>:<user_name>", state.Id.ValueString()),
		)
		return
	}

	_, found, err := findUser(r.Client(), instanceId, userName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListUsers", userTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListUsers", userTypeName).Error(),
		)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.InstanceId = types.StringValue(instanceId)
	state.UserName = types.StringValue(userName)
	// ListUsers returns only the user name, never SuperUser, so the value
	// already in state (from create or a prior read) is preserved here.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute carries RequiresReplace, so any change
// routes through Delete+Create. The method exists only to satisfy the resource
// interface and keeps state in sync with the plan.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, userName, ok := parseUserId(state.Id.ValueString())
	if !ok {
		return
	}
	action := fmt.Sprintf("/api/v1/instances/%s/dropUser", instanceId)
	body := map[string]interface{}{"userName": userName}
	_, err := r.Client().RoaPost("Hologram", "2022-06-01", action, nil, nil, body, true)
	if err != nil && !isUserNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DropUser", userTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DropUser", userTypeName).Error(),
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseUserId splits <instance_id>:<user_name>. The user name may not contain a
// colon, so a strict two-part split is enough.
func parseUserId(id string) (instanceId, userName string, ok bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// findUser lists the users of an instance via ListUsers (the only available read
// API; it is private-visibility but gateway-callable) and returns the item whose
// name matches. Hologram returns the per-user name as "useName" in the response,
// not "userName" — see ACC taskId=11808 — so the match reads "useName".
func findUser(client *connectivity.AliyunClient, instanceId, userName string) (map[string]interface{}, bool, error) {
	users, err := listUsers(client, instanceId)
	if err != nil {
		return nil, false, err
	}
	for _, u := range users {
		item, _ := u.(map[string]interface{})
		if fmt.Sprint(item["useName"]) == userName {
			return item, true, nil
		}
	}
	return nil, false, nil
}

// listUsers paginates ListUsers using nextToken. Each page returns userList, an
// array of objects carrying only useName.
func listUsers(client *connectivity.AliyunClient, instanceId string) ([]interface{}, error) {
	action := fmt.Sprintf("/api/v1/instances/%s/users", instanceId)
	var all []interface{}
	query := map[string]*string{"maxResults": strPointer("100")}
	for {
		response, err := client.RoaGet("Hologram", "2022-06-01", action, query, nil, nil)
		if err != nil {
			return nil, err
		}
		if raw, ok := response["userList"].([]interface{}); ok {
			all = append(all, raw...)
		}
		token, _ := response["nextToken"].(string)
		if token == "" {
			break
		}
		query["nextToken"] = strPointer(token)
	}
	return all, nil
}

func strPointer(v string) *string { return &v }

// isUserNotExist reports whether an error means the user is already gone, so a
// delete of a vanished resource stays idempotent.
func isUserNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is not found")
}
