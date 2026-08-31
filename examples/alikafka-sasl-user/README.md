### ALIKAFKA sasl user Example

The example launches ALIKAFKA sasl user. The parameter in variables.tf can let you specify the sasl user.

### Get up and running

* Planning phase

		terraform plan 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.username
  				Enter a value: {var.username} 
  			var.password
                Enter a value: {var.password}
	    

* Apply phase

		terraform apply 
	    	var.instance_id
	    		Enter a value: {var.instance_id}
			var.username
  				Enter a value: {var.username} 
            var.password
                Enter a value: {var.password}

* Destroy 

		terraform destroy