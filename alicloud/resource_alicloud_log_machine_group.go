package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/aliyun-log-go-sdk"
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
	invoker := NewInvoker()
	invoker.AddCatcher(SlsClientTimeoutCatcher)
	if err := invoker.Run(func() error {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateMachineGroup(d.Get("project").(string), &sls.MachineGroup{
				Name:          d.Get("name").(string),
				MachineIDType: d.Get("identify_type").(string),
				MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
				Attribute: sls.MachinGroupAttribute{
					TopicName: d.Get("topic").(string),
				},
			})
		})
		return err
	}); err != nil {
		return fmt.Errorf("CreateLogMachineGroup got an error: %s.", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)

	group, err := logService.DescribeLogMachineGroup(split[0], split[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeLogMachineGroup got an error: %#v.", err)
	}

	d.Set("project", split[0])
	d.Set("name", group.Name)
	d.Set("identify_type", group.MachineIDType)
	d.Set("identify_list", group.MachineIDList)
	d.Set("topic", group.Attribute.TopicName)

	return nil
}

func resourceAlicloudLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("identify_type") || d.HasChange("identify_list") || d.HasChange("topic") {
		split := strings.Split(d.Id(), COLON_SEPARATED)

		client := meta.(*connectivity.AliyunClient)
		invoker := NewInvoker()
		invoker.AddCatcher(SlsClientTimeoutCatcher)
		if err := invoker.Run(func() error {
			_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.UpdateMachineGroup(split[0], &sls.MachineGroup{
					Name:          split[1],
					MachineIDType: d.Get("identify_type").(string),
					MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
					Attribute: sls.MachinGroupAttribute{
						TopicName: d.Get("topic").(string),
					},
				})
			})
			return err
		}); err != nil {
			return fmt.Errorf("UpdateLogMachineGroup %s got an error: %#v.", split[1], err)
		}
	}

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteMachineGroup(split[0], split[1])
		})
		if err != nil {
			if IsExceptedErrors(err, []string{LogClientTimeout}) {
				return resource.RetryableError(fmt.Errorf("Timeout. DeleteMachineGroup %s got an error: %#v", split[1], err))
			}
			return resource.NonRetryableError(fmt.Errorf("DeleteMachineGroup %s got an error: %#v", split[1], err))
		}

		if _, err := logService.DescribeLogMachineGroup(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Deleting log machine group %s timeout.", split[1]))
	})
}
