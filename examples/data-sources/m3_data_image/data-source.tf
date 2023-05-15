
data "m3_data_image" "dim" {
  owner = "some@gmail.com"
  name = "LINUX (can be regexp)"
  os_type = "l or w"
  alias = "CentOS7_64-bit"
}

data "m3_data_image" "dim" {
  owner = "some@gmail.com"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  name = "LINUX (can be regexp)"
  os_type = "l or w"
  alias = "CentOS7_64-bit"
}

data "m3_data_image" "dim" {
  only_system_images = true
  name = "LINUX (can be regexp)"
  os_type = "l or w"
  alias = "CentOS7_64-bit"
}

data "m3_data_image" "dim" {
  only_system_images = true
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  name = "LINUX (can be regexp)"
  os_type = "l or w"
  alias = "CentOS7_64-bit"
}
