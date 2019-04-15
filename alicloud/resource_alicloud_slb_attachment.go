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
	loadBalancer, err := slbService.DescribeLoadBalancerAttribute(d.Get("load_balancer_id").(string))
	if err != nil {
		return WrapError(err)
	}

	d.SetId(loadBalancer.LoadBalancerId)

	return resourceAliyunSlbAttachmentUpdate(d, meta)
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	loadBalancer, err := slbService.DescribeLoadBalancerAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	backendServerType := loadBalancer.BackendServers
	servers := backendServerType.BackendServer
	instanceIds := make([]string, 0, len(servers))
	var weight int
	if len(servers) > 0 {
		weight = servers[0].Weight
		for _, e := range servers {
			instanceIds = append(instanceIds, e.ServerId)
		}
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("load_balancer_id", loadBalancer.LoadBalancerId)
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
			req := slb.CreateAddBackendServersRequest()
			req.LoadBalancerId = d.Id()
			req.BackendServers = expandBackendServersToString(ns.Difference(os).List(), weight)
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddBackendServers(req)
				})
				if err != nil {
					if IsExceptedErrors(err, SlbIsBusy) {
						return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
					}
					if IsExceptedErrors(err, BackendServerNotReadyStatus) {
						return resource.RetryableError(WrapErrorf(err, DefaultDebugMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
					}
					return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
				}
				addDebug(req.GetActionName(), raw)
				return nil
			}); err != nil {
				return WrapError(err)
			}
		}
		if len(remove) > 0 {
			if err := removeBackendServers(d, meta, remove); err != nil {
				return WrapError(err)
			}
		}

		if len(add) < 1 && len(remove) < 1 {
			update = true
		}
		d.SetPartial("instance_ids")
	}

	if update {
		req := slb.CreateSetBackendServersRequest()
		req.LoadBalancerId = d.Id()
		req.BackendServers = expandBackendServersToString(d.Get("instance_ids").(*schema.Set).List(), weight)
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetBackendServers(req)
			})
			if err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
				}
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			addDebug(req.GetActionName(), raw)
			return nil
		}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliyunSlbAttachmentRead(d, meta)

}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	return removeBackendServers(d, meta, d.Get("instance_ids").(*schema.Set).List())
}

func removeBackendServers(d *schema.ResourceData, meta interface{}, servers []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	instanceSet := d.Get("instance_ids").(*schema.Set)
	if len(servers) > 0 {
		req := slb.CreateRemoveBackendServersRequest()
		req.LoadBalancerId = d.Id()
		req.BackendServers = convertListToJsonString(servers)
		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.RemoveBackendServers(req)
			})
			if err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
				}
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			addDebug(req.GetActionName(), raw)
			loadBalancer, err := slbService.DescribeLoadBalancerAttribute(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))

			}

			servers := loadBalancer.BackendServers.BackendServer

			if len(servers) > 0 {
				for _, e := range servers {
					if instanceSet.Contains(e.ServerId) {
						return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), req.GetActionName(), ProviderERROR))
					}
				}
			}
			return nil
		})
	}
	return nil
}
