package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudFCTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCTriggerCreate,
		Read:   resourceAlicloudFCTriggerRead,
		Update: resourceAlicloudFCTriggerUpdate,
		Delete: resourceAlicloudFCTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateStringLengthInRange(1, 128),
			},
			"name_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					// uuid is 26 characters, limit the prefix to 229.
					value := v.(string)
					if len(value) > 122 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 102 characters, name is limited to 128", k))
					}
					return
				},
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew: true,
			},

			"config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateJsonString,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{string(fc.TRIGGER_TYPE_HTTP), string(fc.TRIGGER_TYPE_LOG),
					string(fc.TRIGGER_TYPE_OSS), string(fc.TRIGGER_TYPE_TIMER)}),
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	conn := client.fcconn

	serviceName := d.Get("service").(string)
	fcName := d.Get("function").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	var config interface{}
	if err := json.Unmarshal([]byte(d.Get("config").(string)), &config); err != nil {
		return fmt.Errorf("Unmarshalling config got an error: %#v.", err)
	}

	object := fc.TriggerCreateObject{
		TriggerName:    StringPointer(name),
		SourceARN:      StringPointer(d.Get("source_arn").(string)),
		TriggerType:    StringPointer(d.Get("type").(string)),
		InvocationRole: StringPointer(d.Get("role").(string)),
		TriggerConfig:  config,
	}
	input := &fc.CreateTriggerInput{
		ServiceName:         StringPointer(serviceName),
		FunctionName:        StringPointer(fcName),
		TriggerCreateObject: object,
	}
	var trigger *fc.CreateTriggerOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		out, err := conn.CreateTrigger(input)
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied}) {
				return resource.RetryableError(fmt.Errorf("Error creating function compute service got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error creating function compute trigger got an error: %#v", err))
		}
		trigger = out
		return nil

	}); err != nil {
		return err
	}

	if trigger == nil {
		return fmt.Errorf("Creating function compute trigger got a empty response: %#v.", trigger)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", serviceName, COLON_SEPARATED, fcName, COLON_SEPARATED, *trigger.TriggerName))

	return resourceAlicloudFCTriggerRead(d, meta)
}

func resourceAlicloudFCTriggerRead(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}

	client := meta.(*AliyunClient)

	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) < 3 {
		return fmt.Errorf("Invalid resource ID %s. Please check it and try again.", d.Id())
	}
	trigger, err := client.DescribeFcTrigger(split[0], split[1], split[2])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeFCTrigger %s got an error: %#v", d.Id(), err)
	}

	d.Set("service", split[0])
	d.Set("function", split[1])
	d.Set("name", trigger.TriggerName)
	d.Set("role", trigger.InvocationRole)
	d.Set("source_arn", trigger.SourceARN)
	d.Set("config", trigger.TriggerConfig)
	d.Set("type", trigger.TriggerType)
	d.Set("last_modified", trigger.LastModifiedTime)

	return nil
}

func resourceAlicloudFCTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)

	d.Partial(true)
	updateInput := &fc.UpdateTriggerInput{}

	if d.HasChange("role") {
		updateInput.InvocationRole = StringPointer(d.Get("role").(string))
		d.SetPartial("role")
	}
	if d.HasChange("config") {
		var config interface{}
		if err := json.Unmarshal([]byte(d.Get("config").(string)), &config); err != nil {
			return fmt.Errorf("When updating, unmarshalling config got an error: %#v.", err)
		}
		updateInput.TriggerConfig = config
		d.SetPartial("config")
	}

	if updateInput != nil {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		if len(split) < 3 {
			return fmt.Errorf("Invalid resource ID %s. Please check it and try again.", d.Id())
		}
		updateInput.ServiceName = StringPointer(split[0])
		updateInput.FunctionName = StringPointer(split[1])
		updateInput.TriggerName = StringPointer(split[2])

		if _, err := client.fcconn.UpdateTrigger(updateInput); err != nil {
			return fmt.Errorf("UpdateTrigger %s got an error: %#v.", d.Id(), err)
		}
	}

	d.Partial(false)
	return resourceAlicloudFCTriggerRead(d, meta)
}

func resourceAlicloudFCTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) < 3 {
		return fmt.Errorf("Invalid resource ID %s. Please check it and try again.", d.Id())
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.fcconn.DeleteTrigger(&fc.DeleteTriggerInput{
			ServiceName:  StringPointer(split[0]),
			FunctionName: StringPointer(split[1]),
			TriggerName:  StringPointer(split[2]),
		}); err != nil {
			if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting trigger got an error: %#v.", err))
		}

		if _, err := client.DescribeFcTrigger(split[0], split[1], split[2]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While deleting function trigger, getting trigger %s got an error: %#v.", d.Id(), err))
		}
		return nil
	})

}
