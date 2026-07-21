package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DynamoDB data type descriptors for item attributes
const (
	dynamoDataTypeBinary    = "B"
	dynamoDataTypeBinarySet = "BS"
	dynamoDataTypeBoolean   = "BOOL"
	dynamoDataTypeList      = "L"
	dynamoDataTypeMap       = "M"
	dynamoDataTypeNull      = "NULL"
	dynamoDataTypeNumber    = "N"
	dynamoDataTypeNumberSet = "NS"
	dynamoDataTypeString    = "S"
	dynamoDataTypeStringSet = "SS"
)

// jsonToAttributeValue converts a JSON value to DynamoDB AttributeValue format.
func jsonToAttributeValue(value interface{}) *dynamodb.AttributeValue {
	switch v := value.(type) {
	case map[string]interface{}:
		if len(v) == 1 {
			for key, val := range v {
				switch key {
				case dynamoDataTypeString:
					return &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%v", val))}
				case dynamoDataTypeNumber:
					return &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%v", val))}
				case dynamoDataTypeBinary:
					return &dynamodb.AttributeValue{B: []byte(fmt.Sprintf("%v", val))}
				case dynamoDataTypeBoolean:
					boolVal := fmt.Sprintf("%v", val) == "true"
					return &dynamodb.AttributeValue{BOOL: aws.Bool(boolVal)}
				case dynamoDataTypeNull:
					nullVal := fmt.Sprintf("%v", val) == "true"
					return &dynamodb.AttributeValue{NULL: aws.Bool(nullVal)}
				case dynamoDataTypeList:
					if list, ok := val.([]interface{}); ok {
						items := make([]*dynamodb.AttributeValue, 0, len(list))
						for _, item := range list {
							items = append(items, jsonToAttributeValue(item))
						}
						return &dynamodb.AttributeValue{L: items}
					}
				case dynamoDataTypeMap:
					if m, ok := val.(map[string]interface{}); ok {
						items := make(map[string]*dynamodb.AttributeValue)
						for k, mv := range m {
							items[k] = jsonToAttributeValue(mv)
						}
						return &dynamodb.AttributeValue{M: items}
					}
				case dynamoDataTypeStringSet:
					if ss, ok := val.([]interface{}); ok {
						strs := make([]*string, 0, len(ss))
						for _, s := range ss {
							strs = append(strs, aws.String(fmt.Sprintf("%v", s)))
						}
						return &dynamodb.AttributeValue{SS: strs}
					}
				case dynamoDataTypeNumberSet:
					if ns, ok := val.([]interface{}); ok {
						nums := make([]*string, 0, len(ns))
						for _, n := range ns {
							nums = append(nums, aws.String(fmt.Sprintf("%v", n)))
						}
						return &dynamodb.AttributeValue{NS: nums}
					}
				case dynamoDataTypeBinarySet:
					if bs, ok := val.([]interface{}); ok {
						bins := make([][]byte, 0, len(bs))
						for _, b := range bs {
							bins = append(bins, []byte(fmt.Sprintf("%v", b)))
						}
						return &dynamodb.AttributeValue{BS: bins}
					}
				}
			}
		}
		// Regular map
		items := make(map[string]*dynamodb.AttributeValue)
		for k, val := range v {
			items[k] = jsonToAttributeValue(val)
		}
		return &dynamodb.AttributeValue{M: items}
	case []interface{}:
		items := make([]*dynamodb.AttributeValue, 0, len(v))
		for _, item := range v {
			items = append(items, jsonToAttributeValue(item))
		}
		return &dynamodb.AttributeValue{L: items}
	case string:
		return &dynamodb.AttributeValue{S: aws.String(v)}
	case float64:
		return &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%v", v))}
	case bool:
		return &dynamodb.AttributeValue{BOOL: aws.Bool(v)}
	case nil:
		return &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	default:
		return &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%v", v))}
	}
}

// attributeValueToJSON converts a DynamoDB AttributeValue to JSON-compatible format.
func attributeValueToJSON(av *dynamodb.AttributeValue) interface{} {
	if av == nil {
		return nil
	}
	if av.S != nil {
		return map[string]interface{}{dynamoDataTypeString: *av.S}
	}
	if av.N != nil {
		return map[string]interface{}{dynamoDataTypeNumber: *av.N}
	}
	if av.B != nil {
		return map[string]interface{}{dynamoDataTypeBinary: string(av.B)}
	}
	if av.BOOL != nil {
		return map[string]interface{}{dynamoDataTypeBoolean: *av.BOOL}
	}
	if av.NULL != nil {
		return map[string]interface{}{dynamoDataTypeNull: *av.NULL}
	}
	if av.L != nil {
		items := make([]interface{}, 0, len(av.L))
		for _, item := range av.L {
			items = append(items, attributeValueToJSON(item))
		}
		return map[string]interface{}{dynamoDataTypeList: items}
	}
	if av.M != nil {
		m := make(map[string]interface{})
		for k, v := range av.M {
			m[k] = attributeValueToJSON(v)
		}
		return map[string]interface{}{dynamoDataTypeMap: m}
	}
	if av.SS != nil {
		items := make([]string, 0, len(av.SS))
		for _, s := range av.SS {
			if s != nil {
				items = append(items, *s)
			}
		}
		return map[string]interface{}{dynamoDataTypeStringSet: items}
	}
	if av.NS != nil {
		items := make([]string, 0, len(av.NS))
		for _, n := range av.NS {
			if n != nil {
				items = append(items, *n)
			}
		}
		return map[string]interface{}{dynamoDataTypeNumberSet: items}
	}
	if av.BS != nil {
		items := make([]string, 0, len(av.BS))
		for _, b := range av.BS {
			if b != nil {
				items = append(items, string(b))
			}
		}
		return map[string]interface{}{dynamoDataTypeBinarySet: items}
	}
	return nil
}

func resourceAlicloudPolarDBDynamoItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDynamoItemCreate,
		Read:   resourceAlicloudPolarDBDynamoItemRead,
		Update: resourceAlicloudPolarDBDynamoItemUpdate,
		Delete: resourceAlicloudPolarDBDynamoItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The PolarDB DynamoDB-compatible endpoint URL.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The account name for PolarDB DynamoDB authentication. If not set, the provider's access key will be used.",
			},
			"account_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The authentication password for PolarDB DynamoDB. If not set, the provider's secret key will be used.",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the PolarDB cluster where the DynamoDB table resides.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the DynamoDB-compatible table.",
			},
			"hash_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The partition key (hash key) attribute name of the item.",
			},
			"range_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The sort key (range key) attribute name of the item.",
			},
			"item": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateDynamoItemJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
				Description:      "JSON representation of the DynamoDB item attributes.",
			},
		},
	}
}

func validateDynamoItemJSON(v interface{}, k string) (ws []string, errors []error) {
	jsonStr := v.(string)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		errors = append(errors, fmt.Errorf("invalid JSON format for %q: %w", k, err))
		return
	}
	return
}

func suppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	var oldMap, newMap map[string]interface{}
	if err := json.Unmarshal([]byte(old), &oldMap); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newMap); err != nil {
		return false
	}
	oldJSON, _ := json.Marshal(oldMap)
	newJSON, _ := json.Marshal(newMap)
	return string(oldJSON) == string(newJSON)
}

func resourceAlicloudPolarDBDynamoItemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	endpoint := d.Get("endpoint").(string)
	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	itemJSON := d.Get("item").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	var itemAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &itemAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_dynamo_item", "ParseItem", ProviderERROR)
	}

	// Convert JSON attributes to DynamoDB AttributeValue map
	dynamoItem := make(map[string]*dynamodb.AttributeValue)
	for k, v := range itemAttrs {
		dynamoItem[k] = jsonToAttributeValue(v)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
			TableName:                aws.String(tableName),
			Item:                     dynamoItem,
			ConditionExpression:      aws.String("attribute_not_exists(#hk)"),
			ExpressionAttributeNames: map[string]*string{"#hk": aws.String(hashKey)},
		})
		if err != nil {
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("PutItem", nil, nil, nil)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_dynamo_item", "PutItem", AlibabaCloudSdkGoERROR)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	rangeKey := d.Get("range_key").(string)
	resourceID := buildDynamoItemResourceID(dbClusterId, tableName, hashKey, rangeKey, itemAttrs)
	d.SetId(resourceID)

	return resourceAlicloudPolarDBDynamoItemRead(d, meta)
}

func resourceAlicloudPolarDBDynamoItemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		parts, err = ParseResourceId(d.Id(), 3)
		if err != nil {
			return WrapError(err)
		}
	}

	dbClusterId := parts[0]
	tableName := parts[1]
	hashKeyValue := parts[2]
	var rangeKeyValue string
	if len(parts) > 3 {
		rangeKeyValue = parts[3]
	}

	endpoint := d.Get("endpoint").(string)
	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	// On import only the resource ID is available, so resolve the DynamoDB-compatible
	// endpoint address from the cluster to rebuild the connection info.
	if endpoint == "" {
		polarDBService := PolarDBService{client}
		endpoint, err = polarDBService.DescribePolarDBDynamoEndpointAddress(dbClusterId)
		if err != nil {
			return WrapError(err)
		}
		d.Set("endpoint", endpoint)
	}

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	// On import the key attribute names are unknown (the ID only carries the key
	// values), so resolve them from the table's key schema.
	if hashKey == "" {
		tableOutput, describeErr := dynamoClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if describeErr != nil {
			return WrapErrorf(describeErr, DefaultErrorMsg, d.Id(), "DescribeTable", AlibabaCloudSdkGoERROR)
		}
		if tableOutput.Table != nil {
			for _, ks := range tableOutput.Table.KeySchema {
				switch aws.StringValue(ks.KeyType) {
				case "HASH":
					hashKey = aws.StringValue(ks.AttributeName)
				case "RANGE":
					rangeKey = aws.StringValue(ks.AttributeName)
				}
			}
		}
		d.Set("hash_key", hashKey)
		if rangeKey != "" {
			d.Set("range_key", rangeKey)
		}
	}

	// Build key
	key := map[string]*dynamodb.AttributeValue{
		hashKey: extractAttributeValueTyped(hashKeyValue),
	}
	if rangeKey != "" && rangeKeyValue != "" {
		key[rangeKey] = extractAttributeValueTyped(rangeKeyValue)
	}

	var output *dynamodb.GetItemOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		output, err = dynamoClient.GetItem(&dynamodb.GetItemInput{
			TableName:      aws.String(tableName),
			Key:            key,
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("GetItem", output, nil, nil)
		return nil
	})

	if err != nil {
		if isDynamoNotFoundError(err) {
			log.Printf("[WARN] DynamoDB Item (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetItem", AlibabaCloudSdkGoERROR)
	}

	if len(output.Item) == 0 {
		log.Printf("[WARN] DynamoDB Item (%s) not found in response, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	// Convert AttributeValue map back to JSON
	jsonItem := make(map[string]interface{})
	for k, v := range output.Item {
		jsonItem[k] = attributeValueToJSON(v)
	}

	itemJSON, err := json.Marshal(jsonItem)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "MarshalItem", ProviderERROR)
	}

	d.Set("db_cluster_id", dbClusterId)
	d.Set("table_name", tableName)
	d.Set("item", string(itemJSON))

	newResourceID := buildDynamoItemResourceID(dbClusterId, tableName, hashKey, rangeKey, jsonItem)
	if newResourceID != d.Id() {
		d.SetId(newResourceID)
	}

	return nil
}

func resourceAlicloudPolarDBDynamoItemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if !d.HasChange("item") {
		return nil
	}

	endpoint := d.Get("endpoint").(string)
	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	_, newItemJSON := d.GetChange("item")
	var newAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(newItemJSON.(string)), &newAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseNewItem", ProviderERROR)
	}

	// Use PutItem to replace the entire item
	dynamoItem := make(map[string]*dynamodb.AttributeValue)
	for k, v := range newAttrs {
		dynamoItem[k] = jsonToAttributeValue(v)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      dynamoItem,
		})
		if err != nil {
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PutItem", AlibabaCloudSdkGoERROR)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	newResourceID := buildDynamoItemResourceID(dbClusterId, tableName, hashKey, rangeKey, newAttrs)
	d.SetId(newResourceID)

	return resourceAlicloudPolarDBDynamoItemRead(d, meta)
}

func resourceAlicloudPolarDBDynamoItemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	endpoint := d.Get("endpoint").(string)
	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	itemJSON := d.Get("item").(string)
	var itemAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &itemAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseItem", ProviderERROR)
	}

	key := map[string]*dynamodb.AttributeValue{
		hashKey: jsonToAttributeValue(itemAttrs[hashKey]),
	}
	if rangeKey != "" {
		key[rangeKey] = jsonToAttributeValue(itemAttrs[rangeKey])
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key:       key,
		})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteItem", nil, nil, nil)
		return nil
	})

	if err != nil {
		if isDynamoNotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteItem", AlibabaCloudSdkGoERROR)
	}

	return nil
}

// Helper functions

func buildDynamoItemResourceID(dbClusterId, tableName, hashKey, rangeKey string, attrs map[string]interface{}) string {
	idParts := []string{dbClusterId, tableName}

	if hv, ok := attrs[hashKey]; ok {
		idParts = append(idParts, extractAttributeValueString(hv))
	}

	if rangeKey != "" {
		if rv, ok := attrs[rangeKey]; ok {
			idParts = append(idParts, extractAttributeValueString(rv))
		}
	}

	return strings.Join(idParts, COLON_SEPARATED)
}

func extractAttributeValueTyped(v interface{}) *dynamodb.AttributeValue {
	if m, ok := v.(map[string]interface{}); ok {
		return jsonToAttributeValue(m)
	}
	return &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%v", v))}
}

func extractAttributeValueString(v interface{}) string {
	if m, ok := v.(map[string]interface{}); ok {
		for _, val := range m {
			return fmt.Sprintf("%v", val)
		}
	}
	return fmt.Sprintf("%v", v)
}

func keysEqual(key1, key2 map[string]interface{}) bool {
	if len(key1) != len(key2) {
		return false
	}
	for k, v1 := range key1 {
		v2, ok := key2[k]
		if !ok {
			return false
		}
		s1 := extractAttributeValueString(v1)
		s2 := extractAttributeValueString(v2)
		if s1 != s2 {
			return false
		}
	}
	return true
}
