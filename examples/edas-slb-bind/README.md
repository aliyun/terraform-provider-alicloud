### Bind Slb to EDAS application Example

The example binds a slb to application in EDAS. The variables.tf can let you create specify parameter instances, such as app_id, slb_id, slb_ip etc.

### Get up and running

* Planning phase

		terraform plan 
    		var.app_id
  				Enter a value: {var.app_id}  
	    	var.slb_id
	    		Enter a value: {var.slb_id} 
	    	var.slb_ip
	    		Enter a value: {var.slb_ip} 
	    	....

* Apply phase

		terraform apply 
		    var.app_id
                Enter a value: {var.app_id}  
            var.slb_id
                Enter a value: {var.slb_id} 
            var.slb_ip
                Enter a value: {var.slb_ip} 
            ....

* Destroy 

		terraform destroy