package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcsImageComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsImageComponentCreate,
		Read:   resourceAlicloudEcsImageComponentRead,
		Delete: resourceAlicloudEcsImageComponentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"component_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Build"}, false),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"image_component_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.:-]{1,127}$`), "It must be `2` to `128` characters in length, The description must start with a letter, and can contain letters, digits, underscores (_), and hyphens (-).")),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Linux"}, false),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsImageComponentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateImageComponent"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("component_type"); ok {
		request["ComponentType"] = v
	}
	request["Content"] = d.Get("content")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("image_component_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("system_type"); ok {
		request["SystemType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	request["ClientToken"] = buildClientToken("CreateImageComponent")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_image_component", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ImageComponentId"]))

	return resourceAlicloudEcsImageComponentRead(d, meta)
}
func resourceAlicloudEcsImageComponentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsImageComponent(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_image_component ecsService.DescribeEcsImageComponent Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("component_type", object["ComponentType"])
	d.Set("content", object["Content"])
	d.Set("description", object["Description"])
	d.Set("image_component_name", object["Name"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("system_type", object["SystemType"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	return nil
}
func resourceAlicloudEcsImageComponentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteImageComponent"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ImageComponentId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
