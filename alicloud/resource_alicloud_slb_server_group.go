package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ids": {
							Type:       schema.TypeList,
							Optional:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							MinItems:   1,
							Deprecated: "Field 'server_ids' has been deprecated from provider version 1.93.0. Use 'server_id' replaces it.",
						},
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(ECS),
							ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
						},
						"server_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunSlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateCreateVServerGroupRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := d.GetOk("name"); ok {
		request.VServerGroupName = v.(string)
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateVServerGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_server_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateVServerGroupResponse)
	d.SetId(response.VServerGroupId)

	return resourceAliyunSlbServerGroupUpdate(d, meta)
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
	for _, server := range object.BackendServers.BackendServer {
		s := map[string]interface{}{
			"server_id": server.ServerId,
			"port":      server.Port,
			"weight":    server.Weight,
			"type":      server.Type,
			"server_ip": server.ServerIp,
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
	var removeserverSet, addServerSet, updateServerSet *schema.Set
	serverUpdate := false
	step := 20
	if d.HasChange("servers") {
		o, n := d.GetChange("servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()
		oldIdPort := getIdPortSetFromServers(remove)
		newIdPort := getIdPortSetFromServers(add)
		updateServerSet = oldIdPort.Intersection(newIdPort)
		removeserverSet = oldIdPort.Difference(newIdPort)
		addServerSet = newIdPort.Difference(oldIdPort)
		if removeserverSet.Len() > 0 {
			rmservers := make([]interface{}, 0)
			for _, rmserver := range remove {
				rms := rmserver.(map[string]interface{})
				idPort := fmt.Sprintf("%s:%d", rms["server_id"], rms["port"])
				if removeserverSet.Contains(idPort) {
					rmsm := map[string]interface{}{
						"server_id": rms["server_id"],
						"port":      rms["port"],
						"type":      rms["type"],
						"weight":    rms["weight"],
						"server_ip": rms["server_ip"],
					}
					rmservers = append(rmservers, rmsm)
				}
			}
			request := slb.CreateRemoveVServerGroupBackendServersRequest()
			request.RegionId = client.RegionId
			request.VServerGroupId = d.Id()
			segs := len(rmservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(rmservers) {
					end = len(rmservers)
				}
				request.BackendServers = expandBackendServersWithPortToString(rmservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveVServerGroupBackendServers(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				d.SetPartial("servers")
			}
		}
		if addServerSet.Len() > 0 {
			addservers := make([]interface{}, 0)
			for _, addserver := range add {
				adds := addserver.(map[string]interface{})
				idPort := fmt.Sprintf("%s:%d", adds["server_id"], adds["port"])
				if addServerSet.Contains(idPort) {
					addsm := map[string]interface{}{
						"server_id": adds["server_id"],
						"port":      adds["port"],
						"type":      adds["type"],
						"weight":    adds["weight"],
						"server_ip": adds["server_ip"],
					}
					addservers = append(addservers, addsm)
				}
			}
			request := slb.CreateAddVServerGroupBackendServersRequest()
			request.RegionId = client.RegionId
			request.VServerGroupId = d.Id()
			segs := len(addservers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(addservers) {
					end = len(addservers)
				}

				request.BackendServers = expandBackendServersWithPortToString(addservers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddVServerGroupBackendServers(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				d.SetPartial("servers")
			}
		}
	}
	name := d.Get("name").(string)
	nameUpdate := false
	if d.HasChange("name") {
		nameUpdate = true
	}
	if d.HasChange("servers") {
		serverUpdate = true
	}
	if serverUpdate || nameUpdate {
		request := slb.CreateSetVServerGroupAttributeRequest()
		request.RegionId = client.RegionId
		request.VServerGroupId = d.Id()
		request.VServerGroupName = name
		if serverUpdate {
			servers := make([]interface{}, 0)
			for _, server := range d.Get("servers").(*schema.Set).List() {
				s := server.(map[string]interface{})
				idPort := fmt.Sprintf("%s:%d", s["server_id"], s["port"])
				if updateServerSet.Contains(idPort) {
					sm := map[string]interface{}{
						"server_id": s["server_id"],
						"port":      s["port"],
						"type":      s["type"],
						"weight":    s["weight"],
						"server_ip": s["server_ip"],
					}
					servers = append(servers, sm)
				}
			}
			segs := len(servers)/step + 1
			for i := 0; i < segs; i++ {
				start := i * step
				end := (i + 1) * step
				if end >= len(servers) {
					end = len(servers)
				}
				request.BackendServers = expandBackendServersWithPortToString(servers[start:end])
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.SetVServerGroupAttribute(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				d.SetPartial("servers")
				d.SetPartial("name")
			}
		} else {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetVServerGroupAttribute(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			d.SetPartial("name")
		}
	}
	d.Partial(false)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current VServerGroup's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the group.", lbId))
		}
	}

	request := slb.CreateDeleteVServerGroupRequest()
	request.RegionId = client.RegionId
	request.VServerGroupId = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteVServerGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RspoolVipExist"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(slbService.WaitForSlbServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
