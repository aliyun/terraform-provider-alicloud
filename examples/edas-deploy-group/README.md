### Create EDAS Application Deploy Group Example

The example create an application deploy group in EDAS. The variables.tf can let you create specify parameter instances, such as app_id, group_name.

### Get up and running

* Planning phase

		terraform plan 
    		var.app_id
  				Enter a value: {var.app_id}  
	    	var.group_name
	    		Enter a value: {var.group_name} 

* Apply phase

		terraform apply 
		    var.app_id
                Enter a value: {var.app_id}  
            var.group_name
                Enter a value: {var.group_name}

* Destroy 

		terraform destroy