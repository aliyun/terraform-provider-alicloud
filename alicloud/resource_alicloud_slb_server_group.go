package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-server-group",
			},

			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MaxItems: 20,
							MinItems: 1,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(1, 65535),
						},
						"weight": {
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
	client := meta.(*connectivity.AliyunClient)
	var groupId string
	req := slb.CreateCreateVServerGroupRequest()
	req.LoadBalancerId = d.Get("load_balancer_id").(string)
	req.VServerGroupName = d.Get("name").(string)
	req.BackendServers = expandBackendServersWithPortToString(d.Get("servers").(*schema.Set).List())
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateVServerGroup(req)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "slb_servere_group", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetActionName(), raw)
	group, _ := raw.(*slb.CreateVServerGroupResponse)
	groupId = group.VServerGroupId

	d.SetId(groupId)

	return resourceAliyunSlbServerGroupUpdate(d, meta)
}

func resourceAliyunSlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	group, err := slbService.DescribeSlbVServerGroupAttribute(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", group.VServerGroupName)
	d.Set("load_balancer_id", group.LoadBalancerId)

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
			return WrapError(fmt.Errorf("Convertting port %s to int got an error: %#v.", k[0], e))
		}
		w, e := strconv.Atoi(k[1])
		if e != nil {
			return WrapError(fmt.Errorf("Convertting weight %s to int got an error: %#v.", k[1], e))
		}
		s := map[string]interface{}{
			"server_ids": value,
			"port":       p,
			"weight":     w,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunSlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.RemoveVServerGroupBackendServers(req)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(req.GetActionName(), raw)
		}
		if len(add) > 0 {
			req := slb.CreateAddVServerGroupBackendServersRequest()
			req.VServerGroupId = d.Id()
			req.BackendServers = expandBackendServersWithPortToString(add)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.AddVServerGroupBackendServers(req)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(req.GetActionName(), raw)
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
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetVServerGroupAttribute(req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(req.GetActionName(), raw)
	}

	d.Partial(false)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	req := slb.CreateDeleteVServerGroupRequest()
	req.VServerGroupId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteVServerGroup(req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
				return nil
			}
			if IsExceptedErrors(err, []string{RspoolVipExist}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(req.GetActionName(), raw)

		if _, err := slbService.DescribeSlbVServerGroupAttribute(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), req.GetActionName(), ProviderERROR))
	})
}
