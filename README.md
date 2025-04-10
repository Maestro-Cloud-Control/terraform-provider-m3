Terraform Provider for M3 (pre-alpha)
==================
- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.23 (to build the provider plugin)


Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](https://github.com/terraform-providers/terraform-provider-aws#requirements) before proceeding).

Clone repository to: `$PATH_TO_YOUR_DIRECTORY/`

To test plugin run:
```buildoutcfg
go test ./service -v
```

To build plugin run:
```buildoutcfg
#linux
go build -o terraform-provider-m3_v0.6.7
#windows
go build -o terraform-provider-m3_v0.6.7.exe
```

After that move it to user plugins directory (see https://www.terraform.io/docs/plugins/basics.html)
```buildoutcfg
#linux
mv terraform-provider-m3_v0.6.7 ~/.terraform.d/plugins
#windows
move terraform-provider-m3_v0.6.7.exe %APPDATA%\terraform.d\plugins
```

Then you can work with it like with any other provider.

[Here](https://www.terraform.io/docs/internals/debugging.html) you can find information about enabling logs for terraform

# M3 Provider 

Terraform is a framework for configuration, but one Terraform template cannot be written in a unified way so that it could be used for different cloud providers (the template for AWS will not work if you deploy it in Azure). Maestro3 provides its own custom tool -  Terraform-Provider that gives the users possibility to work with the provided API so that they can deploy the infrastructure with a single template on any cloud provider supported by Maestro3, including private OpenStack region.

Use M3 provider to interact with the Maestro3 resources.

## Supported interaction
|                      | AWS | Azure | Google | Yandex | OpenStack | vCloud | vSphere | Selectel |
|----------------------|:---:|:-----:|:------:|:------:|:---------:|:------:|:-------:|:--------:|
| Instance management  |  +  |    +  |    +   |   +*   |    +      |   +    |    +    |    +     |
| Volume management    |  +* |   +*  |    +*  |        |    +      |   +    |   -**   |   -***   |
| Image management     |  +* |   +*  |    +*  |        |    +      |  -**   |   -**   |    +     |
| Script management    |  +* |   +*  |    +*  |   +*   |    +*     |   +*   |   +*    |    +*    |
| Schedule management  |  +* |   +*  |    +*  |   +*   |    +*     |   +*   |   +*    |    +*    |
| Keypair management   |  +* |   +*  |    +*  |   +*   |    +*     |   -    |    -    |    -     |
| Instance update      |     |       |        |        |           |        |         |          |
| Volume update        |     |       |        |        |           |        |         |          |
| Image update         |     |       |        |        |           |        |         |          |
| Script update        |  +* |   +*  |    +*  |   +*   |    +*     |   +*   |   +*    |    +*    |
| Schedule update      |  +* |   +*  |    +*  |   +*   |    +*     |   +*   |   +*    |    +*    |
| Keypair update       |  +* |   +*  |    +*  |   +*   |    +*     |   -    |    -    |    -     |

- \* - need to be tested
- \** - not implemented on backend
- \*** - used old API

## Provider
| Name               | Version |
|--------------------|---------|
| [m3](#provider_m3) | 0.6.7   |
