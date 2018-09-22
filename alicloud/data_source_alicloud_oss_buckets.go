package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudOssBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssBucketsRead,

		Schema: map[string]*schema.Schema{
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
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
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
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cors_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_methods": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_origins": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"expose_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"max_age_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"website": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"error_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},

						"logging": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},

						"referer_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_empty": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"referers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							MaxItems: 1,
						},

						"lifecycle_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"expiration": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"days": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										MaxItems: 1,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOssBucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	var allBuckets []oss.BucketProperties
	nextMarker := ""
	for {
		var options []oss.Option
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}

		resp, err := client.ossconn.ListBuckets(options...)
		if err != nil {
			return err
		}

		if resp.Buckets == nil || len(resp.Buckets) < 1 {
			break
		}

		allBuckets = append(allBuckets, resp.Buckets...)

		nextMarker = resp.NextMarker
		if nextMarker == "" {
			break
		}
	}

	var filteredBucketsTemp []oss.BucketProperties
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, bucket := range allBuckets {
			if r != nil && !r.MatchString(bucket.Name) {
				continue
			}
			filteredBucketsTemp = append(filteredBucketsTemp, bucket)
		}
	} else {
		filteredBucketsTemp = allBuckets
	}

	if len(filteredBucketsTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_oss_buckets - Bucket found: %#v", filteredBucketsTemp)

	return bucketsDescriptionAttributes(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionAttributes(d *schema.ResourceData, buckets []oss.BucketProperties, meta interface{}) error {
	client := meta.(*AliyunClient)

	var ids []string
	var s []map[string]interface{}
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"name":          bucket.Name,
			"location":      bucket.Location,
			"storage_class": bucket.StorageClass,
			"creation_date": bucket.CreationDate.Format("2006-01-02"),
		}

		// Add additional information
		resp, err := client.ossconn.GetBucketInfo(bucket.Name)
		if err == nil {
			mapping["acl"] = resp.BucketInfo.ACL
			mapping["extranet_endpoint"] = resp.BucketInfo.ExtranetEndpoint
			mapping["intranet_endpoint"] = resp.BucketInfo.IntranetEndpoint
			mapping["owner"] = resp.BucketInfo.Owner.ID
		} else {
			log.Printf("[WARN] Unable to get additional information for the bucket %s: %v", bucket.Name, err)
		}

		// Add CORS rule information
		var ruleMappings []map[string]interface{}
		cors, err := client.ossconn.GetBucketCORS(bucket.Name)
		if err != nil && !IsExceptedErrors(err, []string{NoSuchCORSConfiguration}) {
			log.Printf("[WARN] Unable to get CORS information for the bucket %s: %v", bucket.Name, err)
		} else if err == nil && cors.CORSRules != nil {
			for _, rule := range cors.CORSRules {
				ruleMapping := make(map[string]interface{})
				ruleMapping["allowed_headers"] = rule.AllowedHeader
				ruleMapping["allowed_methods"] = rule.AllowedMethod
				ruleMapping["allowed_origins"] = rule.AllowedOrigin
				ruleMapping["expose_headers"] = rule.ExposeHeader
				ruleMapping["max_age_seconds"] = rule.MaxAgeSeconds
				ruleMappings = append(ruleMappings, ruleMapping)
			}
		}
		mapping["cors_rules"] = ruleMappings

		// Add website configuration
		var websiteMappings []map[string]interface{}
		ws, err := client.ossconn.GetBucketWebsite(bucket.Name)
		if err != nil && !IsExceptedErrors(err, []string{NoSuchWebsiteConfiguration}) {
			log.Printf("[WARN] Unable to get website information for the bucket %s: %v", bucket.Name, err)
		} else if err == nil && &ws != nil {
			websiteMapping := make(map[string]interface{})
			if v := &ws.IndexDocument; v != nil {
				websiteMapping["index_document"] = v.Suffix
			}
			if v := &ws.ErrorDocument; v != nil {
				websiteMapping["error_document"] = v.Key
			}
			websiteMappings = append(websiteMappings, websiteMapping)
		}
		mapping["website"] = websiteMappings

		// Add logging information
		var loggingMappings []map[string]interface{}
		logging, err := client.ossconn.GetBucketLogging(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get logging information for the bucket %s: %v", bucket.Name, err)
		} else if logging.LoggingEnabled.TargetBucket != "" || logging.LoggingEnabled.TargetPrefix != "" {
			loggingMapping := map[string]interface{}{
				"target_bucket": logging.LoggingEnabled.TargetBucket,
				"target_prefix": logging.LoggingEnabled.TargetPrefix,
			}
			loggingMappings = append(loggingMappings, loggingMapping)
		}
		mapping["logging"] = loggingMappings

		// Add referer information
		var refererMappings []map[string]interface{}
		referer, err := client.ossconn.GetBucketReferer(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get referer information for the bucket %s: %v", bucket.Name, err)
		} else {
			refererMapping := map[string]interface{}{
				"allow_empty": referer.AllowEmptyReferer,
				"referers":    referer.RefererList,
			}
			refererMappings = append(refererMappings, refererMapping)
		}
		mapping["referer_config"] = refererMappings

		// Add lifecycle information
		var lifecycleRuleMappings []map[string]interface{}
		lifecycle, err := client.ossconn.GetBucketLifecycle(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get lifecycle information for the bucket %s: %v", bucket.Name, err)
		} else if len(lifecycle.Rules) > 0 {
			for _, lifecycleRule := range lifecycle.Rules {
				ruleMapping := make(map[string]interface{})
				ruleMapping["id"] = lifecycleRule.ID
				ruleMapping["prefix"] = lifecycleRule.Prefix
				if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
					ruleMapping["enabled"] = true
				} else {
					ruleMapping["enabled"] = false
				}

				// Expiration
				expirationMapping := make(map[string]interface{})
				if !lifecycleRule.Expiration.Date.IsZero() {
					expirationMapping["date"] = (lifecycleRule.Expiration.Date).Format("2006-01-02")
				}
				if &lifecycleRule.Expiration.Days != nil {
					expirationMapping["days"] = int(lifecycleRule.Expiration.Days)
				}
				ruleMapping["expiration"] = []map[string]interface{}{expirationMapping}
				lifecycleRuleMappings = append(lifecycleRuleMappings, ruleMapping)
			}
		}
		mapping["lifecycle_rule"] = lifecycleRuleMappings

		log.Printf("[DEBUG] alicloud_oss_buckets - adding bucket mapping: %v", mapping)
		ids = append(ids, bucket.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
