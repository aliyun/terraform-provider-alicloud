package alicloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

// formatDynamoNumber renders a JSON number as a plain decimal string. json.Unmarshal
// decodes every JSON number as float64, and fmt's default (%v/%g) switches to
// scientific notation for large/small magnitudes (1000000 -> "1e+06"), which would
// corrupt a DynamoDB N attribute. strconv.FormatFloat with 'f'/-1 keeps the full
// decimal form. Non-float values (already-string descriptor payloads) pass through.
func formatDynamoNumber(v interface{}) string {
	if f, ok := v.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", v)
}

// jsonToAttributeValue converts a JSON value to DynamoDB AttributeValue format.
func jsonToAttributeValue(value interface{}) types.AttributeValue {
	switch v := value.(type) {
	case map[string]interface{}:
		if len(v) == 1 {
			for key, val := range v {
				switch key {
				case dynamoDataTypeString:
					return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", val)}
				case dynamoDataTypeNumber:
					return &types.AttributeValueMemberN{Value: formatDynamoNumber(val)}
				case dynamoDataTypeBinary:
					decoded, decErr := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", val))
					if decErr != nil {
						return &types.AttributeValueMemberB{Value: []byte(fmt.Sprintf("%v", val))}
					}
					return &types.AttributeValueMemberB{Value: decoded}
				case dynamoDataTypeBoolean:
					boolVal := fmt.Sprintf("%v", val) == "true"
					return &types.AttributeValueMemberBOOL{Value: boolVal}
				case dynamoDataTypeNull:
					nullVal := fmt.Sprintf("%v", val) == "true"
					return &types.AttributeValueMemberNULL{Value: nullVal}
				case dynamoDataTypeList:
					if list, ok := val.([]interface{}); ok {
						items := make([]types.AttributeValue, 0, len(list))
						for _, item := range list {
							items = append(items, jsonToAttributeValue(item))
						}
						return &types.AttributeValueMemberL{Value: items}
					}
				case dynamoDataTypeMap:
					if m, ok := val.(map[string]interface{}); ok {
						items := make(map[string]types.AttributeValue)
						for k, mv := range m {
							items[k] = jsonToAttributeValue(mv)
						}
						return &types.AttributeValueMemberM{Value: items}
					}
				case dynamoDataTypeStringSet:
					if ss, ok := val.([]interface{}); ok {
						strs := make([]string, 0, len(ss))
						for _, s := range ss {
							strs = append(strs, fmt.Sprintf("%v", s))
						}
						return &types.AttributeValueMemberSS{Value: strs}
					}
				case dynamoDataTypeNumberSet:
					if ns, ok := val.([]interface{}); ok {
						nums := make([]string, 0, len(ns))
						for _, n := range ns {
							nums = append(nums, formatDynamoNumber(n))
						}
						return &types.AttributeValueMemberNS{Value: nums}
					}
				case dynamoDataTypeBinarySet:
					if bs, ok := val.([]interface{}); ok {
						bins := make([][]byte, 0, len(bs))
						for _, b := range bs {
							decoded, decErr := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", b))
							if decErr != nil {
								bins = append(bins, []byte(fmt.Sprintf("%v", b)))
							} else {
								bins = append(bins, decoded)
							}
						}
						return &types.AttributeValueMemberBS{Value: bins}
					}
				}
			}
		}
		// Regular map
		items := make(map[string]types.AttributeValue)
		for k, val := range v {
			items[k] = jsonToAttributeValue(val)
		}
		return &types.AttributeValueMemberM{Value: items}
	case []interface{}:
		items := make([]types.AttributeValue, 0, len(v))
		for _, item := range v {
			items = append(items, jsonToAttributeValue(item))
		}
		return &types.AttributeValueMemberL{Value: items}
	case string:
		return &types.AttributeValueMemberS{Value: v}
	case float64:
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(v, 'f', -1, 64)}
	case bool:
		return &types.AttributeValueMemberBOOL{Value: v}
	case nil:
		return &types.AttributeValueMemberNULL{Value: true}
	default:
		return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", v)}
	}
}

