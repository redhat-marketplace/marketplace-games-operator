# Arcade Operator (Marketplace Games Operator)

## Overview

The Arcade Operator is an Operator built to showcase how to build an Operator around an existing application; the Operator was built using the [Operator-SDK](https://sdk.operatorframework.io/) which helped provide most of the scaffolding needed.

The Operator itself can deploy the [`marketplace-games-ui`](https://github.com/redhat-marketplace/marketplace-games-ui) application, a web app built with `Svelte` on the frontend and `Go` for the backend.

## Usage

If you are already connected to a Kubernetes cluster, and logged in as admin, in your local dev environment, iterative testing of the Operators controller logic can be done using the following commands:

- Install Arcade Operators CRDs into cluster:
  - `make install`

- Run operator in local environment:
  - `make run`

Deploy the sample Arcade Spec under
`config/samples/game_v1alpha1_arcade.yaml` and watch the logs in terminal where you ran `make run` to see the Operator reconciliation in action:

```bash
kubectl apply -f config/samples/game_v1alpha1_arcade.yaml
```
 
### Clean-up

- Uninstall Operators CRDs from cluster:
  - `make uninstall`

_**Note**: If are running the project outside of your `$GOPATH` be sure to export `GO111MODULE=on` to activate go module support._

## Testing

Testing for the Operator can be done using `make`

To run the full test-suite run:

```bash
$ make test

  go test ./... -coverprofile cover.out
  ...
  Running Suite: Controller Suite
  ===============================
  ...
```
