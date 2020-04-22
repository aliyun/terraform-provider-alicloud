### Create EDAS cluster Example

The example creates a cluster in EDAS. The variables.tf can let you create specify parameter instances, such as cluster_name, cluster_type, network_mode etc.

### Get up and running

* Planning phase

		terraform plan 
    		var.cluster_name
  				Enter a value: {var.cluster_name}  
	    	var.cluster_type
	    		Enter a value: {var.cluster_type} /*2 for ECS cluster*/
	    	var.network_mode
	    		Enter a value: {var.vswitch_id}
	    	....

* Apply phase

		terraform apply 
		    var.cluster_name
                Enter a value: {var.cluster_name}  
            var.cluster_type
                Enter a value: {var.cluster_type} /*2 for ECS cluster*/
            var.network_mode
                Enter a value: {var.vswitch_id}
            ....

* Destroy 

		terraform destroy