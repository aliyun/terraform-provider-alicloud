### ALIKAFKA topic Example

The example launches ALIKAFKA Instance. The parameter in variables.tf can let you specify the instance.

### Get up and running

* Planning phase

		terraform plan 
		    var.name
          		Enter a value: {var.name} 
    		var.partition_num
  				Enter a value: {var.partition_num} 
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
            var.paid_type
                Enter a value: {var.paid_type}
            var.spec_type
                Enter a value: {var.spec_type}
            var.vswitch_id
                Enter a value: {var.vswitch_id}
	    

* Apply phase

		terraform apply 
		    var.name
          		Enter a value: {var.name} 
    		var.partition_num
  				Enter a value: {var.partition_num} 
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
            var.paid_type
                Enter a value: {var.paid_type}
            var.spec_type
                Enter a value: {var.spec_type}
            var.vswitch_id
                Enter a value: {var.vswitch_id}
	    	    	

* Destroy 

		terraform destroy