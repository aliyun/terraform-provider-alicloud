package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsNetworkInterfaceAttachmentCreate,
		Read:   resourceAliCloudEcsNetworkInterfaceAttachmentRead,
		Update: resourceAliCloudEcsNetworkInterfaceAttachmentUpdate,
		Delete: resourceAliCloudEcsNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trunk_network_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_card_index": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"wait_for_network_configuration_ready": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEcsNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AttachNetworkInterface"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["NetworkInterfaceId"] = d.Get("network_interface_id")
	request["InstanceId"] = d.Get("instance_id")

	if v, ok := d.GetOk("trunk_network_instance_id"); ok {
		request["TrunkNetworkInstanceId"] = v
	}

	if v, ok := d.GetOkExists("network_card_index"); ok {
		request["NetworkCardIndex"] = v
	}

	if v, ok := d.GetOkExists("wait_for_network_configuration_ready"); ok {
		request["WaitForNetworkConfigurationReady"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_network_interface_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["NetworkInterfaceId"], request["InstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(fmt.Sprint(request["NetworkInterfaceId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsNetworkInterfaceAttachmentRead(d, meta)
}

func resourceAliCloudEcsNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_network_interface_attachment ecsService.DescribeNetworkInterfaceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("network_interface_id", object["NetworkInterfaceId"])
	d.Set("instance_id", object["InstanceId"])

	if attachment, ok := object["Attachment"]; ok {
		attachmentArg := attachment.(map[string]interface{})

		if fmt.Sprint(object["Type"]) == "Member" || fmt.Sprint(object["Type"]) == "slave" {
			if instanceId, ok := attachmentArg["InstanceId"]; ok {
				d.Set("instance_id", instanceId)
			}
		}

		if trunkNetworkInterfaceId, ok := attachmentArg["TrunkNetworkInterfaceId"]; ok {
			d.Set("trunk_network_instance_id", trunkNetworkInterfaceId)
		}

		if networkCardIndex, ok := attachmentArg["NetworkCardIndex"]; ok {
			d.Set("network_card_index", networkCardIndex)
		}
	}

	return nil
}

func resourceAliCloudEcsNetworkInterfaceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAliCloudEcsNetworkInterfaceAttachmentRead(d, meta)
}

func resourceAliCloudEcsNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	action := "DetachNetworkInterface"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"NetworkInterfaceId": parts[0],
		"InstanceId":         parts[1],
	}

	if v, ok := d.GetOk("trunk_network_instance_id"); ok {
		request["TrunkNetworkInstanceId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidEcsId.NotFound", "InvalidEniId.NotFound", "InvalidSecurityGroupId.NotFound", "InvalidVSwitchId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
