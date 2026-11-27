package vpc

import (
	"context"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliyunNetworkAclAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliyunNetworkAclAttachmentCreate,
		ReadContext:   resourceAliyunNetworkAclAttachmentRead,
		UpdateContext: resourceAliyunNetworkAclAttachmentUpdate,
		DeleteContext: resourceAliyunNetworkAclAttachmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		DeprecationMessage: "This resource has been deprecated since v1.124.0 and will be removed in the future. Please use 'resources' in 'alicloud_network_acl' instead.",
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func networkAclAttachmentRemaining(deadline, now time.Time) time.Duration {
	return deadline.Sub(now)
}

func networkAclAttachmentRetry(ctx context.Context, deadline time.Time, f func() error, additionalRetryCodes ...string) error {
	transientCodes := append([]string{"NetworkStatus.Modifying", "ResourceStatus.Error", "OperationDenied.NetworkAclAttachmentInMiddleStatus"}, additionalRetryCodes...)
	return retry.RetryContext(ctx, networkAclAttachmentRemaining(deadline, time.Now()), func() *retry.RetryError {
		err := f()
		if err == nil {
			return nil
		}
		if IsExpectedErrors(err, transientCodes) || NeedRetry(err) {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
}

func networkAclAttachmentDeleteIdempotent(err error) bool {
	return IsExpectedErrors(err, []string{"InvalidNetworkAcl.NotFound", "InvalidResource.NotBinding", "InvalidResource.NotFound", "Mismatch.NetworkAclAndVSwitch"})
}

func resourceAliyunNetworkAclAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("network_acl_id").(string) + COLON_SEPARATED + id.UniqueId())

	return resourceAliyunNetworkAclAttachmentUpdateWithDeadline(ctx, d, meta, time.Now().Add(d.Timeout(schema.TimeoutCreate)))
}

func resourceAliyunNetworkAclAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return diag.FromErr(WrapError(err))
	}
	networkAclId := parts[0]
	vpcResource := []vpc.Resource{}
	for _, e := range d.Get("resources").(*schema.Set).List() {
		resourceId := e.(map[string]interface{})["resource_id"]
		resourceType := e.(map[string]interface{})["resource_type"]
		vpcResource = append(vpcResource, vpc.Resource{
			ResourceId:   resourceId.(string),
			ResourceType: resourceType.(string),
		})
	}
	vpcResource, err = vpcService.DescribeNetworkAclAttachment(ctx, networkAclId, vpcResource)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(WrapError(err))
	}
	d.Set("network_acl_id", networkAclId)

	vpcRes := make([]map[string]interface{}, 0)
	for _, v := range vpcResource {
		vpcRes = append(vpcRes, map[string]interface{}{
			"resource_id":   v.ResourceId,
			"resource_type": v.ResourceType,
		})
	}
	d.Set("resources", vpcRes)
	return nil
}

func resourceAliyunNetworkAclAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceAliyunNetworkAclAttachmentUpdateWithDeadline(ctx, d, meta, time.Now().Add(d.Timeout(schema.TimeoutUpdate)))
}