// attributeValueToJSON converts a DynamoDB AttributeValue to JSON-compatible format.
func attributeValueToJSON(av types.AttributeValue) interface{} {
	switch v := av.(type) {
	case *types.AttributeValueMemberS:
		return map[string]interface{}{dynamoDataTypeString: v.Value}
	case *types.AttributeValueMemberN:
		return map[string]interface{}{dynamoDataTypeNumber: v.Value}
	case *types.AttributeValueMemberB:
		return map[string]interface{}{dynamoDataTypeBinary: base64.StdEncoding.EncodeToString(v.Value)}
	case *types.AttributeValueMemberBOOL:
		return map[string]interface{}{dynamoDataTypeBoolean: v.Value}
	case *types.AttributeValueMemberNULL:
		return map[string]interface{}{dynamoDataTypeNull: v.Value}
	case *types.AttributeValueMemberL:
		items := make([]interface{}, 0, len(v.Value))
		for _, item := range v.Value {
			items = append(items, attributeValueToJSON(item))
		}
		return map[string]interface{}{dynamoDataTypeList: items}
	case *types.AttributeValueMemberM:
		m := make(map[string]interface{})
		for k, mv := range v.Value {
			m[k] = attributeValueToJSON(mv)
		}
		return map[string]interface{}{dynamoDataTypeMap: m}
	case *types.AttributeValueMemberSS:
		items := make([]string, 0, len(v.Value))
		items = append(items, v.Value...)
		return map[string]interface{}{dynamoDataTypeStringSet: items}
	case *types.AttributeValueMemberNS:
		items := make([]string, 0, len(v.Value))
		items = append(items, v.Value...)
		return map[string]interface{}{dynamoDataTypeNumberSet: items}
	case *types.AttributeValueMemberBS:
		items := make([]string, 0, len(v.Value))
		for _, b := range v.Value {
			items = append(items, base64.StdEncoding.EncodeToString(b))
		}
		return map[string]interface{}{dynamoDataTypeBinarySet: items}
	default:
		return nil
	}
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The PolarDB DynamoDB-compatible endpoint URL.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The account name for PolarDB DynamoDB authentication. If not set, it is resolved from the cluster's DynamoDB-type account automatically.",
			},
			"account_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The authentication password for PolarDB DynamoDB. If not set, it is resolved from the cluster's DynamoDB-type account automatically.",
			},
			"db_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The ID of the PolarDB cluster where the DynamoDB table resides.",
			},
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The name of the DynamoDB-compatible table.",
			},
			"hash_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The partition key (hash key) attribute name of the item.",
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
				ValidateFunc:     validatePolarDBDynamoItemJSON,
				DiffSuppressFunc: suppressPolarDBDynamoItemJSONDiffs,
				Description:      "JSON representation of the DynamoDB item attributes.",
			},
		},
	}
}

func validatePolarDBDynamoItemJSON(v interface{}, k string) (ws []string, errors []error) {
	jsonStr := v.(string)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		errors = append(errors, fmt.Errorf("invalid JSON format for %q: %w", k, err))
		return
	}
	return
}

func suppressPolarDBDynamoItemJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	eq, err := compareJsonTemplateAreEquivalent(old, new)
	if err != nil {
		return false
	}
	return eq
}

