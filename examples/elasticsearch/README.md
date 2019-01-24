### ECS In VPC Example

The example launches Elasticsearch in VPC, vswitch_id parameter is the vswitch id from VPC. 
### Get up and running

* Planning phase

		terraform plan 
    		var.availability_zones
  				Enter a value: {var.availability_zones}  /*cn-beijing-b*/
	    	var.datacenter
	    		Enter a value: {datacenter}
	    	var.vswitch_id
	    		Enter a value: {vswitch_id}
	    	....

* Apply phase

		terraform apply 
		    var.availability_zones
  				Enter a value: {var.availability_zones}  /*cn-beijing-b*/
	    	var.datacenter
	    		Enter a value: {datacenter}
	    	var.vswitch_id
	    		Enter a value: {vswitch_id}
	    	....

* Destroy 

		terraform destroy