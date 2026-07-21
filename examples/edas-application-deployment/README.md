### Deploy EDAS application Example

The example deploy an application in EDAS. The variables.tf can let you create specify parameter instances, such as app_id, group_id, package_version etc.

### Get up and running

* Planning phase

		terraform plan 
    		var.app_id
  				Enter a value: {var.app_id}  
	    	var.group_id
	    		Enter a value: {var.group_id} 
	    	var.package_version
	    		Enter a value: {var.package_version} 
	    	....

* Apply phase

		terraform apply 
		    var.app_id
                Enter a value: {var.app_id}  
            var.group_id
                Enter a value: {var.group_id} 
            var.package_version
                Enter a value: {var.package_version} 
            ....

* Destroy 

		terraform destroy