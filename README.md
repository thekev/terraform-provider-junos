terraform-provider-junos
========================
![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-junos)
[![Go Status](https://github.com/jeremmfr/terraform-provider-junos/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-junos/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-junos/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-junos)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-junos)
[![Website](https://img.shields.io/badge/doc-website-lightgrey)](https://terraform-provider-junos.jeremm.fr/)
[![Terraform Registry](https://img.shields.io/badge/doc-terraform_registry-lightgrey)](https://registry.terraform.io/providers/jeremmfr/junos)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/jeremmfr/terraform-provider-junos/blob/master/LICENSE)
<br/><br/>
This is an **unofficial** terraform provider for Junos devices with netconf protocol

See [website](https://terraform-provider-junos.jeremm.fr/) or
[terraform registry](https://registry.terraform.io/providers/jeremmfr/junos)
for provider and resources documentation.

Requirements
---
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x or 0.13.x

Optional
---
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Install automatic
---
With terraform >= 0.13, add source information inside the terraform configuration block for automatic provider installation :
```hcl
terraform {
  required_providers {
    junos = {
      source = "jeremmfr/junos"
    }
  }
}
```


Install binary on disk
---
Download latest version in [releases](https://github.com/jeremmfr/terraform-provider-junos/releases)
##### terraform 0.13
```bash
for archive in $(ls terraform-provider-junos*.zip) ; do
  OS_ARCH=$(echo $archive | cut -d'_' -f3-4 | cut -d'.' -f1)
  VERSION=$(echo $archive | cut -d'_' -f2)
  tfPath="${HOME}/.terraform.d/plugins/registry.local/jeremmfr/junos/${VERSION}/${OS_ARCH}/"
  mkdir -p ${tfPath}
  unzip ${archive} -d ${tfPath}
done
```
and add this inside the terraform configuration block :
```hcl
terraform {
  required_providers {
    junos = {
      source = "registry.local/jeremmfr/junos"
    }
  }
}
```
##### terraform 0.12
```bash
tfPath=$(which terraform | rev | cut -d'/' -f2- | rev)
unzip terraform-provider-junos*.zip -d ${tfPath}
```

Building binary provider with latest tag (terraform 0.13)
---
```bash
git clone https://github.com/jeremmfr/terraform-provider-junos.git
cd terraform-provider-junos && git fetch --tags
latestTag=$(git describe --tags `git rev-list --tags --max-count=1`)
git checkout ${latestTag}
tfPath="${HOME}/.terraform.d/plugins/registry.local/jeremmfr/junos/${latestTag:1}/$(go env GOOS)_$(go env GOARCH)/"
mkdir -p ${tfPath}
go build -o ${tfPath}/terraform-provider-junos_${latestTag}
unset latestTag tfPath
```
and add this inside the terraform configuration block :
```hcl
terraform {
  required_providers {
    junos = {
      source = "registry.local/jeremmfr/junos"
    }
  }
}
```

Building binary provider with latest tag (terraform 0.12)
---
```bash
git clone https://github.com/jeremmfr/terraform-provider-junos.git
cd terraform-provider-junos && git fetch --tags
latestTag=$(git describe --tags `git rev-list --tags --max-count=1`)
git checkout ${latestTag}
tfPath=$(which terraform | rev | cut -d'/' -f2- | rev)
go build -o ${tfPath}/terraform-provider-junos_${latestTag}
unset latestTag tfPath
```

Details
---
Some Junos parameters are not included in provider for various reasons (time, utility, understanding, ...)
