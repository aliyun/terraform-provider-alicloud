package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsInstanceCreate,
		Read:   resourceAlicloudOnsInstanceRead,
		Update: resourceAlicloudOnsInstanceUpdate,
		Delete: resourceAlicloudOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.97.0. Use 'instance_name' instead.",
				ConflictsWith: []string{"instance_name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ons.CreateOnsInstanceCreateRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.InstanceName = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "instance_name" must be set one!`))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceCreate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ons.OnsInstanceCreateResponse)
		d.SetId(fmt.Sprintf("%v", response.Data.InstanceId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudOnsInstanceUpdate(d, meta)
}
func resourceAlicloudOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_instance onsService.DescribeOnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_name", object.InstanceName)
	d.Set("name", object.InstanceName)
	d.Set("instance_type", object.InstanceType)
	d.Set("release_time", time.Unix(int64(object.ReleaseTime)/1000, 0).Format("2006-01-02 03:04:05"))
	d.Set("remark", object.Remark)
	d.Set("status", object.InstanceStatus)
	d.Set("instance_status", object.InstanceStatus)

	listTagResourcesObject, err := onsService.ListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tags := make(map[string]string)
	for _, t := range listTagResourcesObject.TagResources {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := ons.CreateOnsInstanceUpdateRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request.InstanceName = d.Get("instance_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request.InstanceName = d.Get("name").(string)
	}
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
		request.Remark = d.Get("remark").(string)
	}
	if update {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceUpdate(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("instance_name")
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudOnsInstanceRead(d, meta)
}
func resourceAlicloudOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ons.CreateOnsInstanceDeleteRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		args := *request
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceDelete(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"INSTANCE_NOT_EMPTY", "Throttling.User"}) {
				wait()
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
	return nil
}
