
resource "m3_volume" "my-volume" {
  name = "name"
  tenant = "EPMC-EOOS"
  region = "COMPANY-OPENSTACK-3"
  size_in_gb = "8"
  instance_id = "instance id for attaching"
}

resource "m3_volume" "my-volume" {
  name = "name"
  size_in_gb = "8"
  instance_id = "instance id for attaching"
}