### Scale EDAS application Example

The example scale an application in EDAS. The variables.tf can let you create specify parameter instances, such as app_id, deploy_group, ecu_info etc.

### Get up and running

* Planning phase

		terraform plan 
    		var.app_id
  				Enter a value: {var.app_id}  
	    	var.deploy_group
	    		Enter a value: {var.deploy_group} 
	    	var.ecu_info
	    		Enter a value: {var.ecu_info} 
	    	....

* Apply phase

		terraform apply 
		    var.app_id
                Enter a value: {var.app_id}  
            var.deploy_group
                Enter a value: {var.deploy_group} 
            var.ecu_info
                Enter a value: {var.ecu_info} 
            ....

* Destroy 

		terraform destroy