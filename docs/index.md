---
page_title: "m3 Provider"
subcategory: ""
description: |-
The m3 provider gives ability to manage resources of different cloud providers in a unified way via Maestro.
---

# m3 Provider

The m3 provider gives ability to manage resources of different cloud providers in a unified way via Maestro.

# Supported clouds

AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK, VSPHERE, VMWARE, YANDEX


## Example Usage

```terraform
provider "m3" {
	url = "http://ip:port/maestro/api/V3"
	access_key = "access_key"
	secret_key = "secret_key"
	user_identifier = "user_identifier"
}

provider "m3" {
	url = "http://ip:port/maestro/api/V3"
	access_key = "access_key"
	secret_key = "secret_key"
	user_identifier = "user_identifier"
	region = "COMPANY-OPENSTACK-3"
	tenant = "EPMC-EOOS"
	cloud = "cloud"
}
```
