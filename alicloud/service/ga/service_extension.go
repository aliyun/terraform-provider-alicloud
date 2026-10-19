package ga

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const serviceExtensionTypeName = "alicloud_ga_service_extension"

var (
	_ resource.Resource                = (*serviceExtensionResource)(nil)
	_ resource.ResourceWithImportState = (*serviceExtensionResource)(nil)
)

func NewServiceExtensionResource() resource.Resource {
	return &serviceExtensionResource{}
}

type serviceExtensionResource struct {
	fwadapt.ResourceBase
}

// serviceExtensionResourceModel is the tfsdk model for the resource. Computed
// attributes that do not move under provider control on an update carry
// UseStateForUnknown so a plan does not demand their recomputation.
type serviceExtensionResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	ServiceExtensionId types.String `tfsdk:"service_extension_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
	ResourceGroupId    types.String `tfsdk:"resource_group_id"`
	Tags               types.Map    `tfsdk:"tags"`
	Components         types.Set    `tfsdk:"components"`
	Resources          types.Set    `tfsdk:"resources"`
	CreateTime         types.String `tfsdk:"create_time"`
	UpdateTime         types.String `tfsdk:"update_time"`
}

type componentModel struct {
	ServiceComponentId types.String `tfsdk:"service_component_id"`
	Priority           types.Int64  `tfsdk:"priority"`
	Timeout            types.Int64  `tfsdk:"timeout"`
	Config             types.String `tfsdk:"config"`
	FailPolicy         types.String `tfsdk:"fail_policy"`
	ComponentName      types.String `tfsdk:"component_name"`
}

type extensionResourceModel struct {
	AcceleratorId types.String `tfsdk:"accelerator_id"`
	ResourceType  types.String `tfsdk:"resource_type"`
	ResourceId    types.String `tfsdk:"resource_id"`
	AssociateId   types.String `tfsdk:"associate_id"`
}

func (r *serviceExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = serviceExtensionTypeName
}

func (r *serviceExtensionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Global Accelerator (GA) Service Extension resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the service extension, same as `service_extension_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_extension_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the service extension.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the service extension.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the service extension.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the service extension.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the service extension.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to assign to the resource.",
			},
			"components": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of service components added to the service extension.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service_component_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the service component.",
						},
						"priority": schema.Int64Attribute{
							Optional:    true,
							Description: "The priority of the service component.",
						},
						"timeout": schema.Int64Attribute{
							Optional:    true,
							Description: "The timeout of the service component. Unit: seconds.",
						},
						"config": schema.StringAttribute{
							Optional:    true,
							Description: "The configuration of the service component.",
						},
						"fail_policy": schema.StringAttribute{
							Optional:    true,
							Description: "The fail policy of the service component.",
						},
						"component_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the service component.",
						},
					},
				},
			},
			"resources": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of resources associated with the service extension.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accelerator_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the accelerator to associate.",
						},
						"resource_type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the resource to associate.",
						},
						"resource_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the resource to associate.",
						},
						"associate_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the association record.",
						},
					},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the service extension was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the service extension was last updated.",
			},
		},
	}
}

