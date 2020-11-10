# Arcade Operator (Marketplace Games Operator)

## Overview

The Arcade operator is an example operator built to showcase how to build an operator around an existing application. The operator deploys the [`marketplace-games-ui`](https://github.com/redhat-marketplace/marketplace-games-ui) application, a web app built with `Svelte` on the frontend and `Go` for the backend.

## Usage

If you are already connected to a k8s cluster in your local dev environment, iterative testing of the operators controller logic can be done using the following commands:

- Install Operators CRDs into cluster:
  - `make install`

- Run operator in local environment:
  - `make run`

### Clean-up

- Uninstall Operators CRDs from cluster:
  - `make uninstall`

_**Note**: If are running the project outside of your `$GOPATH` be sure to export `GO111MODULE=on` to activate go module support._

## Testing

Testing for the operator can be done using `make`

To run the full test-suite run:

```bash
$ make test

  go test ./... -coverprofile cover.out
  ...
  Running Suite: Controller Suite
  ===============================
  ...
```
