package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudFcTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcTriggersRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invocation_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcTriggersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)
	functionName := d.Get("function_name").(string)

	var ids []string
	var triggerMappings []map[string]interface{}
	nextToken := ""
	for {
		args := fc.NewListTriggersInput(serviceName, functionName)
		if nextToken != "" {
			args.NextToken = &nextToken
		}

		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListTriggers(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*fc.ListTriggersOutput)

		if resp.Triggers == nil || len(resp.Triggers) < 1 {
			break
		}

		for _, trigger := range resp.Triggers {
			mapping := map[string]interface{}{
				"id":                     *trigger.TriggerID,
				"name":                   *trigger.TriggerName,
				"source_arn":             *trigger.SourceARN,
				"type":                   *trigger.TriggerType,
				"invocation_role":        *trigger.InvocationRole,
				"config":                 string(trigger.RawTriggerConfig),
				"creation_time":          *trigger.CreatedTime,
				"last_modification_time": *trigger.LastModifiedTime,
			}

			triggerMappings = append(triggerMappings, mapping)

			log.Printf("[DEBUG] alicloud_fc_triggers - adding trigger mapping: %v", mapping)
			ids = append(ids, *trigger.TriggerID)
		}

		nextToken = ""
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		}
		if nextToken == "" {
			break
		}
	}

	var filteredTriggerMappings []map[string]interface{}
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, mapping := range triggerMappings {
			if r != nil && !r.MatchString(mapping["name"].(string)) {
				continue
			}
			filteredTriggerMappings = append(filteredTriggerMappings, mapping)
		}
	} else {
		filteredTriggerMappings = triggerMappings
	}

	if len(filteredTriggerMappings) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	log.Printf("[DEBUG] alicloud_fc_triggers - Triggers found: %#v", triggerMappings)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("triggers", triggerMappings); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), triggerMappings)
	}
	return nil
}
