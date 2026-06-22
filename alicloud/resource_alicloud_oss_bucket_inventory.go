// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketInventory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketInventoryCreate,
		Read:   resourceAliCloudOssBucketInventoryRead,
		Update: resourceAliCloudOssBucketInventoryUpdate,
		Delete: resourceAliCloudOssBucketInventoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inventory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"included_object_versions": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"All", "Current"}, false),
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Daily", "Weekly"}, false),
						},
					},
				},
			},
			"optional_fields": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oss_bucket_destination": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: StringInSlice([]string{"CSV", "ORC", "Parquet"}, false),
									},
									"account_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
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

func buildOssBucketInventoryConfiguration(d *schema.ResourceData) map[string]interface{} {
	inventory := make(map[string]interface{})
	inventory["Id"] = d.Get("inventory_id")
	inventory["IsEnabled"] = d.Get("is_enabled")
	inventory["IncludedObjectVersions"] = d.Get("included_object_versions")

	if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		scheduleMap := v.([]interface{})[0].(map[string]interface{})
		inventory["Schedule"] = map[string]interface{}{"Frequency": scheduleMap["frequency"]}
	}

	if v, ok := d.GetOk("optional_fields"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		optionalFieldsMap := v.([]interface{})[0].(map[string]interface{})
		if fields, ok := optionalFieldsMap["field"].([]interface{}); ok && len(fields) > 0 {
			inventory["OptionalFields"] = map[string]interface{}{"Field": fields}
		}
	}

	if v, ok := d.GetOk("filter"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		filterMap := v.([]interface{})[0].(map[string]interface{})
		filter := make(map[string]interface{})
		if val, ok := filterMap["prefix"]; ok && val.(string) != "" {
			filter["Prefix"] = val
		}
		if len(filter) > 0 {
			inventory["Filter"] = filter
		}
	}

	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		destinationMap := v.([]interface{})[0].(map[string]interface{})
		if dest, ok := destinationMap["oss_bucket_destination"].([]interface{}); ok && len(dest) > 0 && dest[0] != nil {
			destMap := dest[0].(map[string]interface{})
			ossBucketDestination := make(map[string]interface{})
			ossBucketDestination["Format"] = destMap["format"]
			ossBucketDestination["AccountId"] = destMap["account_id"]
			ossBucketDestination["RoleArn"] = destMap["role_arn"]
			ossBucketDestination["Bucket"] = destMap["bucket"]
			if prefix, ok := destMap["prefix"]; ok && prefix.(string) != "" {
				ossBucketDestination["Prefix"] = prefix
			}
			inventory["Destination"] = map[string]interface{}{"OSSBucketDestination": ossBucketDestination}
		}
	}

	return inventory
}

func resourceAliCloudOssBucketInventoryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?inventory")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	query["inventoryId"] = StringPointer(d.Get("inventory_id").(string))

	request["InventoryConfiguration"] = buildOssBucketInventoryConfiguration(d)

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketInventory", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil && !IsExpectedErrors(err, []string{"InventoryConfigurationAlreadyExists"}) {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_inventory", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", *hostMap["bucket"], COLON_SEPARATED, d.Get("inventory_id").(string)))

	return resourceAliCloudOssBucketInventoryRead(d, meta)
}

func resourceAliCloudOssBucketInventoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketInventory(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_inventory DescribeOssBucketInventory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	d.Set("bucket", parts[0])
	d.Set("inventory_id", objectRaw["Id"])
	d.Set("is_enabled", formatBool(objectRaw["IsEnabled"]))
	d.Set("included_object_versions", objectRaw["IncludedObjectVersions"])

	if scheduleRaw, ok := objectRaw["Schedule"].(map[string]interface{}); ok {
		d.Set("schedule", []map[string]interface{}{{"frequency": scheduleRaw["Frequency"]}})
	}

	if optionalFieldsRaw, ok := objectRaw["OptionalFields"].(map[string]interface{}); ok {
		fields := make([]interface{}, 0)
		if fieldRaw := optionalFieldsRaw["Field"]; fieldRaw != nil {
			switch f := fieldRaw.(type) {
			case []interface{}:
				fields = f
			default:
				fields = append(fields, f)
			}
		}
		d.Set("optional_fields", []map[string]interface{}{{"field": fields}})
	}

	if filterRaw, ok := objectRaw["Filter"].(map[string]interface{}); ok {
		filterMap := map[string]interface{}{"prefix": filterRaw["Prefix"]}
		d.Set("filter", []map[string]interface{}{filterMap})
	}

	if destinationRaw, ok := objectRaw["Destination"].(map[string]interface{}); ok {
		if ossDestRaw, ok := destinationRaw["OSSBucketDestination"].(map[string]interface{}); ok {
			destMap := map[string]interface{}{
				"format":     ossDestRaw["Format"],
				"account_id": ossDestRaw["AccountId"],
				"role_arn":   ossDestRaw["RoleArn"],
				"bucket":     ossDestRaw["Bucket"],
				"prefix":     ossDestRaw["Prefix"],
			}
			d.Set("destination", []map[string]interface{}{{"oss_bucket_destination": []map[string]interface{}{destMap}}})
		}
	}

	return nil
}

func resourceAliCloudOssBucketInventoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	action := fmt.Sprintf("/?inventory")
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(parts[0])
	query["inventoryId"] = StringPointer(parts[1])

	// PutBucketInventory rejects an existing configuration with InventoryConfigurationAlreadyExists,
	// so an update must remove the old rule before writing the new one.
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", action), query, nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil && !IsExpectedErrors(err, []string{"NoSuchInventory", "NoSuchBucket"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	request := make(map[string]interface{})
	request["InventoryConfiguration"] = buildOssBucketInventoryConfiguration(d)
	body = request

	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketInventory", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil && !IsExpectedErrors(err, []string{"InventoryConfigurationAlreadyExists"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudOssBucketInventoryRead(d, meta)
}

func resourceAliCloudOssBucketInventoryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	action := fmt.Sprintf("/?inventory")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(parts[0])
	query["inventoryId"] = StringPointer(parts[1])

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NoSuchInventory", "NoSuchBucket"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
