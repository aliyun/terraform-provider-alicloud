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
	request := slb.CreateCreateVServerGroupRequest()
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := d.GetOk("name"); ok {
		request.VServerGroupName = v.(string)
	}
	if v, ok := d.GetOk("servers"); ok {
		request.BackendServers = expandBackendServersWithPortToString(v.(*schema.Set).List())
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateVServerGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_server_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*slb.CreateVServerGroupResponse)
	d.SetId(response.VServerGroupId)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbServerGroup(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.VServerGroupName)
	d.Set("load_balancer_id", object.LoadBalancerId)

	servers := make([]map[string]interface{}, 0)
	portAndWeight := make(map[string][]string)
	for _, server := range object.BackendServers.BackendServer {
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
			return WrapError(e)
		}
		w, e := strconv.Atoi(k[1])
		if e != nil {
			return WrapError(e)
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

	if d.HasChange("name") {
		d.SetPartial("name")
		update = true
	}

	if d.HasChange("servers") {
		o, n := d.GetChange("servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			request := slb.CreateRemoveVServerGroupBackendServersRequest()
			request.VServerGroupId = d.Id()
			request.BackendServers = expandBackendServersWithPortToString(remove)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.RemoveVServerGroupBackendServers(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}
		if len(add) > 0 {
			request := slb.CreateAddVServerGroupBackendServersRequest()
			request.VServerGroupId = d.Id()
			request.BackendServers = expandBackendServersWithPortToString(add)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.AddVServerGroupBackendServers(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}
		if len(add) < 1 && len(remove) < 1 {
			update = true
		}

		d.SetPartial("servers")
	}

	if update {
		request := slb.CreateSetVServerGroupAttributeRequest()
		request.VServerGroupId = d.Id()
		request.VServerGroupName = name
		request.BackendServers = expandBackendServersWithPortToString(d.Get("servers").(*schema.Set).List())
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetVServerGroupAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	d.Partial(false)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteVServerGroupRequest()
	request.VServerGroupId = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteVServerGroup(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{RspoolVipExist}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(slbService.WaitForSlbServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
