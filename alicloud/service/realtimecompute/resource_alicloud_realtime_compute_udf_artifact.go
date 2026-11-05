// Package realtimecompute. Handwritten Framework resource for the ververica
// 2022-07-18 UdfArtifact API. The four operations are synchronous and ROA-style:
//   - CreateUdfArtifact  POST   /api/v2/namespaces/{namespace}/udfartifacts
//   - GetUdfArtifacts    GET    /api/v2/namespaces/{namespace}/udfartifacts  (udfArtifactName query)
//   - UpdateUdfArtifact  PUT    /api/v2/namespaces/{namespace}/udfartifacts/{udfArtifactName}
//   - DeleteUdfArtifact  DELETE /api/v2/namespaces/{namespace}/udfartifacts/{udfArtifactName}
//
// GetUdfArtifacts answers $.data as a list; Service.DescribeRealtimeComputeUdfArtifact
// filters it by name and returns errs.NotFoundError (990301 or no match). The read path
// of Create/Read/Update and the ACC test harness (ResourceCheckInitWithDescribeMethod)
// share that single Describe method, matching the SDKv2 realtimecompute resources. The
// migrated SDKv2 types in this package stay SDKv2; this is the first Framework type and
// embeds fwadapt.ResourceBase so it can coexist in the same product package.
package realtimecompute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const udfArtifactTypeName = "alicloud_realtime_compute_udf_artifact"

var (
	_ resource.Resource                = &udfArtifactResource{}
	_ resource.ResourceWithImportState = &udfArtifactResource{}
)

// udfClassObjectType is the Framework object type for a single udfClasses entry.
// The server derives udfClasses from the uploaded jar/python file, so every
// nested attribute is Computed.
var udfClassObjectType = types.ObjectType{
	AttributeTypes: map[string]attr.Type{
		"class_name":     types.StringType,
		"class_type":     types.StringType,
		"function_names": types.ListType{ElemType: types.StringType},
	},
}

func fwTimeoutTypesUdfArtifact() map[string]attr.Type {
	return fwutil.TimeoutTypes("create", "update", "delete")
}

func fwTimeoutBlockUdfArtifact() schema.SingleNestedBlock {
	return fwutil.TimeoutBlock("create", "update", "delete")
}

func fwTimeoutDurationUdfArtifact(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	return fwutil.TimeoutDuration(value, operation, fallback)
}

// NewRealtimeComputeUdfArtifactResource is the Framework factory registered in
// register.go as the first Framework (conns.Resources) type of this package.
func NewRealtimeComputeUdfArtifactResource() resource.Resource {
	return &udfArtifactResource{}
}

type udfArtifactResource struct {
	fwadapt.ResourceBase
}

type udfArtifactResourceModel struct {
	Id                types.String `tfsdk:"id"`
	FwTimeouts        types.Object `tfsdk:"timeouts"`
	Workspace         types.String `tfsdk:"workspace"`
	Namespace         types.String `tfsdk:"namespace"`
	Name              types.String `tfsdk:"name"`
	ArtifactType      types.String `tfsdk:"artifact_type"`
	JarUrl            types.String `tfsdk:"jar_url"`
	DependencyJarUris types.List   `tfsdk:"dependency_jar_uris"`
	UdfClasses        types.List   `tfsdk:"udf_classes"`
	Creator           types.String `tfsdk:"creator"`
	CreatedAt         types.String `tfsdk:"created_at"`
	ModifiedAt        types.String `tfsdk:"modified_at"`
}

func (r *udfArtifactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = udfArtifactTypeName
}

