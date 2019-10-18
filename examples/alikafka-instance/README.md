### ALIKAFKA topic Example

The example launches ALIKAFKA Instance. The parameter in variables.tf can let you specify the instance.

### Get up and running

* Planning phase

		terraform plan 
		    var.name
          		Enter a value: {var.name} 
    		var.topic_quota
  				Enter a value: {var.topic_quota} 
	    	var.disk_type
	    		Enter a value: {var.disk_type} 
			var.disk_size
  				Enter a value: {var.disk_size} 
			var.deploy_type
  				Enter a value: {var.deploy_type}
	    	var.io_max
	    		Enter a value: {var.io_max}
	    	var.eip_max
                Enter a value: {var.eip_max}
            var.vpc_id
                Enter a value: {var.vpc_id}
            var.vswitch_id
                Enter a value: {var.vswitch_id}
            var.zone_id
                Enter a value: {var.zone_id}
	    

* Apply phase

		terraform apply 
		    var.name
          		Enter a value: {var.name} 
    		var.topic_quota
  				Enter a value: {var.topic_quota} 
	    	var.disk_type
	    		Enter a value: {var.disk_type} 
			var.disk_size
  				Enter a value: {var.disk_size} 
			var.deploy_type
  				Enter a value: {var.deploy_type}
	    	var.io_max
	    		Enter a value: {var.io_max}
	    	var.eip_max
                Enter a value: {var.eip_max}
            var.vpc_id
                Enter a value: {var.vpc_id}
            var.vswitch_id
                Enter a value: {var.vswitch_id}
            var.zone_id
                Enter a value: {var.zone_id}
	    	    	

* Destroy 

		terraform destroy