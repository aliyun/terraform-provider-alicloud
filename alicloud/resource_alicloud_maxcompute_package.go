package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMaxComputePackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputePackageCreate,
		Read:   resourceAliCloudMaxComputePackageRead,
		Update: resourceAliCloudMaxComputePackageUpdate,
		Delete: resourceAliCloudMaxComputePackageDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid resource id, expected format <project_name>:<package_name>, got %q", d.Id())
				}
				if err := d.Set("project_name", parts[0]); err != nil {
					return nil, err
				}
				if err := d.Set("package_name", parts[1]); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_install": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudMaxComputePackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName := d.Get("project_name").(string)
	packageName := d.Get("package_name").(string)
	action := fmt.Sprintf("/api/v1/projects/%s/packages", projectName)
	query := make(map[string]*string)
	if v, ok := d.GetOkExists("is_install"); ok {
		query["isInstall"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	// The package name is sent as the raw request body when creating the package.
	// Rebuild the reader on every retry attempt so the body is not exhausted.
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("MaxCompute", "2022-01-04", action, query, nil, strings.NewReader(packageName), true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, packageName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_package", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s:%s", projectName, packageName))
	// Body is update-only: if the user set it in config, send it via
	// UpdatePackage right after creation.
	return resourceAliCloudMaxComputePackageUpdate(d, meta)
}

func resourceAliCloudMaxComputePackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}
	objectRaw, err := maxComputeServiceV2.DescribeMaxComputePackage(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_maxcompute_package DescribeMaxComputePackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := objectRaw["name"].(string); ok && v != "" {
		if err := d.Set("package_name", v); err != nil {
			return err
		}
	}
	return nil
}

func resourceAliCloudMaxComputePackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName, packageName, err := parseMaxComputePackageId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	action := fmt.Sprintf("/api/v1/projects/%s/packages/%s", projectName, packageName)
	query := make(map[string]*string)
	var response map[string]interface{}
	// The body field carries the package content and is only sent on update.
	// Only send the body when it actually changed.
	if d.HasChange("body") {
		bodyContent := d.Get("body").(string)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, strings.NewReader(bodyContent), true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, bodyContent)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAliCloudMaxComputePackageRead(d, meta)
}

func resourceAliCloudMaxComputePackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName, packageName, err := parseMaxComputePackageId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	action := fmt.Sprintf("/api/v1/projects/%s/packages/%s", projectName, packageName)
	query := make(map[string]*string)
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("MaxCompute", "2022-01-04", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func parseMaxComputePackageId(id string) (projectName, packageName string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id, expected format <project_name>:<package_name>, got %q", id)
	}
	return parts[0], parts[1], nil
}
