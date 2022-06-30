package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcdRamDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdRamDirectoryCreate,
		Read:   resourceAlicloudEcdRamDirectoryRead,
		Delete: resourceAlicloudEcdRamDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"desktop_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "INTERNET", "ANY"}, false),
			},
			"enable_admin_access": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ram_directory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_:-]{1,254}$`), "It must be `2` to `255` characters in length, The name must start with a letter, and can contain letters, digits, colons (:), underscores (_), and hyphens (-).")),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcdRamDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRAMDirectory"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("desktop_access_type"); ok {
		request["DesktopAccessType"] = v
	}
	if v, ok := d.GetOkExists("enable_admin_access"); ok {
		request["EnableAdminAccess"] = v
	}
	if v, ok := d.GetOkExists("enable_internet_access"); ok {
		request["EnableInternetAccess"] = v
	}
	request["DirectoryName"] = d.Get("ram_directory_name")
	request["RegionId"] = client.RegionId
	request["VSwitchId"] = d.Get("vswitch_ids")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_ram_directory", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DirectoryId"]))

	return resourceAlicloudEcdRamDirectoryRead(d, meta)
}
func resourceAlicloudEcdRamDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdRamDirectory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_ram_directory ecdService.DescribeEcdRamDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("desktop_access_type", object["DesktopAccessType"])
	d.Set("enable_admin_access", object["EnableAdminAccess"])
	d.Set("enable_internet_access", object["EnableInternetAccess"])
	d.Set("ram_directory_name", object["Name"])
	d.Set("status", object["Status"])
	d.Set("vswitch_ids", object["VSwitchIds"])
	return nil
}
func resourceAlicloudEcdRamDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDirectories"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": []string{d.Id()},
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
