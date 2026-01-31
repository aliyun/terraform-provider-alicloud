// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNasFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNasFileSystemCreate,
		Read:   resourceAliCloudNasFileSystemRead,
		Update: resourceAliCloudNasFileSystemUpdate,
		Delete: resourceAliCloudNasFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
				ValidateFunc: IntInSlice([]int{0, 1, 2}),
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"standard", "extreme", "cpfs", "cpfsse"}, false),
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
				ValidateFunc: StringInSlice([]string{"NFS", "SMB", "cpfs"}, false),
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
							ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
						},
						"reserved_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 180),
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"Performance", "Capacity", "standard", "advance", "advance_100", "advance_200", "Premium", "economic"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceAliCloudNasFileSystemCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("redundancy_type"); ok {
		request["RedundancyType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
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
		redundancyVSwitchIdsMapsArray := convertToInterfaceArray(v)

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_file_system", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasServiceV2.NasFileSystemStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNasFileSystemUpdate(d, meta)
}

func resourceAliCloudNasFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasServiceV2 := NasServiceV2{client}

	objectRaw, err := nasServiceV2.DescribeNasFileSystem(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_file_system DescribeNasFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	redundancyVSwitchIdRaw, _ := jsonpath.Get("$.RedundancyVSwitchIds.RedundancyVSwitchId", objectRaw)
	d.Set("redundancy_vswitch_ids", redundancyVSwitchIdRaw)

	checkValue00 := d.Get("file_system_type")
	checkValue01 := d.Get("protocol_type")
	if (checkValue00 == "standard") && (checkValue01 == "SMB") {
		objectRaw, err = nasServiceV2.DescribeFileSystemDescribeSmbAcl(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
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
	if checkValue00 == "standard" {
		objectRaw, err = nasServiceV2.DescribeFileSystemGetRecycleBinAttribute(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
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
	objectRaw, err = nasServiceV2.DescribeFileSystemListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	checkValue00 = d.Get("file_system_type")
	checkValue01 = d.Get("protocol_type")
	if (checkValue00 == "standard") && (checkValue01 == "NFS") {
		objectRaw, err = nasServiceV2.DescribeFileSystemDescribeNfsAcl(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
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

func resourceAliCloudNasFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	nasServiceV2 := NasServiceV2{client}
	objectRaw, _ := nasServiceV2.DescribeNasFileSystem(d.Id())

	if d.HasChange("recycle_bin.0.status") {
		objectRaw1, _ := nasServiceV2.DescribeFileSystemGetRecycleBinAttribute(d.Id())
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcGet("NAS", "2017-06-26", action, query, request)
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

		}
	}

	if d.HasChange("nfs_acl.0.enabled") {
		objectRaw2, _ := nasServiceV2.DescribeFileSystemDescribeNfsAcl(d.Id())
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

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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

		}
	}

	if d.HasChange("smb_acl.0.enabled") {
		var err error
		objectRaw3, err := nasServiceV2.DescribeFileSystemDescribeSmbAcl(d.Id())
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
				return WrapError(err)
			}
			request["Options"] = string(optionsJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
	}
	update = false
	objectRaw, _ = nasServiceV2.DescribeNasFileSystem(d.Id())
	enableUpgradeFileSystem1 := false
	checkValue00 := objectRaw["FileSystemType"]
	if !(checkValue00 == "standard") {
		enableUpgradeFileSystem1 = true
	}
	action = "UpgradeFileSystem"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("capacity") {
		update = true
	}
	request["Capacity"] = d.Get("capacity")
	request["DryRun"] = "false"
	if update && enableUpgradeFileSystem1 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
		nasServiceV2 := NasServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasServiceV2.NasFileSystemStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
	}
	update = false
	objectRaw, _ = nasServiceV2.DescribeNasFileSystem(d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcGet("NAS", "2017-06-26", action, query, request)
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
	}
	update = false
	objectRaw, _ = nasServiceV2.DescribeNasFileSystem(d.Id())
	objectRaw3, err := nasServiceV2.DescribeFileSystemDescribeSmbAcl(d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
	}

	if d.HasChange("tags") {
		nasServiceV2 := NasServiceV2{client}
		if err := nasServiceV2.SetResourceTags(d, "filesystem"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudNasFileSystemRead(d, meta)
}

func resourceAliCloudNasFileSystemDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasServiceV2.DescribeAsyncNasFileSystemStateRefreshFunc(d, response, "$.FileSystems", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
