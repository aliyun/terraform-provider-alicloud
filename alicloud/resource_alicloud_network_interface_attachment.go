package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNetworkInterfaceAttachmentCreate,
		Read:   resourceAliyunNetworkInterfaceAttachmentRead,
		Delete: resourceAliyunNetworkInterfaceAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	eniId := d.Get("network_interface_id").(string)
	instanceId := d.Get("instance_id").(string)

	request := ecs.CreateAttachNetworkInterfaceRequest()
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId

	err := resource.Retry(DefaultTimeout*time.Minute, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachNetworkInterface(request)
		})
		if err != nil {
			if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_netWork_interface_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(eniId + COLON_SEPARATED + instanceId)
	if err = ecsService.WaitForNetworkInterface(eniId, InUse, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunNetworkInterfaceAttachmentRead(d, meta)
}

func resourceAliyunNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("network_interface_id", object.NetworkInterfaceId)

	return nil
}

func resourceAliyunNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]

	_, err = ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request := ecs.CreateDetachNetworkInterfaceRequest()
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId

	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachNetworkInterface(request)
		})
		if err != nil {
			if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForNetworkInterface(eniId, Available, DefaultTimeoutMedium))
}
