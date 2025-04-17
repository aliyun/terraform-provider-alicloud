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
				ValidateFunc: StringInSlice([]string{"standard", "extreme", "cpfs"}, false),
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
				ValidateFunc: StringInSlice([]string{"Performance", "Capacity", "standard", "advance", "advance_100", "advance_200", "Premium"}, false),
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

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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
			if NeedRetry(err) || IsExpectedErrors(err, []string{"User.InDebt"}) {
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

	if d.HasChange("recycle_bin.0.status") {
		var err error

		target := d.Get("recycle_bin.0.status").(string)
		enableEnableRecycleBin := false
		checkValue01 := d.Get("file_system_type")
		if checkValue01 == "standard" {
			enableEnableRecycleBin = true
		}
		if enableEnableRecycleBin && target == "Enable" {
			action := "EnableRecycleBin"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["FileSystemId"] = d.Id()

			if v, ok := d.GetOkExists("recycle_bin"); ok {
				jsonPathResult, err := jsonpath.Get("$[0].reserved_days", v)
				if err == nil && jsonPathResult != "" {
					request["ReservedDays"] = jsonPathResult
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
		enableDisableAndCleanRecycleBin := false
		checkValue01 = d.Get("file_system_type")
		if checkValue01 == "standard" {
			enableDisableAndCleanRecycleBin = true
		}
		if enableDisableAndCleanRecycleBin && target == "Disable" {
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
		var err error

		target := d.Get("nfs_acl.0.enabled").(bool)
		enableEnableNfsAcl := false
		checkValue00 := d.Get("file_system_type")
		checkValue02 := d.Get("protocol_type")
		if (checkValue00 == "standard") && (checkValue02 == "NFS") {
			enableEnableNfsAcl = true
		}
		if enableEnableNfsAcl && target == true {
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
		enableDisableNfsAcl := false
		checkValue00 = d.Get("file_system_type")
		checkValue02 = d.Get("protocol_type")
		if (checkValue00 == "standard") && (checkValue02 == "NFS") {
			enableDisableNfsAcl = true
		}
		if enableDisableNfsAcl && target == false {
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

		target := d.Get("smb_acl.0.enabled").(bool)
		enableEnableSmbAcl := false
		checkValue00 := d.Get("file_system_type")
		checkValue01 := d.Get("protocol_type")
		if (checkValue00 == "standard") && (checkValue01 == "SMB") {
			enableEnableSmbAcl = true
		}
		if enableEnableSmbAcl && target == true {
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
		enableDisableSmbAcl := false
		checkValue00 = d.Get("file_system_type")
		checkValue01 = d.Get("protocol_type")
		if (checkValue00 == "standard") && (checkValue01 == "SMB") {
			enableDisableSmbAcl = true
		}
		if enableDisableSmbAcl && target == false {
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
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("options"); v != nil {
			enableOplock1, _ := jsonpath.Get("$[0].enable_oplock", v)
			if enableOplock1 != nil && (d.HasChange("options.0.enable_oplock") || enableOplock1 != "") {
				objectDataLocalMap["EnableOplock"] = enableOplock1
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["Options"] = string(objectDataLocalMapJson)
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
	enableUpgradeFileSystem := false
	checkValue00 := d.Get("file_system_type")
	if !(checkValue00 == "standard") {
		enableUpgradeFileSystem = true
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
	if update && enableUpgradeFileSystem {
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
	enableUpdateRecycleBinAttribute := false
	checkValue00 = d.Get("recycle_bin.0.status")
	checkValue01 := d.Get("file_system_type")
	if (checkValue00 == "Enable") && (checkValue01 == "standard") {
		enableUpdateRecycleBinAttribute = true
	}
	action = "UpdateRecycleBinAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = d.Id()

	if d.HasChange("recycle_bin.0.reserved_days") {
		update = true
	}
	if v, ok := d.GetOk("recycle_bin.0.reserved_days"); ok {
		query["ReservedDays"] = strconv.Itoa(v.(int))
	}

	if update && enableUpdateRecycleBinAttribute {
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
	enableModifySmbAcl := false
	checkValue00 = d.Get("file_system_type")
	checkValue01 = d.Get("protocol_type")
	checkValue02 := d.Get("smb_acl.0.enabled")
	if (checkValue00 == "standard") && (checkValue01 == "SMB") && (checkValue02 == true) {
		enableModifySmbAcl = true
	}
	action = "ModifySmbAcl"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	if d.HasChange("smb_acl.0.super_admin_sid") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].super_admin_sid", d.Get("smb_acl"))
		if err == nil {
			request["SuperAdminSid"] = jsonPathResult
		}
	}

	if d.HasChange("smb_acl.0.home_dir_path") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].home_dir_path", d.Get("smb_acl"))
		if err == nil {
			request["HomeDirPath"] = jsonPathResult1
		}
	}

	if v, ok := d.GetOk("keytab_md5"); ok {
		request["KeytabMd5"] = v
	}
	if d.HasChange("smb_acl.0.enable_anonymous_access") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].enable_anonymous_access", d.Get("smb_acl"))
		if err == nil {
			request["EnableAnonymousAccess"] = jsonPathResult3
		}
	}

	if d.HasChange("smb_acl.0.encrypt_data") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].encrypt_data", d.Get("smb_acl"))
		if err == nil {
			request["EncryptData"] = jsonPathResult4
		}
	}

	if d.HasChange("smb_acl.0.reject_unencrypted_access") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$[0].reject_unencrypted_access", d.Get("smb_acl"))
		if err == nil {
			request["RejectUnencryptedAccess"] = jsonPathResult5
		}
	}

	if v, ok := d.GetOk("keytab"); ok {
		request["Keytab"] = v
	}
	if update && enableModifySmbAcl {
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
