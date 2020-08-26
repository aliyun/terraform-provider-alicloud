package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudNasAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAccessGroupCreate,
		Read:   resourceAlicloudNasAccessGroupRead,
		Update: resourceAlicloudNasAccessGroupUpdate,
		Delete: resourceAlicloudNasAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"access_group_name"},
			},
			"access_group_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"type"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"access_group_type"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "extreme"}, false),
				Default:      "standard",
			},
		},
	}
}

func resourceAlicloudNasAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessGroupRequest()
	if v, ok := d.GetOk("access_group_name"); ok {
		request.AccessGroupName = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.AccessGroupName = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "access_group_name" must be set one!`))
	}

	if v, ok := d.GetOk("access_group_type"); ok {
		request.AccessGroupType = v.(string)
	} else if v, ok := d.GetOk("type"); ok {
		request.AccessGroupType = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "type" or "access_group_type" must be set one!`))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("file_system_type"); ok {
		request.FileSystemType = v.(string)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.CreateAccessGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceTimeout", "ServiceUnavailable", "InternalError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*nas.CreateAccessGroupResponse)
		d.SetId(fmt.Sprintf("%v:%v", response.AccessGroupName, request.FileSystemType))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudNasAccessGroupRead(d, meta)
}
func resourceAlicloudNasAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	object, err := nasService.DescribeNasAccessGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("access_group_name", parts[0])
	d.Set("name", parts[0])
	d.Set("file_system_type", parts[1])
	d.Set("access_group_type", object.AccessGroupType)
	d.Set("type", object.AccessGroupType)
	d.Set("description", object.Description)
	return nil
}
func resourceAlicloudNasAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("description") {
		request := nas.CreateModifyAccessGroupRequest()
		request.AccessGroupName = parts[0]
		request.FileSystemType = parts[1]
		request.Description = d.Get("description").(string)
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyAccessGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudNasAccessGroupRead(d, meta)
}
func resourceAlicloudNasAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := nas.CreateDeleteAccessGroupRequest()
	request.AccessGroupName = parts[0]
	request.FileSystemType = parts[1]
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.DeleteAccessGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
