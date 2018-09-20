### MNS Queue Example

The example launches MNS Queue. the  parameter in variables.tf can let you specify the queue.

### Get up and running

* Planning phase

		terraform plan 
    		var.name
  				Enter a value: {var.name}  /*tf-example-queue*/
	    	var.delay_seconds
	    		Enter a value: {var.delay_seconds} /*0*/
	    	....

* Apply phase

		terraform apply 
		    var.name
  				Enter a value: {var.name}  /*tf-example-queue*/
	    	var.delay_seconds
	    		Enter a value: {var.delay_seconds} /*0*/
	    	....

* Destroy 

		terraform destroy