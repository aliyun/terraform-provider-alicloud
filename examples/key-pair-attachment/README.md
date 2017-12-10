### ECS Example

The example launches a keypair and a keypair attachment, and use them can attach one keypair to launched ECS instances. Besides, the example also launches other resources, such as vpc, vswitch, disk and so on.
The count parameter in variables.tf can let you create specify number ECS instances.

### Get up and running

* Planning phase

		terraform plan

* Apply phase

		terraform apply

* Destroy

		terraform destroy