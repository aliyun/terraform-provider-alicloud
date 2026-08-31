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

func resourceAliCloudNasMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNasMountTargetCreate,
		Read:   resourceAliCloudNasMountTargetRead,
		Update: resourceAliCloudNasMountTargetUpdate,
		Delete: resourceAliCloudNasMountTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dual_stack": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mount_target_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudNasMountTargetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateMountTarget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("file_system_id"); ok {
		request["FileSystemId"] = v
	}

	if v, ok := d.GetOkExists("dual_stack"); ok {
		request["EnableIpv6"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["AccessGroupName"] = d.Get("access_group_name")
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	request["NetworkType"] = string(Vpc)
	if _, ok := d.GetOk("vpc_id"); !ok {
		vswitchId := Trim(d.Get("vswitch_id").(string))
		if vswitchId != "" {
			vpcService := VpcService{client}
			vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
			if err != nil {
				return WrapError(err)
			}
			request["VpcId"] = vsw["VpcId"]
			request["VSwitchId"] = vswitchId
		}
	}

	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.InvalidState"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_mount_target", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["FileSystemId"], response["MountTargetDomain"]))

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasServiceV2.NasMountTargetStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNasMountTargetUpdate(d, meta)
}

func resourceAliCloudNasMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasServiceV2 := NasServiceV2{client}

	objectRaw, err := nasServiceV2.DescribeNasMountTarget(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_mount_target DescribeNasMountTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_group_name", objectRaw["AccessGroup"])
	d.Set("network_type", objectRaw["NetworkType"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VswId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("mount_target_domain", objectRaw["MountTargetDomain"])

	parts := strings.Split(d.Id(), ":")
	d.Set("file_system_id", parts[0])

	return nil
}

func resourceAliCloudNasMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyMountTarget"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["MountTargetDomain"] = parts[1]
	request["FileSystemId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("access_group_name") {
		update = true
	}
	request["AccessGroupName"] = d.Get("access_group_name")
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.InvalidState"}) {
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
		stateConf := BuildStateConf([]string{}, []string{"[]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasServiceV2.DescribeAsyncNasMountTargetStateRefreshFunc(d, response, "$.FileSystems.FileSystem[*].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	return resourceAliCloudNasMountTargetRead(d, meta)
}

func resourceAliCloudNasMountTargetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteMountTarget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["MountTargetDomain"] = parts[1]
	request["FileSystemId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

		if err != nil || IsExpectedErrors(err, []string{"VolumeStatusForbidOperation", "OperationDenied.InvalidState"}) {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil || IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[]"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasServiceV2.DescribeAsyncNasMountTargetStateRefreshFunc(d, response, "$.MountTargets.MountTarget[*]", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
