package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbAttachmentCreate,
		Read:   resourceAliyunSlbAttachmentRead,
		Update: resourceAliyunSlbAttachmentUpdate,
		Delete: resourceAliyunSlbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"slb_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.6.0. New field 'load_balancer_id' replaces it.",
			},

			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instances' has been deprecated from provider version 1.6.0. New field 'instance_ids' replaces it.",
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validateIntegerInRange(0, 100),
			},

			"backend_servers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Get("load_balancer_id").(string))
	if err != nil {
		return WrapError(err)
	}
	d.SetId(object.LoadBalancerId)

	return resourceAliyunSlbAttachmentUpdate(d, meta)
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	servers := object.BackendServers.BackendServer
	instanceIds := make([]string, 0, len(servers))
	var weight int
	if len(servers) > 0 {
		weight = servers[0].Weight
		for _, e := range servers {
			instanceIds = append(instanceIds, e.ServerId)
		}
	}

	d.Set("load_balancer_id", object.LoadBalancerId)
	d.Set("instance_ids", instanceIds)
	d.Set("weight", weight)
	d.Set("backend_servers", strings.Join(instanceIds, ","))

	return nil
}

func resourceAliyunSlbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	update := false
	weight := d.Get("weight").(int)

	if d.HasChange("weight") {
		update = true
		d.SetPartial("weight")
	}
	if d.HasChange("instance_ids") {
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(add) > 0 {
			request := slb.CreateAddBackendServersRequest()
			request.RegionId = client.RegionId
			request.LoadBalancerId = d.Id()
			request.BackendServers = expandBackendServersToString(ns.Difference(os).List(), weight)
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddBackendServers(request)
				})
				if err != nil {
					if IsExceptedErrors(err, SlbIsBusy) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
		if len(remove) > 0 {
			request := slb.CreateRemoveBackendServersRequest()
			request.RegionId = client.RegionId
			request.LoadBalancerId = d.Id()
			request.BackendServers = expandBackendServersToString(os.Difference(ns).List(), weight)
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveBackendServers(request)
				})
				if err != nil {
					if IsExceptedErrors(err, SlbIsBusy) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		if len(add) < 1 && len(remove) < 1 {
			update = true
		}
		d.SetPartial("instance_ids")
	}

	if update {
		request := slb.CreateSetBackendServersRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.BackendServers = expandBackendServersToString(d.Get("instance_ids").(*schema.Set).List(), weight)
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetBackendServers(request)
			})
			if err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliyunSlbAttachmentRead(d, meta)

}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	instanceSet := d.Get("instance_ids").(*schema.Set)
	if len(instanceSet.List()) > 0 {
		request := slb.CreateRemoveBackendServersRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.BackendServers = convertListToJsonString(instanceSet.List())
		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.RemoveBackendServers(request)
			})
			if err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return WrapError(slbService.WaitSlbAttribute(d.Id(), instanceSet, DefaultTimeout))
}
