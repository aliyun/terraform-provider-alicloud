# Example: WordPress + MySQL on Alibaba Cloud Kubernetes

This example is heavily inspired by https://kubernetes.io/docs/tutorials/stateful-application/mysql-wordpress-persistent-volume/

It describes how to use local volume to deploy [WordPress](https://wordpress.org/)
and [MySQL](https://www.mysql.com/) on Alibaba Cloud Kubernetes Cluster.

## Used resources

### Kubernetes Provider

 - `kubernetes_persistent_volume_claim`
 - `kubernetes_persistent_volume`
 - `kubernetes_replication_controller`
 - `kubernetes_secret`
 - `kubernetes_service`

## Prerequsites

### Kubernetes

This example expects you to already have a running K8S cluster
and credentials set up in a config or environment variables.

If you have it not yet, you can use [kubernetes example](https://github.com/terraform-providers/terraform-provider-alicloud/tree/master/examples/kubernetes)
to create a new cluster and it can download automatically kube config into a file by setting `kube_config`, like `~/.kube/config`.

## How to

### Create

First we make sure kubernetes providers is downloaded and available

```sh
terraform init
```

then we carry on by creating the real infrastructure which requires
password for the MySQL server.

```sh
terraform apply -var 'mysql_password=Yourpassword'
```

You may also specify version of WordPress and/or MySQL

```sh
terraform apply \
	-var 'mysql_version=5.6' \
	-var 'wordpress_version=4.7.3' \
	-var 'mysql_password=Yourpassword'
```

After the `apply` operation has finished you should see output
in your console similar to the one below

```
...

Outputs:

slb_ip = 35.197.xx.xx
```

This is the IP address of your public load balancer
which exposes the Apache web server serving WordPress.
Open that IP in your browser to see the welcome page.

```sh
open "http://$(terraform output lb_ip)"
```

### Destroy

```
terraform destroy -var 'mysql_password=Yourpassword'
```