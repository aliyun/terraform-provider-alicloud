package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudHbrPolicyBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHbrPolicyBindingCreate,
		Read:   resourceAliCloudHbrPolicyBindingRead,
		Update: resourceAliCloudHbrPolicyBindingUpdate,
		Delete: resourceAliCloudHbrPolicyBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"advanced_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oss_detail": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ignore_archive_object": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"inventory_cleanup_policy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"inventory_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"udm_detail": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_disk_id_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"disk_id_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"app_consistent": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"snapshot_group": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ram_role_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"pre_script_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"post_script_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable_fs_freeze": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"timeout_in_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"enable_writers": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_account_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cross_account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"SELF_ACCOUNT", "CROSS_ACCOUNT"}, false),
			},
			"cross_account_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"exclude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_binding_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"UDM_ECS", "NAS", "OSS", "File", "ECS_FILE", "OTS"}, false),
			},
			"speed_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudHbrPolicyBindingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicyBindings"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("policy_id"); ok {
		request["PolicyId"] = v
	}

	policyBindingListDataList := make(map[string]interface{})

	if v, ok := d.GetOkExists("disabled"); ok {
		policyBindingListDataList["Disabled"] = v
	}

	if v, ok := d.GetOkExists("include"); ok {
		policyBindingListDataList["Include"] = v
	}

	if v, ok := d.GetOkExists("cross_account_role_name"); ok {
		policyBindingListDataList["CrossAccountRoleName"] = v
	}

	if v, ok := d.GetOkExists("cross_account_user_id"); ok {
		policyBindingListDataList["CrossAccountUserId"] = v
	}

	if v, ok := d.GetOkExists("data_source_id"); ok {
		policyBindingListDataList["DataSourceId"] = v
	}

	if v, ok := d.GetOkExists("source_type"); ok {
		policyBindingListDataList["SourceType"] = v
	}

	if v, ok := d.GetOkExists("policy_binding_description"); ok {
		policyBindingListDataList["PolicyBindingDescription"] = v
	}

	if v, ok := d.GetOkExists("speed_limit"); ok {
		policyBindingListDataList["SpeedLimit"] = v
	}

	if v, ok := d.GetOkExists("source"); ok {
		policyBindingListDataList["Source"] = v
	}

	if v, ok := d.GetOkExists("cross_account_type"); ok {
		policyBindingListDataList["CrossAccountType"] = v
	}

	if v, ok := d.GetOkExists("exclude"); ok {
		policyBindingListDataList["Exclude"] = v
	}

	if v := d.Get("advanced_options"); !IsNil(v) {
		advancedOptions := make(map[string]interface{})
		ossDetail := make(map[string]interface{})
		ignoreArchiveObject1, _ := jsonpath.Get("$[0].oss_detail[0].ignore_archive_object", d.Get("advanced_options"))
		if ignoreArchiveObject1 != nil && ignoreArchiveObject1 != "" {
			ossDetail["IgnoreArchiveObject"] = ignoreArchiveObject1
		}
		inventoryId1, _ := jsonpath.Get("$[0].oss_detail[0].inventory_id", d.Get("advanced_options"))
		if inventoryId1 != nil && inventoryId1 != "" {
			ossDetail["InventoryId"] = inventoryId1
		}
		inventoryCleanupPolicy1, _ := jsonpath.Get("$[0].oss_detail[0].inventory_cleanup_policy", d.Get("advanced_options"))
		if inventoryCleanupPolicy1 != nil && inventoryCleanupPolicy1 != "" {
			ossDetail["InventoryCleanupPolicy"] = inventoryCleanupPolicy1
		}

		if len(ossDetail) > 0 {
			advancedOptions["OssDetail"] = ossDetail
		}
		udmDetail := make(map[string]interface{})
		diskIdList1, _ := jsonpath.Get("$[0].udm_detail[0].disk_id_list", d.Get("advanced_options"))
		if diskIdList1 != nil && diskIdList1 != "" {
			udmDetail["DiskIdList"] = diskIdList1
		}
		excludeDiskIdList1, _ := jsonpath.Get("$[0].udm_detail[0].exclude_disk_id_list", d.Get("advanced_options"))
		if excludeDiskIdList1 != nil && excludeDiskIdList1 != "" {
			udmDetail["ExcludeDiskIdList"] = excludeDiskIdList1
		}
		destinationKmsKeyId1, _ := jsonpath.Get("$[0].udm_detail[0].destination_kms_key_id", d.Get("advanced_options"))
		if destinationKmsKeyId1 != nil && destinationKmsKeyId1 != "" {
			udmDetail["DestinationKmsKeyId"] = destinationKmsKeyId1
		}
		appConsistent1, _ := jsonpath.Get("$[0].udm_detail[0].app_consistent", d.Get("advanced_options"))
		if appConsistent1 != nil && appConsistent1 != "" {
			udmDetail["AppConsistent"] = appConsistent1
		}
		snapshotGroup1, _ := jsonpath.Get("$[0].udm_detail[0].snapshot_group", d.Get("advanced_options"))
		if snapshotGroup1 != nil && snapshotGroup1 != "" {
			udmDetail["SnapshotGroup"] = snapshotGroup1
		}
		ramRoleName1, _ := jsonpath.Get("$[0].udm_detail[0].ram_role_name", d.Get("advanced_options"))
		if ramRoleName1 != nil && ramRoleName1 != "" {
			udmDetail["RamRoleName"] = ramRoleName1
		}
		preScriptPath1, _ := jsonpath.Get("$[0].udm_detail[0].pre_script_path", d.Get("advanced_options"))
		if preScriptPath1 != nil && preScriptPath1 != "" {
			udmDetail["PreScriptPath"] = preScriptPath1
		}
		postScriptPath1, _ := jsonpath.Get("$[0].udm_detail[0].post_script_path", d.Get("advanced_options"))
		if postScriptPath1 != nil && postScriptPath1 != "" {
			udmDetail["PostScriptPath"] = postScriptPath1
		}
		enableFsFreeze1, _ := jsonpath.Get("$[0].udm_detail[0].enable_fs_freeze", d.Get("advanced_options"))
		if enableFsFreeze1 != nil && enableFsFreeze1 != "" {
			udmDetail["EnableFsFreeze"] = enableFsFreeze1
		}
		timeoutInSeconds1, _ := jsonpath.Get("$[0].udm_detail[0].timeout_in_seconds", d.Get("advanced_options"))
		if timeoutInSeconds1 != nil && timeoutInSeconds1 != "" {
			udmDetail["TimeoutInSeconds"] = timeoutInSeconds1
		}
		enableWriters1, _ := jsonpath.Get("$[0].udm_detail[0].enable_writers", d.Get("advanced_options"))
		if enableWriters1 != nil && enableWriters1 != "" {
			udmDetail["EnableWriters"] = enableWriters1
		}

		if len(udmDetail) > 0 {
			advancedOptions["UdmDetail"] = udmDetail
		}

		if len(advancedOptions) > 0 {
			policyBindingListDataList["AdvancedOptions"] = advancedOptions
		}
	}

	PolicyBindingListMap := make([]interface{}, 0)
	PolicyBindingListMap = append(PolicyBindingListMap, policyBindingListDataList)
	policyBindingListDataListJson, err := json.Marshal(PolicyBindingListMap)
	if err != nil {
		return WrapError(err)
	}
	request["PolicyBindingList"] = string(policyBindingListDataListJson)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_policy_binding", action, AlibabaCloudSdkGoERROR)
	}

	PolicyBindingListSourceTypeVar := d.Get("source_type")
	PolicyBindingListDataSourceIdVar := d.Get("data_source_id")
	d.SetId(fmt.Sprintf("%v:%v:%v", request["PolicyId"], PolicyBindingListSourceTypeVar, PolicyBindingListDataSourceIdVar))

	// Capture the user-configured UdmDetail booleans before Read overwrites
	// the ResourceData writer with backend values.
	udmBoolFields := []string{"snapshot_group", "app_consistent", "enable_fs_freeze", "enable_writers"}
	configUdmBools := make(map[string]bool, len(udmBoolFields))
	for _, f := range udmBoolFields {
		configUdmBools[f] = hbrPolicyBindingUdmDetailBool(d, f)
	}
	// advancedOptions was built from the configuration above and stored on the
	// policy binding list payload; keep a reference so the Create-then-Update
	// fallback below can re-apply it verbatim.
	configAdvancedOptions, _ := policyBindingListDataList["AdvancedOptions"].(map[string]interface{})

	if err := resourceAliCloudHbrPolicyBindingRead(d, meta); err != nil {
		return err
	}

	// Create-then-Update fallback: the CreatePolicyBindings backend silently
	// coerces some UdmDetail booleans (notably snapshot_group) from true to
	// false, while UpdatePolicyBinding persists them. When the user configured
	// any such field as true but the backend stored false, re-apply the full
	// AdvancedOptions via UpdatePolicyBinding so the user's configuration takes
	// effect. This fallback can be removed once CreatePolicyBindings persists
	// these fields.
	needReapply := false
	for _, f := range udmBoolFields {
		if configUdmBools[f] && !hbrPolicyBindingUdmDetailBool(d, f) {
			needReapply = true
			break
		}
	}
	if needReapply {
		if err := hbrPolicyBindingReapplyAdvancedOptions(d, meta, configAdvancedOptions); err != nil {
			return err
		}
		return resourceAliCloudHbrPolicyBindingRead(d, meta)
	}
	return nil
}

func resourceAliCloudHbrPolicyBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrServiceV2 := HbrServiceV2{client}

	objectRaw, err := hbrServiceV2.DescribeHbrPolicyBinding(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_policy_binding DescribeHbrPolicyBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("cross_account_role_name", objectRaw["CrossAccountRoleName"])
	d.Set("cross_account_type", objectRaw["CrossAccountType"])
	d.Set("cross_account_user_id", objectRaw["CrossAccountUserId"])
	d.Set("disabled", objectRaw["Disabled"])
	d.Set("exclude", objectRaw["Exclude"])
	d.Set("include", objectRaw["Include"])
	d.Set("policy_binding_description", objectRaw["PolicyBindingDescription"])
	d.Set("source", objectRaw["Source"])
	d.Set("speed_limit", objectRaw["SpeedLimit"])
	d.Set("data_source_id", objectRaw["DataSourceId"])
	d.Set("policy_id", objectRaw["PolicyId"])
	d.Set("source_type", objectRaw["SourceType"])

	advancedOptionsMaps := make([]map[string]interface{}, 0)
	advancedOptionsMap := make(map[string]interface{})
	advancedOptionsRaw := make(map[string]interface{})
	if objectRaw["AdvancedOptions"] != nil {
		advancedOptionsRaw = objectRaw["AdvancedOptions"].(map[string]interface{})
	}
	if len(advancedOptionsRaw) > 0 {

		ossDetailMaps := make([]map[string]interface{}, 0)
		ossDetailMap := make(map[string]interface{})
		ossDetailRaw := make(map[string]interface{})
		if advancedOptionsRaw["OssDetail"] != nil {
			ossDetailRaw = advancedOptionsRaw["OssDetail"].(map[string]interface{})
		}
		if len(ossDetailRaw) > 0 {
			ossDetailMap["ignore_archive_object"] = ossDetailRaw["IgnoreArchiveObject"]
			ossDetailMap["inventory_cleanup_policy"] = ossDetailRaw["InventoryCleanupPolicy"]
			ossDetailMap["inventory_id"] = ossDetailRaw["InventoryId"]

			ossDetailMaps = append(ossDetailMaps, ossDetailMap)
		}
		advancedOptionsMap["oss_detail"] = ossDetailMaps
		udmDetailMaps := make([]map[string]interface{}, 0)
		udmDetailMap := make(map[string]interface{})
		udmDetailRaw := make(map[string]interface{})
		if advancedOptionsRaw["UdmDetail"] != nil {
			udmDetailRaw = advancedOptionsRaw["UdmDetail"].(map[string]interface{})
		}
		if len(udmDetailRaw) > 0 {
			udmDetailMap["destination_kms_key_id"] = udmDetailRaw["DestinationKmsKeyId"]

			diskIdListRaw := make([]interface{}, 0)
			if udmDetailRaw["DiskIdList"] != nil {
				diskIdListRaw = convertToInterfaceArray(udmDetailRaw["DiskIdList"])
			}

			udmDetailMap["disk_id_list"] = diskIdListRaw
			excludeDiskIdListRaw := make([]interface{}, 0)
			if udmDetailRaw["ExcludeDiskIdList"] != nil {
				excludeDiskIdListRaw = convertToInterfaceArray(udmDetailRaw["ExcludeDiskIdList"])
			}

			udmDetailMap["exclude_disk_id_list"] = excludeDiskIdListRaw
			udmDetailMap["app_consistent"] = udmDetailRaw["AppConsistent"]
			udmDetailMap["snapshot_group"] = udmDetailRaw["SnapshotGroup"]
			udmDetailMap["ram_role_name"] = udmDetailRaw["RamRoleName"]
			udmDetailMap["pre_script_path"] = udmDetailRaw["PreScriptPath"]
			udmDetailMap["post_script_path"] = udmDetailRaw["PostScriptPath"]
			udmDetailMap["enable_fs_freeze"] = udmDetailRaw["EnableFsFreeze"]
			udmDetailMap["timeout_in_seconds"] = formatInt(udmDetailRaw["TimeoutInSeconds"])
			udmDetailMap["enable_writers"] = udmDetailRaw["EnableWriters"]
			udmDetailMaps = append(udmDetailMaps, udmDetailMap)
		}
		advancedOptionsMap["udm_detail"] = udmDetailMaps
		advancedOptionsMaps = append(advancedOptionsMaps, advancedOptionsMap)
	}
	if err := d.Set("advanced_options", advancedOptionsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudHbrPolicyBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdatePolicyBinding"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyId"] = parts[0]
	request["DataSourceId"] = parts[2]
	request["SourceType"] = parts[1]

	if d.HasChange("disabled") {
		update = true
		request["Disabled"] = d.Get("disabled")
	}

	if d.HasChange("advanced_options") {
		update = true
		advancedOptions := make(map[string]interface{})

		if v := d.Get("advanced_options"); v != nil {
			ossDetail := make(map[string]interface{})
			ignoreArchiveObject1, _ := jsonpath.Get("$[0].oss_detail[0].ignore_archive_object", d.Get("advanced_options"))
			if ignoreArchiveObject1 != nil && ignoreArchiveObject1 != "" {
				ossDetail["IgnoreArchiveObject"] = ignoreArchiveObject1
			}
			inventoryId1, _ := jsonpath.Get("$[0].oss_detail[0].inventory_id", d.Get("advanced_options"))
			if inventoryId1 != nil && inventoryId1 != "" {
				ossDetail["InventoryId"] = inventoryId1
			}
			inventoryCleanupPolicy1, _ := jsonpath.Get("$[0].oss_detail[0].inventory_cleanup_policy", d.Get("advanced_options"))
			if inventoryCleanupPolicy1 != nil && inventoryCleanupPolicy1 != "" {
				ossDetail["InventoryCleanupPolicy"] = inventoryCleanupPolicy1
			}

			if len(ossDetail) > 0 {
				advancedOptions["OssDetail"] = ossDetail
			}
			udmDetail := make(map[string]interface{})
			diskIdList1, _ := jsonpath.Get("$[0].udm_detail[0].disk_id_list", d.Get("advanced_options"))
			if diskIdList1 != nil && diskIdList1 != "" {
				udmDetail["DiskIdList"] = diskIdList1
			}
			excludeDiskIdList1, _ := jsonpath.Get("$[0].udm_detail[0].exclude_disk_id_list", d.Get("advanced_options"))
			if excludeDiskIdList1 != nil && excludeDiskIdList1 != "" {
				udmDetail["ExcludeDiskIdList"] = excludeDiskIdList1
			}
			destinationKmsKeyId1, _ := jsonpath.Get("$[0].udm_detail[0].destination_kms_key_id", d.Get("advanced_options"))
			if destinationKmsKeyId1 != nil && destinationKmsKeyId1 != "" {
				udmDetail["DestinationKmsKeyId"] = destinationKmsKeyId1
			}
			appConsistent1, _ := jsonpath.Get("$[0].udm_detail[0].app_consistent", d.Get("advanced_options"))
			if appConsistent1 != nil && appConsistent1 != "" {
				udmDetail["AppConsistent"] = appConsistent1
			}
			snapshotGroup1, _ := jsonpath.Get("$[0].udm_detail[0].snapshot_group", d.Get("advanced_options"))
			if snapshotGroup1 != nil && snapshotGroup1 != "" {
				udmDetail["SnapshotGroup"] = snapshotGroup1
			}
			ramRoleName1, _ := jsonpath.Get("$[0].udm_detail[0].ram_role_name", d.Get("advanced_options"))
			if ramRoleName1 != nil && ramRoleName1 != "" {
				udmDetail["RamRoleName"] = ramRoleName1
			}
			preScriptPath1, _ := jsonpath.Get("$[0].udm_detail[0].pre_script_path", d.Get("advanced_options"))
			if preScriptPath1 != nil && preScriptPath1 != "" {
				udmDetail["PreScriptPath"] = preScriptPath1
			}
			postScriptPath1, _ := jsonpath.Get("$[0].udm_detail[0].post_script_path", d.Get("advanced_options"))
			if postScriptPath1 != nil && postScriptPath1 != "" {
				udmDetail["PostScriptPath"] = postScriptPath1
			}
			enableFsFreeze1, _ := jsonpath.Get("$[0].udm_detail[0].enable_fs_freeze", d.Get("advanced_options"))
			if enableFsFreeze1 != nil && enableFsFreeze1 != "" {
				udmDetail["EnableFsFreeze"] = enableFsFreeze1
			}
			timeoutInSeconds1, _ := jsonpath.Get("$[0].udm_detail[0].timeout_in_seconds", d.Get("advanced_options"))
			if timeoutInSeconds1 != nil && timeoutInSeconds1 != "" {
				udmDetail["TimeoutInSeconds"] = timeoutInSeconds1
			}
			enableWriters1, _ := jsonpath.Get("$[0].udm_detail[0].enable_writers", d.Get("advanced_options"))
			if enableWriters1 != nil && enableWriters1 != "" {
				udmDetail["EnableWriters"] = enableWriters1
			}

			if len(udmDetail) > 0 {
				advancedOptions["UdmDetail"] = udmDetail
			}

			advancedOptionsJson, err := json.Marshal(advancedOptions)
			if err != nil {
				return WrapError(err)
			}
			request["AdvancedOptions"] = string(advancedOptionsJson)
		}
	}

	if d.HasChange("include") {
		update = true
		request["Include"] = d.Get("include")
	}

	if d.HasChange("policy_binding_description") {
		update = true
		request["PolicyBindingDescription"] = d.Get("policy_binding_description")
	}

	if d.HasChange("speed_limit") {
		update = true
		request["SpeedLimit"] = d.Get("speed_limit")
	}

	if d.HasChange("source") {
		update = true
		request["Source"] = d.Get("source")
	}

	if d.HasChange("exclude") {
		update = true
		request["Exclude"] = d.Get("exclude")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
	}

	return resourceAliCloudHbrPolicyBindingRead(d, meta)
}

func resourceAliCloudHbrPolicyBindingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePolicyBinding"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyId"] = parts[0]
	request["SourceType"] = parts[1]
	request["DataSourceIds"] = "[\"" + parts[2] + "\"]"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// hbrPolicyBindingUdmDetailBool reads a boolean field from
// advanced_options.0.udm_detail[0].<field> in the current ResourceData.
// Before Read this reflects the configured value; after Read it reflects the
// backend value, because the SDK writer (populated by d.Set during Read)
// takes precedence over the diff.
func hbrPolicyBindingUdmDetailBool(d *schema.ResourceData, field string) bool {
	raw, _ := jsonpath.Get("$[0].udm_detail[0]."+field, d.Get("advanced_options"))
	b, _ := raw.(bool)
	return b
}

// hbrPolicyBindingReapplyAdvancedOptions re-issues UpdatePolicyBinding with the
// full user-configured AdvancedOptions. It is the Create-then-Update fallback
// for UdmDetail booleans that CreatePolicyBindings silently coerces (notably
// snapshot_group=true -> false); UpdatePolicyBinding persists them correctly.
func hbrPolicyBindingReapplyAdvancedOptions(d *schema.ResourceData, meta interface{}, advancedOptions map[string]interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UpdatePolicyBinding"
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	request["PolicyId"] = parts[0]
	request["DataSourceId"] = parts[2]
	request["SourceType"] = parts[1]

	advancedOptionsJson, err := json.Marshal(advancedOptions)
	if err != nil {
		return WrapError(err)
	}
	request["AdvancedOptions"] = string(advancedOptionsJson)

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
