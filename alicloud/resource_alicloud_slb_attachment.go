package alicloud

import (
	"fmt"

	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

			"slb_id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.6.0. New field 'load_balancer_id' replaces it.",
			},

			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instances' has been deprecated from provider version 1.6.0. New field 'instance_ids' replaces it.",
			},

			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"weight": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validateIntegerInRange(0, 100),
			},

			"backend_servers": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(d.Get("load_balancer_id").(string))
	if err != nil {
		return err
	}

	d.SetId(loadBalancer.LoadBalancerId)

	return resourceAliyunSlbAttachmentUpdate(d, meta)
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
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
			return err
		}
	}

	d.Set("load_balancer_id", loadBalancer.LoadBalancerId)
	d.Set("instance_ids", instanceIds)
	d.Set("weight", weight)
	d.Set("backend_servers", strings.Join(instanceIds, ","))

	return nil
}

func resourceAliyunSlbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn
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
		add := expandBackendServers(ns.Difference(os).List(), weight)

		if len(add) > 0 {
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				_, err := slbconn.AddBackendServers(d.Id(), add)
				if err != nil {
					if IsExceptedErrors(err, SlbIsBusy) {
						return resource.RetryableError(fmt.Errorf("Load banalcer adds backend servers timeout and got an error: %#v.", err))
					}
					return resource.NonRetryableError(fmt.Errorf("Add backend servers got an error: %#v", err))
				}
				return nil
			}); err != nil {
				return err
			}
		}
		if len(remove) > 0 {
			if err := removeBackendServers(d, meta, remove); err != nil {
				return err
			}
		}

		if len(add) < 1 && len(remove) < 1 {
			update = true
		}
		d.SetPartial("instance_ids")
	}

	if update {
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			if _, err := slbconn.SetBackendServers(d.Id(), expandBackendServers(d.Get("instance_ids").(*schema.Set).List(), weight)); err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(fmt.Errorf("Load banalcer sets backend servers timeout and got an error: %#v.", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Set backend servers got an error: %#v", err))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceAliyunSlbAttachmentRead(d, meta)

}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	return removeBackendServers(d, meta, d.Get("instance_ids").(*schema.Set).List())
}

func removeBackendServers(d *schema.ResourceData, meta interface{}, servers []interface{}) error {
	client := meta.(*AliyunClient)
	instanceSet := d.Get("instance_ids").(*schema.Set)
	if len(servers) > 0 {

		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.slbconn.RemoveBackendServers(d.Id(), convertArrayInterfaceToArrayString(servers))
			if err != nil {
				if IsExceptedErrors(err, SlbIsBusy) {
					return resource.RetryableError(fmt.Errorf("Load balancer removes backend servers timeout and got an error: %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Remove backend servers got an error: %#v", err))
			}

			loadBalancer, err := client.DescribeLoadBalancerAttribute(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return resource.NonRetryableError(fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err))

			}

			servers := loadBalancer.BackendServers.BackendServer

			if len(servers) > 0 {
				for _, e := range servers {
					if instanceSet.Contains(e.ServerId) {
						return resource.RetryableError(fmt.Errorf("There are still target backend servers in the SLB."))
					}
				}
			}
			return nil
		})
	}
	return nil
}
