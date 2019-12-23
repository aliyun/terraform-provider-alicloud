### ALIKAFKA consumer group Example

The example launches ALIKAFKA consumer group. The parameter in variables.tf can let you specify the consumer group.

### Get up and running

* Planning phase

		terraform plan 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.consumer_id
  				Enter a value: {var.consumer_id} 
	    

* Apply phase

		terraform apply 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.consumer_id
  				Enter a value: {var.consumer_id} 

* Destroy 

		terraform destroy