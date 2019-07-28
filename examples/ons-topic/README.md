### ONS instance Example

The example launches ONS Topic. The parameter in variables.tf can let you specify the topic.

### Get up and running

* Planning phase

		terraform plan 
    		var.name
  				Enter a value: {var.name} 
	    	var.instance_remark
	    		Enter a value: {var.instance_remark} /*tf-example-ons-instance-remark*/
			var.topic
  				Enter a value: {var.topic} 
			var.message_type
  				Enter a value: {var.message_type}
	    	var.topic_remark
	    		Enter a value: {var.topic_remark} /*tf-example-ons-topic-remark*/
	    

* Apply phase

		terraform apply 
		    var.name
  		        Enter a value: {var.name}
	        var.instance_remark
	    	    Enter a value: {var.remark} /*tf-example-ons-instance-remark*/
			var.topic
  				Enter a value: {var.topic} 
			var.message_type
  				Enter a value: {var.message_type}
	    	var.topic_remark
	    		Enter a value: {var.topic_remark} /*tf-example-ons-topic-remark*/
	    	

* Destroy 

		terraform destroy