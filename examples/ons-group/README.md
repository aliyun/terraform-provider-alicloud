### ONS instance Example

The example launches ONS group. The parameter in variables.tf can let you specify the group.

### Get up and running

* Planning phase

		terraform plan 
    		var.name
  				Enter a value: {var.name} 
	    	var.instance_remark
	    		Enter a value: {var.instance_remark} /*tf-example-ons-instance-remark*/
			var.group
  				Enter a value: {var.group} 
	    	var.group_remark
	    		Enter a value: {var.group_remark} /*tf-example-ons-group-remark*/
	    

* Apply phase

		terraform apply 
		    var.name
  		        Enter a value: {var.name}
	        var.instance_remark
	    	    Enter a value: {var.remark} /*tf-example-ons-instance-remark*/
			var.group
  				Enter a value: {var.group} 
	    	var.group_remark
	    		Enter a value: {var.group_remark} /*tf-example-ons-group-remark*/
	    	

* Destroy 

		terraform destroy