func (r *serviceExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceExtensionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Client()

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("CreateServiceExtension"),
		"Name":        plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.ResourceGroupId.IsNull() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}

	response, err := rpcPostRetry(client, "CreateServiceExtension", request, defaultCreateTimeout, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateServiceExtension", serviceExtensionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, serviceExtensionTypeName, "CreateServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	idRaw, err := jsonpath.Get("$.ServiceExtensionId", response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading ServiceExtensionId", serviceExtensionTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, serviceExtensionTypeName, "$.ServiceExtensionId", response).Error(),
		)
		return
	}
	id := fmt.Sprint(idRaw)

	stateConf := buildStateConf([]string{}, []string{"active"}, defaultCreateTimeout, 5*time.Second,
		serviceExtensionStateRefreshFunc(client, id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for active", serviceExtensionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "CreateServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// The SDKv2 resource deferred components/resources to a later Update whose
	// branches were gated by !IsNewResource and skipped on a fresh resource.
	// terraform-plugin-framework rejects an applied state that diverges from the
	// planned Optional+Computed set, so create applies tags, components and
	// associated resources up front instead of leaving them as drift.
	newTags := tagsToMapFromState(ctx, plan.Tags)
	if err := applyServiceExtensionTags(client, id, map[string]string{}, newTags); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: applying tags", serviceExtensionTypeName),
			err.Error(),
		)
		return
	}

	emptyComponents, _ := componentSet(ctx, nil)
	if err := updateServiceExtensionComponents(ctx, client, id, emptyComponents, plan.Components); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: adding service components", serviceExtensionTypeName),
			err.Error(),
		)
		return
	}
	emptyResources, _ := resourceSet(ctx, nil)
	if err := updateServiceExtensionResources(ctx, client, id, emptyResources, plan.Resources); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: associating resources", serviceExtensionTypeName),
			err.Error(),
		)
		return
	}

	state, diags := readServiceExtensionState(ctx, client, id, plan.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serviceExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceExtensionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Client()

	refreshed, diags := readServiceExtensionState(ctx, client, state.Id.ValueString(), state.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if refreshed.Id.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *serviceExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state serviceExtensionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Client()
	id := state.Id.ValueString()

	// tags
	if !plan.Tags.Equal(state.Tags) {
		if err := applyServiceExtensionTags(client, id, tagsToMapFromState(ctx, state.Tags), tagsToMapFromState(ctx, plan.Tags)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: applying tags", serviceExtensionTypeName),
				err.Error(),
			)
			return
		}
	}

	// name/description via UpdateServiceExtension
	if !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) {
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ClientToken":        helper.BuildClientToken("UpdateServiceExtension"),
			"ServiceExtensionId": id,
			"Name":               plan.Name.ValueString(),
		}
		if !plan.Description.IsNull() {
			request["Description"] = plan.Description.ValueString()
		}
		if _, err := rpcPostRetry(client, "UpdateServiceExtension", request, defaultUpdateTimeout, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateServiceExtension", serviceExtensionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "UpdateServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		// UpdateServiceExtension is asynchronous: converge updating -> active.
		stateConf := buildStateConf([]string{"updating"}, []string{"active"}, defaultUpdateTimeout, 5*time.Second,
			serviceExtensionStateRefreshFunc(client, id, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for active after UpdateServiceExtension", serviceExtensionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "UpdateServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	if !plan.Components.Equal(state.Components) {
		if err := updateServiceExtensionComponents(ctx, client, id, state.Components, plan.Components); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: reconciling components", serviceExtensionTypeName),
				err.Error(),
			)
			return
		}
	}

	if !plan.Resources.Equal(state.Resources) {
		if err := updateServiceExtensionResources(ctx, client, id, state.Resources, plan.Resources); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: reconciling resources", serviceExtensionTypeName),
				err.Error(),
			)
			return
		}
	}

	if !plan.ResourceGroupId.Equal(state.ResourceGroupId) {
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"ClientToken":  helper.BuildClientToken("ChangeResourceGroup"),
			"ResourceId":   id,
			"ResourceType": serviceExtensionResourceType,
		}
		if !plan.ResourceGroupId.IsNull() {
			request["NewResourceGroupId"] = plan.ResourceGroupId.ValueString()
		}
		if _, err := rpcPostRetry(client, "ChangeResourceGroup", request, defaultUpdateTimeout, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ChangeResourceGroup", serviceExtensionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "ChangeResourceGroup", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	refreshed, diags := readServiceExtensionState(ctx, client, id, plan.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *serviceExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceExtensionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Client()
	id := state.Id.ValueString()

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ClientToken":        helper.BuildClientToken("DeleteServiceExtension"),
		"ServiceExtensionId": id,
	}
	_, err := rpcPostRetry(client, "DeleteServiceExtension", request, defaultDeleteTimeout, true)
	if err != nil {
		if errs.NotFoundError(err) || isExpectedErrors(err, []string{"NotExist.ServiceExtension"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteServiceExtension", serviceExtensionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Wait for the resource to disappear.
	stateConf := buildStateConf([]string{}, []string{}, defaultDeleteTimeout, 5*time.Second,
		serviceExtensionStateRefreshFunc(client, id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for removal", serviceExtensionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteServiceExtension", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *serviceExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readServiceExtensionState fetches the service extension and maps it onto the
// resource model. A missing resource is signalled by a null Id so callers can
// remove it from state. Tagging returns Invalid.ResourceType for service
// extensions; when that happens the prior tags are preserved rather than
// overwritten with an empty set, mirroring the SDKv2 resource which skipped
// d.Set("tags") on that error and so kept the configured tags in state.
func readServiceExtensionState(ctx context.Context, client *connectivity.AliyunClient, id string, priorTags types.Map) (serviceExtensionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	object, err := describeServiceExtension(client, id)
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource %s DescribeGaServiceExtension not found: %s", serviceExtensionTypeName, err)
			return serviceExtensionResourceModel{Id: types.StringNull()}, diags
		}
		diags.AddError(
			fmt.Sprintf("Reading %s: calling GetServiceExtension", serviceExtensionTypeName),
			errs.WrapError(err).Error(),
		)
		return serviceExtensionResourceModel{}, diags
	}

	model := serviceExtensionResourceModel{
		Id:                 types.StringValue(id),
		ServiceExtensionId: stringValue(object["ServiceExtensionId"]),
		Name:               stringValue(object["Name"]),
		Description:        stringValue(object["Description"]),
		Type:               stringValue(object["Type"]),
		State:              stringValue(object["State"]),
		ResourceGroupId:    stringValue(object["ResourceGroupId"]),
		CreateTime:         stringValue(object["CreateTime"]),
		UpdateTime:         stringValue(object["UpdateTime"]),
	}

	if v, ok := object["Components"].([]interface{}); ok {
		components, d := componentSet(ctx, v)
		model.Components = components
		diags.Append(d...)
	} else {
		set, d := componentSet(ctx, nil)
		model.Components = set
		diags.Append(d...)
	}

	if v, ok := object["Resources"].([]interface{}); ok {
		resources, d := resourceSet(ctx, v)
		model.Resources = resources
		diags.Append(d...)
	} else {
		set, d := resourceSet(ctx, nil)
		model.Resources = set
		diags.Append(d...)
	}

	tags, terr := listServiceExtensionTags(client, id)
	if terr != nil {
		if isExpectedErrors(terr, []string{"Invalid.ResourceType"}) {
			log.Printf("[WARN] Resource %s tagging is not supported by the backend (Invalid.ResourceType), preserving prior tags for %s: %s", serviceExtensionTypeName, id, terr)
			model.Tags = priorTags
		} else {
			diags.AddError(
				fmt.Sprintf("Reading %s: calling ListTagResources", serviceExtensionTypeName),
				errs.WrapError(terr).Error(),
			)
			return model, diags
		}
	} else {
		model.Tags = tagsMapToValue(ctx, tags, &diags)
	}

	return model, diags
}

// updateServiceExtensionComponents reconciles the old and new component sets by
// service_component_id. Items differing in any non-identity field are removed
// and re-added, matching SDKv2 TypeSet semantics. Each asynchronous call waits
// for state to converge from updating back to active.
func updateServiceExtensionComponents(ctx context.Context, client *connectivity.AliyunClient, id string, oldSet, newSet types.Set) error {
	oldItems := componentsFromSet(ctx, oldSet)
	newItems := componentsFromSet(ctx, newSet)

	oldByID := indexComponents(oldItems)
	newByID := indexComponents(newItems)

	var toRemove, toAdd []componentModel
	for id2, old := range oldByID {
		new, ok := newByID[id2]
		if !ok || !componentsEqual(old, new) {
			toRemove = append(toRemove, old)
			if ok {
				toAdd = append(toAdd, new)
			}
		}
	}
	for id2, new := range newByID {
		if _, ok := oldByID[id2]; !ok {
			toAdd = append(toAdd, new)
		}
	}

	if len(toRemove) > 0 {
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ClientToken":        helper.BuildClientToken("RemoveServiceComponentsFromServiceExtension"),
			"ServiceExtensionId": id,
		}
		for i, c := range toRemove {
			request[fmt.Sprintf("ServiceComponentIds.%d", i+1)] = c.ServiceComponentId.ValueString()
		}
		if _, err := rpcPostRetry(client, "RemoveServiceComponentsFromServiceExtension", request, defaultUpdateTimeout, true); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "RemoveServiceComponentsFromServiceExtension", errs.AlibabaCloudSdkGoERROR)
		}
		if err := waitActive(client, id, defaultUpdateTimeout); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ClientToken":        helper.BuildClientToken("AddServiceComponentsToServiceExtension"),
			"ServiceExtensionId": id,
		}
		for i, c := range toAdd {
			request[fmt.Sprintf("ServiceComponents.%d.ServiceComponentId", i+1)] = c.ServiceComponentId.ValueString()
			request[fmt.Sprintf("ServiceComponents.%d.Priority", i+1)] = c.Priority.ValueInt64()
			request[fmt.Sprintf("ServiceComponents.%d.Timeout", i+1)] = c.Timeout.ValueInt64()
			if v := c.Config.ValueString(); v != "" {
				request[fmt.Sprintf("ServiceComponents.%d.Config", i+1)] = v
			}
			if v := c.FailPolicy.ValueString(); v != "" {
				request[fmt.Sprintf("ServiceComponents.%d.FailPolicy", i+1)] = v
			}
		}
		if _, err := rpcPostRetry(client, "AddServiceComponentsToServiceExtension", request, defaultUpdateTimeout, true); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "AddServiceComponentsToServiceExtension", errs.AlibabaCloudSdkGoERROR)
		}
		if err := waitActive(client, id, defaultUpdateTimeout); err != nil {
			return err
		}
	}
	return nil
}

