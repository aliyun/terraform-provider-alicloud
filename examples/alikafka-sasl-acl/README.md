### ALIKAFKA sasl user Example

The example launches ALIKAFKA sasl user. The parameter in variables.tf can let you specify the sasl user.

### Get up and running

* Planning phase

		terraform plan 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.username
  				Enter a value: {var.username} 
  			var.acl_resource_type
                Enter a value: {var.acl_resource_type}
            var.acl_resource_name
                Enter a value: {var.acl_resource_name}
            var.acl_resource_pattern_type
                Enter a value: {var.acl_resource_pattern_type}
            var.acl_operation_type
                Enter a value: {var.acl_operation_type}
	    

* Apply phase

		terraform apply 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.username
  				Enter a value: {var.username} 
  			var.acl_resource_type
                Enter a value: {var.acl_resource_type}
            var.acl_resource_name
                Enter a value: {var.acl_resource_name}
            var.acl_resource_pattern_type
                Enter a value: {var.acl_resource_pattern_type}
            var.acl_operation_type
                Enter a value: {var.acl_operation_type}
* Destroy 

		terraform destroy