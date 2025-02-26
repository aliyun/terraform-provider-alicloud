package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFcLayerVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFcLayerVersionCreate,
		Read:   resourceAlicloudFcLayerVersionRead,
		Update: resourceAlicloudFcLayerVersionUpdate,
		Delete: resourceAlicloudFcLayerVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"arn": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"code_check_sum": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"oss_bucket_name": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"zip_file"},
			},
			"oss_object_name": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"zip_file"},
			},
			"zip_file": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"oss_object_name", "oss_bucket_name"},
			},
			"compatible_runtime": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"skip_destroy": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"description": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"layer_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"version": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudFcLayerVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	layerName := d.Get("layer_name")
	action := fmt.Sprintf("/2021-04-06/layers/%s/versions", layerName)

	body := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}

	compatibleRuntimes := make([]string, 0)
	for _, v := range d.Get("compatible_runtime").(*schema.Set).List() {
		compatibleRuntimes = append(compatibleRuntimes, v.(string))
	}
	body["compatibleRuntime"] = compatibleRuntimes

	codeMaps := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("zip_file"); ok && fmt.Sprint(v) != "" && v.(string) != "" {
		file, err := loadFileContent(v.(string))
		if err != nil {
			return WrapError(err)
		}
		codeMaps["zipFile"] = file
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok && fmt.Sprint(v) != "" {
		codeMaps["ossBucketName"] = v
	}
	if v, ok := d.GetOk("oss_object_name"); ok && fmt.Sprint(v) != "" {
		codeMaps["ossObjectName"] = v
	}
	body["Code"] = codeMaps

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("FC-Open", "2021-04-06", action, nil, nil, body, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, body)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_layer_version", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(layerName, ":", response["version"]))

	return resourceAlicloudFcLayerVersionRead(d, meta)
}

func resourceAlicloudFcLayerVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcOpenService := FcOpenService{client}

	object, err := fcOpenService.DescribeFcLayerVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fc_layer_version fcOpenService.DescribeFcLayerVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("layer_name", parts[0])
	d.Set("version", parts[1])
	d.Set("acl", object["acl"])
	d.Set("arn", object["arn"])
	d.Set("code_check_sum", object["codeChecksum"])
	d.Set("compatible_runtime", object["compatibleRuntime"])
	d.Set("description", object["description"])
	return nil
}

func resourceAlicloudFcLayerVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] The property in this resource cannot be updated.")
	return nil
}

func resourceAlicloudFcLayerVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if v, ok := d.GetOkExists("skip_destroy"); ok && v.(bool) {
		log.Printf("[INFO] Setting `skip_destroy` to `true` means that the Alicloud Provider will not destroy any layer version, even when running `terraform destroy`. Layer versions are thus intentional dangling resources that are not managed by Terraform and may incur extra expense in your Alicloud account.")
		return nil
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	var response map[string]interface{}
	action := fmt.Sprintf("/2021-04-06/layers/%s/versions/%s", parts[0], parts[1])
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RoaDelete("FC-Open", "2021-04-06", action, nil, nil, nil, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
