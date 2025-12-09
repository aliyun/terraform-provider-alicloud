// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudOpenApiExplorerApiMcpServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOpenApiExplorerApiMcpServerCreate,
		Read:   resourceAliCloudOpenApiExplorerApiMcpServerRead,
		Update: resourceAliCloudOpenApiExplorerApiMcpServerUpdate,
		Delete: resourceAliCloudOpenApiExplorerApiMcpServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"additional_api_descriptions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_output_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"api_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"const_parameters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"api_override_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.ValidateJsonString,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"product": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"execute_cli_command": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"apis": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
						"selectors": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"assume_role_extra_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"assume_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_assume_role": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_custom_vpc_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instructions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"EN_US", "ZH_CN"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oauth_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prompts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arguments": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"required": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"public_access": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off", "follow"}, false),
			},
			"system_tools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"terraform_tools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"async": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"destroy_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"NEVER", "ALWAYS", "ON_FAILURE"}, false),
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vpc_whitelists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOpenApiExplorerApiMcpServerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/apimcpserver")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("additional_api_descriptions"); ok {
		additionalApiDescriptionsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["apiOverrideJson"] = dataLoopTmp["api_override_json"]
			dataLoopMap["executeCliCommand"] = dataLoopTmp["execute_cli_command"]
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["const_parameters"]
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["constParameters"] = localMaps
			dataLoopMap["apiVersion"] = dataLoopTmp["api_version"]
			dataLoopMap["product"] = dataLoopTmp["product"]
			dataLoopMap["apiName"] = dataLoopTmp["api_name"]
			dataLoopMap["enableOutputSchema"] = dataLoopTmp["enable_output_schema"]
			additionalApiDescriptionsMapsArray = append(additionalApiDescriptionsMapsArray, dataLoopMap)
		}
		request["additionalApiDescriptions"] = additionalApiDescriptionsMapsArray
	}

	if v, ok := d.GetOk("terraform_tools"); ok {
		terraformToolsMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["async"] = dataLoop2Tmp["async"]
			dataLoop2Map["code"] = dataLoop2Tmp["code"]
			dataLoop2Map["description"] = dataLoop2Tmp["description"]
			dataLoop2Map["destroyPolicy"] = dataLoop2Tmp["destroy_policy"]
			dataLoop2Map["name"] = dataLoop2Tmp["name"]
			terraformToolsMapsArray = append(terraformToolsMapsArray, dataLoop2Map)
		}
		request["terraformTools"] = terraformToolsMapsArray
	}

	if v, ok := d.GetOk("apis"); ok {
		apisMapsArray := make([]interface{}, 0)
		for _, dataLoop3 := range convertToInterfaceArray(v) {
			dataLoop3Tmp := dataLoop3.(map[string]interface{})
			dataLoop3Map := make(map[string]interface{})
			dataLoop3Map["apiVersion"] = dataLoop3Tmp["api_version"]
			dataLoop3Map["product"] = dataLoop3Tmp["product"]
			dataLoop3Map["selectors"] = convertToInterfaceArray(dataLoop3Tmp["selectors"])
			apisMapsArray = append(apisMapsArray, dataLoop3Map)
		}
		request["apis"] = apisMapsArray
	}

	if v, ok := d.GetOk("language"); ok {
		request["language"] = v
	}
	if v, ok := d.GetOk("assume_role_extra_policy"); ok {
		request["assumeRoleExtraPolicy"] = v
	}
	if v, ok := d.GetOk("prompts"); ok {
		promptsMapsArray := make([]interface{}, 0)
		for _, dataLoop4 := range convertToInterfaceArray(v) {
			dataLoop4Tmp := dataLoop4.(map[string]interface{})
			dataLoop4Map := make(map[string]interface{})
			localMaps1 := make([]interface{}, 0)
			localData5 := dataLoop4Tmp["arguments"]
			for _, dataLoop5 := range convertToInterfaceArray(localData5) {
				dataLoop5Tmp := dataLoop5.(map[string]interface{})
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["description"] = dataLoop5Tmp["description"]
				dataLoop5Map["required"] = dataLoop5Tmp["required"]
				dataLoop5Map["name"] = dataLoop5Tmp["name"]
				localMaps1 = append(localMaps1, dataLoop5Map)
			}
			dataLoop4Map["arguments"] = localMaps1
			dataLoop4Map["name"] = dataLoop4Tmp["name"]
			dataLoop4Map["content"] = dataLoop4Tmp["content"]
			dataLoop4Map["description"] = dataLoop4Tmp["description"]
			promptsMapsArray = append(promptsMapsArray, dataLoop4Map)
		}
		request["prompts"] = promptsMapsArray
	}

	if v, ok := d.GetOk("vpc_whitelists"); ok {
		vpcWhitelistsMapsArray := convertToInterfaceArray(v)

		request["vpcWhitelists"] = vpcWhitelistsMapsArray
	}

	if v, ok := d.GetOk("system_tools"); ok {
		systemToolsMapsArray := convertToInterfaceArray(v)

		request["systemTools"] = systemToolsMapsArray
	}

	if v, ok := d.GetOkExists("enable_custom_vpc_whitelist"); ok {
		request["enableCustomVpcWhitelist"] = v
	}
	if v, ok := d.GetOk("oauth_client_id"); ok {
		request["oauthClientId"] = v
	}
	if v, ok := d.GetOk("instructions"); ok {
		request["instructions"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("public_access"); ok {
		request["publicAccess"] = v
	}
	request["name"] = d.Get("name")
	if v, ok := d.GetOkExists("enable_assume_role"); ok {
		request["enableAssumeRole"] = v
	}
	if v, ok := d.GetOk("assume_role_name"); ok {
		request["assumeRoleName"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("OpenAPIExplorer", "2024-11-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_open_api_explorer_api_mcp_server", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["id"]))

	return resourceAliCloudOpenApiExplorerApiMcpServerRead(d, meta)
}

func resourceAliCloudOpenApiExplorerApiMcpServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	openApiExplorerServiceV2 := OpenApiExplorerServiceV2{client}

	objectRaw, err := openApiExplorerServiceV2.DescribeOpenApiExplorerApiMcpServer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_open_api_explorer_api_mcp_server DescribeOpenApiExplorerApiMcpServer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("assume_role_extra_policy", objectRaw["assumeRoleExtraPolicy"])
	d.Set("assume_role_name", objectRaw["assumeRoleName"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("enable_assume_role", objectRaw["enableAssumeRole"])
	d.Set("enable_custom_vpc_whitelist", objectRaw["enableCustomVpcWhitelist"])
	d.Set("instructions", objectRaw["instructions"])
	d.Set("language", objectRaw["language"])
	d.Set("name", objectRaw["name"])
	d.Set("oauth_client_id", objectRaw["oauthClientId"])
	d.Set("public_access", objectRaw["publicAccess"])

	additionalApiDescriptionsRaw := objectRaw["additionalApiDescriptions"]
	additionalApiDescriptionsMaps := make([]map[string]interface{}, 0)
	if additionalApiDescriptionsRaw != nil {
		for _, additionalApiDescriptionsChildRaw := range convertToInterfaceArray(additionalApiDescriptionsRaw) {
			additionalApiDescriptionsMap := make(map[string]interface{})
			additionalApiDescriptionsChildRaw := additionalApiDescriptionsChildRaw.(map[string]interface{})
			additionalApiDescriptionsMap["api_name"] = additionalApiDescriptionsChildRaw["apiName"]
			additionalApiDescriptionsMap["api_override_json"] = additionalApiDescriptionsChildRaw["apiOverrideJson"]
			additionalApiDescriptionsMap["api_version"] = additionalApiDescriptionsChildRaw["apiVersion"]
			additionalApiDescriptionsMap["enable_output_schema"] = additionalApiDescriptionsChildRaw["enableOutputSchema"]
			additionalApiDescriptionsMap["execute_cli_command"] = additionalApiDescriptionsChildRaw["executeCliCommand"]
			additionalApiDescriptionsMap["product"] = additionalApiDescriptionsChildRaw["product"]

			constParametersRaw := additionalApiDescriptionsChildRaw["constParameters"]
			constParametersMaps := make([]map[string]interface{}, 0)
			if constParametersRaw != nil {
				for _, constParametersChildRaw := range convertToInterfaceArray(constParametersRaw) {
					constParametersMap := make(map[string]interface{})
					constParametersChildRaw := constParametersChildRaw.(map[string]interface{})
					constParametersMap["key"] = constParametersChildRaw["key"]
					constParametersMap["value"] = constParametersChildRaw["value"]

					constParametersMaps = append(constParametersMaps, constParametersMap)
				}
			}
			additionalApiDescriptionsMap["const_parameters"] = constParametersMaps
			additionalApiDescriptionsMaps = append(additionalApiDescriptionsMaps, additionalApiDescriptionsMap)
		}
	}
	if err := d.Set("additional_api_descriptions", additionalApiDescriptionsMaps); err != nil {
		return err
	}
	apisRaw := objectRaw["apis"]
	apisMaps := make([]map[string]interface{}, 0)
	if apisRaw != nil {
		for _, apisChildRaw := range convertToInterfaceArray(apisRaw) {
			apisMap := make(map[string]interface{})
			apisChildRaw := apisChildRaw.(map[string]interface{})
			apisMap["api_version"] = apisChildRaw["apiVersion"]
			apisMap["product"] = apisChildRaw["product"]

			selectorsRaw := make([]interface{}, 0)
			if apisChildRaw["selectors"] != nil {
				selectorsRaw = convertToInterfaceArray(apisChildRaw["selectors"])
			}

			apisMap["selectors"] = selectorsRaw
			apisMaps = append(apisMaps, apisMap)
		}
	}
	if err := d.Set("apis", apisMaps); err != nil {
		return err
	}
	promptsRaw := objectRaw["prompts"]
	promptsMaps := make([]map[string]interface{}, 0)
	if promptsRaw != nil {
		for _, promptsChildRaw := range convertToInterfaceArray(promptsRaw) {
			promptsMap := make(map[string]interface{})
			promptsChildRaw := promptsChildRaw.(map[string]interface{})
			promptsMap["content"] = promptsChildRaw["content"]
			promptsMap["description"] = promptsChildRaw["description"]
			promptsMap["name"] = promptsChildRaw["name"]

			argumentsRaw := promptsChildRaw["arguments"]
			argumentsMaps := make([]map[string]interface{}, 0)
			if argumentsRaw != nil {
				for _, argumentsChildRaw := range convertToInterfaceArray(argumentsRaw) {
					argumentsMap := make(map[string]interface{})
					argumentsChildRaw := argumentsChildRaw.(map[string]interface{})
					argumentsMap["description"] = argumentsChildRaw["description"]
					argumentsMap["name"] = argumentsChildRaw["name"]
					argumentsMap["required"] = argumentsChildRaw["required"]

					argumentsMaps = append(argumentsMaps, argumentsMap)
				}
			}
			promptsMap["arguments"] = argumentsMaps
			promptsMaps = append(promptsMaps, promptsMap)
		}
	}
	if err := d.Set("prompts", promptsMaps); err != nil {
		return err
	}
	systemToolsRaw := make([]interface{}, 0)
	if objectRaw["systemTools"] != nil {
		systemToolsRaw = convertToInterfaceArray(objectRaw["systemTools"])
	}

	d.Set("system_tools", systemToolsRaw)
	terraformToolsRaw := objectRaw["terraformTools"]
	terraformToolsMaps := make([]map[string]interface{}, 0)
	if terraformToolsRaw != nil {
		for _, terraformToolsChildRaw := range convertToInterfaceArray(terraformToolsRaw) {
			terraformToolsMap := make(map[string]interface{})
			terraformToolsChildRaw := terraformToolsChildRaw.(map[string]interface{})
			terraformToolsMap["async"] = terraformToolsChildRaw["async"]
			terraformToolsMap["code"] = terraformToolsChildRaw["code"]
			terraformToolsMap["description"] = terraformToolsChildRaw["description"]
			terraformToolsMap["destroy_policy"] = terraformToolsChildRaw["destroyPolicy"]
			terraformToolsMap["name"] = terraformToolsChildRaw["name"]

			terraformToolsMaps = append(terraformToolsMaps, terraformToolsMap)
		}
	}
	if err := d.Set("terraform_tools", terraformToolsMaps); err != nil {
		return err
	}
	vpcWhitelistsRaw := make([]interface{}, 0)
	if objectRaw["vpcWhitelists"] != nil {
		vpcWhitelistsRaw = convertToInterfaceArray(objectRaw["vpcWhitelists"])
	}

	d.Set("vpc_whitelists", vpcWhitelistsRaw)

	return nil
}

func resourceAliCloudOpenApiExplorerApiMcpServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/apimcpserver")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["id"] = StringPointer(d.Id())

	if d.HasChange("additional_api_descriptions") {
		update = true
	}
	if v, ok := d.GetOk("additional_api_descriptions"); ok || d.HasChange("additional_api_descriptions") {
		additionalApiDescriptionsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["apiOverrideJson"] = dataLoopTmp["api_override_json"]
			dataLoopMap["executeCliCommand"] = dataLoopTmp["execute_cli_command"]
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["const_parameters"]
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["constParameters"] = localMaps
			dataLoopMap["apiVersion"] = dataLoopTmp["api_version"]
			dataLoopMap["product"] = dataLoopTmp["product"]
			dataLoopMap["apiName"] = dataLoopTmp["api_name"]
			dataLoopMap["enableOutputSchema"] = dataLoopTmp["enable_output_schema"]
			additionalApiDescriptionsMapsArray = append(additionalApiDescriptionsMapsArray, dataLoopMap)
		}
		request["additionalApiDescriptions"] = additionalApiDescriptionsMapsArray
	}

	if d.HasChange("terraform_tools") {
		update = true
	}
	if v, ok := d.GetOk("terraform_tools"); ok || d.HasChange("terraform_tools") {
		terraformToolsMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["async"] = dataLoop2Tmp["async"]
			dataLoop2Map["code"] = dataLoop2Tmp["code"]
			dataLoop2Map["description"] = dataLoop2Tmp["description"]
			dataLoop2Map["destroyPolicy"] = dataLoop2Tmp["destroy_policy"]
			dataLoop2Map["name"] = dataLoop2Tmp["name"]
			terraformToolsMapsArray = append(terraformToolsMapsArray, dataLoop2Map)
		}
		request["terraformTools"] = terraformToolsMapsArray
	}

	if d.HasChange("apis") {
		update = true
	}
	if v, ok := d.GetOk("apis"); ok || d.HasChange("apis") {
		apisMapsArray := make([]interface{}, 0)
		for _, dataLoop3 := range convertToInterfaceArray(v) {
			dataLoop3Tmp := dataLoop3.(map[string]interface{})
			dataLoop3Map := make(map[string]interface{})
			dataLoop3Map["apiVersion"] = dataLoop3Tmp["api_version"]
			dataLoop3Map["product"] = dataLoop3Tmp["product"]
			dataLoop3Map["selectors"] = convertToInterfaceArray(dataLoop3Tmp["selectors"])
			apisMapsArray = append(apisMapsArray, dataLoop3Map)
		}
		request["apis"] = apisMapsArray
	}

	if d.HasChange("language") {
		update = true
	}
	if v, ok := d.GetOk("language"); ok || d.HasChange("language") {
		request["language"] = v
	}
	if d.HasChange("assume_role_extra_policy") {
		update = true
	}
	if v, ok := d.GetOk("assume_role_extra_policy"); ok || d.HasChange("assume_role_extra_policy") {
		request["assumeRoleExtraPolicy"] = v
	}
	if d.HasChange("prompts") {
		update = true
	}
	if v, ok := d.GetOk("prompts"); ok || d.HasChange("prompts") {
		promptsMapsArray := make([]interface{}, 0)
		for _, dataLoop4 := range convertToInterfaceArray(v) {
			dataLoop4Tmp := dataLoop4.(map[string]interface{})
			dataLoop4Map := make(map[string]interface{})
			localMaps1 := make([]interface{}, 0)
			localData5 := dataLoop4Tmp["arguments"]
			for _, dataLoop5 := range convertToInterfaceArray(localData5) {
				dataLoop5Tmp := dataLoop5.(map[string]interface{})
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["description"] = dataLoop5Tmp["description"]
				dataLoop5Map["required"] = dataLoop5Tmp["required"]
				dataLoop5Map["name"] = dataLoop5Tmp["name"]
				localMaps1 = append(localMaps1, dataLoop5Map)
			}
			dataLoop4Map["arguments"] = localMaps1
			dataLoop4Map["name"] = dataLoop4Tmp["name"]
			dataLoop4Map["content"] = dataLoop4Tmp["content"]
			dataLoop4Map["description"] = dataLoop4Tmp["description"]
			promptsMapsArray = append(promptsMapsArray, dataLoop4Map)
		}
		request["prompts"] = promptsMapsArray
	}

	if d.HasChange("vpc_whitelists") {
		update = true
	}
	if v, ok := d.GetOk("vpc_whitelists"); ok || d.HasChange("vpc_whitelists") {
		vpcWhitelistsMapsArray := convertToInterfaceArray(v)

		request["vpcWhitelists"] = vpcWhitelistsMapsArray
	}

	if d.HasChange("system_tools") {
		update = true
	}
	if v, ok := d.GetOk("system_tools"); ok || d.HasChange("system_tools") {
		systemToolsMapsArray := convertToInterfaceArray(v)

		request["systemTools"] = systemToolsMapsArray
	}

	if d.HasChange("enable_custom_vpc_whitelist") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_custom_vpc_whitelist"); ok || d.HasChange("enable_custom_vpc_whitelist") {
		request["enableCustomVpcWhitelist"] = v
	}
	if d.HasChange("oauth_client_id") {
		update = true
	}
	if v, ok := d.GetOk("oauth_client_id"); ok || d.HasChange("oauth_client_id") {
		request["oauthClientId"] = v
	}
	if d.HasChange("instructions") {
		update = true
	}
	if v, ok := d.GetOk("instructions"); ok || d.HasChange("instructions") {
		request["instructions"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("public_access") {
		update = true
	}
	if v, ok := d.GetOk("public_access"); ok || d.HasChange("public_access") {
		request["publicAccess"] = v
	}
	if d.HasChange("enable_assume_role") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_assume_role"); ok || d.HasChange("enable_assume_role") {
		request["enableAssumeRole"] = v
	}
	if d.HasChange("assume_role_name") {
		update = true
	}
	if v, ok := d.GetOk("assume_role_name"); ok || d.HasChange("assume_role_name") {
		request["assumeRoleName"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("OpenAPIExplorer", "2024-11-30", action, query, header, body, true)
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
	}

	return resourceAliCloudOpenApiExplorerApiMcpServerRead(d, meta)
}

func resourceAliCloudOpenApiExplorerApiMcpServerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/apimcpserver")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["id"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("OpenAPIExplorer", "2024-11-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"InvalidApiMcpServer.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
