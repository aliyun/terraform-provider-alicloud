### ALIKAFKA topic Example

The example launches ALIKAFKA Topic. The parameter in variables.tf can let you specify the topic.

### Get up and running

* Planning phase

		terraform plan 
		    var.instance_id
          		Enter a value: {var.instance_id} 
    		var.topic
  				Enter a value: {var.topic} 
	    	var.local_topic
	    		Enter a value: {var.local_topic} 
			var.compact_topic
  				Enter a value: {var.compact_topic} 
			var.partition_num
  				Enter a value: {var.partition_num}
	    	var.remark
	    		Enter a value: {var.remark} /*tf-example-alikafka-topic-remark*/
	    

* Apply phase

		terraform apply 
            var.instance_id
                Enter a value: {var.instance_id} 
    		var.topic
  				Enter a value: {var.topic} 
	    	var.local_topic
	    		Enter a value: {var.local_topic} 
			var.compact_topic
  				Enter a value: {var.compact_topic} 
			var.partition_num
  				Enter a value: {var.partition_num}
	    	var.remark
	    		Enter a value: {var.remark} /*tf-example-alikafka-topic-remark*/
	    	

* Destroy 

		terraform destroy