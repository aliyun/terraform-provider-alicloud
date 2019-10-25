package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAlikafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaInstanceCreate,
		Read:   resourceAlicloudAlikafkaInstanceRead,
		Update: resourceAlicloudAlikafkaInstanceUpdate,
		Delete: resourceAlicloudAlikafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAlikafkaInstanceNameLen,
			},
			"topic_quota": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deploy_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"io_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}

	regionId := client.RegionId
	topicQuota := d.Get("topic_quota").(int)
	diskType := d.Get("disk_type").(int)
	diskSize := d.Get("disk_size").(int)
	deployType := d.Get("deploy_type").(int)
	ioMax := d.Get("io_max").(int)
	vswitchId := d.Get("vswitch_id").(string)

	// Get vswitch info by vswitchId
	vsw, err := vpcService.DescribeVSwitch(vswitchId)
	if err != nil {
		return WrapError(err)
	}

	// 1. Create post-pay order
	createOrderReq := alikafka.CreateCreatePostPayOrderRequest()
	createOrderReq.RegionId = regionId
	createOrderReq.TopicQuota = requests.NewInteger(topicQuota)
	createOrderReq.DiskType = strconv.Itoa(diskType)
	createOrderReq.DiskSize = requests.NewInteger(diskSize)
	createOrderReq.DeployType = requests.NewInteger(deployType)
	createOrderReq.IoMax = requests.NewInteger(ioMax)
	if v, ok := d.GetOk("eip_max"); ok {
		createOrderReq.EipMax = requests.NewInteger(v.(int))
	}

	var createOrderResp *alikafka.CreatePostPayOrderResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreatePostPayOrder(createOrderReq)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(createOrderReq.GetActionName(), raw, createOrderReq.RpcRequest, createOrderReq)
		v, _ := raw.(*alikafka.CreatePostPayOrderResponse)
		createOrderResp = v
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", createOrderReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAlikafkaInstanceByOrderId(createOrderResp.OrderId, 60)

	if err != nil {
		return WrapError(err)
	}

	instanceId := alikafkaInstanceVO.InstanceId
	d.SetId(instanceId)

	// 3. Start instance
	startInstanceReq := alikafka.CreateStartInstanceRequest()
	startInstanceReq.RegionId = regionId
	startInstanceReq.InstanceId = instanceId
	startInstanceReq.VpcId = vsw.VpcId
	startInstanceReq.VSwitchId = vswitchId
	startInstanceReq.ZoneId = vsw.ZoneId
	if _, ok := d.GetOk("eip_max"); ok {
		startInstanceReq.IsEipInner = requests.NewBoolean(true)
		startInstanceReq.DeployModule = "eip"
	}
	if v, ok := d.GetOk("name"); ok {
		startInstanceReq.Name = v.(string)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.StartInstance(startInstanceReq)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(startInstanceReq.GetActionName(), raw, startInstanceReq.RpcRequest, startInstanceReq)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// 3. wait until running
	err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, DefaultLongTimeout)

	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudAlikafkaInstanceRead(d, meta)
}

func resourceAlicloudAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("topic_quota", object.TopicNumLimit)
	d.Set("disk_type", object.DiskType)
	d.Set("disk_size", object.DiskSize)
	d.Set("deploy_type", object.DeployType)
	d.Set("io_max", object.IoMax)
	d.Set("eip_max", object.EipMax)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("zone_id", object.ZoneId)

	return nil
}

func resourceAlicloudAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	// Process change instance name.
	if d.HasChange("name") {
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		modifyInstanceNameReq := alikafka.CreateModifyInstanceNameRequest()
		modifyInstanceNameReq.RegionId = client.RegionId
		modifyInstanceNameReq.InstanceId = d.Id()
		modifyInstanceNameReq.InstanceName = name

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ModifyInstanceName(modifyInstanceNameReq)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyInstanceNameReq.GetActionName(), raw, modifyInstanceNameReq.RpcRequest, modifyInstanceNameReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceNameReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	attributeUpdate := false

	upgradeReq := alikafka.CreateUpgradePostPayOrderRequest()
	upgradeReq.RegionId = client.RegionId
	upgradeReq.InstanceId = d.Id()
	upgradeReq.TopicQuota = requests.NewInteger(d.Get("topic_quota").(int))
	upgradeReq.DiskSize = requests.NewInteger(d.Get("disk_size").(int))
	upgradeReq.IoMax = requests.NewInteger(d.Get("io_max").(int))

	if d.HasChange("topic_quota") || d.HasChange("disk_size") || d.HasChange("io_max") {
		attributeUpdate = true
	}
	eipMax := 0
	if v, ok := d.GetOk("eip_max"); ok {
		eipMax = v.(int)
	}
	if d.HasChange("eip_max") {

		if v, ok := d.GetOk("eip_max"); ok {
			eipMax = v.(int)
		}
		upgradeReq.EipMax = requests.NewInteger(eipMax)
		attributeUpdate = true
	}

	if attributeUpdate {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.UpgradePostPayOrder(upgradeReq)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(upgradeReq.GetActionName(), raw, upgradeReq.RpcRequest, upgradeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), upgradeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	err := alikafkaService.WaitForAlikafkaInstanceUpdated(d.Id(), d.Get("topic_quota").(int),
		d.Get("disk_size").(int), d.Get("io_max").(int), eipMax, DefaultTimeoutMedium)

	if err != nil {
		return WrapError(err)
	}

	err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, 2000)

	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudAlikafkaInstanceRead(d, meta)
}

func resourceAlicloudAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId
	request.ReleaseIgnoreTime = requests.NewBoolean(true)
	request.ForceDeleteInstance = requests.NewBoolean(true)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.ReleaseInstance(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAllAlikafkaNodeRelease(d.Id(), DefaultTimeoutMedium))
}
