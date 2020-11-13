# Nomad Job Concourse Resource

A concurse resource for jobs on [Hashicorp Nomad](https://www.nomadproject.io/).

## Source Configuration
```yaml
resource_types:
  - name: nomad
    type: registry-image
    source:
      repository: cioplenu/concourse-nomad-resource
      tag: latest

resources:
  - name: test-job
    type: nomad
    source:
      url: http://nomad-cluster.local:4646
      name: test
      token: XXXX
```

* `url`: (required) URL of the nomad api on a node of your cluster
* `name`: (required) name of the job
* `token`: Nomad ACL token

## Behavior

### `check`: Check for new [versions](https://www.nomadproject.io/api-docs/jobs#list-job-versions) of the job

Checks if there are new versions of the job. It uses the `nomad job history` command of the [nomad
CLI](https://www.nomadproject.io/docs/commands/job/history) and not all versions since the first
deployment might be available.

### `in`: noop

### `out`: Run a new version of the job

Deploys a new version of the job to the cluster using `nomad job run` by reading a HCL job file.
[Golang template variables](https://golang.org/pkg/text/template/) can be used to insert information
like a version on the fly.

#### Parameters
* `job_path`: (required) Path of the HCL job file to run
* `vars`: { [key: string]: string } dictionary of variables to substitute in the job file. Each key
  should be represented in the job file as `{{.key}}`
* `var_files`: { [key: string]: string } dictionary of paths to files to read to get variable
  values. Each key should be represented in the job file as `{{.key}}` and the values should be path
  to text files which content will be used as the variable value. Whitespace and trailing newlines
  will be trimmed from the value.

## Example

```yaml
resource_types:
  - name: nomad
    type: registry-image
    source:
      repository: cioplenu/concourse-nomad-resource
      tag: latest

resources:
  - name: sample-job
    type: nomad
    source:
      url: http://10.4.0.4:4646
      name: sample
      token: ((nomad-cluster-token))

  - name: sample-repo
    type: git
    source:
      uri: git@github.com:cioplenu/sample.git
      branch: main
      private_key: ((private-repo-key))

jobs:
  - name: deploy
    plan:
      - get: sample-repo
      - put: sample-job
        params:
          job_path: sample-repo/sample.nomad
          vars:
            registry_token: ((registry-token))
          var_files:
            version: sample-repo/ci/version
```

with a job file like:

```hcl
job "sample" {
  region      = "de"
  datacenters = ["dc1"]
  type        = "service"

  group "sample" {
    task "sample" {
      driver = "docker"

      config {
        image = "cioplenu/sample:{{.version}}"
        auth  = {
          username = "some-username"
          password = "{{.registry_token}}"
        }
        force_pull = true
      }

      resources {
        cpu    = 50
        memory = 50
        network {
          mbits = 1
        }
      }
    }
  }
}
```
and a version file like:
```text
9.1.4
```
