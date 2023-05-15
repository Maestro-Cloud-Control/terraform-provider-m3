
resource "m3_image" "my-image" {
  name = "ImageFromTf"
  source_instance_id = "ecs00100019F"
  description = "Here is image description"
}

resource "m3_image" "my-image" {
  tenant = "EPMC-EOOS"
  region = "COMPANY-OPENSTACK-3"
  name = "ImageFromTf"
  source_instance_id = "ecs00100019F"
  description = "Here is image description"
}
