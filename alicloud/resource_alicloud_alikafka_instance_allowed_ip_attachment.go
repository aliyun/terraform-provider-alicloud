package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAliKafkaInstanceAllowedIpAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAliKafkaInstanceAllowedIpAttachmentCreate,
		Read:   resourceAliCloudAliKafkaInstanceAllowedIpAttachmentRead,
		Delete: resourceAliCloudAliKafkaInstanceAllowedIpAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allowed_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"vpc", "internet"}, false),
			},
			"port_range": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"9092/9092", "9093/9093", "9094/9094", "9095/9095"}, false),
			},
			"allowed_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAliKafkaInstanceAllowedIpAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UpdateAllowedIp"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["UpdateType"] = "add"
	request["InstanceId"] = d.Get("instance_id")
	request["AllowedListType"] = d.Get("allowed_type")
	request["PortRange"] = d.Get("port_range")
	request["AllowedListIp"] = d.Get("allowed_ip")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
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

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["InstanceId"], request["AllowedListType"], request["PortRange"], request["AllowedListIp"]))

	return resourceAliCloudAliKafkaInstanceAllowedIpAttachmentRead(d, meta)
}

func resourceAliCloudAliKafkaInstanceAllowedIpAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAliKafkaInstanceAllowedIpAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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

	d.Set("instance_id", parts[0])
	d.Set("allowed_type", parts[1])
	d.Set("port_range", object["PortRange"])
	d.Set("allowed_ip", parts[3])

	return nil
}

func resourceAliCloudAliKafkaInstanceAllowedIpAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateAllowedIp"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"UpdateType":      "delete",
		"InstanceId":      parts[0],
		"AllowedListType": parts[1],
		"PortRange":       parts[2],
		"AllowedListIp":   parts[3],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
