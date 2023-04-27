package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEbsEnterpriseSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEbsEnterpriseSnapshotPolicyCreate,
		Read:   resourceAlicloudEbsEnterpriseSnapshotPolicyRead,
		Update: resourceAlicloudEbsEnterpriseSnapshotPolicyUpdate,
		Delete: resourceAlicloudEbsEnterpriseSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"cross_region_copy_info": {
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"regions": {
							Optional: true,
							Type:     schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"retain_days": {
										Optional: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
					},
				},
			},
			"desc": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"enterprise_snapshot_policy_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"retain_rule": {
				MaxItems: 1,
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeInt,
						},
						"time_interval": {
							Optional: true,
							Type:     schema.TypeInt,
						},
						"time_unit": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"schedule": {
				Required: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_expression": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED"}, false),
			},
			"storage_rule": {
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_immediate_access": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"target_type": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudEbsEnterpriseSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewEbsClient()
	if err != nil {
		return WrapError(err)
	}
	request["Name"] = d.Get("enterprise_snapshot_policy_name")
	request["TargetType"] = d.Get("target_type")
	request["State"] = d.Get("status")
	request["Desc"] = d.Get("desc")

	request["Schedule"], err = convertMaptoJsonString(map[string]interface{}{
		"CronExpression": d.Get("schedule").([]interface{})[0].(map[string]interface{})["cron_expression"],
	})
	if err != nil {
		return WrapError(err)
	}

	retainRuleObj := d.Get("retain_rule").([]interface{})[0].(map[string]interface{})
	retainRuleMap := make(map[string]interface{}, 0)
	if v, ok := retainRuleObj["number"]; ok && formatInt(v) != 0 {
		retainRuleMap["Number"] = v
	}
	if v, ok := retainRuleObj["time_interval"]; ok {
		retainRuleMap["TimeInterval"] = v
	}
	if v, ok := retainRuleObj["time_unit"]; ok {
		retainRuleMap["TimeUnit"] = v
	}
	request["RetainRule"], err = convertMaptoJsonString(retainRuleMap)
	if err != nil {
		return WrapError(err)
	}

	storageRuleMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("storage_rule"); ok && len(v.([]interface{})) > 0 {
		storageRuleObg := v.([]interface{})[0].(map[string]interface{})
		storageRuleMap["EnableImmediateAccess"] = storageRuleObg["enable_immediate_access"]
		request["StorageRule"], err = convertMaptoJsonString(storageRuleMap)
		if err != nil {
			return WrapError(err)
		}
	}

	crossRegionCopyInfoMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("cross_region_copy_info"); ok && v.(*schema.Set).Len() > 0 {
		crossRegionCopyInfoObj := v.(*schema.Set).List()[0].(map[string]interface{})
		crossRegionCopyInfoMap["Enabled"] = crossRegionCopyInfoObj["enabled"]

		regionsMaps := make([]map[string]interface{}, 0)
		if v, ok := crossRegionCopyInfoObj["regions"]; ok {
			for _, vv := range v.(*schema.Set).List() {
				regionsObj := vv.(map[string]interface{})

				regionsMaps = append(regionsMaps, map[string]interface{}{
					"RegionId":   regionsObj["region_id"],
					"RetainDays": regionsObj["retain_days"],
				})
			}
		}
		crossRegionCopyInfoMap["Regions"] = regionsMaps
	}
	request["CrossRegionCopyInfo"], err = convertMaptoJsonString(crossRegionCopyInfoMap)
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}

	request["ClientToken"] = buildClientToken("CreateEnterpriseSnapshotPolicy")
	var response map[string]interface{}
	action := "CreateEnterpriseSnapshotPolicy"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_enterprise_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PolicyId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ebs_enterprise_snapshot_policy")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudEbsEnterpriseSnapshotPolicyRead(d, meta)
}

func resourceAlicloudEbsEnterpriseSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}

	object, err := ebsService.DescribeEbsEnterpriseSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_enterprise_snapshot_policy ebsService.DescribeEbsEnterpriseSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	createTime57 := object["CreateTime"]
	d.Set("create_time", createTime57)
	crossRegionCopyInfo15Maps := make([]map[string]interface{}, 0)
	crossRegionCopyInfo15Map := make(map[string]interface{})
	crossRegionCopyInfo15Raw := object["CrossRegionCopyInfo"].(map[string]interface{})
	crossRegionCopyInfo15Map["enabled"] = crossRegionCopyInfo15Raw["Enabled"]
	if regions15Raw, ok := crossRegionCopyInfo15Raw["Regions"]; ok {
		regions15Maps := make([]map[string]interface{}, 0)
		for _, value1 := range regions15Raw.([]interface{}) {
			regions15 := value1.(map[string]interface{})
			regions15Map := make(map[string]interface{})
			regions15Map["region_id"] = regions15["RegionId"]
			regions15Map["retain_days"] = regions15["RetainDays"]
			regions15Maps = append(regions15Maps, regions15Map)
		}
		crossRegionCopyInfo15Map["regions"] = regions15Maps
	}
	crossRegionCopyInfo15Maps = append(crossRegionCopyInfo15Maps, crossRegionCopyInfo15Map)
	d.Set("cross_region_copy_info", crossRegionCopyInfo15Maps)

	desc71 := object["Desc"]
	d.Set("desc", desc71)

	name39 := object["Name"]
	d.Set("enterprise_snapshot_policy_name", name39)

	resourceGroupId30 := object["ResourceGroupId"]
	d.Set("resource_group_id", resourceGroupId30)
	retainRule13Maps := make([]map[string]interface{}, 0)
	retainRule13Map := make(map[string]interface{})
	retainRule13Raw := object["RetainRule"].(map[string]interface{})
	retainRule13Map["number"] = retainRule13Raw["Number"]
	retainRule13Map["time_interval"] = retainRule13Raw["TimeInterval"]
	retainRule13Map["time_unit"] = retainRule13Raw["TimeUnit"]
	retainRule13Maps = append(retainRule13Maps, retainRule13Map)
	d.Set("retain_rule", retainRule13Maps)
	schedule0Maps := make([]map[string]interface{}, 0)
	schedule0Map := make(map[string]interface{})
	schedule0Raw := object["Schedule"].(map[string]interface{})
	schedule0Map["cron_expression"] = schedule0Raw["CronExpression"]
	schedule0Maps = append(schedule0Maps, schedule0Map)
	d.Set("schedule", schedule0Maps)

	state59 := object["State"]
	d.Set("status", state59)
	storageRule20Maps := make([]map[string]interface{}, 0)
	storageRule20Map := make(map[string]interface{})
	storageRule20Raw := object["StorageRule"].(map[string]interface{})
	storageRule20Map["enable_immediate_access"] = storageRule20Raw["EnableImmediateAccess"]
	storageRule20Maps = append(storageRule20Maps, storageRule20Map)
	d.Set("storage_rule", storageRule20Maps)

	targetType83 := object["TargetType"]
	d.Set("target_type", targetType83)

	tagsRaw, err := ebsService.ListTagResources(d.Id(), "EnterpriseSnapshotPolicy")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(tagsRaw))

	return nil
}

func resourceAlicloudEbsEnterpriseSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}
	conn, err := client.NewEbsClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"PolicyId": d.Id(),
		"RegionId": client.RegionId,
	}

	if d.HasChange("desc") {
		update = true
		if v, ok := d.GetOk("desc"); ok {
			request["Desc"] = v
		}
	}
	if d.HasChange("enterprise_snapshot_policy_name") {
		update = true
		if v, ok := d.GetOk("enterprise_snapshot_policy_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("status") {
		update = true
		if v, ok := d.GetOk("status"); ok {
			request["State"] = v
		}
	}

	if d.HasChange("schedule") {
		update = true
	}

	request["Schedule"], err = convertMaptoJsonString(map[string]interface{}{
		"CronExpression": d.Get("schedule").([]interface{})[0].(map[string]interface{})["cron_expression"],
	})
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("retain_rule") {
		update = true
	}
	retainRuleObj := d.Get("retain_rule").([]interface{})[0].(map[string]interface{})
	retainRuleMap := make(map[string]interface{}, 0)
	if v, ok := retainRuleObj["number"]; ok {
		retainRuleMap["Number"] = v
	}
	if v, ok := retainRuleObj["time_interval"]; ok {
		retainRuleMap["TimeInterval"] = v
	}
	if v, ok := retainRuleObj["time_unit"]; ok {
		retainRuleMap["TimeUnit"] = v
	}
	request["RetainRule"], err = convertMaptoJsonString(retainRuleMap)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("storage_rule") {
		update = true
	}
	storageRuleMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("storage_rule"); ok && len(v.([]interface{})) > 0 {
		storageRuleObg := v.([]interface{})[0].(map[string]interface{})
		storageRuleMap["EnableImmediateAccess"] = storageRuleObg["enable_immediate_access"]
	}
	request["StorageRule"], err = convertMaptoJsonString(storageRuleMap)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("cross_region_copy_info") {
		update = true
		crossRegionCopyInfoMap := make(map[string]interface{}, 0)
		if v, ok := d.GetOk("cross_region_copy_info"); ok && v.(*schema.Set).Len() > 0 {
			crossRegionCopyInfoObj := v.(*schema.Set).List()[0].(map[string]interface{})
			crossRegionCopyInfoMap["Enabled"] = crossRegionCopyInfoObj["enabled"]

			regionsMaps := make([]map[string]interface{}, 0)
			if v, ok := crossRegionCopyInfoObj["regions"]; ok {
				for _, vv := range v.(*schema.Set).List() {
					regionsObj := vv.(map[string]interface{})

					regionsMaps = append(regionsMaps, map[string]interface{}{
						"RegionId":   regionsObj["region_id"],
						"RetainDays": regionsObj["retain_days"],
					})
				}
			}
			crossRegionCopyInfoMap["Regions"] = regionsMaps
		}
		request["CrossRegionCopyInfo"], err = convertMaptoJsonString(crossRegionCopyInfoMap)
		if err != nil {
			return WrapError(err)
		}
	}

	if update {
		action := "UpdateEnterpriseSnapshotPolicy"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("desc")
		d.SetPartial("enterprise_snapshot_policy_name")
		d.SetPartial("status")
		d.SetPartial("schedule")
		d.SetPartial("cross_region_copy_info")
		d.SetPartial("storage_rule")
		d.SetPartial("retain_rule")
	}

	if err := ebsService.SetResourceTags(d, "EnterpriseSnapshotPolicy"); err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAlicloudEbsEnterpriseSnapshotPolicyRead(d, meta)
}

func resourceAlicloudEbsEnterpriseSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEbsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"PolicyId": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteEnterpriseSnapshotPolicy")
	action := "DeleteEnterpriseSnapshotPolicy"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
