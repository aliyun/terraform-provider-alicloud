package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Optional: true,
				ForceNew: true,
			},
			"base_image_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMAGE", "IMAGE_FAMILY", "OSS"}, false),
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
			// lintignore: S006
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
			"test_content": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_family": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"repair_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nvme_support": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"advanced_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retain_cloud_assistant": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"image_name_suffix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"import_image_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"boot_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"license_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"retain_imported_image": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"features": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nvme_support": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 17,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oss_bucket": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"oss_object": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"format": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"disk_image_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// expandImagePipelineCreateRequest builds the CreateImagePipeline RPC request map from the
// ResourceData schema. RegionId and ClientToken are intentionally excluded (they depend on the
// caller/client) and set by Create. Extracted to enable no-cloud unit testing of the full
// parameter mapping and disk_device_mappings ordering without invoking the SDK or a real API.
func expandImagePipelineCreateRequest(d *schema.ResourceData, regionId string) map[string]interface{} {
	request := make(map[string]interface{})
	if v, ok := d.GetOk("add_account"); ok {
		request["AddAccount"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("base_image"); ok {
		request["BaseImage"] = v
	}
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
	request["RegionId"] = regionId
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
	if v, ok := d.GetOk("test_content"); ok {
		request["TestContent"] = v
	}
	if v, ok := d.GetOk("image_family"); ok {
		request["ImageFamily"] = v
	}
	if v, ok := d.GetOk("repair_mode"); ok {
		request["RepairMode"] = v
	}
	if v, ok := d.GetOk("nvme_support"); ok {
		request["NvmeSupport"] = v
	}
	if v, ok := d.GetOk("advanced_options"); ok && len(v.([]interface{})) > 0 {
		if retain, ok := d.GetOkExists("advanced_options.0.retain_cloud_assistant"); ok {
			request["AdvancedOptions.RetainCloudAssistant"] = retain
		}
		if suffix, ok := d.GetOkExists("advanced_options.0.image_name_suffix"); ok {
			request["AdvancedOptions.ImageNameSuffix"] = suffix
		}
	}
	if v, ok := d.GetOk("import_image_options"); ok && len(v.([]interface{})) > 0 {
		if architecture, ok := d.GetOk("import_image_options.0.architecture"); ok {
			request["ImportImageOptions.Architecture"] = architecture
		}
		if platform, ok := d.GetOk("import_image_options.0.platform"); ok {
			request["ImportImageOptions.Platform"] = platform
		}
		if bootMode, ok := d.GetOk("import_image_options.0.boot_mode"); ok {
			request["ImportImageOptions.BootMode"] = bootMode
		}
		if licenseType, ok := d.GetOk("import_image_options.0.license_type"); ok {
			request["ImportImageOptions.LicenseType"] = licenseType
		}
		if osType, ok := d.GetOk("import_image_options.0.os_type"); ok {
			request["ImportImageOptions.OSType"] = osType
		}
		if retain, ok := d.GetOkExists("import_image_options.0.retain_imported_image"); ok {
			request["ImportImageOptions.RetainImportedImage"] = retain
		}
		if features, ok := d.GetOk("import_image_options.0.features"); ok && len(features.([]interface{})) > 0 {
			if nvme, ok := d.GetOk("import_image_options.0.features.0.nvme_support"); ok {
				request["ImportImageOptions.Features.NvmeSupport"] = nvme
			}
		}
		if mappings, ok := d.GetOk("import_image_options.0.disk_device_mappings"); ok {
			for i, raw := range mappings.([]interface{}) {
				mapping := raw.(map[string]interface{})
				index := i + 1
				if v, ok := mapping["oss_bucket"]; ok && fmt.Sprint(v) != "" {
					request[fmt.Sprintf("ImportImageOptions.DiskDeviceMappings.%d.OSSBucket", index)] = v
				}
				if v, ok := mapping["oss_object"]; ok && fmt.Sprint(v) != "" {
					request[fmt.Sprintf("ImportImageOptions.DiskDeviceMappings.%d.OSSObject", index)] = v
				}
				if v, ok := mapping["format"]; ok && fmt.Sprint(v) != "" {
					request[fmt.Sprintf("ImportImageOptions.DiskDeviceMappings.%d.Format", index)] = v
				}
				if v, ok := mapping["disk_image_size"]; ok && fmt.Sprint(v) != "0" {
					request[fmt.Sprintf("ImportImageOptions.DiskDeviceMappings.%d.DiskImageSize", index)] = v
				}
			}
		}
	}
	return request
}

func resourceAlicloudEcsImagePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateImagePipeline"
	request := expandImagePipelineCreateRequest(d, client.RegionId)
	var err error
	request["ClientToken"] = buildClientToken("CreateImagePipeline")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
	return flattenImagePipelineRead(d, object)
}

// flattenImagePipelineRead writes the DescribeImagePipelines API response into the ResourceData
// schema, including the nested AdvancedOptions, ImportImageOptions.Features and
// ImportImageOptions.DiskDeviceMappings blocks. Extracted from Read so the flatten/round-trip
// (the same path the Importer uses via ImportStatePassthrough → Read) is unit-testable without
// invoking the SDK or a real API.
func flattenImagePipelineRead(d *schema.ResourceData, object map[string]interface{}) error {
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
	d.Set("test_content", object["TestContent"])
	d.Set("image_family", object["ImageFamily"])
	d.Set("repair_mode", object["RepairMode"])
	d.Set("nvme_support", object["NvmeSupport"])

	advancedOptionsList := make([]map[string]interface{}, 0)
	if advancedOptions, ok := object["AdvancedOptions"].(map[string]interface{}); ok {
		m := map[string]interface{}{}
		if v, ok := advancedOptions["RetainCloudAssistant"]; ok {
			m["retain_cloud_assistant"] = v
		}
		if v, ok := advancedOptions["ImageNameSuffix"]; ok {
			m["image_name_suffix"] = v
		}
		advancedOptionsList = append(advancedOptionsList, m)
	}
	d.Set("advanced_options", advancedOptionsList)

	importImageOptionsList := make([]map[string]interface{}, 0)
	if importImageOptions, ok := object["ImportImageOptions"].(map[string]interface{}); ok {
		m := map[string]interface{}{
			"architecture":          importImageOptions["Architecture"],
			"platform":              importImageOptions["Platform"],
			"boot_mode":             importImageOptions["BootMode"],
			"license_type":          importImageOptions["LicenseType"],
			"os_type":               importImageOptions["OSType"],
			"retain_imported_image": importImageOptions["RetainImportedImage"],
		}
		featuresList := make([]map[string]interface{}, 0)
		if features, ok := importImageOptions["Features"].(map[string]interface{}); ok {
			fm := map[string]interface{}{}
			if v, ok := features["NvmeSupport"]; ok {
				fm["nvme_support"] = v
			}
			featuresList = append(featuresList, fm)
		}
		m["features"] = featuresList
		diskDeviceMappingsList := make([]map[string]interface{}, 0)
		if diskDeviceMappings, ok := importImageOptions["DiskDeviceMappings"].(map[string]interface{}); ok {
			if diskDeviceMapping, ok := diskDeviceMappings["DiskDeviceMapping"].([]interface{}); ok {
				for _, item := range diskDeviceMapping {
					mapping, ok := item.(map[string]interface{})
					if !ok {
						continue
					}
					dm := map[string]interface{}{
						"oss_bucket": mapping["OSSBucket"],
						"oss_object": mapping["OSSObject"],
						"format":     mapping["Format"],
					}
					if v, ok := mapping["DiskImageSize"]; ok && fmt.Sprint(v) != "0" {
						dm["disk_image_size"] = formatInt(v)
					}
					diskDeviceMappingsList = append(diskDeviceMappingsList, dm)
				}
			}
		}
		m["disk_device_mappings"] = diskDeviceMappingsList
		importImageOptionsList = append(importImageOptionsList, m)
	}
	d.Set("import_image_options", importImageOptionsList)
	return nil
}
func resourceAlicloudEcsImagePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteImagePipeline"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ImagePipelineId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
