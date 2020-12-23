## Build:

Build and place executable into a conforming *kustomize* exec plugin directory:

```
cd situ-vault-kustomize
go build -o testdata/kustomize/plugin/situ-vault/v1/secretgenerator/SecretGenerator .
```

## Run locally:

Configure *kustomize* plugin directory:

```
export XDG_CONFIG_HOME=$(pwd)/testdata
```

Provide password, e.g. by environment variable:

```
export SITU_VAULT_PASSWORD="test-pw"
```

Build and show resulting configuration (during development):

```
kustomize build --enable_alpha_plugins testdata/example
```

Build configuration and apply it directly:

```
kustomize build --enable_alpha_plugins testdata/example | kubectl apply -f - 
kubectl get all
```
