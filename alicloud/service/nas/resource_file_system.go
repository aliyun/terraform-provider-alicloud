package nas

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileSystemCreate,
		Read:   resourceFileSystemRead,
		Update: resourceFileSystemUpdate,
		Delete: resourceFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_upgrade_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_used_ratio": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.IntBetween(1, 99),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"step": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypt_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 1, 2}),
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"standard", "extreme", "cpfs", "cpfsse"}, false),
			},
			"keytab": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keytab_md5": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nfs_acl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_oplock": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"NFS", "SMB", "cpfs"}, false),
			},
			"recycle_bin": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"Enable", "Disable"}, false),
						},
						"reserved_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.IntBetween(0, 180),
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if v, ok := d.GetOk("recycle_bin.0.status"); ok && v.(string) == "Enable" {
									return false
								}
								return true
							},
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"secondary_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"redundancy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"LRS", "ZRS"}, false),
			},
			"redundancy_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"smb_acl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"home_dir_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_anonymous_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"super_admin_sid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"reject_unencrypted_access": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"encrypt_data": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Performance", "Capacity", "standard", "advance", "advance_100", "advance_200", "Premium", "Agentic"}, false),
			},
			"tags": tags.TagsSchema(),
			"vswitch_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: fileSystemVpcDiffSuppressFunc,
			},
			"vpc_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: fileSystemVpcDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

// fileSystemVpcDiffSuppressFunc suppresses spurious diffs on vpc_id and vswitch_id
// for standard and extreme NAS file systems. The NAS DescribeFileSystem API does not
// return VpcId/QuorumVswId for these file system types, so the value read back into
// state is empty while the user's config still holds the value supplied at create
// time. Without suppression this produces a permanent non-empty diff that, combined
// with ForceNew, forces a destructive replace on every plan. The fields are immutable
// at the API level (ModifyFileSystem does not accept them), so suppressing the diff
// for types that never surface the value is safe.
func fileSystemVpcDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if fsType := d.Get("file_system_type").(string); fsType == "standard" || fsType == "extreme" {
		return true
	}
	return false
}

func resourceFileSystemCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("redundancy_type"); ok {
		request["RedundancyType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOkExists("capacity"); ok {
		request["Capacity"] = v
	}
	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("redundancy_vswitch_ids"); ok {
		redundancyVSwitchIdsMapsArray := convert.ConvertToInterfaceArray(v)

		request["RedundancyVSwitchIds"] = redundancyVSwitchIdsMapsArray
	}

	if v, ok := d.GetOkExists("encrypt_type"); ok {
		request["EncryptType"] = v
	}
	request["DryRun"] = "false"
	request["ProtocolType"] = d.Get("protocol_type")
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"User.InDebt"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_nas_file_system", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))

	service := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, service.FileSystemStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceFileSystemUpdate(d, meta)
}

func resourceFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeFileSystem(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_file_system DescribeFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("capacity", objectRaw["Capacity"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("encrypt_type", objectRaw["EncryptType"])
	d.Set("file_system_type", objectRaw["FileSystemType"])
	d.Set("kms_key_id", objectRaw["KMSKeyId"])
	d.Set("protocol_type", objectRaw["ProtocolType"])
	d.Set("redundancy_type", objectRaw["RedundancyType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("vswitch_id", objectRaw["QuorumVswId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	autoUpgradeConfigMaps := make([]map[string]interface{}, 0)
	if autoUpgradeConfigRaw, ok := objectRaw["AutoUpgradeConfig"].(map[string]interface{}); ok && len(autoUpgradeConfigRaw) > 0 {
		autoUpgradeConfigMap := make(map[string]interface{})
		if v, ok := autoUpgradeConfigRaw["enabled"]; ok {
			autoUpgradeConfigMap["enabled"] = v
		}
		if v, ok := autoUpgradeConfigRaw["capacityUsedRatio"]; ok {
			capacityUsedRatio, err := convert.FormatInt(v)
			if err != nil {
				return errs.WrapError(err)
			}
			autoUpgradeConfigMap["capacity_used_ratio"] = capacityUsedRatio
		}
		if v, ok := autoUpgradeConfigRaw["step"]; ok {
			step, err := convert.FormatInt(v)
			if err != nil {
				return errs.WrapError(err)
			}
			autoUpgradeConfigMap["step"] = step
		}
		if v, ok := autoUpgradeConfigRaw["time"]; ok {
			upgradeTime, err := convert.FormatInt(v)
			if err != nil {
				return errs.WrapError(err)
			}
			autoUpgradeConfigMap["time"] = upgradeTime
		}
		autoUpgradeConfigMaps = append(autoUpgradeConfigMaps, autoUpgradeConfigMap)
	}
	if err := d.Set("auto_upgrade_config", autoUpgradeConfigMaps); err != nil {
		return errs.WrapError(err)
	}

	optionsMaps := make([]map[string]interface{}, 0)
	optionsMap := make(map[string]interface{})
	optionsRaw := make(map[string]interface{})
	if objectRaw["Options"] != nil {
		optionsRaw = objectRaw["Options"].(map[string]interface{})
	}
	if len(optionsRaw) > 0 {
		optionsMap["enable_oplock"] = optionsRaw["EnableOplock"]

		optionsMaps = append(optionsMaps, optionsMap)
	}
	if err := d.Set("options", optionsMaps); err != nil {
		return err
	}
	if objectRaw["RedundancyType"] == "ZRS" {
		redundancyVSwitchIdRaw, _ := jsonpath.Get("$.RedundancyVSwitchIds.RedundancyVSwitchId", objectRaw)
		d.Set("redundancy_vswitch_ids", redundancyVSwitchIdRaw)
	}

	// Agentic file systems are standard-type file systems without the general-purpose
	// NAS recycle bin and ACL features: skip those cross-contract reads and keep the
	// corresponding attributes empty in state instead of recording ineffective values.
	isAgenticStorage := d.Get("storage_type") == "Agentic"
	if isAgenticStorage {
		if err := d.Set("smb_acl", []interface{}{}); err != nil {
			return errs.WrapError(err)
		}
		if err := d.Set("recycle_bin", []interface{}{}); err != nil {
			return errs.WrapError(err)
		}
		if err := d.Set("nfs_acl", []interface{}{}); err != nil {
			return errs.WrapError(err)
		}
	}
	checkValue00 := d.Get("file_system_type")
	checkValue01 := d.Get("protocol_type")
	if !isAgenticStorage && (checkValue00 == "standard") && (checkValue01 == "SMB") {
		objectRaw, err = service.DescribeSmbAcl(d.Id())
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		smbAclMaps := make([]map[string]interface{}, 0)
		smbAclMap := make(map[string]interface{})

		smbAclMap["enable_anonymous_access"] = objectRaw["EnableAnonymousAccess"]
		smbAclMap["enabled"] = objectRaw["Enabled"]
		smbAclMap["encrypt_data"] = objectRaw["EncryptData"]
		smbAclMap["home_dir_path"] = objectRaw["HomeDirPath"]
		smbAclMap["reject_unencrypted_access"] = objectRaw["RejectUnencryptedAccess"]
		smbAclMap["super_admin_sid"] = objectRaw["SuperAdminSid"]

		smbAclMaps = append(smbAclMaps, smbAclMap)
		if err := d.Set("smb_acl", smbAclMaps); err != nil {
			return err
		}

	}
	checkValue00 = d.Get("file_system_type")
	if checkValue00 == "standard" && !isAgenticStorage {
		objectRaw, err = service.GetRecycleBinAttribute(d.Id())
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		recycleBinMaps := make([]map[string]interface{}, 0)
		recycleBinMap := make(map[string]interface{})

		recycleBinMap["enable_time"] = objectRaw["EnableTime"]
		recycleBinMap["reserved_days"] = objectRaw["ReservedDays"]
		recycleBinMap["secondary_size"] = objectRaw["SecondarySize"]
		recycleBinMap["size"] = objectRaw["Size"]
		recycleBinMap["status"] = objectRaw["Status"]

		recycleBinMaps = append(recycleBinMaps, recycleBinMap)
		if err := d.Set("recycle_bin", recycleBinMaps); err != nil {
			return err
		}

	}
	objectRaw, err = service.ListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	checkValue00 = d.Get("file_system_type")
	checkValue01 = d.Get("protocol_type")
	if !isAgenticStorage && (checkValue00 == "standard") && (checkValue01 == "NFS") {
		objectRaw, err = service.DescribeNfsAcl(d.Id())
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		nfsAclMaps := make([]map[string]interface{}, 0)
		nfsAclMap := make(map[string]interface{})

		nfsAclMap["enabled"] = objectRaw["Enabled"]

		nfsAclMaps = append(nfsAclMaps, nfsAclMap)
		if err := d.Set("nfs_acl", nfsAclMaps); err != nil {
			return err
		}

	}

	return nil
}

func resourceFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	service := NewService(client)
	objectRaw, _ := service.DescribeFileSystem(d.Id())

	// Agentic file systems are standard-type file systems, but the recycle bin and ACL
	// features below are general-purpose NAS features and their APIs do not apply to
	// them; reject the change up front instead of silently skipping it.
	isAgenticStorage := d.Get("storage_type") == "Agentic"

	if d.HasChange("recycle_bin.0.status") && isAgenticStorage {
		return errs.WrapError(fmt.Errorf("recycle_bin is not supported for Agentic file systems and cannot be modified on %s", d.Id()))
	}
	if d.HasChange("recycle_bin.0.status") {
		objectRaw1, _ := service.GetRecycleBinAttribute(d.Id())
		var err error
		target := d.Get("recycle_bin.0.status").(string)
		enableEnableRecycleBinEnable := false
		checkValue00 := objectRaw1["Status"]
		checkValue01 := objectRaw["FileSystemType"]
		if (checkValue00 == "Disable") && (checkValue01 == "standard") {
			enableEnableRecycleBinEnable = true
		}
		if enableEnableRecycleBinEnable && target == "Enable" {
			action := "EnableRecycleBin"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			if v, ok := d.GetOkExists("recycle_bin"); ok {
				recycleBinReservedDaysJsonPath, err := jsonpath.Get("$[0].reserved_days", v)
				if err == nil && recycleBinReservedDaysJsonPath != "" {
					request["ReservedDays"] = recycleBinReservedDaysJsonPath
				}
			}
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
		enableDisableAndCleanRecycleBinDisable := false

		checkValue01 = objectRaw["FileSystemType"]
		if (checkValue00 == "Enable") && (checkValue01 == "standard") {
			enableDisableAndCleanRecycleBinDisable = true
		}
		if enableDisableAndCleanRecycleBinDisable && target == "Disable" {
			action := "DisableAndCleanRecycleBin"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["FileSystemId"] = d.Id()

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcGet("NAS", "2017-06-26", action, query, request)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	if d.HasChange("nfs_acl.0.enabled") && isAgenticStorage {
		return errs.WrapError(fmt.Errorf("nfs_acl is not supported for Agentic file systems and cannot be modified on %s", d.Id()))
	}
	if d.HasChange("nfs_acl.0.enabled") {
		objectRaw2, _ := service.DescribeNfsAcl(d.Id())
		var err error
		target := d.Get("nfs_acl.0.enabled").(bool)
		enableEnableNfsAclTrue := false
		checkValue00 := objectRaw["FileSystemType"]
		checkValue01 := objectRaw2["Enabled"]
		checkValue02 := objectRaw["ProtocolType"]
		if (checkValue00 == "standard") && (checkValue01 == false) && (checkValue02 == "NFS") {
			enableEnableNfsAclTrue = true
		}
		if enableEnableNfsAclTrue && target == true {
			action := "EnableNfsAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
		enableDisableNfsAclFalse := false
		checkValue00 = objectRaw["FileSystemType"]

		checkValue02 = objectRaw["ProtocolType"]
		if (checkValue00 == "standard") && (checkValue01 == true) && (checkValue02 == "NFS") {
			enableDisableNfsAclFalse = true
		}
		if enableDisableNfsAclFalse && target == false {
			action := "DisableNfsAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	if d.HasChange("smb_acl.0.enabled") && isAgenticStorage {
		return errs.WrapError(fmt.Errorf("smb_acl is not supported for Agentic file systems and cannot be modified on %s", d.Id()))
	}
	if d.HasChange("smb_acl.0.enabled") {
		var err error
		objectRaw3, err := service.DescribeSmbAcl(d.Id())
		target := d.Get("smb_acl.0.enabled").(bool)
		enableEnableSmbAclTrue := false
		checkValue00 := objectRaw["FileSystemType"]
		checkValue01 := objectRaw["ProtocolType"]
		checkValue02 := objectRaw3["Enabled"]
		if (checkValue00 == "standard") && (checkValue01 == "SMB") && (checkValue02 == false) {
			enableEnableSmbAclTrue = true
		}
		if enableEnableSmbAclTrue && target == true {
			action := "EnableSmbAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			request["KeytabMd5"] = d.Get("keytab_md5")
			request["Keytab"] = d.Get("keytab")
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
		enableDisableSmbAclFalse := false
		checkValue00 = objectRaw["FileSystemType"]
		checkValue01 = objectRaw["ProtocolType"]

		if (checkValue00 == "standard") && (checkValue01 == "SMB") && (checkValue02 == true) {
			enableDisableSmbAclFalse = true
		}
		if enableDisableSmbAclFalse && target == false {
			action := "DisableSmbAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	var err error
	action := "ModifyFileSystem"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	if d.HasChange("options") {
		update = true
		options := make(map[string]interface{})

		if v := d.Get("options"); v != nil {
			enableOplock1, _ := jsonpath.Get("$[0].enable_oplock", v)
			if enableOplock1 != nil && enableOplock1 != "" {
				options["EnableOplock"] = enableOplock1
			}

			optionsJson, err := json.Marshal(options)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Options"] = string(optionsJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("auto_upgrade_config") {
		// auto_upgrade_config and its sub-attributes are Optional+Computed: the
		// state echo of a value that is no longer in the configuration is
		// indistinguishable from an explicitly configured value through
		// GetOkExists. Read the raw configuration instead so only attributes the
		// user actually configured are sent to ModifyFileSystem; omitted ones are
		// left to the server (computed) instead of being re-sent from state.
		autoUpgradeConfig := make(map[string]interface{})
		rawConfig := d.GetRawConfig()
		if !rawConfig.IsNull() {
			rawAutoUpgradeConfig := rawConfig.GetAttr("auto_upgrade_config")
			if !rawAutoUpgradeConfig.IsNull() && rawAutoUpgradeConfig.IsKnown() && rawAutoUpgradeConfig.LengthInt() > 0 {
				autoUpgradeConfigRaw := rawAutoUpgradeConfig.Index(cty.NumberIntVal(0))
				if !autoUpgradeConfigRaw.IsNull() && autoUpgradeConfigRaw.IsKnown() {
					if v := autoUpgradeConfigRaw.GetAttr("enabled"); !v.IsNull() && v.IsKnown() {
						autoUpgradeConfig["enabled"] = v.True()
					}
					if v := autoUpgradeConfigRaw.GetAttr("capacity_used_ratio"); !v.IsNull() && v.IsKnown() {
						capacityUsedRatio, _ := v.AsBigFloat().Int64()
						autoUpgradeConfig["capacityUsedRatio"] = int(capacityUsedRatio)
					}
					if v := autoUpgradeConfigRaw.GetAttr("step"); !v.IsNull() && v.IsKnown() {
						step, _ := v.AsBigFloat().Int64()
						autoUpgradeConfig["step"] = int(step)
					}
					if v := autoUpgradeConfigRaw.GetAttr("time"); !v.IsNull() && v.IsKnown() {
						upgradeTime, _ := v.AsBigFloat().Int64()
						autoUpgradeConfig["time"] = int(upgradeTime)
					}
				}
			}
		}
		// A removed block has no explicitly configured attributes; do not send a
		// stale object and do not issue an empty ModifyFileSystem request.
		if len(autoUpgradeConfig) > 0 {
			update = true
			autoUpgradeConfigJson, err := json.Marshal(autoUpgradeConfig)
			if err != nil {
				return errs.WrapError(err)
			}
			request["AutoUpgradeConfig"] = string(autoUpgradeConfigJson)
		}
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	objectRaw, _ = service.DescribeFileSystem(d.Id())
	enableUpgradeFileSystem1 := false
	checkValue00 := objectRaw["FileSystemType"]
	if !(checkValue00 == "standard") {
		enableUpgradeFileSystem1 = true
	}
	action = "UpgradeFileSystem"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("capacity") {
		update = true
	}
	request["Capacity"] = d.Get("capacity")
	request["DryRun"] = "false"
	if update && enableUpgradeFileSystem1 {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.FileSystemStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "filesystem"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	objectRaw, _ = service.DescribeFileSystem(d.Id())
	enableUpdateRecycleBinAttribute1 := false

	checkValue01 := objectRaw["FileSystemType"]
	if (checkValue00 == "Enable") && (checkValue01 == "standard") {
		enableUpdateRecycleBinAttribute1 = true
	}
	action = "UpdateRecycleBinAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = d.Id()

	if d.HasChange("recycle_bin.0.reserved_days") {
		update = true
	}
	if v, ok := d.GetOkExists("recycle_bin.0.reserved_days"); ok {
		query["ReservedDays"] = strconv.Itoa(v.(int))
	}

	if update && enableUpdateRecycleBinAttribute1 {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcGet("NAS", "2017-06-26", action, query, request)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	objectRaw, _ = service.DescribeFileSystem(d.Id())
	var objectRaw3 map[string]interface{}
	if !isAgenticStorage {
		objectRaw3, _ = service.DescribeSmbAcl(d.Id())
	}
	enableModifySmbAcl1 := false
	checkValue00 = objectRaw["FileSystemType"]
	checkValue01 = objectRaw["ProtocolType"]
	checkValue02 := objectRaw3["Enabled"]
	if (checkValue00 == "standard") && (checkValue01 == "SMB") && (checkValue02 == true) {
		enableModifySmbAcl1 = true
	}
	action = "ModifySmbAcl"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	if d.HasChange("smb_acl.0.super_admin_sid") {
		update = true
		smbAclSuperAdminSidJsonPath, err := jsonpath.Get("$[0].super_admin_sid", d.Get("smb_acl"))
		if err == nil {
			request["SuperAdminSid"] = smbAclSuperAdminSidJsonPath
		}
	}

	if d.HasChange("smb_acl.0.home_dir_path") {
		update = true
		smbAclHomeDirPathJsonPath, err := jsonpath.Get("$[0].home_dir_path", d.Get("smb_acl"))
		if err == nil {
			request["HomeDirPath"] = smbAclHomeDirPathJsonPath
		}
	}

	if v, ok := d.GetOk("keytab_md5"); ok {
		request["KeytabMd5"] = v
	}
	if d.HasChange("smb_acl.0.enable_anonymous_access") {
		update = true
		smbAclEnableAnonymousAccessJsonPath, err := jsonpath.Get("$[0].enable_anonymous_access", d.Get("smb_acl"))
		if err == nil {
			request["EnableAnonymousAccess"] = smbAclEnableAnonymousAccessJsonPath
		}
	}

	if d.HasChange("smb_acl.0.encrypt_data") {
		update = true
		smbAclEncryptDataJsonPath, err := jsonpath.Get("$[0].encrypt_data", d.Get("smb_acl"))
		if err == nil {
			request["EncryptData"] = smbAclEncryptDataJsonPath
		}
	}

	if d.HasChange("smb_acl.0.reject_unencrypted_access") {
		update = true
		smbAclRejectUnencryptedAccessJsonPath, err := jsonpath.Get("$[0].reject_unencrypted_access", d.Get("smb_acl"))
		if err == nil {
			request["RejectUnencryptedAccess"] = smbAclRejectUnencryptedAccessJsonPath
		}
	}

	if v, ok := d.GetOk("keytab"); ok {
		request["Keytab"] = v
	}
	if update && enableModifySmbAcl1 {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		service := NewService(client)
		if err := service.SetResourceTags(d, "filesystem"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceFileSystemRead(d, meta)
}

func resourceFileSystemDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	service := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, service.DescribeAsyncFileSystemStateRefreshFunc(d, response, "$.FileSystems", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return nil
}
