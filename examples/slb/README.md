### SLB Example

The example create SLB and additional listener, the listener parameter following:

### SLB Listener parameter describe
listener parameter | support protocol | value range | remark |
------------- | ------------- | ------------- |  ------------- |
backend_port | http & https & tcp & udp | 1-65535 | the ecs instance port |
frontend_port | http & https & tcp & udp | 1-65535 | the slb linstener port |
protocol | http & https & tcp & udp | http or https or tcp or udp | |
bandwidth | http & https & tcp & udp | -1 / 1-1000 | |
scheduler | http & https & tcp & udp | wrr or wlc | |
sticky_session | http & https | on or off | |
sticky_session_type | http & https | insert or server | if sticky_session is on, the value must have|
cookie_timeout | http & https | 1-86400  | if sticky_session is on and sticky_session_type is insert, the value must have|
cookie | http & https |   | if sticky_session is on and sticky_session_type is server, the value must have|
persistence_timeout | tcp & udp | 0-3600 | |
health_check | http & https | on or off | | TCP and UDP listener's HealthCheck is always on
health_check_type | tcp | tcp or http | if health_check is on, the value must have |
health_check_domain | http & https & tcp | | one string which length is 1-80 and only allow letters, digits, '-' and '.' characters. When it is not set or empty,  Server Load Balancer uses the private network IP address of each backend server as Domain used for health check  |
health_check_uri | http & https & tcp |  | example: /aliyun. if health_check is on, the value must have . Default to "/"|
health_check_connect_port | http & https & tcp & udp | 1-65535 | If the parameter is not set, the backend server port (BackendServerPort) will be used. |
healthy_threshold | http & https & tcp & udp | 1-10 | default to 3 when the health_check is on |
unhealthy_threshold | http & https & tcp & udp | 1-10 | default to 3 when the health_check is on |
health_check_timeout | http & https & tcp & udp | 1-300 | default to 5 when the health_check is on |
health_check_interval | http & https & tcp & udp | 1-50 | default to 2 when the health_check is on |
health_check_http_code | http & https & tcp | http_2xx,http_3xx,http_4xx,http_5xx | default to http_2xx when the health_check is on |
ssl_certificate_id | https |  |  |
acl_status | http & https & tcp & udp | on or off | default to on |
acl_type   | http & https & tcp & udp | white or black |  |
acl_id     | http & https & tcp & udp | the id of resource alicloud_slb_acl|  |
established_timeout | tcp       | 10-900|tcp listener's EstablishedTimeout for established connection idle timeout.|
idle_timeout |http & https      | 1-60  | http/https listener's IdleTimeout for established connection idle timeout. defaut to 15.|
request_timeout |http & https   | 1-180 | http/https listener's RequestTimeout for request which does not get response from backend timeout. defaut to 60.|
enable_http2    |https          | on or off | default to on|
tls_cipher_policy |https        |  tls_cipher_policy_1_0, tls_cipher_policy_1_1, tls_cipher_policy_1_2, tls_cipher_policy_1_2_strict | default to tls_cipher_policy_1_0 |

### Get up and running

* Planning phase

		terraform plan 

* Apply phase

		terraform apply 


* Destroy 

		terraform destroy