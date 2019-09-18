### ALIKAFKA instance Example

The example launches ALIKAFKA group. The parameter in variables.tf can let you specify the group.

### Get up and running

* Planning phase

		terraform plan 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.group
  				Enter a value: {var.consumer_id} 
	    

* Apply phase

		terraform apply 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.group
  				Enter a value: {var.consumer_id} 

* Destroy 

		terraform destroy