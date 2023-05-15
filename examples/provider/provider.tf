
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