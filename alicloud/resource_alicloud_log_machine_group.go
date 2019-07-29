package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudLogMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogMachineGroupCreate,
		Read:   resourceAlicloudLogMachineGroupRead,
		Update: resourceAlicloudLogMachineGroupUpdate,
		Delete: resourceAlicloudLogMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identify_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  sls.MachineIDTypeIP,
				ValidateFunc: validateAllowedStringValue([]string{
					string(sls.MachineIDTypeIP),
					string(sls.MachineIDTypeUserDefined),
				}),
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identify_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MinItems: 1,
			},
		},
	}
}

func resourceAlicloudLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateMachineGroup(d.Get("project").(string), &sls.MachineGroup{
				Name:          d.Get("name").(string),
				MachineIDType: d.Get("identify_type").(string),
				MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
				Attribute: sls.MachinGroupAttribute{
					TopicName: d.Get("topic").(string),
				},
			})
		})
		if err != nil {
			if IsExceptedError(err, LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateMachineGroup", raw)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_machine_group", "CreateMachineGroup", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogMachineGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("identify_type", object.MachineIDType)
	d.Set("identify_list", object.MachineIDList)
	d.Set("topic", object.Attribute.TopicName)

	return nil
}

func resourceAlicloudLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("identify_type") || d.HasChange("identify_list") || d.HasChange("topic") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}

		client := meta.(*connectivity.AliyunClient)

		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.UpdateMachineGroup(parts[0], &sls.MachineGroup{
					Name:          parts[1],
					MachineIDType: d.Get("identify_type").(string),
					MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
					Attribute: sls.MachinGroupAttribute{
						TopicName: d.Get("topic").(string),
					},
				})
			})
			if err != nil {
				if IsExceptedError(err, LogClientTimeout) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("UpdateMachineGroup", raw)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateMachineGroup", AliyunLogGoSdkERROR)
		}
	}

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteMachineGroup(parts[0], parts[1])
		})
		if err != nil {
			if IsExceptedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteMachineGroup", raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "ListShards", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogMachineGroup(d.Id(), Deleted, DefaultTimeout))
}
