### ONS instance Example

The example launches ONS Instance. The parameter in variables.tf can let you specify the instance.

### Get up and running

* Planning phase

		terraform plan 
    		var.name
  				Enter a value: {var.name} 
	    	var.remark
	    		Enter a value: {var.remark} /*tf-example-ons-instance-remark*/
	    

* Apply phase

		terraform apply 
		    var.name
  		        Enter a value: {var.name}
	        var.remark
	    	    Enter a value: {var.remark} /*tf-example-ons-instance-remark*/
	    	

* Destroy 

		terraform destroy