// updateServiceExtensionResources reconciles the old and new associated resource
// sets by (accelerator_id, resource_type, resource_id). Dissociate takes the
// AssociateId recorded in state; associate sends the triple and a fresh
// AssociateId is read back. Each call waits for state convergence.
func updateServiceExtensionResources(ctx context.Context, client *connectivity.AliyunClient, id string, oldSet, newSet types.Set) error {
	oldItems := resourcesFromSet(ctx, oldSet)
	newItems := resourcesFromSet(ctx, newSet)

	oldByKey := indexResources(oldItems)
	newByKey := indexResources(newItems)

	var toRemove, toAdd []extensionResourceModel
	for k, old := range oldByKey {
		if _, ok := newByKey[k]; !ok {
			toRemove = append(toRemove, old)
		}
	}
	for k, new := range newByKey {
		if _, ok := oldByKey[k]; !ok {
			toAdd = append(toAdd, new)
		}
	}

	if len(toRemove) > 0 {
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ClientToken":        helper.BuildClientToken("DissociateResourcesFromServiceExtension"),
			"ServiceExtensionId": id,
		}
		for i, r := range toRemove {
			if associateId := r.AssociateId.ValueString(); associateId != "" {
				request[fmt.Sprintf("AssociateIds.%d", i+1)] = associateId
			}
		}
		if _, err := rpcPostRetry(client, "DissociateResourcesFromServiceExtension", request, defaultUpdateTimeout, true); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DissociateResourcesFromServiceExtension", errs.AlibabaCloudSdkGoERROR)
		}
		if err := waitActive(client, id, defaultUpdateTimeout); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ClientToken":        helper.BuildClientToken("AssociateResourcesWithServiceExtension"),
			"ServiceExtensionId": id,
		}
		for i, r := range toAdd {
			request[fmt.Sprintf("Resources.%d.AcceleratorId", i+1)] = r.AcceleratorId.ValueString()
			request[fmt.Sprintf("Resources.%d.ResourceType", i+1)] = r.ResourceType.ValueString()
			request[fmt.Sprintf("Resources.%d.ResourceId", i+1)] = r.ResourceId.ValueString()
		}
		if _, err := rpcPostRetry(client, "AssociateResourcesWithServiceExtension", request, defaultUpdateTimeout, true); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "AssociateResourcesWithServiceExtension", errs.AlibabaCloudSdkGoERROR)
		}
		if err := waitActive(client, id, defaultUpdateTimeout); err != nil {
			return err
		}
	}
	return nil
}

// waitActive waits for a service extension to converge from updating to active.
func waitActive(client *connectivity.AliyunClient, id string, timeout time.Duration) error {
	stateConf := buildStateConf([]string{"updating"}, []string{"active"}, timeout, 5*time.Second,
		serviceExtensionStateRefreshFunc(client, id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "WaitForState", errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
