resource "m3_volume" "my-volume" {
  volume_name = "name"
  tenant_name = "EPMC-EOOS"
  region_name = "COMPANY-OPENSTACK-3"
  size_in_gb = "8"
  instance_id = "instance id for attaching"
}

//////////////////////////////////////////////////////////////////

resource "m3_volume" "my-volume" {
  volume_name = "name"
  size_in_gb = "8"
  instance_id = "instance id for attaching"
}