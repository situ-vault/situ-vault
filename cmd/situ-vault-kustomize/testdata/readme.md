

## Run locally:

Configure *kustomize* plugin directory:
```
export XDG_CONFIG_HOME=$(pwd)/testdata
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
