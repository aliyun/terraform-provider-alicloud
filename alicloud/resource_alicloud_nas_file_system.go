// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
				ForceNew:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2}),
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"standard", "extreme", "cpfs"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				ForceNew: true,
				Computed: true,
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
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}
	if v, ok := d.GetOkExists("capacity"); ok {
		request["Capacity"] = v
	}
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["ProtocolType"] = d.Get("protocol_type")
	if v, ok := d.GetOkExists("encrypt_type"); ok {
		request["EncryptType"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	request["DryRun"] = "false"
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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

	if objectRaw["Capacity"] != nil {
		d.Set("capacity", objectRaw["Capacity"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["EncryptType"] != nil {
		d.Set("encrypt_type", objectRaw["EncryptType"])
	}
	if objectRaw["FileSystemType"] != nil {
		d.Set("file_system_type", objectRaw["FileSystemType"])
	}
	if objectRaw["KMSKeyId"] != nil {
		d.Set("kms_key_id", objectRaw["KMSKeyId"])
	}
	if objectRaw["ProtocolType"] != nil {
		d.Set("protocol_type", objectRaw["ProtocolType"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["StorageType"] != nil {
		d.Set("storage_type", objectRaw["StorageType"])
	}
	if objectRaw["QuorumVswId"] != nil {
		d.Set("vswitch_id", objectRaw["QuorumVswId"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}
	if objectRaw["ZoneId"] != nil {
		d.Set("zone_id", objectRaw["ZoneId"])
	}

	checkValue00 := d.Get("file_system_type")
	if checkValue00 == "standard" {

		objectRaw, err = nasServiceV2.DescribeGetRecycleBinAttribute(d.Id())
		if err != nil {
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
	objectRaw, err = nasServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	checkValue01 := d.Get("protocol_type")
	if checkValue00 == "standard" && checkValue01 == "NFS" {

		objectRaw, err = nasServiceV2.DescribeDescribeNfsAcl(d.Id())
		if err != nil {
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
	nasServiceV2 := NasServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyFileSystem"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
	action = "UpgradeFileSystem"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("capacity") {
		update = true
	}
	request["Capacity"] = d.Get("capacity")
	request["DryRun"] = "false"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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

	if d.HasChange("recycle_bin.0.status") {
		objectRaw, err := nasServiceV2.DescribeGetRecycleBinAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("recycle_bin.0.status").(string)

		if objectRaw["Status"].(string) != target {
			if target == "Enable" {
				action := "EnableRecycleBin"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FileSystemId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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

			if target == "Disable" {
				action := "DisableAndCleanRecycleBin"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FileSystemId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcGet("NAS", "2017-06-26", action, request, nil)
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

	}

	if d.HasChange("recycle_bin.0.reserved_days") {
		action = "UpdateRecycleBinAttribute"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["FileSystemId"] = d.Id()
		request["RegionId"] = client.RegionId

		if reservedDays, ok := d.GetOkExists("recycle_bin.0.reserved_days"); ok {
			request["ReservedDays"] = reservedDays
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcGet("NAS", "2017-06-26", action, request, nil)
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

	if d.HasChange("nfs_acl.0.enabled") {
		objectRaw, err := nasServiceV2.DescribeDescribeNfsAcl(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("nfs_acl.0.enabled").(bool)

		if objectRaw["Enabled"].(bool) != target {
			if target == true {
				action := "EnableNfsAcl"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FileSystemId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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

			if target == false {
				action := "DisableNfsAcl"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["FileSystemId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasServiceV2.NasFileSystemStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
