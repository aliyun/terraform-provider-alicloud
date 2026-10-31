// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudCmsAddonRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAddonReleaseCreate,
		Read:   resourceAliCloudCmsAddonReleaseRead,
		Update: resourceAliCloudCmsAddonReleaseUpdate,
		Delete: resourceAliCloudCmsAddonReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"addon_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"addon_release_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"addon_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aliyun_lang": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `dry_run` has been deprecated from provider version 2.2.0.",
			},
			"entity_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"entity_types": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"region_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip_match_rule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_field_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip_cidr": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"in", "notIn", "equals", "startWith", "endWith", "inAll", "notExist"}, false),
									},
									"tag_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"field_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"in", "notIn", "equals", "startWith", "endWith", "inAll", "notExist"}, false),
									},
									"field_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"field_values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"env_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"integration_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsAddonReleaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	policyId := d.Get("integration_policy_id")
	action := fmt.Sprintf("/integration-policies/%s/addon-releases", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("addon_release_name"); ok {
		request["releaseName"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		request["values"] = v
	}
	if v, ok := d.GetOk("env_type"); ok {
		request["envType"] = v
	}
	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["aliyunLang"] = v
	}
	request["version"] = d.Get("addon_version")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["dryRun"] = v
	}
	request["addonName"] = d.Get("addon_name")
	if v, ok := d.GetOk("entity_rules"); ok {
		entityRulesMap := map[string]interface{}{}
		for _, entityRules := range convert.ConvertToInterfaceArray(v) {
			if entityRulesArg, ok := entityRules.(map[string]interface{}); ok {
				if resourceGroupId, ok := entityRulesArg["resource_group_id"]; ok {
					entityRulesMap["resourceGroupId"] = resourceGroupId
				}

				if instanceIds, ok := entityRulesArg["instance_ids"]; ok {
					entityRulesMap["instanceIds"] = convert.ConvertToInterfaceArray(instanceIds)
				}

				if entityTypes, ok := entityRulesArg["entity_types"]; ok {
					entityRulesMap["entityTypes"] = convert.ConvertToInterfaceArray(entityTypes)
				}

				if regionIds, ok := entityRulesArg["region_ids"]; ok {
					entityRulesMap["regionIds"] = convert.ConvertToInterfaceArray(regionIds)
				}

				if v, ok := entityRulesArg["ip_match_rule"]; ok {
					ipMatchRuleMap := map[string]interface{}{}
					for _, ipMatchRule := range convert.ConvertToInterfaceArray(v) {
						if ipMatchRuleArg, ok := ipMatchRule.(map[string]interface{}); ok {
							if ipFieldKey, ok := ipMatchRuleArg["ip_field_key"]; ok {
								ipMatchRuleMap["ipFieldKey"] = ipFieldKey
							}

							if ipCIDR, ok := ipMatchRuleArg["ip_cidr"]; ok {
								ipMatchRuleMap["ipCIDR"] = ipCIDR
							}
						}
					}

					entityRulesMap["ipMatchRule"] = ipMatchRuleMap
				}

				if v, ok := entityRulesArg["tags"]; ok {
					tagsMaps := make([]map[string]interface{}, 0)
					for _, tagsList := range convert.ConvertToInterfaceArray(v) {
						tagsMap := map[string]interface{}{}
						if tagsArg, ok := tagsList.(map[string]interface{}); ok {
							if op, ok := tagsArg["op"]; ok {
								tagsMap["op"] = op
							}

							if tagKey, ok := tagsArg["tag_key"]; ok {
								tagsMap["tagKey"] = tagKey
							}

							if tagValues, ok := tagsArg["tag_values"]; ok {
								tagsMap["tagValues"] = convert.ConvertToInterfaceArray(tagValues)
							}
						}

						tagsMaps = append(tagsMaps, tagsMap)
					}

					entityRulesMap["tags"] = tagsMaps
				}

				if v, ok := entityRulesArg["field_rules"]; ok {
					fieldRulesMaps := make([]map[string]interface{}, 0)
					for _, fieldRulesList := range convert.ConvertToInterfaceArray(v) {
						fieldRulesMap := map[string]interface{}{}
						if fieldRulesArg, ok := fieldRulesList.(map[string]interface{}); ok {
							if op, ok := fieldRulesArg["op"]; ok {
								fieldRulesMap["op"] = op
							}

							if fieldKey, ok := fieldRulesArg["field_key"]; ok {
								fieldRulesMap["fieldKey"] = fieldKey
							}

							if fieldValues, ok := fieldRulesArg["field_values"]; ok {
								fieldRulesMap["fieldValues"] = convert.ConvertToInterfaceArray(fieldValues)
							}
						}

						fieldRulesMaps = append(fieldRulesMaps, fieldRulesMap)
					}

					entityRulesMap["fieldRules"] = fieldRulesMaps
				}
			}
		}

		request["entityRules"] = entityRulesMap
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_addon_release", action, errs.AlibabaCloudSdkGoERROR)
	}

	releaseenvironmentIdVar, _ := jsonpath.Get("$.release.environmentId", response)
	releasereleaseNameVar, _ := jsonpath.Get("$.release.releaseName", response)
	d.SetId(fmt.Sprintf("%v:%v", releaseenvironmentIdVar, releasereleaseNameVar))

	return resourceAliCloudCmsAddonReleaseRead(d, meta)
}

func resourceAliCloudCmsAddonReleaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAddonRelease(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_addon_release DescribeCmsAddonRelease Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("addon_name", objectRaw["addonName"])
	d.Set("addon_version", objectRaw["version"])
	d.Set("aliyun_lang", objectRaw["language"])
	d.Set("config", objectRaw["config"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("env_type", objectRaw["envType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("addon_release_name", objectRaw["releaseName"])
	d.Set("integration_policy_id", objectRaw["policyId"])

	if entityRulesObj, ok := objectRaw["entityRules"]; ok {
		if entityRulesObjMap, ok := entityRulesObj.(map[string]interface{}); ok {
			if entityRules, ok := entityRulesObjMap["entityRules"]; ok {
				entityRulesMaps := make([]map[string]interface{}, 0)
				entityRulesMap := map[string]interface{}{}

				if entityRulesArg, ok := entityRules.(map[string]interface{}); ok {
					if resourceGroupId, ok := entityRulesArg["resourceGroupId"]; ok {
						entityRulesMap["resource_group_id"] = resourceGroupId
					}

					if instanceIds, ok := entityRulesArg["instanceIds"]; ok {
						entityRulesMap["instance_ids"] = instanceIds
					}

					if entityTypes, ok := entityRulesArg["entityTypes"]; ok {
						entityRulesMap["entity_types"] = entityTypes
					}

					if regionIds, ok := entityRulesArg["regionIds"]; ok {
						entityRulesMap["region_ids"] = regionIds
					}

					if ipMatchRule, ok := entityRulesArg["ipMatchRule"]; ok {
						ipMatchRuleMaps := make([]map[string]interface{}, 0)
						ipMatchRuleMap := map[string]interface{}{}
						if ipMatchRuleArg, ok := ipMatchRule.(map[string]interface{}); ok {
							if ipFieldKey, ok := ipMatchRuleArg["ipFieldKey"]; ok {
								ipMatchRuleMap["ip_field_key"] = ipFieldKey
							}

							if ipCIDR, ok := ipMatchRuleArg["ipCIDR"]; ok {
								ipMatchRuleMap["ip_cidr"] = ipCIDR
							}
						}

						ipMatchRuleMaps = append(ipMatchRuleMaps, ipMatchRuleMap)

						entityRulesMap["ip_match_rule"] = ipMatchRuleMaps
					}

					if v, ok := entityRulesArg["tags"]; ok {
						tagsMaps := make([]map[string]interface{}, 0)
						for _, tagsList := range convert.ConvertToInterfaceArray(v) {
							tagsMap := map[string]interface{}{}
							if tagsArg, ok := tagsList.(map[string]interface{}); ok {
								if op, ok := tagsArg["op"]; ok {
									tagsMap["op"] = op
								}

								if tagKey, ok := tagsArg["tagKey"]; ok {
									tagsMap["tag_key"] = tagKey
								}

								if tagValues, ok := tagsArg["tagValues"]; ok {
									tagsMap["tag_values"] = tagValues
								}
							}

							tagsMaps = append(tagsMaps, tagsMap)
						}

						entityRulesMap["tags"] = tagsMaps
					}

					if v, ok := entityRulesArg["fieldRules"]; ok {
						fieldRulesMaps := make([]map[string]interface{}, 0)
						for _, fieldRulesList := range convert.ConvertToInterfaceArray(v) {
							fieldRulesMap := map[string]interface{}{}
							if fieldRulesArg, ok := fieldRulesList.(map[string]interface{}); ok {
								if op, ok := fieldRulesArg["op"]; ok {
									fieldRulesMap["op"] = op
								}

								if fieldKey, ok := fieldRulesArg["fieldKey"]; ok {
									fieldRulesMap["field_key"] = fieldKey
								}

								if fieldValues, ok := fieldRulesArg["fieldValues"]; ok {
									fieldRulesMap["field_values"] = fieldValues
								}
							}

							fieldRulesMaps = append(fieldRulesMaps, fieldRulesMap)
						}

						entityRulesMap["field_rules"] = fieldRulesMaps
					}
				}

				entityRulesMaps = append(entityRulesMaps, entityRulesMap)

				d.Set("entity_rules", entityRulesMaps)
			}
		}
	}

	return nil
}

func resourceAliCloudCmsAddonReleaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	releaseName := parts[1]
	policyId := parts[0]
	action := fmt.Sprintf("/integration-policies/%s/addon-releases/%s", policyId, releaseName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("config") {
		update = true
	}
	if v, ok := d.GetOk("config"); ok || d.HasChange("config") {
		request["values"] = v
	}
	if d.HasChange("addon_version") {
		update = true
	}
	request["addonVersion"] = d.Get("addon_version")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["dryRun"] = v
	}

	if d.HasChange("entity_rules") {
		update = true
	}
	if v, ok := d.GetOk("entity_rules"); ok {
		entityRulesMap := map[string]interface{}{}
		for _, entityRules := range convert.ConvertToInterfaceArray(v) {
			if entityRulesArg, ok := entityRules.(map[string]interface{}); ok {
				if resourceGroupId, ok := entityRulesArg["resource_group_id"]; ok {
					entityRulesMap["resourceGroupId"] = resourceGroupId
				}

				if instanceIds, ok := entityRulesArg["instance_ids"]; ok {
					entityRulesMap["instanceIds"] = convert.ConvertToInterfaceArray(instanceIds)
				}

				if entityTypes, ok := entityRulesArg["entity_types"]; ok {
					entityRulesMap["entityTypes"] = convert.ConvertToInterfaceArray(entityTypes)
				}

				if regionIds, ok := entityRulesArg["region_ids"]; ok {
					entityRulesMap["regionIds"] = convert.ConvertToInterfaceArray(regionIds)
				}

				if v, ok := entityRulesArg["ip_match_rule"]; ok {
					ipMatchRuleMap := map[string]interface{}{}
					for _, ipMatchRule := range convert.ConvertToInterfaceArray(v) {
						if ipMatchRuleArg, ok := ipMatchRule.(map[string]interface{}); ok {
							if ipFieldKey, ok := ipMatchRuleArg["ip_field_key"]; ok {
								ipMatchRuleMap["ipFieldKey"] = ipFieldKey
							}

							if ipCIDR, ok := ipMatchRuleArg["ip_cidr"]; ok {
								ipMatchRuleMap["ipCIDR"] = ipCIDR
							}
						}
					}

					entityRulesMap["ipMatchRule"] = ipMatchRuleMap
				}

				if v, ok := entityRulesArg["tags"]; ok {
					tagsMaps := make([]map[string]interface{}, 0)
					for _, tagsList := range convert.ConvertToInterfaceArray(v) {
						tagsMap := map[string]interface{}{}
						if tagsArg, ok := tagsList.(map[string]interface{}); ok {
							if op, ok := tagsArg["op"]; ok {
								tagsMap["op"] = op
							}

							if tagKey, ok := tagsArg["tag_key"]; ok {
								tagsMap["tagKey"] = tagKey
							}

							if tagValues, ok := tagsArg["tag_values"]; ok {
								tagsMap["tagValues"] = convert.ConvertToInterfaceArray(tagValues)
							}
						}

						tagsMaps = append(tagsMaps, tagsMap)
					}

					entityRulesMap["tags"] = tagsMaps
				}

				if v, ok := entityRulesArg["field_rules"]; ok {
					fieldRulesMaps := make([]map[string]interface{}, 0)
					for _, fieldRulesList := range convert.ConvertToInterfaceArray(v) {
						fieldRulesMap := map[string]interface{}{}
						if fieldRulesArg, ok := fieldRulesList.(map[string]interface{}); ok {
							if op, ok := fieldRulesArg["op"]; ok {
								fieldRulesMap["op"] = op
							}

							if fieldKey, ok := fieldRulesArg["field_key"]; ok {
								fieldRulesMap["fieldKey"] = fieldKey
							}

							if fieldValues, ok := fieldRulesArg["field_values"]; ok {
								fieldRulesMap["fieldValues"] = convert.ConvertToInterfaceArray(fieldValues)
							}
						}

						fieldRulesMaps = append(fieldRulesMaps, fieldRulesMap)
					}

					entityRulesMap["fieldRules"] = fieldRulesMaps
				}
			}
		}

		request["entityRules"] = entityRulesMap
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsAddonReleaseRead(d, meta)
}

func resourceAliCloudCmsAddonReleaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	policyId := parts[0]
	action := fmt.Sprintf("/integration-policies/%s/addon-releases", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["releaseName"] = helper.StringPointer(parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"404"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cmsServiceV2 := CmsServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cmsServiceV2.CmsAddonReleaseStateRefreshFunc(d.Id(), "releaseName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
