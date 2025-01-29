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
	query["FileSystemId"] = d.Get("file_system_id")
	query["InputRegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_vsc_mount_point", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["FileSystemId"], response["MountPointId"]))

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

	instances1Raw := objectRaw["Instances"]
	instancesMaps := make([]map[string]interface{}, 0)
	if instances1Raw != nil {
		for _, instancesChild1Raw := range instances1Raw.([]interface{}) {
			instancesMap := make(map[string]interface{})
			instancesChild1Raw := instancesChild1Raw.(map[string]interface{})
			instancesMap["instance_id"] = instancesChild1Raw["InstanceId"]
			instancesMap["status"] = instancesChild1Raw["Status"]

			vscs1Raw := instancesChild1Raw["Vscs"]
			vscsMaps := make([]map[string]interface{}, 0)
			if vscs1Raw != nil {
				for _, vscsChild1Raw := range vscs1Raw.([]interface{}) {
					vscsMap := make(map[string]interface{})
					vscsChild1Raw := vscsChild1Raw.(map[string]interface{})
					vscsMap["vsc_id"] = vscsChild1Raw["VscId"]
					vscsMap["vsc_status"] = vscsChild1Raw["VscStatus"]
					vscsMap["vsc_type"] = vscsChild1Raw["VscType"]

					vscsMaps = append(vscsMaps, vscsMap)
				}
			}
			instancesMap["vscs"] = vscsMaps
			instancesMaps = append(instancesMaps, instancesMap)
		}
	}
	d.Set("instances", instancesMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("file_system_id", parts[0])
	d.Set("mount_point_id", parts[1])

	return nil
}

func resourceAliCloudDfsVscMountPointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyVscMountPoint"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = parts[0]
	query["MountPointId"] = parts[1]
	query["InputRegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

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
	query["MountPointId"] = parts[1]
	query["FileSystemId"] = parts[0]
	query["InputRegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

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

	return nil
}