func (r *udfArtifactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlockUdfArtifact()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The resource id in format <workspace>:<namespace>:<udfArtifactName>, built from the user-provided identity parts.",
			},
			"workspace": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The workspace ID of the Realtime Compute (Ververica) instance, sent in the request header.",
			},
			"namespace": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The project namespace name, part of the resource path.",
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The UDF artifact name (the jar or python file name). It is the third identity part and the path segment of update/delete.",
			},
			"artifact_type": schema.StringAttribute{
				Required:    true,
				Description: "The artifact type. Valid values: ARTIFACT_TYPE_JAVA, ARTIFACT_TYPE_PYTHON.",
			},
			"jar_url": schema.StringAttribute{
				Required:    true,
				Description: "The OSS storage path of the UDF artifact, e.g. oss://bucket/udfCollection.jar.",
			},
			"dependency_jar_uris": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The list of additional dependency jar storage paths.",
			},
			"udf_classes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"class_name": schema.StringAttribute{Computed: true},
						"class_type": schema.StringAttribute{Computed: true},
						"function_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
				Description: "The UDF classes the server derives from the uploaded artifact.",
			},
			"creator":     schema.StringAttribute{Computed: true, Description: "The creator of the artifact."},
			"created_at":  schema.StringAttribute{Computed: true, Description: "The creation time (Unix seconds)."},
			"modified_at": schema.StringAttribute{Computed: true, Description: "The last modification time (Unix seconds)."},
		},
	}
}

