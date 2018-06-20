package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identify_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  sls.MachineIDTypeIP,
				ValidateFunc: validateAllowedStringValue([]string{
					string(sls.MachineIDTypeIP),
					string(sls.MachineIDTypeUserDefined),
				}),
			},
			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"identify_list": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MinItems: 1,
			},
		},
	}
}

func resourceAlicloudLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {

	if err := meta.(*AliyunClient).logconn.CreateMachineGroup(d.Get("project").(string), &sls.MachineGroup{
		Name:          d.Get("name").(string),
		MachineIDType: d.Get("identify_type").(string),
		MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
		Attribute: sls.MachinGroupAttribute{
			TopicName: d.Get("topic").(string),
		},
	}); err != nil {
		return fmt.Errorf("CreateLogMachineGroup got an error: %#v.", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), COLON_SEPARATED)

	group, err := meta.(*AliyunClient).DescribeLogMachineGroup(split[0], split[1])
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
	split := strings.Split(d.Id(), COLON_SEPARATED)
	d.Partial(true)

	update := false
	if d.HasChange("identify_type") {
		update = true
		d.SetPartial("identify_type")
	}
	if d.HasChange("identify_list") {
		update = true
		d.SetPartial("identify_list")
	}
	if d.HasChange("topic") {
		update = true
		d.SetPartial("topic")
	}

	if update {
		if err := meta.(*AliyunClient).logconn.UpdateMachineGroup(split[0], &sls.MachineGroup{
			Name:          split[1],
			MachineIDType: d.Get("identify_type").(string),
			MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
			Attribute: sls.MachinGroupAttribute{
				TopicName: d.Get("topic").(string),
			},
		}); err != nil {
			return fmt.Errorf("UpdateLogMachineGroup %s got an error: %#v.", split[1], err)
		}
	}
	d.Partial(false)

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	split := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		if err := client.logconn.DeleteMachineGroup(split[0], split[1]); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting log machine group %s got an error: %#v", split[1], err))
		}

		if _, err := client.DescribeLogMachineGroup(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Deleting log machine group %s timeout.", split[1]))
	})
}
