package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSlbMasterSlaveServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbMasterSlaveServerGroupCreate,
		Read:   resourceAliyunSlbMasterSlaveServerGroupRead,
		Delete: resourceAliyunSlbMasterSlaveServerGroupDelete,
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
							ValidateFunc: validateIntegerInRange(1, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validateIntegerInRange(0, 100),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(ECS),
							ValidateFunc: validateAllowedStringValue([]string{string(ENI), string(ECS)}),
						},
						"server_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue([]string{string("Master"), string("Slave")}),
						},
						"is_backup": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateAllowedIntValue([]int{0, 1}),
						},
					},
				},
				MaxItems: 2,
				MinItems: 2,
			},
		},
	}
}

func resourceAliyunSlbMasterSlaveServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateCreateMasterSlaveServerGroupRequest()
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_master_slave_server_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*slb.CreateMasterSlaveServerGroupResponse)
	d.SetId(response.MasterSlaveServerGroupId)

	return resourceAliyunSlbMasterSlaveServerGroupRead(d, meta)
}

func resourceAliyunSlbMasterSlaveServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbMasterSlaveServerGroup(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
			"is_backup":   server.IsBackup,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunSlbMasterSlaveServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteMasterSlaveServerGroupRequest()
	request.MasterSlaveServerGroupId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteMasterSlaveServerGroup(request)
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
		if IsExceptedErrors(err, []string{MasterSlaveServerGroupNotFoundMessage, InvalidParameter}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbMasterSlaveServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
