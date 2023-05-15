
resource "m3_script" "my_script1" {
  name = "script"
  content = <<EOF
	#! /bin/bash
	sudo apt-get update
sudo apt-get install -y apache2
sudo systemctl start apache2
sudo systemctl enable apache2
echo "<h1>Deployed via Terraform</h1>" | sudo tee /var/www/html/index.html
EOF
  extension = "for example: .cmd"
  tenant = "tenant"
  region = "region"
  cloud = "cloud"
}

resource "m3_script" "my_script1" {
  name = "script"
  content = "some script"
  extension = "for example: .cmd"
  tenant = "tenant"
  region = "region"
  cloud = "cloud"
}

resource "m3_script" "my_script1" {
  name = "script"
  content = <<EOF
	#! /bin/bash
	sudo apt-get update
sudo apt-get install -y apache2
sudo systemctl start apache2
sudo systemctl enable apache2
echo "<h1>Deployed via Terraform</h1>" | sudo tee /var/www/html/index.html
EOF
  extension = "for example: .cmd"
}

resource "m3_script" "my_script1" {
  name = "script"
  content = "some script"
  extension = "for example: .cmd"
}