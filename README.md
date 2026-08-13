# Byteplus Cloud Control Provider

The Byteplus Cloud Control Provider enables interaction with various Byteplus-supported resources through the Cloud Control API. Prior to usage, you must configure the provider with appropriate credentials.

Use the left navigation panel to explore available resource documentation. If you cannot find the desired resource, please submit an issue report for assistance.

> **NOTE**：
> 1. The Byteplus Cloud Control provider requires the use of Terraform 1.0.7 or later.
> 2. This guide requires an available Byteplus account or sub-account to create resources.

## Example Usage

Terraform 1.0.7 and later:

```shell
# Configure the Bytepluscc Provider
terraform {
  required_providers {
    bytepluscc = {
      source  = "byteplus-sdk/bytepluscc"
      version = "~> 0.0.1"
    }
  }
}

# Create VPC
resource "bytepluscc_vpc_vpc" "VPCDemo" {
  cidr_block            = "192.168.0.0/24"
  support_ipv_4_gateway = true
  vpc_name              = "vpc-demo"
  description           = "VpcDemo Example"
  dns_servers = [
    "12.3.x.x"
  ]
  associate_cens = [
    {
      cen_id       = "cen-3re8cx4vwdibk5zsk2xxxx"
      cen_owner_id = "200000xxx"
      cen_status   = "Attaching"
    }
  ]
  nat_gateway_ids = [
    "ngw-2d6tp1y8zq41s58ozfdxxx"
  ]
  route_table_ids = [
    "vtb-29mkf8ft83l6o1e1hgixxx"
  ]
  security_group_ids = [
    "sg-29mkx39jb80741e1hgj2bxxx"
  ]
  project_name = "default"
  tags = [
    {
      key   = "env"
      value = "test"
    }
  ]
}
# Query VPC
data "bytepluscc_vpc_vpc" "VpcVpcDataSource" {
  id = bytepluscc_vpc_vpc.VPCDemo.id
}
```

## Authentication

The Bytepluscc provider first resolves source credentials and then optionally uses those credentials to call STS AssumeRole.

### Static credentials

Static credentials can be provided by adding an public_key and private_key in-line in the bytepluscc provider block:

> **Warning**:
> Hard-coded credentials are not recommended in any Terraform configuration and risks secret leakage should this file ever be committed to a public version control system.

**Authentication Resolution and Requirements:**

- **AK/SK Authentication**: Complete `access_key` and `secret_key` credentials have the highest source priority; `session_token` is forwarded when temporary credentials are used.
- **Profile Authentication**: `profile` or `file_path` selects BytePlus CLI credentials when a complete AK/SK pair is absent.
- **Default Credential Provider**: When neither complete AK/SK nor Profile inputs are configured, the BytePlus SDK default credential chain is used.
- **AssumeRole**: When `assume_role` is configured, it wraps the selected AK/SK, Profile, or default-chain source and supplies refreshable target-role credentials.
- **Validation Rules**:
    - If complete AK/SK and Profile inputs are both present, AK/SK takes precedence and the provider returns a warning
    - A partial AK/SK pair is not a usable source and does not override Profile inputs
    - The `file_path` parameter is optional. If not specified, the default file `~/.byteplus/config.json` is used

Usage:

```terraform
provider "bytepluscc" {
  region    = "ap-southeast-1"
  profile   = "platform-admin"
  file_path = "/path/to/.byteplus/config.json"

  assume_role {
    assume_role_trn  = "trn:iam::222222222222:role/terraform-execution"
    duration_seconds = 3600
  }
}
```

### Environment variables

Set `BYTEPLUS_REGION` and choose AK/SK, Profile, or the default credential chain. If both complete AK/SK and Profile inputs are configured, AK/SK takes precedence and the provider returns a warning:
```shell
provider "bytepluscc" {

}
```

Usage:

```shell
$ export BYTEPLUS_REGION="ap-southeast-1"

# Option 1: AK/SK credentials
$ export BYTEPLUS_ACCESS_KEY="your_public_key"
$ export BYTEPLUS_SECRET_KEY="your_private_key"
$ export BYTEPLUS_SESSION_TOKEN="your_session_token" # optional for temporary credentials

# Option 2: Profile credentials (do not set AK/SK at the same time)
# export BYTEPLUS_PROFILE="your_profile"
# export BYTEPLUS_FILE_PATH="your_file_path" # defaults to ~/.byteplus/config.json
```

## Authenticated Cloud Control proxy

Use `proxy_url` together with `proxy_authorization` when Cloud Control API requests must pass through an authenticated HTTP proxy. `proxy_authorization` is the complete value of the `Proxy-Authorization` header and is marked sensitive. For HTTP proxies, keep credentials out of `proxy_url` and use this attribute instead. `proxy_authorization` supports HTTP and HTTPS proxy URLs; SOCKS5 and SOCKS5H URLs remain available without this attribute.

For a ZTI proxy, configure the same values used by curl:

```shell
PSM="<your-psm>"
PROXY="<your-proxy-domain>"
AUTH_TOKEN=$(echo -n "ZTI_$PSM:$(cat "$SEC_TOKEN_PATH")" | base64 -w0)

export BYTEPLUS_PROXY_URL="http://${PROXY}:8080"
export BYTEPLUS_PROXY_AUTHORIZATION="Basic ${AUTH_TOKEN}"
```

