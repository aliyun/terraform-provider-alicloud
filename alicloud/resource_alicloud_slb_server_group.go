package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbServerGroupCreate,
		Read:   resourceAliyunSlbServerGroupRead,
		Update: resourceAliyunSlbServerGroupUpdate,
		Delete: resourceAliyunSlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-server-group",
			},

			"servers": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ids": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MaxItems: 20,
							MinItems: 1,
						},
						"port": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(1, 65535),
						},
						"weight": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validateIntegerInRange(0, 100),
						},
					},
				},
				MaxItems: 20,
				MinItems: 0,
			},
		},
	}
}

func resourceAliyunSlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	var groupId string
	req := slb.CreateCreateVServerGroupRequest()
	req.LoadBalancerId = d.Get("load_balancer_id").(string)
	req.VServerGroupName = d.Get("name").(string)
	req.BackendServers = expandBackendServersWithPortToString(d.Get("servers").(*schema.Set).List())
	if group, err := meta.(*AliyunClient).slbconn.CreateVServerGroup(req); err != nil {
		return fmt.Errorf("CreateVServerGroup got an error: %#v", err)
	} else {
		groupId = group.VServerGroupId
	}

	d.SetId(groupId)

	return resourceAliyunSlbServerGroupUpdate(d, meta)
}

func resourceAliyunSlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	group, err := meta.(*AliyunClient).DescribeSlbVServerGroupAttribute(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", group.VServerGroupName)
	d.Set("load_balancer_id", d.Get("load_balancer_id").(string))

	servers := make([]map[string]interface{}, 0)
	portAndWeight := make(map[string][]string)
	for _, server := range group.BackendServers.BackendServer {
		key := fmt.Sprintf("%d%s%d", server.Port, COLON_SEPARATED, server.Weight)
		if v, ok := portAndWeight[key]; !ok {
			portAndWeight[key] = []string{server.ServerId}
		} else {
			v = append(v, server.ServerId)
			portAndWeight[key] = v
		}
	}
	for key, value := range portAndWeight {
		k := strings.Split(key, COLON_SEPARATED)
		p, e := strconv.Atoi(k[0])
		if e != nil {
			return fmt.Errorf("Convertting port %s to int got an error: %#v.", k[0], e)
		}
		w, e := strconv.Atoi(k[1])
		if e != nil {
			return fmt.Errorf("Convertting weight %s to int got an error: %#v.", k[1], e)
		}
		s := map[string]interface{}{
			"server_ids": value,
			"port":       p,
			"weight":     w,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return err
	}

	return nil
}

func resourceAliyunSlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn

	d.Partial(true)

	name := d.Get("name").(string)
	update := false

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		update = true
	}

	if d.HasChange("servers") && !d.IsNewResource() {
		o, n := d.GetChange("servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			req := slb.CreateRemoveVServerGroupBackendServersRequest()
			req.VServerGroupId = d.Id()
			req.BackendServers = expandBackendServersWithPortToString(remove)
			if _, err := slbconn.RemoveVServerGroupBackendServers(req); err != nil {
				return fmt.Errorf("RemoveVServerGroupBackendServers got an error: %#v", err)
			}
		}
		if len(add) > 0 {
			req := slb.CreateAddVServerGroupBackendServersRequest()
			req.VServerGroupId = d.Id()
			req.BackendServers = expandBackendServersWithPortToString(add)
			if _, err := slbconn.AddVServerGroupBackendServers(req); err != nil {
				return fmt.Errorf("AddVServerGroupBackendServers got an error: %#v", err)
			}
		}
		if len(add) < 1 && len(remove) < 1 {
			update = true
		}

		d.SetPartial("servers")
	}

	if update {
		req := slb.CreateSetVServerGroupAttributeRequest()
		req.VServerGroupId = d.Id()
		req.VServerGroupName = name
		req.BackendServers = expandBackendServersWithPortToString(d.Get("servers").(*schema.Set).List())
		if _, err := slbconn.SetVServerGroupAttribute(req); err != nil {
			return fmt.Errorf("SetVServerGroupAttribute got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	req := slb.CreateDeleteVServerGroupRequest()
	req.VServerGroupId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.slbconn.DeleteVServerGroup(req); err != nil {
			if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
				return nil
			}
			if IsExceptedErrors(err, []string{RspoolVipExist}) {
				return resource.RetryableError(fmt.Errorf("DeleteVServerGroup got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}

		if _, err := meta.(*AliyunClient).DescribeSlbVServerGroupAttribute(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While deleting VServer Group, DescribeVServerGroupAttribute got an error: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("DeleteVServerGroup %s timeout.", d.Id()))
	})
}
