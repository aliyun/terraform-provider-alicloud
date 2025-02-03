// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNasAccessPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNasAccessPointCreate,
		Read:   resourceAliCloudNasAccessPointRead,
		Update: resourceAliCloudNasAccessPointUpdate,
		Delete: resourceAliCloudNasAccessPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_point_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_point_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_ram": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"posix_user": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"posix_group_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"posix_user_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"posix_secondary_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"root_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"root_path_permission": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner_user_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"owner_group_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudNasAccessPointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["FileSystemId"] = d.Get("file_system_id")

	request["AccessGroup"] = d.Get("access_group")
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("access_point_name"); ok {
		request["AccessPointName"] = v
	}
	if v, ok := d.GetOkExists("enabled_ram"); ok {
		request["EnabledRam"] = v
	}
	if v, ok := d.GetOk("root_path_permission"); ok {
		jsonPathResult4, err := jsonpath.Get("$[0].owner_user_id", v)
		if err == nil && jsonPathResult4 != "" {
			request["OwnerUserId"] = jsonPathResult4
		}
	}
	if v, ok := d.GetOk("root_path_permission"); ok {
		jsonPathResult5, err := jsonpath.Get("$[0].owner_group_id", v)
		if err == nil && jsonPathResult5 != "" {
			request["OwnerGroupId"] = jsonPathResult5
		}
	}
	if v, ok := d.GetOk("root_path_permission"); ok {
		jsonPathResult6, err := jsonpath.Get("$[0].permission", v)
		if err == nil && jsonPathResult6 != "" {
			request["Permission"] = jsonPathResult6
		}
	}
	if v, ok := d.GetOk("posix_user"); ok {
		jsonPathResult7, err := jsonpath.Get("$[0].posix_user_id", v)
		if err == nil && jsonPathResult7 != "" {
			request["PosixUserId"] = jsonPathResult7
		}
	}
	if v, ok := d.GetOk("posix_user"); ok {
		jsonPathResult8, err := jsonpath.Get("$[0].posix_group_id", v)
		if err == nil && jsonPathResult8 != "" {
			request["PosixGroupId"] = jsonPathResult8
		}
	}
	request["VswId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("root_path"); ok {
		request["RootDirectory"] = v
	}
	if v, ok := d.GetOk("posix_user"); ok {
		jsonPathResult11, err := jsonpath.Get("$[0].posix_secondary_group_ids", v)
		if err == nil && jsonPathResult11 != "" {
			request["PosixSecondaryGroupIds"] = jsonPathResult11
		}
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_point", action, AlibabaCloudSdkGoERROR)
	}

	AccessPointAccessPointId, _ := jsonpath.Get("$.AccessPoint.AccessPointId", response)
	d.SetId(fmt.Sprintf("%v:%v", query["FileSystemId"], AccessPointAccessPointId))

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasServiceV2.NasAccessPointStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNasAccessPointRead(d, meta)
}

func resourceAliCloudNasAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasServiceV2 := NasServiceV2{client}

	objectRaw, err := nasServiceV2.DescribeNasAccessPoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_access_point DescribeNasAccessPoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_group", objectRaw["AccessGroup"])
	d.Set("access_point_name", objectRaw["AccessPointName"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enabled_ram", objectRaw["EnabledRam"])
	d.Set("root_path", objectRaw["RootPath"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("access_point_id", objectRaw["AccessPointId"])
	d.Set("file_system_id", objectRaw["FileSystemId"])

	posixUserMaps := make([]map[string]interface{}, 0)
	posixUserMap := make(map[string]interface{})
	posixUser1Raw := make(map[string]interface{})
	if objectRaw["PosixUser"] != nil {
		posixUser1Raw = objectRaw["PosixUser"].(map[string]interface{})
	}
	if len(posixUser1Raw) > 0 {
		posixUserMap["posix_group_id"] = posixUser1Raw["PosixGroupId"]
		posixUserMap["posix_user_id"] = posixUser1Raw["PosixUserId"]

		posixSecondaryGroupIds1Raw := make([]interface{}, 0)
		if posixUser1Raw["PosixSecondaryGroupIds"] != nil {
			posixSecondaryGroupIds1Raw = posixUser1Raw["PosixSecondaryGroupIds"].([]interface{})
		}

		posixUserMap["posix_secondary_group_ids"] = posixSecondaryGroupIds1Raw
		posixUserMaps = append(posixUserMaps, posixUserMap)
	}
	d.Set("posix_user", posixUserMaps)
	rootPathPermissionMaps := make([]map[string]interface{}, 0)
	rootPathPermissionMap := make(map[string]interface{})
	rootPathPermission1Raw := make(map[string]interface{})
	if objectRaw["RootPathPermission"] != nil {
		rootPathPermission1Raw = objectRaw["RootPathPermission"].(map[string]interface{})
	}
	if len(rootPathPermission1Raw) > 0 {
		rootPathPermissionMap["owner_group_id"] = rootPathPermission1Raw["OwnerGroupId"]
		rootPathPermissionMap["owner_user_id"] = rootPathPermission1Raw["OwnerUserId"]
		rootPathPermissionMap["permission"] = rootPathPermission1Raw["Permission"]

		rootPathPermissionMaps = append(rootPathPermissionMaps, rootPathPermissionMap)
	}
	d.Set("root_path_permission", rootPathPermissionMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("file_system_id", parts[0])
	d.Set("access_point_id", parts[1])

	return nil
}

func resourceAliCloudNasAccessPointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyAccessPoint"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = parts[0]
	query["AccessPointId"] = parts[1]
	if d.HasChange("access_point_name") {
		update = true
		request["AccessPointName"] = d.Get("access_point_name")
	}

	if d.HasChange("access_group") {
		update = true
	}
	request["AccessGroup"] = d.Get("access_group")
	if d.HasChange("enabled_ram") {
		update = true
		request["EnabledRam"] = d.Get("enabled_ram")
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudNasAccessPointRead(d, meta)
}

func resourceAliCloudNasAccessPointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccessPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["FileSystemId"] = parts[0]
	query["AccessPointId"] = parts[1]

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasServiceV2.NasAccessPointStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