func (r *udfArtifactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan udfArtifactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta()
	operationTimeout, timeoutDiagnostics := fwTimeoutDurationUdfArtifact(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace := plan.Workspace.ValueString()
	namespace := plan.Namespace.ValueString()
	name := plan.Name.ValueString()
	body := map[string]interface{}{
		"name":              name,
		"namespace":         namespace,
		"jarUrl":            plan.JarUrl.ValueString(),
		"artifactType":      plan.ArtifactType.ValueString(),
		"dependencyJarUris": dependencyJarUrisSlice(plan.DependencyJarUris),
	}
	action := fmt.Sprintf("/api/v2/namespaces/%s/udfartifacts", namespace)
	query := map[string]*string{}
	header := map[string]*string{"workspace": helper.StringPointer(workspace)}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	callErr := retry.RetryContext(ctx, operationTimeout, func() *retry.RetryError {
		var err error
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if callErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUdfArtifact", udfArtifactTypeName),
			errs.WrapErrorf(callErr, errs.DefaultErrorMsg, udfArtifactID(workspace, namespace, name), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id := udfArtifactID(workspace, namespace, name)
	plan.Id = types.StringValue(id)
	// CreateUdfArtifact is synchronous, but the list index may lag a moment; the read path
	// (Service.DescribeRealtimeComputeUdfArtifact) retries the GET internally, so by the time
	// it returns the artifact is either visible or genuinely absent.
	object, err := NewService(client).DescribeRealtimeComputeUdfArtifact(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: reading back after create", udfArtifactTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetUdfArtifacts", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back after create", udfArtifactTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, udfArtifactState(object, &plan))...)
}

func (r *udfArtifactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state udfArtifactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	object, err := NewService(r.Meta()).DescribeRealtimeComputeUdfArtifact(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetUdfArtifacts", udfArtifactTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, udfArtifactState(object, &state))...)
}

func (r *udfArtifactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state udfArtifactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta()
	operationTimeout, timeoutDiagnostics := fwTimeoutDurationUdfArtifact(state.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace := state.Workspace.ValueString()
	namespace := state.Namespace.ValueString()
	name := state.Name.ValueString()
	body := map[string]interface{}{
		"name":              name,
		"namespace":         namespace,
		"jarUrl":            plan.JarUrl.ValueString(),
		"artifactType":      plan.ArtifactType.ValueString(),
		"dependencyJarUris": dependencyJarUrisSlice(plan.DependencyJarUris),
	}
	action := fmt.Sprintf("/api/v2/namespaces/%s/udfartifacts/%s", namespace, name)
	query := map[string]*string{}
	header := map[string]*string{"workspace": helper.StringPointer(workspace)}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	callErr := retry.RetryContext(ctx, operationTimeout, func() *retry.RetryError {
		var err error
		response, err = client.RoaPut("ververica", "2022-07-18", action, query, header, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if callErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateUdfArtifact", udfArtifactTypeName),
			errs.WrapErrorf(callErr, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(client).DescribeRealtimeComputeUdfArtifact(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading back after update", udfArtifactTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, udfArtifactState(object, &plan))...)
}

func (r *udfArtifactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state udfArtifactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta()
	operationTimeout, timeoutDiagnostics := fwTimeoutDurationUdfArtifact(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, namespace, name, err := parseUdfArtifactID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s: parsing resource id", udfArtifactTypeName), err.Error())
		return
	}
	action := fmt.Sprintf("/api/v2/namespaces/%s/udfartifacts/%s", namespace, name)
	query := map[string]*string{}
	header := map[string]*string{"workspace": helper.StringPointer(workspace)}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	callErr := retry.RetryContext(ctx, operationTimeout, func() *retry.RetryError {
		var err error
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)
	if callErr != nil {
		// 990301 means the artifact (or its namespace) is already gone.
		if errs.IsExpectedErrors(callErr, []string{"990301"}) || errs.NotFoundError(callErr) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteUdfArtifact", udfArtifactTypeName),
			errs.WrapErrorf(callErr, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	// Some ververica responses surface NotFound as success=false in the body rather
	// than a Go error; treat that as already deleted too.
	if fmt.Sprint(response["success"]) == "false" && fmt.Sprint(response["errorCode"]) == "990301" {
		return
	}
}

func (r *udfArtifactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// udfArtifactID builds the composite id <workspace>:<namespace>:<udfArtifactName>.
func udfArtifactID(workspace, namespace, name string) string {
	return fmt.Sprintf("%s:%s:%s", workspace, namespace, name)
}

// parseUdfArtifactID splits the composite id into its three identity parts.
func parseUdfArtifactID(id string) (workspace, namespace, name string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid resource id %q: expected <workspace>:<namespace>:<udfArtifactName>", id)
	}
	return parts[0], parts[1], parts[2], nil
}

// udfArtifactState maps a ververica UdfArtifact response object onto the resource
// model. workspace is a ForceNew identity part sent in the request header rather than
// echoed in the list response, and the composite id is not part of the response, so both
// are carried over from the prior plan/state; the remaining fields come from the server.
func udfArtifactState(object map[string]interface{}, prior *udfArtifactResourceModel) udfArtifactResourceModel {
	return udfArtifactResourceModel{
		Id:                prior.Id,
		FwTimeouts:        prior.FwTimeouts,
		Workspace:         prior.Workspace,
		Namespace:         strField(object, "namespace"),
		Name:              strField(object, "name"),
		ArtifactType:      strField(object, "artifactType"),
		JarUrl:            strField(object, "jarUrl"),
		DependencyJarUris: strListField(object, "dependencyJarUris"),
		UdfClasses:        udfClassesField(object, "udfClasses"),
		Creator:           strField(object, "creator"),
		CreatedAt:         strField(object, "createdAt"),
		ModifiedAt:        strField(object, "modifiedAt"),
	}
}

// strField reads a string field; a present-but-empty string stays "" instead of null.
func strField(object map[string]interface{}, key string) types.String {
	v, ok := object[key]
	if !ok || v == nil {
		return types.StringNull()
	}
	return types.StringValue(fmt.Sprint(v))
}

// strListField reads a string list field into a types.List of StringType.
func strListField(object map[string]interface{}, key string) types.List {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.ListNull(types.StringType)
	}
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return types.ListNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(arr))
	for _, v := range arr {
		elems = append(elems, types.StringValue(fmt.Sprint(v)))
	}
	list, _ := types.ListValue(types.StringType, elems)
	return list
}

// udfClassesField reads the server-derived udfClasses list into a types.List of
// the udfClass object type.
func udfClassesField(object map[string]interface{}, key string) types.List {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.ListNull(udfClassObjectType)
	}
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return types.ListNull(udfClassObjectType)
	}
	elems := make([]attr.Value, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		obj, _ := types.ObjectValue(udfClassObjectType.AttributeTypes, map[string]attr.Value{
			"class_name":     strField(m, "className"),
			"class_type":     strField(m, "classType"),
			"function_names": strListField(m, "functionNames"),
		})
		elems = append(elems, obj)
	}
	if len(elems) == 0 {
		return types.ListNull(udfClassObjectType)
	}
	list, _ := types.ListValue(udfClassObjectType, elems)
	return list
}

// dependencyJarUrisSlice converts a plan dependency_jar_uris list into the body
// slice; a null/empty list yields an empty slice so PUT carries the full state.
func dependencyJarUrisSlice(l types.List) []interface{} {
	out := []interface{}{}
	if l.IsNull() || l.IsUnknown() {
		return out
	}
	for _, e := range l.Elements() {
		if s, ok := e.(types.String); ok {
			out = append(out, s.ValueString())
		}
	}
	return out
}
