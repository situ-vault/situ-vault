<a id="situ-vault-logo" href="#situ-vault-logo">
<img align="right" width="206" height="206" src="https://raw.githubusercontent.com/situ-vault/situ-vault/main/gui/situ-vault/icon/icon.png" alt="situ-vault app lock icon">
</a>

# situ-vault

> Handling secrets like a piece of cake!

Simple toolbox for working with encrypted secrets.

This repository contains tools to locally encrypt and decrypt secrets stored in a text based format.

* Symmetric encryption using current algorithms:
    * 🔐 **Salted key derivation** from a user supplied password (PBKDF2, Argon2id or Scrypt)
    * ✅ **Authenticated encryption** with state-of-the-art ciphers (AES-256-GCM, NaCl Secretbox or XChaCha20-Poly1305)
* Slim text based output:
    * 🔣 Various encodings that enable **easy copy 'n' paste** like Base32 or Base62 (but also Base64)
    * 📦 Selected cipher suite stored as a prefix to the ciphertext to allow decryption without further configuration
* Two user interfaces:
    * ⌨️ CLI: Simplistic **command line interface**, reads from files or flags and writes to stdout
    * 🖱️ GUI: Cross-platform **graphical user interface** built with ``fyne.io``
* Application interfaces:
    * 🐵 Concise Go pkg
    * 📃 Text based output and **only standard algorithms** allow straightforward implementation with other languages
* Cloud native tools integration:
    * ⛅ Usage of encrypted secrets via a **kustomize exec plugin** to safely store encrypted secrets in a repository but directly use them for Kubernetes application deployments

The name is inspired from the latin ***In situ*** which can mean ***in place***. As ``situ-vault`` only works with the
secrets locally and does not depend on a remote system like Hashicorp Vault, AWS KMS or Azure KV this quite nicely
captures its unique differentiator.

