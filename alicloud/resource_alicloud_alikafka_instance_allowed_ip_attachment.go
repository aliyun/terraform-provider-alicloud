package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAliKafkaInstanceAllowedIpAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAliKafkaInstanceAllowedIpAttachmentCreate,
		Read:   resourceAlicloudAliKafkaInstanceAllowedIpAttachmentRead,
		Delete: resourceAlicloudAliKafkaInstanceAllowedIpAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allowed_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allowed_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"vpc"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_range": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"9092/9092"}, false),
			},
		},
	}
}

func resourceAlicloudAliKafkaInstanceAllowedIpAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UpdateAllowedIp"
	request := make(map[string]interface{})
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	request["AllowedListIp"] = d.Get("allowed_ip")
	request["AllowedListType"] = d.Get("allowed_type")
	request["InstanceId"] = d.Get("instance_id")
	request["PortRange"] = d.Get("port_range")
	request["RegionId"] = client.RegionId
	request["UpdateType"] = "add"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance_allowed_ip_attachment", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["AllowedListType"], ":", request["PortRange"], ":", request["AllowedListIp"]))

	return resourceAlicloudAliKafkaInstanceAllowedIpAttachmentRead(d, meta)
}
func resourceAlicloudAliKafkaInstanceAllowedIpAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	_, err := alikafkaService.DescribeAliKafkaInstanceAllowedIpAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_instance_allowed_ip_attachment alikafkaService.DescribeAliKafkaInstanceAllowedIpAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("allowed_ip", parts[3])
	d.Set("allowed_type", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("port_range", parts[2])
	return nil
}
func resourceAlicloudAliKafkaInstanceAllowedIpAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateAllowedIp"
	var response map[string]interface{}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AllowedListIp":   parts[3],
		"AllowedListType": parts[1],
		"InstanceId":      parts[0],
		"PortRange":       parts[2],
	}

	request["RegionId"] = client.RegionId
	request["UpdateType"] = "delete"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
