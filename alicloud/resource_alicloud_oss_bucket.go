package alicloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudOssBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOssBucketCreate,
		Read:   resourceAlicloudOssBucketRead,
		Update: resourceAlicloudOssBucketUpdate,
		Delete: resourceAlicloudOssBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateOssBucketName,
				Default:      resource.PrefixedUniqueId("tf-oss-bucket-"),
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      oss.ACLPrivate,
				Optional:     true,
				ValidateFunc: validateOssBucketAcl,
			},

			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				MaxItems: 10,
			},

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Required: true,
						},

						"error_document": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},

			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},

			"logging_isenable": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated from 1.37.0. When `logging` is set, the bucket logging will be able.",
			},

			"referer_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_empty": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"referers": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				MaxItems: 1,
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateOssBucketLifecycleRuleId,
						},
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      expirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateOssBucketDateTimestamp,
									},
									"days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
				MaxItems: 1000,
			},

			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Default:  oss.StorageStandard,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(oss.StorageStandard),
					string(oss.StorageIA),
					string(oss.StorageArchive),
				}),
			},
			"server_side_encryption_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sse_algorithm": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateAllowedStringValue([]string{
								ServerSideEncryptionAes256,
								ServerSideEncryptionKMS,
							}),
						},
					},
				},
				MaxItems: 1,
			},

			"tags": tagsSchema(),

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateAllowedStringValue([]string{
								"Enabled",
								"Suspended",
							}),
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceAlicloudOssBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bucket := d.Get("bucket").(string)
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.IsBucketExist(bucket)
	})
	if err != nil {
		return err
	}
	isExist, _ := raw.(bool)
	if isExist {
		return fmt.Errorf("[ERROR] The specified bucket name: %#v is not available. The bucket namespace is shared by all users of the OSS system. Please select a different name and try again.", bucket)
	}

	_, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.CreateBucket(bucket, oss.StorageClass(oss.StorageClassType(d.Get("storage_class").(string))))
	})
	if err != nil {
		return fmt.Errorf("Error creating OSS bucket: %#v", err)
	}

	retryErr := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.IsBucketExist(bucket)
		})

		if err != nil {
			return resource.NonRetryableError(err)
		}
		isExist, _ := raw.(bool)
		if !isExist {
			return resource.RetryableError(fmt.Errorf("Trying to ensure new OSS bucket %#v has been created successfully.", bucket))
		}

		return nil
	})

	if retryErr != nil {
		return fmt.Errorf("Error creating OSS bucket: %#v", retryErr)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)

	return resourceAlicloudOssBucketUpdate(d, meta)
}

func resourceAlicloudOssBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(d.Id())
	})
	if err != nil {
		if ossNotFoundError(err) {
			return nil
		}
		return err
	}
	info, _ := raw.(oss.GetBucketInfoResult)
	d.Set("bucket", d.Id())

	d.Set("acl", info.BucketInfo.ACL)
	d.Set("creation_date", info.BucketInfo.CreationDate.Format("2016-01-01"))
	d.Set("extranet_endpoint", info.BucketInfo.ExtranetEndpoint)
	d.Set("intranet_endpoint", info.BucketInfo.IntranetEndpoint)
	d.Set("location", info.BucketInfo.Location)
	d.Set("owner", info.BucketInfo.Owner.ID)
	d.Set("storage_class", info.BucketInfo.StorageClass)

	if &info.BucketInfo.SseRule != nil {
		if len(info.BucketInfo.SseRule.SSEAlgorithm) > 0 && info.BucketInfo.SseRule.SSEAlgorithm != "None" {
			rule := make(map[string]interface{})
			rule["sse_algorithm"] = info.BucketInfo.SseRule.SSEAlgorithm
			data := make([]map[string]interface{}, 0)
			data = append(data, rule)
			d.Set("server_side_encryption_rule", data)
		}
	}

	if info.BucketInfo.Versioning != "" {
		data := map[string]interface{}{
			"status": info.BucketInfo.Versioning,
		}
		versioning := make([]map[string]interface{}, 0)
		versioning = append(versioning, data)
		d.Set("versioning", versioning)
	}

	// Read the CORS
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketCORS(d.Id())
	})
	if err != nil && !IsExceptedErrors(err, []string{NoSuchCORSConfiguration}) {
		return err
	}
	cors, _ := raw.(oss.GetBucketCORSResult)
	rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
	for _, r := range cors.CORSRules {
		rule := make(map[string]interface{})
		rule["allowed_headers"] = r.AllowedHeader
		rule["allowed_methods"] = r.AllowedMethod
		rule["allowed_origins"] = r.AllowedOrigin
		rule["expose_headers"] = r.ExposeHeader
		rule["max_age_seconds"] = r.MaxAgeSeconds

		rules = append(rules, rule)
	}
	if err := d.Set("cors_rule", rules); err != nil {
		return err
	}

	// Read the website configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketWebsite(d.Id())
	})
	if err != nil && !IsExceptedErrors(err, []string{NoSuchWebsiteConfiguration}) {
		return fmt.Errorf("Error getting bucket website: %#v", err)
	}
	ws, _ := raw.(oss.GetBucketWebsiteResult)
	if err == nil && &ws != nil {
		var websites []map[string]interface{}
		w := make(map[string]interface{})

		if v := &ws.IndexDocument; v != nil {
			w["index_document"] = v.Suffix
		}

		if v := &ws.ErrorDocument; v != nil {
			w["error_document"] = v.Key
		}
		websites = append(websites, w)
		if err := d.Set("website", websites); err != nil {
			return err
		}
	}

	// Read the logging configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketLogging(d.Id())
	})
	if err != nil {
		return fmt.Errorf("Error getting bucket logging: %#v", err)
	}
	logging, _ := raw.(oss.GetBucketLoggingResult)
	if &logging != nil {
		enable := logging.LoggingEnabled
		if &enable != nil {
			lgs := make([]map[string]interface{}, 0)
			tb := logging.LoggingEnabled.TargetBucket
			tp := logging.LoggingEnabled.TargetPrefix
			if tb != "" || tp != "" {
				lgs = append(lgs, map[string]interface{}{
					"target_bucket": tb,
					"target_prefix": tp,
				})
			}
			if err := d.Set("logging", lgs); err != nil {
				return err
			}
		}
	}

	// Read the bucket referer
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketReferer(d.Id())
	})
	referers := make([]map[string]interface{}, 0)
	if err != nil {
		return fmt.Errorf("Error getting bucket referer: %#v", err)
	}
	referer, _ := raw.(oss.GetBucketRefererResult)
	if len(referer.RefererList) > 0 {
		referers = append(referers, map[string]interface{}{
			"allow_empty": referer.AllowEmptyReferer,
			"referers":    referer.RefererList,
		})
		if err := d.Set("referer_config", referers); err != nil {
			return err
		}
	}

	// Read the lifecycle rule configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketLifecycle(d.Id())
	})
	if err != nil {
		if ossNotFoundError(err) {
			log.Printf("[WARN] OSS bucket: %s, no lifecycle could be found.", d.Id())
			return nil
		}
		return fmt.Errorf("Error getting bucket lifecycle: %#v", err)
	}
	lifecycle, _ := raw.(oss.GetBucketLifecycleResult)
	if len(lifecycle.Rules) > 0 {
		rules := make([]map[string]interface{}, 0, len(lifecycle.Rules))

		for _, lifecycleRule := range lifecycle.Rules {
			rule := make(map[string]interface{})
			rule["id"] = lifecycleRule.ID
			rule["prefix"] = lifecycleRule.Prefix
			if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
				rule["enabled"] = true
			} else {
				rule["enabled"] = false
			}

			// expiration
			if &lifecycleRule.Expiration != nil {
				e := make(map[string]interface{})
				if lifecycleRule.Expiration.Date != "" {
					t, err := time.Parse("2006-01-02T15:04:05.000Z", lifecycleRule.Expiration.Date)
					if err != nil {
						return err
					}
					e["date"] = t.Format("2006-01-02")
				}
				if &lifecycleRule.Expiration.Days != nil {
					e["days"] = int(lifecycleRule.Expiration.Days)
				}
				rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
			}
			rules = append(rules, rule)
		}

		if err := d.Set("lifecycle_rule", rules); err != nil {
			return err
		}
	}

	// Read Policy
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		params := map[string]interface{}{}
		params["policy"] = nil
		return ossClient.Conn.Do("GET", d.Id(), "", params, nil, nil, 0, nil)
	})

	if err != nil {
		if ossNotFoundError(err) {
			log.Printf("[WARN] OSS bucket: %s, no policy could be found.", d.Id())
			return nil
		}
		return fmt.Errorf("Error getting bucket policy: %#v", err)
	}

	rawResp := raw.(*oss.Response)
	defer rawResp.Body.Close()

	if err == nil {
		rawData, err := ioutil.ReadAll(rawResp.Body)
		if err != nil {
			return err
		}
		err = d.Set("policy", string(rawData))
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Read tags
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketTagging(d.Id())
	})
	if err != nil {
		if ossNotFoundError(err) {
			log.Printf("[WARN] OSS bucket: %s, no tagging could be found.", d.Id())
			return nil
		}
		return fmt.Errorf("Error getting bucket tagging: %#v", err)
	}

	tagging, _ := raw.(oss.GetBucketTaggingResult)
	if len(tagging.Tags) > 0 {
		tagsMap := make(map[string]string)
		for _, t := range tagging.Tags {
			tagsMap[t.Key] = t.Value
		}
		if err := d.Set("tags", tagsMap); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlicloudOssBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("acl") {
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.SetBucketACL(d.Id(), oss.ACLType(d.Get("acl").(string)))
		})
		if err != nil {
			return fmt.Errorf("Error setting OSS bucket ACL: %#v", err)
		}
		d.SetPartial("acl")
	}

	if d.HasChange("cors_rule") {
		if err := resourceAlicloudOssBucketCorsUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("cors_rule")
	}

	if d.HasChange("website") {
		if err := resourceAlicloudOssBucketWebsiteUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("website")
	}

	if d.HasChange("logging") {
		if err := resourceAlicloudOssBucketLoggingUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("logging")
	}

	if d.HasChange("referer_config") {
		if err := resourceAlicloudOssBucketRefererUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("referer_config")
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceAlicloudOssBucketLifecycleRuleUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("lifecycle_rule")
	}

	if d.HasChange("policy") {
		if err := resourceAlicloudOssBucketPolicyUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("policy")
	}

	if d.HasChange("server_side_encryption_rule") {
		if err := resourceAlicloudOssBucketEncryptionUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("server_side_encryption_rule")
	}

	if d.HasChange("tags") {
		if err := resourceAlicloudOssBucketTaggingUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	if d.HasChange("versioning") {
		if err := resourceAlicloudOssBucketVersioningUpdate(client, d); err != nil {
			return err
		}
		d.SetPartial("versioning")
	}

	d.Partial(false)
	return resourceAlicloudOssBucketRead(d, meta)
}

func resourceAlicloudOssBucketCorsUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	cors := d.Get("cors_rule").([]interface{})
	if cors == nil || len(cors) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketCORS(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket cors_rule: %#v", err)
		}
		return nil
	}
	// Put CORS
	rules := make([]oss.CORSRule, 0, len(cors))
	for _, c := range cors {
		corsMap := c.(map[string]interface{})
		rule := oss.CORSRule{}
		for k, v := range corsMap {
			log.Printf("[DEBUG] OSS bucket: %s, put CORS: %#v, %#v", d.Id(), k, v)
			if k == "max_age_seconds" {
				rule.MaxAgeSeconds = v.(int)
			} else {
				rMap := make([]string, len(v.([]interface{})))
				for i, vv := range v.([]interface{}) {
					rMap[i] = vv.(string)
				}
				switch k {
				case "allowed_headers":
					rule.AllowedHeader = rMap
				case "allowed_methods":
					rule.AllowedMethod = rMap
				case "allowed_origins":
					rule.AllowedOrigin = rMap
				case "expose_headers":
					rule.ExposeHeader = rMap
				}
			}
		}
		rules = append(rules, rule)
	}

	log.Printf("[DEBUG] Oss bucket: %s, put CORS: %#v", d.Id(), cors)
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketCORS(d.Id(), rules)
	})
	if err != nil {
		return fmt.Errorf("Error putting oss CORS: %s", err)
	}

	return nil
}
func resourceAlicloudOssBucketWebsiteUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	ws := d.Get("website").([]interface{})
	if ws == nil || len(ws) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketWebsite(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket logging: %#v", err)
		}
		return nil
	}

	var index_document, error_document string
	w := ws[0].(map[string]interface{})

	if v, ok := w["index_document"]; ok {
		index_document = v.(string)
	}
	if v, ok := w["error_document"]; ok {
		error_document = v.(string)
	}
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketWebsite(d.Id(), index_document, error_document)
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS bucket website: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketLoggingUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	logging := d.Get("logging").([]interface{})
	if logging == nil || len(logging) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketLogging(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket logging: %#v", err)
		}
		return nil
	}

	c := logging[0].(map[string]interface{})
	var target_bucket, target_prefix string
	if v, ok := c["target_bucket"]; ok {
		target_bucket = v.(string)
	}
	if v, ok := c["target_prefix"]; ok {
		target_prefix = v.(string)
	}
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketLogging(d.Id(), target_bucket, target_prefix, target_bucket != "" || target_prefix != "")
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS bucket logging: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketRefererUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	config := d.Get("referer_config").([]interface{})
	if config == nil || len(config) < 1 {
		log.Printf("[DEBUG] OSS set bucket referer as nil")
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.SetBucketReferer(d.Id(), nil, true)
		})
		if err != nil {
			return fmt.Errorf("Error deleting OSS website: %#v", err)
		}
		return nil
	}

	c := config[0].(map[string]interface{})

	var allow bool
	var referers []string
	if v, ok := c["allow_empty"]; ok {
		allow = v.(bool)
	}
	if v, ok := c["referers"]; ok {
		for _, referer := range v.([]interface{}) {
			referers = append(referers, referer.(string))
		}
	}
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketReferer(d.Id(), referers, allow)
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS bucket referer configuration: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketLifecycleRuleUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	bucket := d.Id()
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	if lifecycleRules == nil || len(lifecycleRules) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketLifecycle(bucket)
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket lifecycle: %#v", err)
		}
		return nil
	}

	rules := make([]oss.LifecycleRule, 0, len(lifecycleRules))

	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})

		rule := oss.LifecycleRule{
			Prefix: r["prefix"].(string),
		}

		// ID
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = val
		}

		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = string(ExpirationStatusEnabled)
		} else {
			rule.Status = string(ExpirationStatusDisabled)
		}

		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			e := expiration[0].(map[string]interface{})
			i := oss.LifecycleExpiration{}
			valDate, _ := e["date"].(string)
			valDays, _ := e["days"].(int)

			if (valDate != "" && valDays > 0) || (valDate == "" && valDays <= 0) {
				return fmt.Errorf("'date' conflicts with 'days'. One and only one of them can be specified in one expiration configuration.")
			}

			if valDate != "" {
				i.Date = fmt.Sprintf("%sT00:00:00.000Z", valDate)
			}
			if valDays > 0 {
				i.Days = valDays
			}
			rule.Expiration = &i
		}
		rules = append(rules, rule)
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.SetBucketLifecycle(bucket, rules)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS lifecycle rule: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketPolicyUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	bucket := d.Id()
	policy := d.Get("policy").(string)

	if len(policy) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				params := map[string]interface{}{}
				params["policy"] = nil
				return ossClient.Conn.Do("DELETE", bucket, "", params, nil, nil, 0, nil)
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket policy: %#v", err)
		}
		return nil
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			params := map[string]interface{}{}
			params["policy"] = nil

			buffer := new(bytes.Buffer)
			buffer.Write([]byte(policy))
			return ossClient.Conn.Do("PUT", bucket, "", params, nil, buffer, 0, nil)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS lifecycle rule: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketEncryptionUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	encryption_rule := d.Get("server_side_encryption_rule").([]interface{})
	if encryption_rule == nil || len(encryption_rule) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketEncryption(d.Id())
				return nil, nil
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket encryption: %#v", err)
		}
		return nil
	}

	var sseRule oss.ServerEncryptionRule
	c := encryption_rule[0].(map[string]interface{})
	if v, ok := c["sse_algorithm"]; ok {
		sseRule.SSEDefault.SSEAlgorithm = v.(string)
	}

	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketEncryption(d.Id(), sseRule)
	})
	if err != nil {
		return fmt.Errorf("Error putting OSS bucket encryption: %#v", err)
	}

	return nil
}

func resourceAlicloudOssBucketTaggingUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	tagsMap := d.Get("tags").(map[string]interface{})
	if tagsMap == nil || len(tagsMap) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return nil, ossClient.DeleteBucketTagging(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error removing OSS bucket tagging: %#v", err)
		}
		return nil
	}

	// Put tagging
	var bTagging oss.Tagging
	for k, v := range tagsMap {
		bTagging.Tags = append(bTagging.Tags, oss.Tag{
			Key:   k,
			Value: v.(string),
		})
	}
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketTagging(d.Id(), bTagging)
	})
	if err != nil {
		return fmt.Errorf("Error putting oss tagging: %s", err)
	}

	return nil
}

func resourceAlicloudOssBucketVersioningUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	versioning := d.Get("versioning").([]interface{})
	if len(versioning) == 1 {
		var status string
		c := versioning[0].(map[string]interface{})
		if v, ok := c["status"]; ok {
			status = v.(string)
		}

		versioningCfg := oss.VersioningConfig{}
		versioningCfg.Status = status
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.SetBucketVersioning(d.Id(), versioningCfg)
		})

		if err != nil {
			return fmt.Errorf("Error putting oss versioning: %s", err)
		}
	}

	return nil
}

func resourceAlicloudOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossService := OssService{client}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.IsBucketExist(d.Id())
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("OSS delete bucket got an error: %#v", err))
		}
		exist, _ := raw.(bool)
		if !exist {
			return nil
		}

		_, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.DeleteBucket(d.Id())
		})
		if err != nil {
			if err.(oss.ServiceError).Code == "BucketNotEmpty" {
				if d.Get("force_destroy").(bool) {
					log.Printf("[DEBUG] oss Bucket attempting to forceDestroy %+v", err)
					_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
						bucket, _ := ossClient.Bucket(d.Get("bucket").(string))
						lor, err := bucket.ListObjectVersions()
						if err == nil {
							objectsToDelete := make([]oss.DeleteObject, 0)
							for _, object := range lor.ObjectDeleteMarkers {
								objectsToDelete = append(objectsToDelete, oss.DeleteObject{
									Key:       object.Key,
									VersionId: object.VersionId,
								})
							}

							for _, object := range lor.ObjectVersions {
								objectsToDelete = append(objectsToDelete, oss.DeleteObject{
									Key:       object.Key,
									VersionId: object.VersionId,
								})
							}

							_, err = bucket.DeleteObjectVersions(objectsToDelete)
						}
						return nil, err
					})

					if err != nil {
						return resource.NonRetryableError(fmt.Errorf("When force_destroy OSS bucket, got an error: %#v", err))
					}
				}
			}
			return resource.RetryableError(fmt.Errorf("OSS Bucket %s is in use - trying again while it is deleted.", d.Id()))
		}
		bucket, err := ossService.QueryOssBucketById(d.Id())
		if err != nil {
			// Verify the error is what we want
			if IsExceptedErrors(err, []string{OssBucketNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting OSS bucket, describing it got an error: %#v", err))
		}
		if bucket.Name != "" {
			return resource.RetryableError(fmt.Errorf("Deleting OSS Bucket %s timeout.", d.Id()))
		}

		return nil
	})
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}
