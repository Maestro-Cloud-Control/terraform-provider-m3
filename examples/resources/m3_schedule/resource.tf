
resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  action = "start or stop"
  cloud = "cloud"
  cron = "0 1 0 * * ? *"
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  action = "start or stop"
  cloud = "cloud"
  cron = "0 1 0 * * ? *"
  tag_key = "key of tag for select instances"
  tag_value = "value of tag for select instances"
  # also, we can create schedule without tag_value
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  action = "start or stop"
  cloud = "cloud"
  cron = "0 1 0 * * ? *"
  tag_key = "key of tag for select instances"
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  region = "COMPANY-OPENSTACK-3"
  tenant = "EPMC-EOOS"
  action = "start or stop"
  cloud = "cloud"
  cron = "0 1 0 * * ? *"
  instances = ["instance ID", "instance ID"]
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  action = "start or stop"
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  action = "start or stop"
  cron = "0 1 0 * * ? *"
  tag_key = "key of tag for select instances"
  tag_value = "value of tag for select instances"
  # also, we can create schedule without tag_value
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  action = "start or stop"
  cron = "0 1 0 * * ? *"
  tag_key = "key of tag for select instances"
}

resource "m3_schedule" "my-schedule" {
  name = "name"
  description = "description"
  action = "start or stop"
  cron = "0 1 0 * * ? *"
  instances = ["instance ID", "instance ID"]
}