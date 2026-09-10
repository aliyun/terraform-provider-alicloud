package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOosDefaultPatchBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosDefaultPatchBaselineCreate,
		Read:   resourceAlicloudOosDefaultPatchBaselineRead,
		Delete: resourceAlicloudOosDefaultPatchBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"patch_baseline_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"patch_baseline_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudOosDefaultPatchBaselineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("patch_baseline_name"); ok {
		request["Name"] = v
	}

	var response map[string]interface{}
	action := "RegisterDefaultPatchBaseline"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_default_patch_baseline", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PatchBaseline.Name", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_oos_default_patch_baseline")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudOosDefaultPatchBaselineRead(d, meta)
}

func resourceAlicloudOosDefaultPatchBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}

	object, err := oosService.DescribeOosDefaultPatchBaseline(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_default_patch_baseline oosService.DescribeOosDefaultPatchBaseline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("patch_baseline_id", object["Id"])
	d.Set("patch_baseline_name", object["Name"])

	return nil
}

func resourceAlicloudOosDefaultPatchBaselineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	oosService := OosService{client}
	object, err := oosService.DescribeOosPatchBaseline(d.Id())
	if err != nil {
		return WrapError(err)
	}

	defaultPatchBaselineID := FindDefaultDefaultPatchBaselineIDForOS(object["OperationSystem"].(string))
	if defaultPatchBaselineID == "" {
		return fmt.Errorf("failed to find the default default patch baseline ID")
	}

	request := map[string]interface{}{
		"Name":     defaultPatchBaselineID,
		"RegionId": client.RegionId,
	}

	action := "RegisterDefaultPatchBaseline"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func FindDefaultDefaultPatchBaselineIDForOS(operationSystem string) string {
	switch operationSystem {
	case "AliyunLinux":
		return "ACS-AliyunLinux-DefaultPatchBaseline"
	case "AlmaLinux":
		return "ACS-AlmaLinux-DefaultPatchBaseline"
	case "Anolis":
		return "ACS-Anolis-DefaultPatchBaseline"
	case "CentOS":
		return "ACS-CentOS-DefaultPatchBaseline"
	case "Debian":
		return "ACS-Debian-DefaultPatchBaseline"
	case "RedhatEnterpriseLinux":
		return "ACS-RedhatEnterpriseLinux-DefaultPatchBaseline"
	case "Ubuntu":
		return "ACS-Ubuntu-DefaultPatchBaseline"
	case "Windows":
		return "ACS-Windows-DefaultPatchBaseline"
	}
	return ""
}
