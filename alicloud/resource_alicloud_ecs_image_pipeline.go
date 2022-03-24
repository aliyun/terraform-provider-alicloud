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

func resourceAlicloudEcsImagePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsImagePipelineCreate,
		Read:   resourceAlicloudEcsImagePipelineRead,
		Delete: resourceAlicloudEcsImagePipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"add_account": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 20,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"base_image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"base_image_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMAGE", "IMAGE_FAMILY"}, false),
			},
			"build_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delete_instance_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"image_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9:._-]{1,63}$`), "It must be `2` to `64` characters in length, The description must start with a letter, and can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).")),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9:._-]{1,127}$`), "It must be `2` to `128` characters in length, The description must start with a letter, and can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).")),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(20, 500),
				ForceNew:     true,
				Computed:     true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"to_region_id": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 20,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsImagePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateImagePipeline"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("add_account"); ok {
		request["AddAccount"] = v.(*schema.Set).List()
	}
	request["BaseImage"] = d.Get("base_image")
	request["BaseImageType"] = d.Get("base_image_type")
	if v, ok := d.GetOk("build_content"); ok {
		request["BuildContent"] = v
	}
	if v, ok := d.GetOkExists("delete_instance_on_failure"); ok {
		request["DeleteInstanceOnFailure"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("image_name"); ok {
		request["ImageName"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("system_disk_size"); ok {
		request["SystemDiskSize"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("to_region_id"); ok {
		request["ToRegionId"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateImagePipeline")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_image_pipeline", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ImagePipelineId"]))

	return resourceAlicloudEcsImagePipelineRead(d, meta)
}
func resourceAlicloudEcsImagePipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsImagePipeline(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_image_pipeline ecsService.DescribeEcsImagePipeline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("base_image", object["BaseImage"])
	d.Set("base_image_type", object["BaseImageType"])
	d.Set("delete_instance_on_failure", object["DeleteInstanceOnFailure"])
	d.Set("description", object["Description"])
	d.Set("image_name", object["ImageName"])
	d.Set("instance_type", object["InstanceType"])
	if v, ok := object["InternetMaxBandwidthOut"]; ok && fmt.Sprint(v) != "0" {
		d.Set("internet_max_bandwidth_out", formatInt(v))
	}
	d.Set("name", object["Name"])
	if v, ok := object["SystemDiskSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("system_disk_size", formatInt(v))
	}
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	addAccountsList := make([]string, 0)
	if addAccounts, ok := object["AddAccounts"]; ok {
		addAccountsMap := addAccounts.(map[string]interface{})
		if addAccount, ok := addAccountsMap["AddAccount"]; ok {
			for _, item := range addAccount.([]interface{}) {
				addAccountsList = append(addAccountsList, item.(string))
			}
		}
	}
	d.Set("add_account", addAccountsList)

	toRegionIdList := make([]string, 0)
	if toRegionIds, ok := object["ToRegionIds"]; ok {
		toRegionIdsMap := toRegionIds.(map[string]interface{})
		if toRegionId, ok := toRegionIdsMap["ToRegionId"]; ok {
			for _, item := range toRegionId.([]interface{}) {
				toRegionIdList = append(toRegionIdList, item.(string))
			}
		}
	}
	d.Set("to_region_id", toRegionIdList)
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	return nil
}
func resourceAlicloudEcsImagePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteImagePipeline"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ImagePipelineId": d.Id(),
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
