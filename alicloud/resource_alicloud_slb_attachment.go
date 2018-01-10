package alicloud

import (
	"fmt"

	"strings"
	"time"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbAttachmentCreate,
		Read:   resourceAliyunSlbAttachmentRead,
		Update: resourceAliyunSlbAttachmentUpdate,
		Delete: resourceAliyunSlbAttachmentDelete,

		Schema: map[string]*schema.Schema{

			"slb_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
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

	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(d.Get("slb_id").(string))
	if err != nil {
		return err
	}

	if loadBalancer == nil {
		d.SetId("")
		return fmt.Errorf("Special SLB Id %s is not found in %#v.", d.Get("slb_id").(string), getRegion(d, meta))
	}
	d.SetId(loadBalancer.LoadBalancerId)

	return resourceAliyunSlbAttachmentUpdate(d, meta)
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(d.Get("slb_id").(string))
	if err != nil {
		return err
	}

	if loadBalancer == nil {
		d.SetId("")
		return nil
	}

	backendServerType := loadBalancer.BackendServers
	servers := backendServerType.BackendServer
	instanceIds := make([]string, 0, len(servers))
	if len(servers) > 0 {
		for _, e := range servers {
			instanceIds = append(instanceIds, e.ServerId)
		}
		if err != nil {
			return err
		}
	}

	d.Set("slb_id", d.Id())
	d.Set("instances", instanceIds)
	d.Set("backend_servers", strings.Join(instanceIds, ","))

	return nil
}

func resourceAliyunSlbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn
	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandBackendServers(os.Difference(ns).List())
		add := expandBackendServers(ns.Difference(os).List())

		if len(add) > 0 {
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				_, err := slbconn.AddBackendServers(d.Id(), add)
				if err != nil {
					if IsExceptedError(err, ServiceIsConfiguring) {
						return resource.RetryableError(fmt.Errorf("Load banalcer adds backend servers timeout and got an error: %#v.", err))
					}
					return resource.NonRetryableError(fmt.Errorf("Add backend servers got an error: %#v", err))
				}
				return nil
			}); err != nil {
				return err
			}
		}
		if err := removeBackendServers(d, meta, remove); err != nil {
			return err
		}
	}

	return resourceAliyunSlbAttachmentRead(d, meta)

}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	o := d.Get("instances")
	os := o.(*schema.Set)
	remove := expandBackendServers(os.List())

	return removeBackendServers(d, meta, remove)
}

func removeBackendServers(d *schema.ResourceData, meta interface{}, servers []slb.BackendServerType) error {
	slbconn := meta.(*AliyunClient).slbconn
	if len(servers) > 0 {
		removeBackendServers := make([]string, 0, len(servers))
		for _, e := range servers {
			removeBackendServers = append(removeBackendServers, e.ServerId)
		}
		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := slbconn.RemoveBackendServers(d.Id(), removeBackendServers)
			if err != nil {
				if IsExceptedError(err, BackendServerconfiguring) {
					return resource.RetryableError(fmt.Errorf("Load balancer removes backend servers timeout and got an error: %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Remove backend servers got an error: %#v", err))
			}
			return nil
		})
	}
	return nil
}