*Side Note:* If the exact opposite is actually desired, thus not using local symmetric keys,
then [sops](https://github.com/getsops/sops) might be worth a look. However, the greater flexibility that sops offers
also results in more visible complexity, as apparent by the lengthy Yaml or Json structured files for the results as
well as its configuration.

## Usage

### CLI

With flags: (not recommended)

```
echo test-pw | read -s PW
echo test-data | read -s DATA
situ-vault encrypt -password=$PW -cleartext=$DATA
SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32#LB:NO##IYKEB5WQVTPEQ===##I5VS45LGEXJXLZYNU7SYDC3ROJSDPGR2VG7KQSF2##END

situ-vault decrypt -password=$PW -ciphertext="SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32#LB:NO##IYKEB5WQVTPEQ===##I5VS45LGEXJXLZYNU7SYDC3ROJSDPGR2VG7KQSF2##END"
test-data
```

With files:

```
situ-vault encrypt -password="file://./pw.txt" -cleartext="file://./data.txt" > ./data.enc.txt

situ-vault decrypt -password="file://./pw.txt" -ciphertext="file://./data.enc.txt" > ./data.dec.txt
```

Specify a custom vault mode:

```
situ-vault encrypt -password=test-pw -cleartext=test-data -vaultmode="C:XCHACHA20_POLY1305#KDF:ARGON2ID_T1_M65536_C4#SALT:R32B#ENC:BASE62#LB:CH80"
SITU_VAULT_V1##C:XCHACHA20_POLY1305#KDF:ARGON2ID_T1_M65536_C4#SALT:R32B#ENC:BASE62#LB:CH80##lsDfYPcXuqspleYN0yYMw1EJu6mFfYMyP4X1L0HpZRf##
4YVoCE4cXfMxQasx7UsqnIOA6DtsOJswSk##END
```

The direction to a file currently adds a newline after the end of the decrypted content, which might be a problem for
some inputs!

Surrounding whitespace around ciphertexts is cleaned before parsing.

### GUI

#### Encrypt

<a id="situ-vault-encrypt" href="#situ-vault-encrypt">
<img alt="Encrypt GUI situ-vault" max-height="500px" src="https://raw.githubusercontent.com/situ-vault/situ-vault/main/doc/screens/encrypt.png">
</a>

#### Decrypt

<a id="situ-vault-decrypt" href="#situ-vault-decrypt">
<img alt="Decrypt GUI situ-vault" max-height="500px" src="https://raw.githubusercontent.com/situ-vault/situ-vault/main/doc/screens/decrypt.png">
</a>

#### Custom Mode

The GUI also allows to select the algorithms and other parameters for a custom vault mode:

<a id="situ-vault-mode" href="#situ-vault-mode">
<img alt="Encrypt GUI custom mode situ-vault" max-height="500px" src="https://raw.githubusercontent.com/situ-vault/situ-vault/main/doc/screens/mode.png">
</a>

### Kustomize

For the moment, this is documented in the specific readme: [readme.md](./cmd/situ-vault-kustomize/testdata/readme.md)

## Format

The output of ``situ-vault`` is formatted as text, called a ``message`` or ``vaultmessage``:

```
# whole vaultmessage: (encoding of salt and ciphertext depending on vaultmode)
SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32#LB:NO##IYKEB5WQVTPEQ===##I5VS45LGEXJXLZYNU7SYDC3ROJSDPGR2VG7KQSF2##END
<fix-version>##<vaultmode>##<salt>##<ciphertext>##<fix-end>

# with vaultmode:
##<code>:<value>#<code>:<value>#<...>##
```

A text based format instead of a binary one was chosen to enable easy diffing for version control and simple copying.
A self describing format enables decryption without further configuration by the user and allows the reuse of a mode for
further ciphertexts. The clearly readable ``##END`` allows the user to see if a message was copied completely.

The authentication tags are directly part of the ciphertext and not separated by ``##`` as this is the format most often
used by the crypto libraries. However, the nonces are not stored as a prefix of the ciphertext, as these are taken from
the key derivation function too and thus not needed in the message.

## Security

None of the actual cryptographic algorithms are re-implemented in this repository. Only the implementations from
the [crypto package](https://pkg.go.dev/crypto) of the Go standard library and its supplementary extensions
from [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) are used.

One of the main design assumptions for the CLI and GUI are that the local machine is kind of safe enough to at least
temporarily store and work with secrets in plain. The goal is rather to have a secret encrypted when it is stored in
another not as safe place or when in transit.

### Choice

Sometimes choice in cryptography algorithms is seen as undesired by library or tool authors.
However, choice also enables the user to comply to existing requirements or ease crypto agility when new algorithms
become necessary. Thus, in ``situ-vault`` the user is able to build a cipher suite, called ``mode`` or ``vaultmode``,
based on preferences, compliance requirements or tooling interoperability needs.

At least at this point in time, all the provided options are seen as suitable variants. The HKDF should only be used
when the input password is already a strong secret key e.g. when taken from a random source. It is still included here
as it offers a lightweight alternative in these specific cases.

Currently, the mode and other metadata are not authenticated, only the actual input cleartext is authenticated.
The ``vaultmode`` text might be added as additional data (where an AEAD is used) in subsequent versions of situ-vault.

#### Constructs:

Name | Notes | Value (``C``)
--- | --- | ---
AES-256-GCM | Key: 32 byte, Nonce: 12 byte, Tag: 16 byte | ``"AES256_GCM"``
NaCl Secretbox / XSalsa20-Poly1305 | Key: 32 byte, Nonce: 24 byte, Tag: 16 byte | ``"NACL_SECRETBOX"``
XChaCha20-Poly1305 | Key: 32 byte, Nonce: 24 byte, Tag: 16 byte | ``"XCHACHA20_POLY1305"``

#### Key Derivation Functions:

Name | Notes | Value (``KDF``)
--- | --- | ---
PBKDF2 | With SHA-256 and 10000 iterations | ``"PBKDF2_SHA256_I10K"``
Argon2id | With time=1, memory=64*1024 and threads=4 | ``"ARGON2ID_T1_M65536_C4"``
scrypt | With N=32768, r=8 and p=1 | ``"SCRYPT_N32768_R8_P1"``
HKDF | With SHA-256 and no info value | ``"HKDF_SHA256_NOINFO"``

#### Salts:

Name | Notes | Value (``SALT``)
--- | --- | ---
Random 8 bytes | n/a | ``"R8B"``
Random 16 bytes | n/a | ``"R16B"``
Random 24 bytes | n/a | ``"R24B"``
Random 32 bytes | n/a | ``"R32B"``

#### Encodings:

Name | Notes | Value (``ENC``)
--- | --- | ---
Hex | ``[0-9A-F]`` Base16 | ``"HEX"``
Base32 | ``[2-9A-Z]`` Base32, no ``0`` or ``1`` | ``"BASE32"``
Base62 | ``[0-9A-Za-z]`` Base64 without the special characters | ``"BASE62"`` 
Base64 | ``[0-9A-Za-z\+\/]`` Base64 standard | ``"BASE64"``
Base64 URL | ``[0-9A-Za-z\-\_]`` Base64 URL safe variant: ``-``, ``_`` instead of ``+``, ``/`` | ``"BASE64URL"``

#### Linebreaks:

Name | Notes | Value (``LB``)
--- | --- | ---
No linebreaks | n/a | ``"NO"``
After 80 characters | n/a | ``"CH80"``
After 100 characters | n/a | ``"CH100"``
After 120 characters | n/a | ``"CH120"``

### Overview

Schematic overview over the steps during encryption:

```mermaid
graph TB

    %% user inputs:
    pw[/Password/]
    cleartext[/Cleartext/]
    mode[/Vaultmode/]

    style pw stroke:#333,stroke-width:5px
    style cleartext stroke:#333,stroke-width:5px
    style mode stroke:#333,stroke-width:5px
    
    mode -.- rnd
    mode -.- kdf
    mode -.- ae
    mode -.- encode
    mode -.- lb
    
    %% mode influences subprocesses: (link styling not for ids)
    linkStyle 0 stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5
    linkStyle 1 stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5
    linkStyle 2 stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5
    linkStyle 3 stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5
    linkStyle 4 stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5
    
    %% key derivation:
    rnd[[Random Generator]]
    salt[/Salt/]
    rnd --> salt
    kdf[[Key Derivation Function]]
    pw --> kdf
    salt --> kdf
    key[/Key/]
    iv[/IV/]
    kdf --> key
    kdf --> iv
    
    %% encryption:
    ae[[Authenticated Encryption]]
    key --> ae
    iv --> ae
    cleartext-->ae
    ae --> ciphertext
    ciphertext[/Ciphertext & Tag/]
    
    %% encoding:
    encode[[Text Encoding]]
    salt --> encode
    ciphertext --> encode
    encode --> ciphertextEnc
    encode --> saltEnc
    ciphertextEnc[/Ciphertext & Tag Encoded/]
    saltEnc[/Salt Encoded/]
    
    %% line wrap:
    lb[Line Breaking]
    ciphertextWrapped[/Ciphertext & Tag Text Lines/]
    ciphertextEnc --> lb
    lb --> ciphertextWrapped
    
    %% concat message:
    version[/Version Prefix/]
    concat[Concatenation]
    message[/Vaultmessage/]
    style message stroke:#333,stroke-width:5px

    version --> concat
    mode --> concat
    saltEnc --> concat
    ciphertextWrapped --> concat

    concat --> message

```

### Comparison

The combination of algorithms in use can be compared to other well-established tools:

* [openssl-enc](https://www.openssl.org/docs/man1.1.1/man1/enc.html) symmetric cipher routines as of version 1.1.1:
    * algorithms:
        * Various ciphers supported including legacy ciphers, but no AEAD by choice
        * e.g. AES-256-CTR (in some other implementations also GCM; but there the authentication tag is discarded)
        * If enabled PBKDF2 can be used, defaults use SHA-256 and 10000 iterations with a salt (8 bytes) to derive the key (32 bytes) and an IV (16 bytes)
    * differences: no authenticated encryption; selected cipher suite or version information is not included in the output
    * result structure, binary or all together Base64 encoded:
        * a prefix: ``Salted__``
        * the salt
        * the ciphertext
* [Ansible Vault](https://docs.ansible.com/ansible/2.10/user_guide/vault.html#format-of-files-encrypted-with-ansible-vault) payload format 1.2:
    * algorithms:
        * AES-256-CTR
        * PBKDF2 with SHA-256 using 10000 iterations with a salt (32 bytes) to derive the AES key (32 bytes), HMAC key (32 bytes) and an IV (16 bytes)
    * differences: plaintext is padded to block size before encryption; separate HMAC is used as AES mode is CTR instead of GCM
    * result structure, separated by newlines:
        * a header (e.g. ``$ANSIBLE_VAULT;1.2;AES256;vault-id-label``)
        * hex encoded salt
        * hex encoded HMAC of the ciphertext
        * hex encoded ciphertext (newlines after 80 characters)

## Development

### Run all tests

```shell
for directory in pkg cmd gui
do
  ( cd "$directory" && go test ./... )
done
```

### Run all checks

```shell
for directory in pkg cmd gui
do
  ( cd "$directory" && go vet -tags ci ./... && staticcheck ./... )
done
```

### Run commands during development

```shell
go run ./cmd/situ-vault encrypt -password=hello -cleartext=world
go run ./gui/situ-vault
```

### Upgrade dependencies

```shell
for directory in pkg cmd gui
do
  ( 
    cd "$directory" 
    go get -u ./...
    echo 'resetting pkg replace version'
    sed -i '' 's:situ-vault/pkg v[0-9]*.[0-9]*.[0-9]*$:situ-vault/pkg v0.0.0:g' go.mod
    echo 'cleaning checksums'
    rm go.sum
    go mod tidy
    go test ./... 
  )
done
```

## Attribution

The situ-vault app icon emoji is from [OpenMoji](https://openmoji.org/). Icon license: [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/#)