func resourceAlicloudPolarDBDynamoItemCreate(d *schema.ResourceData, meta interface{}) error {
	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	itemJSON := d.Get("item").(string)

	conn, err := resolvePolarDBDynamoConn(d, meta, d.Get("db_cluster_id").(string))
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	var itemAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &itemAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_dynamo_item", "ParseItem", ProviderERROR)
	}

	// Convert JSON attributes to DynamoDB AttributeValue map
	dynamoItem := make(map[string]types.AttributeValue)
	for k, v := range itemAttrs {
		dynamoItem[k] = jsonToAttributeValue(v)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName:                aws.String(tableName),
			Item:                     dynamoItem,
			ConditionExpression:      aws.String("attribute_not_exists(#hk)"),
			ExpressionAttributeNames: map[string]string{"#hk": hashKey},
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
	parts, err := ParseResourceIdWithEscaped(d.Id(), 4)
	if err != nil {
		parts, err = ParseResourceIdWithEscaped(d.Id(), 3)
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

	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)

	conn, err := resolvePolarDBDynamoConn(d, meta, dbClusterId)
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	// Prefer the typed key values from the item JSON in state; the raw ID segments
	// carry no type information, so a Number/Binary key rebuilt as a plain string
	// would miss the row and wrongly drop the resource from state.
	var itemAttrs map[string]interface{}
	if itemJSON := d.Get("item").(string); itemJSON != "" {
		if err := json.Unmarshal([]byte(itemJSON), &itemAttrs); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseItem", ProviderERROR)
		}
	}

	// On import the key attribute names (and the item itself) are unknown, so
	// resolve the key names and key attribute types from the table schema.
	keyAttrTypes := make(map[string]types.ScalarAttributeType)
	if hashKey == "" || itemAttrs == nil {
		tableOutput, describeErr := dynamoClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if describeErr != nil {
			return WrapErrorf(describeErr, DefaultErrorMsg, d.Id(), "DescribeTable", AlibabaCloudSdkGoERROR)
		}
		if tableOutput.Table != nil {
			for _, ks := range tableOutput.Table.KeySchema {
				switch ks.KeyType {
				case types.KeyTypeHash:
					if hashKey == "" {
						hashKey = aws.ToString(ks.AttributeName)
						d.Set("hash_key", hashKey)
					}
				case types.KeyTypeRange:
					if rangeKey == "" {
						rangeKey = aws.ToString(ks.AttributeName)
						d.Set("range_key", rangeKey)
					}
				}
			}
			for _, ad := range tableOutput.Table.AttributeDefinitions {
				keyAttrTypes[aws.ToString(ad.AttributeName)] = ad.AttributeType
			}
		}
	}

	buildKeyValue := func(name, rawValue string) types.AttributeValue {
		if itemAttrs != nil {
			if v, ok := itemAttrs[name]; ok {
				return extractAttributeValueTyped(v)
			}
		}
		switch keyAttrTypes[name] {
		case types.ScalarAttributeTypeN:
			return &types.AttributeValueMemberN{Value: rawValue}
		case types.ScalarAttributeTypeB:
			return &types.AttributeValueMemberB{Value: []byte(rawValue)}
		default:
			return &types.AttributeValueMemberS{Value: rawValue}
		}
	}

	// Build key
	key := map[string]types.AttributeValue{
		hashKey: buildKeyValue(hashKey, hashKeyValue),
	}
	if rangeKey != "" && rangeKeyValue != "" {
		key[rangeKey] = buildKeyValue(rangeKey, rangeKeyValue)
	}

	var output *dynamodb.GetItemOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		output, err = dynamoClient.GetItem(context.Background(), &dynamodb.GetItemInput{
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
	if !d.HasChange("item") {
		return nil
	}

	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)

	conn, err := resolvePolarDBDynamoConn(d, meta, d.Get("db_cluster_id").(string))
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	oldItemJSON, newItemJSON := d.GetChange("item")
	var newAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(newItemJSON.(string)), &newAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseNewItem", ProviderERROR)
	}
	var oldAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(oldItemJSON.(string)), &oldAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseOldItem", ProviderERROR)
	}

	// Use PutItem to replace the entire item
	dynamoItem := make(map[string]types.AttributeValue)
	for k, v := range newAttrs {
		dynamoItem[k] = jsonToAttributeValue(v)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
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

	// When the key values inside "item" change, PutItem writes a brand-new row and
	// the previous row would be orphaned outside Terraform's tracking, so delete it.
	oldKey := map[string]interface{}{hashKey: oldAttrs[hashKey]}
	newKey := map[string]interface{}{hashKey: newAttrs[hashKey]}
	if rangeKey != "" {
		oldKey[rangeKey] = oldAttrs[rangeKey]
		newKey[rangeKey] = newAttrs[rangeKey]
	}
	if !keysEqual(oldKey, newKey) {
		deleteKey := map[string]types.AttributeValue{
			hashKey: jsonToAttributeValue(oldAttrs[hashKey]),
		}
		if rangeKey != "" {
			deleteKey[rangeKey] = jsonToAttributeValue(oldAttrs[rangeKey])
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = dynamoClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key:       deleteKey,
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
		if err != nil && !isDynamoNotFoundError(err) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteOldItem", AlibabaCloudSdkGoERROR)
		}
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	newResourceID := buildDynamoItemResourceID(dbClusterId, tableName, hashKey, rangeKey, newAttrs)
	d.SetId(newResourceID)

	return resourceAlicloudPolarDBDynamoItemRead(d, meta)
}

func resourceAlicloudPolarDBDynamoItemDelete(d *schema.ResourceData, meta interface{}) error {
	tableName := d.Get("table_name").(string)
	hashKey := d.Get("hash_key").(string)
	rangeKey := d.Get("range_key").(string)

	conn, err := resolvePolarDBDynamoConn(d, meta, d.Get("db_cluster_id").(string))
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	itemJSON := d.Get("item").(string)
	var itemAttrs map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &itemAttrs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ParseItem", ProviderERROR)
	}

	key := map[string]types.AttributeValue{
		hashKey: jsonToAttributeValue(itemAttrs[hashKey]),
	}
	if rangeKey != "" {
		key[rangeKey] = jsonToAttributeValue(itemAttrs[rangeKey])
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
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
		idParts = append(idParts, EscapeColons(extractAttributeValueString(hv)))
	}

	if rangeKey != "" {
		if rv, ok := attrs[rangeKey]; ok {
			idParts = append(idParts, EscapeColons(extractAttributeValueString(rv)))
		}
	}

	return strings.Join(idParts, COLON_SEPARATED)
}

func extractAttributeValueTyped(v interface{}) types.AttributeValue {
	if m, ok := v.(map[string]interface{}); ok {
		return jsonToAttributeValue(m)
	}
	return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", v)}
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
