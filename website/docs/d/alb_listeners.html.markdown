---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listeners"
sidebar_current: "docs-alicloud-datasource-alb-listeners"
description: |- 
    Provides a list of Application Load Balancer (ALB) Listeners to the user.
---

# alicloud\_alb\_listeners

This data source provides the Application Load Balancer (ALB) Listeners of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_listeners" "ids" {
  ids = ["example_id"]
}
output "alb_listener_id_1" {
  value = data.alicloud_alb_listeners.ids.listeners.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Listener IDs.
* `listener_ids` - (Optional, ForceNew) The listener ids.
* `listener_protocol` - (Optional, ForceNew) Snooping Protocols. Valid Values: `HTTP`, `HTTPS` Or `QUIC`.
* `load_balancer_ids` - (Optional, ForceNew) The load balancer ids.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The state of the listener. Valid Values: `Running` Or `Stopped`. `Running`: The listener is running. `Stopped`: The listener is stopped.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `max_results` - This Request Returned by the Maximum Number of Records.
* `next_token` - The Current Call Returns to the Position of the Set to Null Represents the Data Has Been Read to the End of.
* `listeners` - A list of Alb Listeners. Each element contains the following attributes:
    * `access_log_record_customized_headers_enabled` - Indicates whether the access log has a custom header field. Valid values: true and false. Default value: false.

    -> **NOTE:** Only Instances outside the Security Group to Access the Log Switch **accesslogenabled** Open, in Order to Set This Parameter to the **True**.
    * `access_log_tracing_config` - Xtrace Configuration Information.
    * `tracing_sample` - Xtrace Sampling Rate. Value: **1~10000**.

    -> **NOTE:** This attribute is valid when **tracingenabled** is **true**.
    * `tracing_type` - Xtrace Type Value Is **Zipkin**.

    -> **NOTE:** This attribute is valid when **tracingenabled** is **true**.
    * `tracing_enabled` - Xtrace Function. Value: True Or False. Default Value: False.
  
    -> **NOTE:** Only Instances outside the Security Group to Access the Log Switch **accesslogenabled** Open, in Order to Set This Parameter to the **True**.

    * `acl_config` - The configurations of the access control lists (ACLs).
        * `acl_relations` - The ACLs that are associated with the listener.
            * `acl_id` - Snooping Binding of the Access Policy Group ID List.
            * `status` - The association status between the ACL and the listener.  Valid values: `Associating`, `Associated` Or `Dissociating`. `Associating`: The ACL is being associated with the listener. `Associated`: The ACL is associated with the listener. `Dissociating`: The ACL is being disassociated from the listener.
        * `acl_type` - The type of the ACL. Valid values: `White` Or `Black`. `White`: specifies the ACL as a whitelist. Only requests from the IP addresses or CIDR blocks in the ACL are forwarded. Whitelists apply to scenarios where only specific IP addresses are allowed to access an application. Risks may occur if the whitelist is improperly set. After you set a whitelist for an Application Load Balancer (ALB) listener, only requests from IP addresses that are added to the whitelist are distributed by the listener. If the whitelist is enabled without IP addresses specified, the ALB listener does not forward requests. `Black`: All requests from the IP addresses or CIDR blocks in the ACL are denied. The blacklist is used to prevent specified IP addresses from accessing an application. If the blacklist is enabled but the corresponding ACL does not contain IP addresses, the ALB listener forwards all requests.
    * `certificates` - The Certificate List.
        * `certificate_id` - The ID of the Certificate.
    * `default_actions` - The Default Rule Action List. 		
        * `forward_group_config` - The configuration of the forwarding rule action. This parameter is required if the Type parameter is set to FowardGroup.
            *  `server_group_tuples` - The destination server group to which requests are forwarded.
                * `server_group_id` - The ID of the destination server group to which requests are forwarded.
        * `type` - Action Type. The value is set to ForwardGroup. It indicates that requests are forwarded to multiple vServer groups.	
    * `gzip_enabled` - Whether to Enable Gzip Compression, as a Specific File Type on a Compression. Valid Values: `True` Or `False`. Default Value: `True`. 	
    * `http2_enabled` - Whether to Enable HTTP/2 Features. Valid Values: `True` Or `False`. Default Value: `True`.

    -> **NOTE:** The attribute is valid when the attribute `ListenerProtocol` is `HTTPS`.
    * `id` - The ID of the Listener. 	
    * `idle_timeout` - Specify the Connection Idle Timeout Value: `1` to `60`. Unit: Seconds.
    * `listener_description` - Set the IP Address of the Listened Description. Length Is from 2 to 256 Characters. 	
    * `listener_id` - on Behalf of the Resource Level Id of the Resources Property Fields. 	
    * `listener_port` - The ALB Instance Front-End, and Those of the Ports Used. Value: `1~65535`. 	
    * `listener_protocol` - Snooping Protocols. Valid Values: `HTTP`, `HTTPS` Or `QUIC`. 	
    * `load_balancer_id` - The ALB Instance Id. 	
    * `quic_config` - Configuration Associated with the QuIC Listening. 		
    * `quic_listener_id` - The ID of the QUIC listener to be associated. If QuicUpgradeEnabled is set to true, this parameter is required. Only HTTPS listeners support this parameter. 		
    * `quic_upgrade_enabled` - Indicates whether quic upgrade is enabled. Valid values: true and false. Default value: false. 	
    * `request_timeout` - The Specified Request Timeout Time. Value: `1` to `180`. Unit: Seconds. Default Value: 60. If the Timeout Time Within the Back-End Server Has Not Answered the ALB Will Give up Waiting, the Client Returns the HTTP 504 Error Code. 	
    * `security_policy_id` - Security Policy.
  
    -> **NOTE:** The attribute is valid when the attribute `ListenerProtocol` is `HTTPS`.
    * `xforwarded_for_config` - xforwardfor Related Attribute Configuration. 		
        * `xforwardedforclientcert_issuerdnenabled` - Indicates Whether the `X-Forwarded-Clientcert-issuerdn` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate after the Manifests Are Signed, the Publisher Information. 		
        * `xforwardedforclientcertclientverifyalias` - The Custom Header Field Names Only When `xforwardedforclientcertclientverifyenabled` Has a Value of True, this Value Will Not Take Effect until.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits. 		
        * `xforwardedforclientcertclientverifyenabled` - Indicates Whether the `X-Forwarded-Clientcert-clientverify` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate to Verify the Results. 		
        * `xforwardedforclientcertfingerprintalias` - The Custom Header Field Names Only When `xforwardedforclientcertfingerprintenabled`, Which Evaluates to True When the Entry into Force of.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits. 		
        * `xforwardedforclientcertfingerprintenabled` - Indicates Whether the `X-Forwarded-Clientcert-fingerprint` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate Fingerprint Value. 		
        * `xforwardedforclientcertsubjectdnalias` - The name of the custom header. This parameter is valid only if `xforwardedforclientcertsubjectdnenabled` is set to true. The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits. 		
        * `xforwardedforclientcertsubjectdnenabled` - Specifies whether to use the `X-Forwarded-Clientcert-subjectdn` header field to obtain information about the owner of the ALB client certificate. Valid values: true and false. Default value: false. 	
        * `xforwardedforclientcert_issuerdnalias` - The Custom Header Field Names Only When `xforwardedforclientcert_issuerdnenabled`, Which Evaluates to True When the Entry into Force of. 		
        * `xforwardedforprotoenabled` - Indicates Whether the X-Forwarded-Proto Header Field Is Used to Obtain the Server Load Balancer Instance Snooping Protocols. 		
        * `xforwardedforenabled` - Indicates whether the X-Forwarded-For header field is used to obtain the real IP address of tqhe client. Valid values: true and false. Default value: true. 		
        * `xforwardedforslbidenabled` - Indicates whether the SLB-ID header field is used to obtain the ID of the ALB instance. Valid values: true and false. Default value: false. 		
        * `xforwardedforslbportenabled` - Indicates Whether the X-Forwarded-Port Header Field Is Used to Obtain the Server Load Balancer Instance Listening Port. 		
        * `xforwardedforclientsrcportenabled` - Indicates Whether the X-Forwarded-Client-Port Header Field Is Used to Obtain Access to Server Load Balancer Instances to the Client, and Those of the Ports.