func resourceAliyunNetworkAclAttachmentUpdateWithDeadline(ctx context.Context, d *schema.ResourceData, meta interface{}, deadline time.Time) diag.Diagnostics {
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return diag.FromErr(WrapError(err))
	}
	networkAclId := parts[0]
	if d.HasChange("resources") {
		oraw, nraw := d.GetChange("resources")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		remove := o.Difference(n).List()
		create := n.Difference(o).List()

		if len(remove) > 0 {
			request := vpc.CreateUnassociateNetworkAclRequest()
			request.RegionId = client.RegionId
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.UnassociateNetworkAclResource
			removed := make([]vpc.Resource, 0, len(remove))
			for _, t := range remove {
				s := t.(map[string]interface{})
				var resourceId, resourceType string
				if v, ok := s["resource_id"]; ok {
					resourceId = v.(string)
				}
				if v, ok := s["resource_type"]; ok {
					resourceType = v.(string)
				}
				resources = append(resources, vpc.UnassociateNetworkAclResource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
				removed = append(removed, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			err = networkAclAttachmentRetry(ctx, deadline, func() error {
				raw, err := client.VpcSdkRpc(request, vpc.CreateUnassociateNetworkAclResponse())
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return err
			})
			if err != nil && !networkAclAttachmentDeleteIdempotent(err) {
				return diag.FromErr(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			if err := vpcService.WaitForNetworkAclAttachment(ctx, request.NetworkAclId, removed, Deleted, deadline); err != nil {
				return diag.FromErr(WrapError(err))
			}
		}

		if len(create) > 0 {
			request := vpc.CreateAssociateNetworkAclRequest()
			request.RegionId = client.RegionId
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.AssociateNetworkAclResource
			for _, t := range create {
				s := t.(map[string]interface{})
				var resourceId, resourceType string
				if v, ok := s["resource_id"]; ok {
					resourceId = v.(string)
				}
				if v, ok := s["resource_type"]; ok {
					resourceType = v.(string)
				}
				resources = append(resources, vpc.AssociateNetworkAclResource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			err = networkAclAttachmentRetry(ctx, deadline, func() error {
				raw, err := client.VpcSdkRpc(request, vpc.CreateAssociateNetworkAclResponse())
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return err
			})
			if err != nil {
				return diag.FromErr(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			desired := make([]vpc.Resource, 0)
			for _, t := range n.List() {
				s := t.(map[string]interface{})
				resourceId, _ := s["resource_id"].(string)
				resourceType, _ := s["resource_type"].(string)
				desired = append(desired, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			if err := vpcService.WaitForNetworkAclAttachment(ctx, request.NetworkAclId, desired, Available, deadline); err != nil {
				return diag.FromErr(WrapError(err))
			}
		}
	}

	return resourceAliyunNetworkAclAttachmentRead(ctx, d, meta)
}

func resourceAliyunNetworkAclAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	deadline := time.Now().Add(d.Timeout(schema.TimeoutDelete))
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return diag.FromErr(WrapError(err))
	}
	networkAclId := parts[0]

	vpcResource := make([]vpc.Resource, 0)
	for _, e := range d.Get("resources").(*schema.Set).List() {
		s := e.(map[string]interface{})
		resourceId, _ := s["resource_id"].(string)
		resourceType, _ := s["resource_type"].(string)
		vpcResource = append(vpcResource, vpc.Resource{
			ResourceId:   resourceId,
			ResourceType: resourceType,
		})
	}

	object, err := vpcService.DescribeNetworkAclWithContext(ctx, networkAclId)
	if err != nil {
		if !NotFoundError(err) {
			return diag.FromErr(WrapError(err))
		}
		object = nil
	}

	// The network acl can be shared by several attachment instances, so only the resources recorded
	// in this instance's state may be unbound; unassociating every resource of the acl would release
	// the bindings that belong to the sibling attachments.
	resources, err := networkAclAttachmentResourcesToUnbind(object, vpcResource)
	if err != nil {
		return diag.FromErr(WrapError(err))
	}
	if len(resources) > 0 {
		request := vpc.CreateUnassociateNetworkAclRequest()
		request.RegionId = client.RegionId
		request.NetworkAclId = networkAclId
		request.ClientToken = buildClientToken(request.GetActionName())
		request.Resource = &resources
		err = networkAclAttachmentRetry(ctx, deadline, func() error {
			raw, err := client.VpcSdkRpc(request, vpc.CreateUnassociateNetworkAclResponse())
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return err
		}, "TaskConflict")
		if err != nil && !networkAclAttachmentDeleteIdempotent(err) {
			return diag.FromErr(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
	}

	return diag.FromErr(vpcService.WaitForNetworkAclAttachment(ctx, networkAclId, vpcResource, Deleted, deadline))
}
