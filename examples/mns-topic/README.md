### MNS topic Example

The example launches MNS topic. the  parameter in variables.tf can let you specify the topic.

### Get up and running

* Planning phase

		terraform plan 
    		var.name
  				Enter a value: {var.name}  /*tf-example-topic*/
	    	var.maximum_message_size
	    		Enter a value: {var.maximum_message_size} /*65536*/
	    	....

* Apply phase

		terraform apply 
		    var.name
  				Enter a value: {var.name}  /*tf-example-topic*/
	    	var.maximum_message_size
	    		Enter a value: {var.maximum_message_size} /*65536*/
	    	....

* Destroy 

		terraform destroy