The equivalent Terraform configuration is:

```hcl
variable "proxy_authorization" {
  type      = string
  sensitive = true
}

provider "bytepluscc" {
  proxy_url           = "http://<your-proxy-domain>:8080"
  proxy_authorization = var.proxy_authorization
}
```

Standard `no_proxy` rules identify destinations that connect directly instead of using `proxy_url`:

```hcl
provider "bytepluscc" {
  proxy_url           = "http://<your-proxy-domain>:8080"
  proxy_authorization = var.proxy_authorization
  no_proxy            = "localhost,127.0.0.1,.internal.example.com"
}
```

To use the proxy only for specified destinations, configure the inverse allowlist with `proxy_include_domains`:

```hcl
provider "bytepluscc" {
  proxy_url           = "http://<your-proxy-domain>:8080"
  proxy_authorization = var.proxy_authorization
  proxy_include_domains = [
    "cloudcontrol.ap-southeast-1.byteplusapi.com",
  ]
}
```

Replace the example host with the Cloud Control hostname for the configured region, or use the host from `endpoints.cloudcontrolapi` when a custom endpoint is configured. A domain without a leading dot matches the domain and its subdomains; a leading dot matches subdomains only. IP addresses, CIDR ranges, optional ports, and `*` follow standard `NO_PROXY` matching. `no_proxy` and `proxy_include_domains` cannot be configured together.

The equivalent environment variables are `BYTEPLUS_NO_PROXY` (with `NO_PROXY` and `no_proxy` as fallbacks) and comma-separated `BYTEPLUS_PROXY_INCLUDE_DOMAINS`. Explicit provider attributes take precedence over the environment variables. These settings apply to Cloud Control API traffic; STS calls made by `assume_role` use a separate client.

<!-- schema generated by tfplugindocs -->

## Schema

### Optional

- `access_key` (String) The Access Key for Byteplus Provider. It can also be sourced from the `BYTEPLUS_ACCESS_KEY` environment variable
- `secret_key` (String) The Secret Key for Byteplus Provider. It can also be sourced from the `BYTEPLUS_SECRET_KEY` environment variable
- `session_token` (String) The Session Token for Byteplus Provider, used together with temporary AK/SK obtained from STS. It can also be sourced from the `BYTEPLUS_SESSION_TOKEN` environment variable
- `profile` (String) The Profile for Byteplus Provider. It can be sourced from the `BYTEPLUS_PROFILE` environment variable. Complete AccessKey and SecretKey credentials take precedence when both sources are configured
- `file_path` (String) The File Path for Byteplus Provider. It specifies the path to the profile configuration file. If not specified, the default file `~/.byteplus/config.json` will be used, and can also be sourced from the `BYTEPLUS_FILE_PATH` environment variable
- `assume_role` (Attributes) An `assume_role` block that uses the selected source credentials to obtain target-role credentials. Only one `assume_role` block may be in the configuration. (see [below for nested schema](#nestedatt--assume_role))
- `customer_headers` (String) CUSTOMER HEADERS for Byteplus Provider. The customer_headers field uses commas (,) to separate multiple headers, and colons (:) to separate each header key from its corresponding value.
- `disable_ssl` (Boolean) Disable SSL for Byteplus Provider
- `endpoints` (Attributes) An `endpoints` block (documented below). Only one `endpoints` block may be in the configuration. (see [below for nested schema](#nestedatt--endpoints))
- `no_proxy` (String) Comma-separated hosts, domain suffixes, IP addresses, or CIDR ranges that bypass proxy_url. It follows standard NO_PROXY matching and can be sourced from BYTEPLUS_NO_PROXY, NO_PROXY, or no_proxy.
- `proxy_authorization` (String, Sensitive) Value of the Proxy-Authorization header for Cloud Control API proxy requests, for example `Basic <token>`. It can also be sourced from the `BYTEPLUS_PROXY_AUTHORIZATION` environment variable.
- `proxy_include_domains` (Set of String) Hosts, domain suffixes, IP addresses, or CIDR ranges that use proxy_url while all other destinations connect directly. It can be sourced as a comma-separated list from BYTEPLUS_PROXY_INCLUDE_DOMAINS and cannot be combined with no_proxy.
- `proxy_url` (String) HTTP, HTTPS, SOCKS5, or SOCKS5H proxy URL for Cloud Control API requests. It can also be sourced from the `BYTEPLUS_PROXY_URL` environment variable.
- `region` (String) The Region for Byteplus Provider. It must be provided, but it can also be sourced from the `BYTEPLUS_REGION` environment variable


<a id="nestedatt--assume_role"></a>

### Nested Schema for `assume_role`

Required:

- `assume_role_trn` (String) he TRN of the role to assume.

Optional:

- `duration_seconds` (Number) The duration of the session when making the AssumeRole call. Its value ranges from 900 to 43200(seconds), and default is 3600 seconds.
- `policy` (String) A more restrictive policy when making the AssumeRole call

<a id="nestedatt--endpoints"></a>

### Nested Schema for `endpoints`

Optional:

- `cloudcontrolapi` (String) Use this to override the default Cloud Control API service endpoint URL
- `sts` (String) Use this to override the default STS service endpoint URL

## Security and privacy

This project takes security seriously.
For vulnerability reporting and supported versions, see [SECURITY.md](SECURITY.md)
