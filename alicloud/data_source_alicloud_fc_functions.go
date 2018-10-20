package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudFcFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcFunctionsRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"functions": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_checksum": {
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

func dataSourceAlicloudFcFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)

	var ids []string
	var functionMappings []map[string]interface{}
	nextToken := ""
	for {
		args := fc.NewListFunctionsInput(serviceName)
		if nextToken != "" {
			args.NextToken = &nextToken
		}

		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListFunctions(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*fc.ListFunctionsOutput)

		if resp.Functions == nil || len(resp.Functions) < 1 {
			break
		}

		for _, function := range resp.Functions {
			mapping := map[string]interface{}{
				"id":                     *function.FunctionID,
				"name":                   *function.FunctionName,
				"description":            *function.Description,
				"runtime":                *function.Runtime,
				"handler":                *function.Handler,
				"timeout":                *function.Timeout,
				"memory_size":            *function.MemorySize,
				"code_size":              *function.CodeSize,
				"code_checksum":          *function.CodeChecksum,
				"creation_time":          *function.CreatedTime,
				"last_modification_time": *function.LastModifiedTime,
			}

			functionMappings = append(functionMappings, mapping)
			log.Printf("[DEBUG] alicloud_fc_functions - adding function mapping: %v", mapping)
			ids = append(ids, *function.FunctionID)
		}

		nextToken = ""
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		}
		if nextToken == "" {
			break
		}
	}

	var filteredFunctionMappings []map[string]interface{}
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, mapping := range functionMappings {
			if r != nil && !r.MatchString(mapping["name"].(string)) {
				continue
			}
			filteredFunctionMappings = append(filteredFunctionMappings, mapping)
		}
	} else {
		filteredFunctionMappings = functionMappings
	}

	if len(filteredFunctionMappings) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	log.Printf("[DEBUG] alicloud_fc_functions - Functions found: %#v", filteredFunctionMappings)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("functions", filteredFunctionMappings); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), filteredFunctionMappings)
	}
	return nil
}
