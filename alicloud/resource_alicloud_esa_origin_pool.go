// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaOriginPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaOriginPoolCreate,
		Read:   resourceAliCloudEsaOriginPoolRead,
		Update: resourceAliCloudEsaOriginPoolUpdate,
		Delete: resourceAliCloudEsaOriginPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"origin_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_pool_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"origins": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"header": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"origin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_conf": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"region": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"auth_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"access_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ip_version_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, n := d.GetChange("origins")
					if o == nil || n == nil {
						return old == new
					}
					oldList := o.(*schema.Set).List()
					newList := n.(*schema.Set).List()
					if len(oldList) != len(newList) {
						return false
					}

					normalizeOrigins := func(items []interface{}) []map[string]interface{} {
						result := make([]map[string]interface{}, 0, len(items))
						for _, item := range items {
							m := item.(map[string]interface{})
							normalized := map[string]interface{}{
								"address": m["address"],
								"type":    m["type"],
								"name":    m["name"],
								"weight":  m["weight"],
								"enabled": m["enabled"],
								"header":  m["header"],
							}
							if authConf, ok := m["auth_conf"]; ok && authConf != nil {
								authList := authConf.([]interface{})
								if len(authList) > 0 {
									auth := authList[0].(map[string]interface{})
									normalized["auth_conf"] = map[string]interface{}{
										"access_key": auth["access_key"],
										"auth_type":  auth["auth_type"],
										"region":     auth["region"],
										"version":    auth["version"],
									}
								}
							}
							result = append(result, normalized)
						}
						sort.Slice(result, func(i, j int) bool {
							return fmt.Sprint(result[i]["name"]) < fmt.Sprint(result[j]["name"])
						})
						return result
					}

					oldNormalized := normalizeOrigins(oldList)
					newNormalized := normalizeOrigins(newList)

					for i := range oldNormalized {
						oldH := fmt.Sprint(oldNormalized[i]["header"])
						newH := fmt.Sprint(newNormalized[i]["header"])
						if oldH != newH {
							if equal, _ := compareJsonTemplateAreEquivalent(oldH, newH); equal {
								oldNormalized[i]["header"] = ""
								newNormalized[i]["header"] = ""
							}
						}
					}

					return reflect.DeepEqual(oldNormalized, newNormalized)
				},
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaOriginPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOriginPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("origins"); ok {
		originsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Address"] = dataLoopTmp["address"]
			if header, ok := dataLoopTmp["header"]; ok && fmt.Sprint(header) != "" {
				dataLoopMap["Header"] = parseHeader(header.(string))
			}
			dataLoopMap["Type"] = dataLoopTmp["type"]
			localData1 := make(map[string]interface{})
			version1, _ := jsonpath.Get("$[0].version", dataLoopTmp["auth_conf"])
			if version1 != nil && version1 != "" {
				localData1["Version"] = version1
			}
			authType1, _ := jsonpath.Get("$[0].auth_type", dataLoopTmp["auth_conf"])
			if authType1 != nil && authType1 != "" {
				localData1["AuthType"] = authType1
			}
			accessKey1, _ := jsonpath.Get("$[0].access_key", dataLoopTmp["auth_conf"])
			if accessKey1 != nil && accessKey1 != "" {
				localData1["AccessKey"] = accessKey1
			}
			secretKey1, _ := jsonpath.Get("$[0].secret_key", dataLoopTmp["auth_conf"])
			if secretKey1 != nil && secretKey1 != "" {
				localData1["SecretKey"] = secretKey1
			}
			region1, _ := jsonpath.Get("$[0].region", dataLoopTmp["auth_conf"])
			if region1 != nil && region1 != "" {
				localData1["Region"] = region1
			}
			dataLoopMap["AuthConf"] = localData1
			if ipVersionPolicy, ok := dataLoopTmp["ip_version_policy"]; ok {
				dataLoopMap["IpVersionPolicy"] = ipVersionPolicy
			}
			dataLoopMap["Weight"] = dataLoopTmp["weight"]
			if enabled, ok := dataLoopTmp["enabled"]; ok {
				dataLoopMap["Enabled"] = enabled
			}
			dataLoopMap["Name"] = dataLoopTmp["name"]
			originsMapsArray = append(originsMapsArray, dataLoopMap)
		}
		originsMapsJson, err := json.Marshal(originsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Origins"] = string(originsMapsJson)
	}

	request["Name"] = d.Get("origin_pool_name")
	if v, ok := d.GetOkExists("enabled"); ok {
		request["Enabled"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_origin_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["Id"]))

	return resourceAliCloudEsaOriginPoolRead(d, meta)
}

func resourceAliCloudEsaOriginPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaOriginPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_origin_pool DescribeEsaOriginPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enabled", objectRaw["Enabled"])
	d.Set("origin_pool_name", objectRaw["Name"])

	if v, ok := objectRaw["Id"]; ok {
		d.Set("origin_pool_id", fmt.Sprint(v))
	}

	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", fmt.Sprint(v))
	}

	secretKeyMap := make(map[string]interface{})
	if origins, ok := d.GetOk("origins"); ok {
		for _, dataLoop := range convertToInterfaceArray(origins) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			originName := fmt.Sprint(dataLoopTmp["name"])
			if !IsNil(dataLoopTmp["auth_conf"]) {
				secretKey1, _ := jsonpath.Get("$[0].secret_key", dataLoopTmp["auth_conf"])
				if secretKey1 != nil && secretKey1 != "" {
					secretKeyMap[originName] = secretKey1
				}
			}
		}
	}
	originsRaw := objectRaw["Origins"]
	originsMaps := make([]map[string]interface{}, 0)
	if originsRaw != nil {
		for _, originsChildRaw := range convertToInterfaceArray(originsRaw) {
			originsMap := make(map[string]interface{})
			originsChildRaw := originsChildRaw.(map[string]interface{})
			originsMap["address"] = originsChildRaw["Address"]
			originsMap["enabled"] = originsChildRaw["Enabled"]
			if header, ok := originsChildRaw["Header"]; ok {
				originsMap["header"] = convertObjectToJsonString(header)
			}
			originsMap["ip_version_policy"] = originsChildRaw["IpVersionPolicy"]
			originsMap["name"] = originsChildRaw["Name"]
			originsMap["origin_id"] = fmt.Sprint(originsChildRaw["Id"])
			originsMap["type"] = originsChildRaw["Type"]
			originsMap["weight"] = originsChildRaw["Weight"]

			authConfMaps := make([]map[string]interface{}, 0)
			authConfMap := make(map[string]interface{})
			authConfRaw := make(map[string]interface{})
			if originsChildRaw["AuthConf"] != nil {
				authConfRaw = originsChildRaw["AuthConf"].(map[string]interface{})
			}
			if len(authConfRaw) > 0 {
				authConfMap["access_key"] = authConfRaw["AccessKey"]
				authConfMap["auth_type"] = authConfRaw["AuthType"]
				authConfMap["region"] = authConfRaw["Region"]
				authConfMap["secret_key"] = authConfRaw["SecretKey"]
				originName := fmt.Sprint(originsMap["name"])
				if sk, ok := secretKeyMap[originName]; ok {
					authConfMap["secret_key"] = sk
				}
				authConfMap["version"] = authConfRaw["Version"]

				authConfMaps = append(authConfMaps, authConfMap)
			}
			originsMap["auth_conf"] = authConfMaps
			originsMaps = append(originsMaps, originsMap)
		}
	}
	if err := d.Set("origins", originsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEsaOriginPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateOriginPool"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	if d.HasChange("origins") {
		update = true
		if v, ok := d.GetOk("origins"); ok || d.HasChange("origins") {
			originsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Address"] = dataLoopTmp["address"]
				if header, ok := dataLoopTmp["header"]; ok && fmt.Sprint(header) != "" {
					dataLoopMap["Header"] = parseHeader(header.(string))
				}
				dataLoopMap["Type"] = dataLoopTmp["type"]
				if !IsNil(dataLoopTmp["auth_conf"]) {
					localData1 := make(map[string]interface{})
					version1, _ := jsonpath.Get("$[0].version", dataLoopTmp["auth_conf"])
					if version1 != nil && version1 != "" {
						localData1["Version"] = version1
					}
					authType1, _ := jsonpath.Get("$[0].auth_type", dataLoopTmp["auth_conf"])
					if authType1 != nil && authType1 != "" {
						localData1["AuthType"] = authType1
					}
					accessKey1, _ := jsonpath.Get("$[0].access_key", dataLoopTmp["auth_conf"])
					if accessKey1 != nil && accessKey1 != "" {
						localData1["AccessKey"] = accessKey1
					}
					secretKey1, _ := jsonpath.Get("$[0].secret_key", dataLoopTmp["auth_conf"])
					if secretKey1 != nil && secretKey1 != "" {
						localData1["SecretKey"] = secretKey1
					}
					region1, _ := jsonpath.Get("$[0].region", dataLoopTmp["auth_conf"])
					if region1 != nil && region1 != "" {
						localData1["Region"] = region1
					}
					dataLoopMap["AuthConf"] = localData1
				}
				if ipVersionPolicy, ok := dataLoopTmp["ip_version_policy"]; ok {
					dataLoopMap["IpVersionPolicy"] = ipVersionPolicy
				}
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				if enabled, ok := dataLoopTmp["enabled"]; ok {
					dataLoopMap["Enabled"] = enabled
				}
				dataLoopMap["Name"] = dataLoopTmp["name"]
				originsMapsArray = append(originsMapsArray, dataLoopMap)
			}
			originsMapsJson, err := json.Marshal(originsMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["Origins"] = string(originsMapsJson)
		}
	}

	if d.HasChange("enabled") {
		update = true
		request["Enabled"] = d.Get("enabled")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return resourceAliCloudEsaOriginPoolRead(d, meta)
}

func resourceAliCloudEsaOriginPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteOriginPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["Id"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func parseHeader(input string) interface{} {
	if b, err := strconv.ParseBool(input); err == nil {
		return b
	}
	if i, err := strconv.Atoi(input); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(input), &obj); err == nil {
		return obj
	}
	var arr []interface{}
	if err := json.Unmarshal([]byte(input), &arr); err == nil {
		return arr
	}
	return input
}

func comparePrefixSuffix(s1, s2 string, length int) bool {
	if len(s1) < length || len(s2) < length {
		return false
	}
	return s1[:length] == s2[:length] && s1[len(s1)-length:] == s2[len(s2)-length:]
}
