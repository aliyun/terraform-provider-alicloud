### Cluster expansion Example

The example expand a cluster in EDAS. The variables.tf can let you create specify parameter instances, such as cluster_id and instance_ids.

### Get up and running

* Planning phase

		terraform plan 
    		var.cluster_id
  				Enter a value: {var.cluster_id}  
	    	var.instance_ids
	    		Enter a value: {var.instance_ids} /*[instanceid1, instanceid2]*/
	  

* Apply phase

		terraform apply 
		    var.cluster_id
                Enter a value: {var.cluster_id}  
            var.instance_ids
                Enter a value: {var.instance_ids} /*[instanceid1, instanceid2]*/

* Destroy 

		terraform destroy