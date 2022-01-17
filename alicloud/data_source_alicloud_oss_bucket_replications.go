package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOssBucketReplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssBucketReplicationsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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
						"redundancy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_region_replication": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transfer_acceleration": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"replication_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prefix_set": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefixes": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"location": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"transfer_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"historical_object_replication": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sync_role": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudOssBucketReplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *oss.Client
	var allBuckets []oss.BucketProperties
	nextMarker := ""
	for {
		var options []oss.Option
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}

		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return ossClient.ListBuckets(options...)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oss_bucket", "CreateBucket", AliyunOssGoSdk)
		}
		if debugOn() {
			addDebug("ListBuckets", raw, requestInfo, map[string]interface{}{"options": options})
		}
		response, _ := raw.(oss.ListBucketsResult)

		if response.Buckets == nil || len(response.Buckets) < 1 {
			break
		}

		allBuckets = append(allBuckets, response.Buckets...)

		nextMarker = response.NextMarker
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
	return bucketsDescriptionReplication(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionReplication(d *schema.ResourceData, buckets []oss.BucketProperties, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var s []map[string]interface{}
	var names []string
	var requestInfo *oss.Client
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"name":          bucket.Name,
			"location":      bucket.Location,
			"storage_class": bucket.StorageClass,
			"creation_date": bucket.CreationDate.Format("2006-01-02"),
		}

		// Add additional information
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return ossClient.GetBucketInfo(bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketInfo", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			response, _ := raw.(oss.GetBucketInfoResult)
			mapping["acl"] = response.BucketInfo.ACL
			mapping["extranet_endpoint"] = response.BucketInfo.ExtranetEndpoint
			mapping["intranet_endpoint"] = response.BucketInfo.IntranetEndpoint
			mapping["owner"] = response.BucketInfo.Owner.ID
			mapping["redundancy_type"] = response.BucketInfo.RedundancyType
			mapping["cross_region_replication"] = response.BucketInfo.CrossRegionReplication
			mapping["transfer_acceleration"] = response.BucketInfo.TransferAcceleration

		} else {
			log.Printf("[WARN] Unable to get additional information for the bucket %s: %v", bucket.Name, err)
		}

		// Add replication information
		var replicationRuleMappings []map[string]interface{}
		replicationRuleMapping := make(map[string]interface{})
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return ossClient.GetBucketReplication(bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketReplication", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			replication, _ := raw.(oss.GetBucketReplicationResult)

			replicationRuleMapping["id"] = replication.Rule.ID
			replicationRuleMapping["action"] = replication.Rule.Action
			replicationRuleMapping["status"] = replication.Rule.Status
			replicationRuleMapping["historical_object_replication"] = replication.Rule.HistoricalObjectReplication
			replicationRuleMapping["sync_role"] = replication.Rule.SyncRole

			// Destination
			destinationMapping := make(map[string]interface{})
			if replication.Rule.Destination.Bucket != "" {
				destinationMapping["bucket"] = replication.Rule.Destination.Bucket
			}
			if replication.Rule.Destination.Location != "" {
				destinationMapping["location"] = replication.Rule.Destination.Location
			}
			if replication.Rule.Destination.TransferType != "" {
				destinationMapping["transfer_type"] = replication.Rule.Destination.TransferType
			}
			replicationRuleMapping["destination"] = []map[string]interface{}{destinationMapping}

			// PrefixSet
			prefixSetMapping := make(map[string]interface{})
			if replication.Rule.PrefixSet != nil && len(replication.Rule.PrefixSet.Prefix) != 0 {
				prefixSetMapping["prefixes"] = replication.Rule.PrefixSet.Prefix
			}
			replicationRuleMapping["prefix_set"] = []map[string]interface{}{prefixSetMapping}

			replicationRuleMappings = append(replicationRuleMappings, replicationRuleMapping)
		} else {
			log.Printf("[WARN] Unable to get replication information for the bucket %s: %v", bucket.Name, err)
		}

		mapping["replication_rule"] = replicationRuleMappings

		ids = append(ids, bucket.Name)
		s = append(s, mapping)
		names = append(names, bucket.Name)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
