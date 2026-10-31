package slb

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMasterSlaveServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceMasterSlaveServerGroupCreate,
		Read:   resourceMasterSlaveServerGroupRead,
		Update: resourceMasterSlaveServerGroupUpdate,
		Delete: resourceMasterSlaveServerGroupDelete,
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
				Required: true,
				ForceNew: true,
			},

			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(ECS),
							ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
						},
						"server_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Master", "Slave"}, false),
						},
					},
				},
				MaxItems: 2,
				MinItems: 2,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceMasterSlaveServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateCreateMasterSlaveServerGroupRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := d.GetOk("name"); ok {
		request.MasterSlaveServerGroupName = v.(string)
	}
	if v, ok := d.GetOk("servers"); ok {
		request.MasterSlaveBackendServers = expandMasterSlaveBackendServersToString(v.(*schema.Set).List())
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateMasterSlaveServerGroup(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_master_slave_server_group", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateMasterSlaveServerGroupResponse)
	d.SetId(response.MasterSlaveServerGroupId)

	return resourceMasterSlaveServerGroupRead(d, meta)
}

func resourceMasterSlaveServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	object, err := slbService.DescribeSlbMasterSlaveServerGroup(d.Id())

	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", object.MasterSlaveServerGroupName)
	d.Set("load_balancer_id", object.LoadBalancerId)

	servers := make([]map[string]interface{}, 0)

	for _, server := range object.MasterSlaveBackendServers.MasterSlaveBackendServer {
		s := map[string]interface{}{
			"server_id":   server.ServerId,
			"port":        server.Port,
			"weight":      server.Weight,
			"type":        server.Type,
			"server_type": server.ServerType,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func resourceMasterSlaveServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMasterSlaveServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return errs.WrapError(fmt.Errorf("Current master-slave server group's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteMasterSlaveServerGroupRequest()
	request.RegionId = client.RegionId
	request.MasterSlaveServerGroupId = d.Id()

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteMasterSlaveServerGroup(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"RspoolVipExist"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return errs.WrapError(slbService.WaitForSlbMasterSlaveServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
