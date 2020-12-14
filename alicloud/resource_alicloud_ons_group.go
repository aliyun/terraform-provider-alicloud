package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsGroupCreate,
		Read:   resourceAlicloudOnsGroupRead,
		Update: resourceAlicloudOnsGroupUpdate,
		Delete: resourceAlicloudOnsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateOnsGroupId,
				ConflictsWith: []string{"group_id"},
			},
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateOnsGroupId,
				Deprecated:    "Field 'group_id' has been deprecated from version 1.98.0. Use 'group_name' instead.",
				ConflictsWith: []string{"group_name"},
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "tcp"}, false),
				Default:      "tcp",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"read_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ons.CreateOnsGroupCreateRequest()
	if v, ok := d.GetOk("group_name"); ok {
		request.GroupId = v.(string)
	} else if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "group_id" or "group_name" must be set one!`))
	}

	if v, ok := d.GetOk("group_type"); ok {
		request.GroupType = v.(string)
	}

	request.InstanceId = d.Get("instance_id").(string)
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupCreate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(fmt.Sprintf("%v:%v", request.InstanceId, request.GroupId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudOnsGroupUpdate(d, meta)
}
func resourceAlicloudOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_group onsService.DescribeOnsGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("group_name", parts[1])
	d.Set("group_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("group_type", object.GroupType)
	d.Set("remark", object.Remark)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		tags[t.Key] = t.Value
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTagsForGroup(d, "GROUP"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("read_enable") {
		request := ons.CreateOnsGroupConsumerUpdateRequest()
		request.GroupId = parts[1]
		request.InstanceId = parts[0]
		request.ReadEnable = requests.NewBoolean(d.Get("read_enable").(bool))
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupConsumerUpdate(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("read_enable")
	}
	d.Partial(false)
	return resourceAlicloudOnsGroupRead(d, meta)
}
func resourceAlicloudOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := ons.CreateOnsGroupDeleteRequest()
	request.GroupId = parts[1]
	request.InstanceId = parts[0]
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupDelete(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
