package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunHaVipAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcHaVipAttachmentCreate,
		Read:   resourceAlicloudVpcHaVipAttachmentRead,
		Update: resourceAlicloudVpcHaVipAttachmentUpdate,
		Delete: resourceAlicloudVpcHaVipAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"havip_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"force": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"instance_type": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudVpcHaVipAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("havip_id"); ok {
		request["HaVipId"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	var response map[string]interface{}
	action := "AssociateHaVip"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("AssociateHaVip")
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "IncorrectHaVipStatus", "InvalidVip.Status", "OperationConflict", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "IncorrectStatus.%s", "SystemBusy", "ServiceUnavailable", "IncorrectInstanceStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_havip_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["HaVipId"], ":", request["InstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcHaVipAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudVpcHaVipAttachmentRead(d, meta)
}

func resourceAlicloudVpcHaVipAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpcHaVipAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ha_vip_attachment vpcService.DescribeVpcHaVipAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("havip_id", parts[0])
	d.Set("instance_id", parts[1])

	associatedInstanceType := object["AssociatedInstanceType"]
	d.Set("instance_type", associatedInstanceType)

	status := object["Status"]
	d.Set("status", status)

	return nil
}

func resourceAlicloudVpcHaVipAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudVpcHaVipAttachmentRead(d, meta)
}

func resourceAlicloudVpcHaVipAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"HaVipId":    parts[0],
		"InstanceId": parts[1],
		"RegionId":   client.RegionId,
	}

	if v, ok := d.GetOk("force"); ok {
		request["Force"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	action := "UnassociateHaVip"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("UnassociateHaVip")
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectHaVipStatus", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "IncorrectStatus.%s", "SystemBusy", "ServiceUnavailable", "IncorrectInstanceStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
