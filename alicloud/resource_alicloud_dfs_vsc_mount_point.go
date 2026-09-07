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

func resourceAliCloudDfsVscMountPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDfsVscMountPointCreate,
		Read:   resourceAliCloudDfsVscMountPointRead,
		Update: resourceAliCloudDfsVscMountPointUpdate,
		Delete: resourceAliCloudDfsVscMountPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alias_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vscs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vsc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vsc_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vsc_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"mount_point_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDfsVscMountPointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVscMountPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("file_system_id"); ok {
		request["FileSystemId"] = v
	}
	request["InputRegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_vsc_mount_point", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["FileSystemId"], response["MountPointId"]))

	return resourceAliCloudDfsVscMountPointUpdate(d, meta)
}

func resourceAliCloudDfsVscMountPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsServiceV2 := DfsServiceV2{client}

	objectRaw, err := dfsServiceV2.DescribeDfsVscMountPoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_vsc_mount_point DescribeDfsVscMountPoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("mount_point_id", objectRaw["MountPointId"])

	instancesRaw := objectRaw["Instances"]
	instancesMaps := make([]map[string]interface{}, 0)
	if instancesRaw != nil {
		for _, instancesChildRaw := range instancesRaw.([]interface{}) {
			instancesMap := make(map[string]interface{})
			instancesChildRaw := instancesChildRaw.(map[string]interface{})
			instancesMap["instance_id"] = instancesChildRaw["InstanceId"]
			instancesMap["status"] = instancesChildRaw["Status"]

			vscsRaw := instancesChildRaw["Vscs"]
			vscsMaps := make([]map[string]interface{}, 0)
			if vscsRaw != nil {
				for _, vscsChildRaw := range vscsRaw.([]interface{}) {
					vscsMap := make(map[string]interface{})
					vscsChildRaw := vscsChildRaw.(map[string]interface{})
					vscsMap["vsc_id"] = vscsChildRaw["VscId"]
					vscsMap["vsc_status"] = vscsChildRaw["VscStatus"]
					vscsMap["vsc_type"] = vscsChildRaw["VscType"]

					vscsMaps = append(vscsMaps, vscsMap)
				}
			}
			instancesMap["vscs"] = vscsMaps
			instancesMaps = append(instancesMaps, instancesMap)
		}
	}
	if err := d.Set("instances", instancesMaps); err != nil {
		return err
	}

	if objectRaw["MountPointAlias"] != nil {
		if fmt.Sprint(objectRaw["MountPointAlias"]) != "" {
			aliasPrefix := strings.Split(fmt.Sprint(objectRaw["MountPointAlias"]), "--")[0]
			d.Set("alias_prefix", aliasPrefix)
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("file_system_id", parts[0])

	return nil
}

func resourceAliCloudDfsVscMountPointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "BindVscMountPointAlias"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["MountPointId"] = parts[1]
	request["FileSystemId"] = parts[0]
	request["InputRegionId"] = client.RegionId
	if d.HasChange("alias_prefix") {
		update = true
	}
	request["AliasPrefix"] = d.Get("alias_prefix")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)
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
	parts = strings.Split(d.Id(), ":")
	action = "ModifyVscMountPoint"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["MountPointId"] = parts[1]
	request["FileSystemId"] = parts[0]
	request["InputRegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudDfsVscMountPointRead(d, meta)
}

func resourceAliCloudDfsVscMountPointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteVscMountPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["MountPointId"] = parts[1]
	request["FileSystemId"] = parts[0]
	request["InputRegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
