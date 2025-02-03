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

func resourceAlicloudNasMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasMountTargetCreate,
		Read:   resourceAlicloudNasMountTargetRead,
		Update: resourceAlicloudNasMountTargetUpdate,
		Delete: resourceAlicloudNasMountTargetDelete,
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
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Vpc", "Classic"}, false),
				Computed:     true,
			},
			"mount_target_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudNasMountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMountTarget"
	request := make(map[string]interface{})
	var err error
	request["FileSystemId"] = d.Get("file_system_id")
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	} else {
		request["NetworkType"] = string(Vpc)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
		if v, ok := d.GetOk("vswitch_id"); ok {
			request["VSwitchId"] = v
		}
	} else {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.InvalidState"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_mount_target", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["FileSystemId"], ":", response["MountTargetDomain"]))
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{"Pending"}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, nasService.NasMountTargetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudNasMountTargetUpdate(d, meta)
}
func resourceAlicloudNasMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	object, err := nasService.DescribeNasMountTarget(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_mount_target nasService.DescribeNasMountTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("file_system_id", parts[0])
	d.Set("access_group_name", object["AccessGroup"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VswId"])
	d.Set("mount_target_domain", object["MountTargetDomain"])
	d.Set("network_type", object["NetworkType"])

	return nil
}
func resourceAlicloudNasMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FileSystemId":      parts[0],
		"MountTargetDomain": parts[1],
	}
	if !d.IsNewResource() && d.HasChange("access_group_name") {
		update = true
		request["AccessGroupName"] = d.Get("access_group_name")
	}
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}
	if update {
		action := "ModifyMountTarget"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.InvalidState"}) {
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
		nasService := NasService{client}
		stateConf := BuildStateConf([]string{"Pending"}, []string{"Active", "Inactive"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, nasService.NasMountTargetStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudNasMountTargetRead(d, meta)
}
func resourceAlicloudNasMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteMountTarget"
	var response map[string]interface{}
	request := map[string]interface{}{
		"FileSystemId":      parts[0],
		"MountTargetDomain": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"VolumeStatusForbidOperation", "OperationDenied.InvalidState"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, nasService.NasMountTargetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
