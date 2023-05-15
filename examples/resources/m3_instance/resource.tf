
resource "m3_instance" "my-server" {
  image = "CentOS7_64-bit"
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
}

resource "m3_instance" "my-server" {
  image = "CentOS7_64-bit"
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  stop can be less then terminate
  #  max 720
  stop_after = 5
  #  max 720
  terminate_after = 6
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  stop can be less then terminate
  #  max 720
  stop_after = 5
  #  max 720
  terminate_after = 6
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  terminate_after = 6
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  terminate_after = 6
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  stop_after = 6
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  stop_after = 6
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = "CentOS7_64-bit"
  name  = "test"
  shape = "MINI"
  key = "sshkey"
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  stop can be less then terminate
  #  max 720
  stop_after = 5
  #  max 720
  terminate_after = 6
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  stop can be less then terminate
  #  max 720
  stop_after = 5
  #  max 720
  terminate_after = 6
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  terminate_after = 6
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  stop_after = 5
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  stop_after = 5
  lock_termination = true
}

resource "m3_instance" "my-server" {
  image = data.m3_data_image.dim.id
  name  = "test"
  shape = "MINI"
  key = "sshkey"
  #  max 720
  terminate_after = 6
  lock_termination = true
}