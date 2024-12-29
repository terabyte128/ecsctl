# ecsctl

Make ECS behave like Kubernetes. Currently focuses on application management, since for my use case the actual cluster resources are managed by Terraform. Can autocomplete resource names to save time.

## Features

* `get`: clusters, services, task definitions, tasks
* `logs`: tail logs for a given task
* `exec`: run an interactive command in a task

## Requirements

* [Session Manager plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) for using `ecsctl exec` to run commands in containers. The AWS CLI itself is not required.

## Installation

You'll need to install [Go](https://go.dev/) (sorry).

```shell
git clone https://github.com/terabyte128/ecsctl
cd ecsctl
go build
```

Copy the resulting binary somewhere in your `$PATH`.

### Shell Completion

Follow instructions for your shell:

```shell
ecsctl completion
```
