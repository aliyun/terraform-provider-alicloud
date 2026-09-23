### Create EDAS application Example

The example creates an application in EDAS. The variables.tf can let you create specify parameter instances, such as application_name, cluster_id, package_type etc.

### Get up and running

* Planning phase

		terraform plan 
    		var.application_name
  				Enter a value: {var.application_name}  
	    	var.cluster_id
	    		Enter a value: {var.cluster_id} 
	    	var.package_type
	    		Enter a value: {var.package_type} /* war or jar*/
	    	....

* Apply phase

		terraform apply 
		    var.application_name
                Enter a value: {var.application_name}  
            var.cluster_id
                Enter a value: {var.cluster_id} 
            var.package_type
                Enter a value: {var.package_type} /* war or jar*/
            ....

* Destroy 

		terraform destroy