### kubernetes cluster-autoscaler Example

The example help you to configure and deploy cluster-autoscaler to specific kubernetes cluster.

### Preconditions  
* Create a kuberentes cluster firstly (managed or dedicated).
* Create scaling group with proper configuration before.
    * CentOS7 base image 
    * Same size of instancesTypes. 
    * Add Policy to RAM role of the node to deploy cluster-autoscaler (https://www.alibabacloud.com/help/doc-detail/119099.htm)
    * Configure RAM role of autoscaling ndoes if necessary. 
* Change the main.tf based on examples/cluster-autoscaler.

### Get up and running

* Planning phase

		terraform plan 
		    var.cluster_id
          		Enter a value: {var.cluster_id} 
    		var.utilization
  				Enter a value: {var.utilization} 
	    	var.cool_down_duration
	    		Enter a value: {var.cool_down_duration} 
			var.defer_scale_in_duration
  				Enter a value: {var.defer_scale_in_duration} 	
	    

* Apply phase

		terraform apply     	
		    var.cluster_id
          		Enter a value: {var.cluster_id} 
    		var.utilization
  				Enter a value: {var.utilization} 
	    	var.cool_down_duration
	    		Enter a value: {var.cool_down_duration} 
			var.defer_scale_in_duration
  				Enter a value: {var.defer_scale_in_duration} 
* Destroy 

		terraform destroy