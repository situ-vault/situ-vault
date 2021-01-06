# situ-vault kustomize exec plugin

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

Service should be running at, e.g: [http://localhost:8666/]

# Ingress TLS setup

Create nginx ingress controller for a local setup:

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.43.0/deploy/static/provider/cloud/deploy.yaml
kubectl --namespace=ingress-nginx get all
```

To create new certificates:

```
mkdir testdata/cert
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout testdata/cert/key.pem -out testdata/cert/cert.pem -subj "/CN=localhost"
kubectl create secret tls secret-tls --key testdata/cert/key.pem --cert testdata/cert/cert.pem
```

Encrypt new files:

```
go install ./../../cmd/...
base64 -i testdata/cert/cert.pem -o testdata/cert/cert.b64.txt
base64 -i testdata/cert/key.pem -o testdata/cert/key.b64.txt

situ-vault encrypt -password=test-pw -vaultmode="C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE62#LB:CH100" -cleartext=file://./testdata/cert/cert.pem > ./testdata/cert/cert.enc.txt
situ-vault encrypt -password=test-pw -vaultmode="C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE62#LB:CH100" -cleartext=file://./testdata/cert/key.pem > ./testdata/cert/key.enc.txt
```

Test request:

```
curl -vk https://localhost:443/
```
