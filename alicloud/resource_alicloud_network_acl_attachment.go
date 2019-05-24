package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunNetworkAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNetworkAclAttachmentCreate,
		Read:   resourceAliyunNetworkAclAttachmentRead,
		Update: resourceAliyunNetworkAclAttachmentUpdate,
		Delete: resourceAliyunNetworkAclAttachmentDelete,

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

func resourceAliyunNetworkAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("network_acl_id").(string) + COLON_SEPARATED + resource.UniqueId())

	return resourceAliyunNetworkAclAttachmentUpdate(d, meta)
}

func resourceAliyunNetworkAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	networkAclService := NetworkAclService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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
	err = networkAclService.DescribeNetworkAclAttachment(networkAclId, vpcResource)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("network_acl_id", networkAclId)
	d.Set("resources", vpcResource)
	return nil
}

func resourceAliyunNetworkAclAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	networkAclService := NetworkAclService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.UnassociateNetworkAclResource
			vpcResource := []vpc.Resource{}
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
				vpcResource = append(vpcResource, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.UnassociateNetworkAcl(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName, raw)
			if err := networkAclService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, DefaultTimeout); err != nil {
				return WrapError(err)
			}
		}

		if len(create) > 0 {
			request := vpc.CreateAssociateNetworkAclRequest()
			request.NetworkAclId = networkAclId
			request.ClientToken = buildClientToken(request.GetActionName())
			var resources []vpc.AssociateNetworkAclResource
			vpcResource := []vpc.Resource{}
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
				vpcResource = append(vpcResource, vpc.Resource{
					ResourceId:   resourceId,
					ResourceType: resourceType,
				})
			}
			request.Resource = &resources
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.AssociateNetworkAcl(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName, raw)
			if err := networkAclService.WaitForNetworkAclAttachment(request.NetworkAclId, vpcResource, Available, DefaultTimeout); err != nil {
				return WrapError(err)
			}
		}
		d.SetPartial("resources")
	}

	return resourceAliyunNetworkAclAttachmentRead(d, meta)
}

func resourceAliyunNetworkAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	networkAclService := NetworkAclService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	networkAclId := parts[0]

	resources := []vpc.UnassociateNetworkAclResource{}
	object, err := networkAclService.DescribeNetworkAcl(networkAclId)
	vpcResource := []vpc.Resource{}
	request := vpc.CreateUnassociateNetworkAclRequest()
	request.NetworkAclId = networkAclId
	request.ClientToken = buildClientToken(request.GetActionName())
	for _, e := range object.Resources.Resource {

		resources = append(resources, vpc.UnassociateNetworkAclResource{
			ResourceId:   e.ResourceId,
			ResourceType: e.ResourceType,
		})
		vpcResource = append(vpcResource, vpc.Resource{
			ResourceId:   e.ResourceId,
			ResourceType: e.ResourceType,
		})
	}
	request.Resource = &resources
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateNetworkAcl(request)
		})
		//Waiting for unassociate the network acl
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(err)
			}
		}
		addDebug(request.GetActionName, raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return networkAclService.WaitForNetworkAclAttachment(networkAclId, vpcResource, Deleted, DefaultTimeout)
}
