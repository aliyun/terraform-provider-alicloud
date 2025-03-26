---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_website"
description: |-
  Provides a Alicloud OSS Bucket Website resource.
---

# alicloud_oss_bucket_website

Provides a OSS Bucket Website resource.

the static website configuration and mirror configuration of the bucket.

For information about OSS Bucket Website and how to use it, see [What is Bucket Website](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketwebsite).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_website&exampleId=49cb70e8-a10e-3389-bebc-3698726a8a22625c0946&activeTab=example&spm=docs.r.oss_bucket_website.0.49cb70e8a1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_uuid" "default" {
}

resource "alicloud_oss_bucket" "defaultnVj9x3" {
  bucket        = "${var.name}-${random_uuid.default.result}"
  storage_class = "Standard"
  lifecycle {
    ignore_changes = [website]
  }
}

resource "alicloud_oss_bucket_website" "default" {
  index_document {
    suffix          = "index.html"
    support_sub_dir = "true"
    type            = "0"
  }

  error_document {
    key         = "error.html"
    http_status = "404"
  }

  bucket = alicloud_oss_bucket.defaultnVj9x3.bucket
  routing_rules {
    routing_rule {
      rule_number = "1"
      condition {
        http_error_code_returned_equals = "404"
      }

      redirect {
        protocol           = "https"
        http_redirect_code = "305"
        redirect_type      = "AliCDN"
        host_name          = "www.alicdn-master.com"
      }

      lua_config {
        script = "example.lua"
      }

    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `error_document` - (Optional, List) The container that holds the error page configuration information. See [`error_document`](#error_document) below.
* `index_document` - (Optional, List) Static Website Default Home Page Configuration See [`index_document`](#index_document) below.
* `routing_rules` - (Optional, List) The container that holds the jump rule or the mirroring back-to-origin rule. See [`routing_rules`](#routing_rules) below.

### `error_document`

The error_document supports the following:
* `http_status` - (Optional, Int) The HTTP status code when the error page is returned. The default 404.
* `key` - (Optional) The error page file. If the Object accessed does not exist, this error page is returned.

### `index_document`

The index_document supports the following:
* `suffix` - (Optional) The default home page.
* `support_sub_dir` - (Optional) Whether to jump to the default home page of a subdirectory when accessing a subdirectory.
* `type` - (Optional) After the default homepage is set, the behavior when an Object that ends with a non-forward slash (/) is accessed and the Object does not exist.

### `routing_rules`

The routing_rules supports the following:
* `routing_rule` - (Optional, List) Specify a jump rule or a mirroring back-to-origin rule, with a maximum of 20 routing rules. See [`routing_rule`](#routing_rules-routing_rule) below.

### `routing_rules-routing_rule`

The routing_rules-routing_rule supports the following:
* `condition` - (Optional, List) Save the criteria that the rule needs to match. See [`condition`](#routing_rules-routing_rule-condition) below.
* `lua_config` - (Optional, List) The Lua script configuration to be executed. See [`lua_config`](#routing_rules-routing_rule-lua_config) below.
* `redirect` - (Optional, List) Specifies the action to perform after this rule is matched. See [`redirect`](#routing_rules-routing_rule-redirect) below.
* `rule_number` - (Optional, Int) The sequence number of the matching and executing jump rules. OSS matches rules according to this sequence number. If the match is successful, the rule is executed and subsequent rules are not executed.

### `routing_rules-routing_rule-condition`

The routing_rules-routing_rule-condition supports the following:
* `http_error_code_returned_equals` - (Optional) When the specified Object is accessed, this status is returned to match this rule. This field must be 404 when the jump rule is mirrored back to the source.
* `include_headers` - (Optional, List) This rule can only be matched if the request contains the specified Header and the value is the specified value. You can specify up to 10 containers. See [`include_headers`](#routing_rules-routing_rule-condition-include_headers) below.
* `key_prefix_equals` - (Optional) Only objects that match this prefix can match this rule.
* `key_suffix_equals` - (Optional) Only objects that match this suffix can match this rule.

### `routing_rules-routing_rule-lua_config`

The routing_rules-routing_rule-lua_config supports the following:
* `script` - (Optional) The Lua script name.

### `routing_rules-routing_rule-redirect`

The routing_rules-routing_rule-redirect supports the following:
* `enable_replace_prefix` - (Optional) If this field is set to true, the prefix of Object is replaced with the value specified by ReplaceKeyPrefixWith. If this field is not specified or is blank, the Object prefix is truncated.
* `host_name` - (Optional) The domain name during the jump. The domain name must comply with the domain name specification.
* `http_redirect_code` - (Optional) The status code returned during the jump. It takes effect only when the RedirectType is set to External or AliCDN.
* `mirror_allow_get_image_info` - (Optional) Image back-to-source allows getting Image information
* `mirror_allow_head_object` - (Optional) Whether to allow HeadObject in image back-to-source
* `mirror_allow_video_snapshot` - (Optional) Mirror back-to-source allows support for video frame truncation
* `mirror_async_status` - (Optional, Int) The status code of the mirror back-to-source trigger asynchronous pull mode.
* `mirror_auth` - (Optional, List) Image back Source station authentication information See [`mirror_auth`](#routing_rules-routing_rule-redirect-mirror_auth) below.
* `mirror_check_md5` - (Optional) Whether to check the MD5 of the source body. It takes effect only when the RedirectType is set to Mirror.
* `mirror_dst_region` - (Optional) Mirrored back-to-source high-speed Channel vpregion
* `mirror_dst_slave_vpc_id` - (Optional) Mirroring back-to-source high-speed Channel standby station VpcId
* `mirror_dst_vpc_id` - (Optional) Mirror back-to-source high-speed Channel VpcId
* `mirror_follow_redirect` - (Optional) If the result of the image back-to-source acquisition is 3xx, whether to continue to jump to the specified Location to obtain data. It takes effect only when the RedirectType is set to Mirror.
* `mirror_headers` - (Optional, List) Specifies the Header carried when the image returns to the source. It takes effect only when the RedirectType is set to Mirror. See [`mirror_headers`](#routing_rules-routing_rule-redirect-mirror_headers) below.
* `mirror_is_express_tunnel` - (Optional) Whether it is a mirror back-to-source high-speed Channel
* `mirror_multi_alternates` - (Optional, Computed, List) Mirror back-to-source multi-source station configuration container. **NOTE:**: If you want to clean one configuration, you must set the configuration to empty value, removing from code cannot make effect. See [`mirror_multi_alternates`](#routing_rules-routing_rule-redirect-mirror_multi_alternates) below.
* `mirror_pass_original_slashes` - (Optional) Transparent transmission/to source Station
* `mirror_pass_query_string` - (Optional) Same as PassQueryString and takes precedence over PassQueryString. It takes effect only when the RedirectType is set to Mirror.
* `mirror_proxy_pass` - (Optional) Whether mirroring back to source does not save data
* `mirror_return_headers` - (Optional, Computed, List) The container that saves the image back to the source and returns the response header rule. **NOTE:**: If you want to clean one configuration, you must set the configuration to empty value, removing from code cannot make effect. See [`mirror_return_headers`](#routing_rules-routing_rule-redirect-mirror_return_headers) below.
* `mirror_role` - (Optional) Roles used when mirroring back-to-source
* `mirror_save_oss_meta` - (Optional) Mirror back-to-source back-to-source OSS automatically saves user metadata
* `mirror_sni` - (Optional) Transparent transmission of SNI
* `mirror_switch_all_errors` - (Optional) It is used to judge the status of active-standby switching. The judgment logic of active-standby switching is that the source station returns an error. If MirrorSwitchAllErrors is true, it is considered a failure except the following status code: 200,206,301,302,303,307,404; If false, only the source Station Returns 5xx or times out is considered a failure.
* `mirror_taggings` - (Optional, Computed, List) Save the label according to the parameters when saving the file from the mirror back to the source. **NOTE:**: If you want to clean one configuration, you must set the configuration to empty value, removing from code cannot make effect. See [`mirror_taggings`](#routing_rules-routing_rule-redirect-mirror_taggings) below.
* `mirror_tunnel_id` - (Optional) Mirror back-to-source leased line back-to-source tunnel ID
* `mirror_url` - (Optional) The address of the origin of the image. It takes effect only when the RedirectType is set to Mirror. The origin address must start with http:// or https:// and end with a forward slash (/). OSS takes the Object name after the Origin address to form the origin URL.
* `mirror_url_probe` - (Optional) Mirror back-to-source Master-backup back-to-source switching decision URL
* `mirror_url_slave` - (Optional) Mirror back-to-source primary backup back-to-source backup station URL
* `mirror_user_last_modified` - (Optional) Whether the source station LastModifiedTime is used for the image back-to-source save file.
* `mirror_using_role` - (Optional) Whether to use role for mirroring back to source
* `pass_query_string` - (Optional) Whether to carry the request parameters when executing the jump or mirror back-to-source rule. Did the user carry the request parameters when requesting OSS? a = B & c = d, and set PassQueryString to true. If the rule is a 302 jump, this request parameter is added to the Location header of the jump. For example Location:example.com? a = B & c = d, and the jump type is mirrored back-to-origin, this request parameter is also carried in the back-to-origin request initiated. Values: true, false (default)
* `protocol` - (Optional) The protocol at the time of the jump. It takes effect only when the RedirectType is set to External or AliCDN.
* `redirect_type` - (Optional) Specifies the type of jump. The value range is as follows: Mirror: Mirror back to the source. External: External redirects, that is, OSS returns a 3xx request to redirect to another address. AliCDN: Alibaba Cloud CDN jump, mainly used for Alibaba Cloud CDN. Unlike External, OSS adds an additional Header. After recognizing this Header, Alibaba Cloud CDN redirects the data to the specified address and returns the obtained data to the user instead of returning the 3xx Redirection request to the user.
* `replace_key_prefix_with` - (Optional) The prefix of the Object name will be replaced with this value during Redirect. If the prefix is empty, this string is inserted in front of the Object name.
* `replace_key_with` - (Optional) During redirection, the Object name is replaced with the value specified by ReplaceKeyWith. You can set variables in ReplaceKeyWith. Currently, the supported variable is ${key}, which indicates the name of the Object in the request.
* `transparent_mirror_response_codes` - (Optional) Mirror back-to-source transparent source station response code list

### `routing_rules-routing_rule-redirect-mirror_auth`

The routing_rules-routing_rule-redirect-mirror_auth supports the following:
* `access_key_id` - (Optional) Mirror back-to-source source Station back-to-source AK
* `access_key_secret` - (Optional) Mirroring back to the source station back to the source SK will be automatically desensitized when obtaining the configuration.
* `auth_type` - (Optional) Authentication type of mirror return Source
* `region` - (Optional) Signature Region

### `routing_rules-routing_rule-redirect-mirror_headers`

The routing_rules-routing_rule-redirect-mirror_headers supports the following:
* `pass` - (Optional, List) Pass through the specified Header to the source site. It takes effect only when the RedirectType is set to Mirror. Each Header is up to 1024 bytes in length and has A character set of 0 to 9, a to Z, A to z, and dashes (-).
* `pass_all` - (Optional) Indicates whether other headers except the following headers are transmitted to the source site. It takes effect only when the RedirectType is set to Mirror. content-length, authorization2, authorization, range, date, and other headers Headers whose names start with oss-/x-oss-/x-drs-
* `remove` - (Optional, List) Do not pass the specified Header to the source site. It takes effect only when the RedirectType is set to Mirror. Each Header is up to 1024 bytes in length and has A character set of 0 to 9, a to Z, A to z, and dashes (-).
* `set` - (Optional, List) Set a Header to send to the source site. Regardless of whether the request contains the specified Header, these headers will be set when returning to the source site. It takes effect only when the RedirectType is set to Mirror. See [`set`](#routing_rules-routing_rule-redirect-mirror_headers-set) below.

### `routing_rules-routing_rule-redirect-mirror_multi_alternates`

The routing_rules-routing_rule-redirect-mirror_multi_alternates supports the following:
* `mirror_multi_alternate` - (Optional, List) Mirror back-to-source multi-source station configuration list See [`mirror_multi_alternate`](#routing_rules-routing_rule-redirect-mirror_multi_alternates-mirror_multi_alternate) below.

### `routing_rules-routing_rule-redirect-mirror_return_headers`

The routing_rules-routing_rule-redirect-mirror_return_headers supports the following:
* `return_header` - (Optional, List) The list of response header rules for mirroring back-to-source return. See [`return_header`](#routing_rules-routing_rule-redirect-mirror_return_headers-return_header) below.

### `routing_rules-routing_rule-redirect-mirror_taggings`

The routing_rules-routing_rule-redirect-mirror_taggings supports the following:
* `taggings` - (Optional, List) Image back-to-source save label rule list See [`taggings`](#routing_rules-routing_rule-redirect-mirror_taggings-taggings) below.

### `routing_rules-routing_rule-redirect-mirror_taggings-taggings`

The routing_rules-routing_rule-redirect-mirror_taggings-taggings supports the following:
* `key` - (Optional) The tag key corresponding to the current rule
* `value` - (Optional) Rules for Saving Label Values

### `routing_rules-routing_rule-redirect-mirror_return_headers-return_header`

The routing_rules-routing_rule-redirect-mirror_return_headers-return_header supports the following:
* `key` - (Optional) Response header corresponding to the current rule
* `value` - (Optional) Rules that return response header values

### `routing_rules-routing_rule-redirect-mirror_multi_alternates-mirror_multi_alternate`

The routing_rules-routing_rule-redirect-mirror_multi_alternates-mirror_multi_alternate supports the following:
* `mirror_multi_alternate_dst_region` - (Optional) Mirroring back-to-source multi-station Region
* `mirror_multi_alternate_number` - (Optional, Int) Image back-to-source multi-source station serial number
* `mirror_multi_alternate_url` - (Optional) Mirroring back-to-source multi-source site URL
* `mirror_multi_alternate_vpc_id` - (Optional) Mirroring back-to-source multi-source VpcId

### `routing_rules-routing_rule-redirect-mirror_headers-set`

The routing_rules-routing_rule-redirect-mirror_headers-set supports the following:
* `key` - (Optional) Set the key of the Header, up to 1024 bytes, and the character set is the same as that of Pass. It takes effect only when the RedirectType is set to Mirror.
* `value` - (Optional) Set the value of the Header to 1024 bytes at most. \r\n. It takes effect only when the RedirectType is set to Mirror.

### `routing_rules-routing_rule-condition-include_headers`

The routing_rules-routing_rule-condition-include_headers supports the following:
* `ends_with` - (Optional) This rule can only be matched if the request contains the Header specified by Key and the value ends with this value.
* `equals` - (Optional) This rule can only be matched if the request contains the Header specified by Key and the value is the specified value.
* `key` - (Optional) This rule can only be matched if the request contains this Header and the value meets the conditions.
* `starts_with` - (Optional) This rule can only be matched if the request contains the Header specified by Key and the value starts with this value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Website.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Website.
* `update` - (Defaults to 5 mins) Used when update the Bucket Website.

## Import

OSS Bucket Website can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_website.example <id>
```