# situ-vault

> Handling secrets like a piece of cake!

Simple toolbox for working with encrypted secrets.

This repository contains tools to locally encrypt and decrypt secrets stored in a text based format.

* Symmetric encryption using current algorithms:
    * 🔐 **Salted key derivation** from a user supplied password (PBKDF2, Argon2id or Scrypt)
    * ✅ **Authenticated encryption** with state-of-the-art ciphers (AES-256-GCM, NaCl Secretbox or XChaCha20-Poly1305)
* Slim text based output:
    * 🔣 Various encodings that enable **easy copy 'n' paste** like Base32 or Base62 (but also the ubiquitous Base64)
    * 📦 Selected suite of ciphers stored as a prefix to the ciphertext to allow easy decryption without further configuration
* Two user interfaces that should suit everyone:
    * ⌨️ CLI: Simplistic **command line interface**, reads from files or flags and outputs to stdout
    * 🖱️ GUI: Cross-platform **graphical user interface** built with ``fyne.io``
* Application interfaces:
    * 🐵 Easy to use Go pkg
    * 📃 Text based output and **only standard algorithms** allow straightforward implementation with other languages and tools
* Cloud native integration:
    * ⛅ Usage of encrypted secrets via a **kustomize exec plugin** to safely store encrypted secrets in a repository but directly use them for Kubernetes application deployments

The name is inspired from the latin ***In situ*** which can mean ***in place***. As ``situ-vault`` only works with the
secrets locally and does not depend on a remote system like Hashicorp Vault, AWS KMS or Azure KV this quite nicely
captures its unique differentiator.

*Side Note:* If the exact opposite is actually desired, thus not using local symmetric keys,
then [sops](https://github.com/mozilla/sops) might be worth a look. However, the greater flexibility that sops offers
also results in more visible complexity, as apparent by the lengthy Yaml or Json structured files for the results as
well as its configuration.

## Usage

### CLI

With flags: (not recommended)

```
situ-vault encrypt -password=test-pw -cleartext=test-data
SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##IYKEB5WQVTPEQ===##I5VS45LGEXJXLZYNU7SYDC3ROJSDPGR2VG7KQSF2

situ-vault decrypt -password=test-pw -ciphertext="SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##IYKEB5WQVTPEQ===##I5VS45LGEXJXLZYNU7SYDC3ROJSDPGR2VG7KQSF2"
test-data
```

With files:

```
situ-vault encrypt -password="file://./pw.txt" -cleartext="file://./data.txt" > ./data.enc.txt

situ-vault decrypt -password="file://./pw.txt" -ciphertext="file://./data.enc.txt" > ./data.dec.txt
```

This direction to a file currently adds a newline after the end of the decrypted content, which might be a problem for
some inputs!

Surrounding whitespace around ciphertexts is cleaned before parsing.

## Security

None of the actual cryptographic algorithms are re-implemented in this repository. Only the implementations from
the [crypto package](https://pkg.go.dev/crypto) of the Go standard library and its supplementary extensions
from [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) are used.

One of the main design assumptions for the CLI and GUI are that the local machine is kind of safe enough to at least
temporarily store and work with secrets in plain. The goal is rather to have a secret encrypted when it is stored in
another not as safe place or when in transit.

### Choice

The user is able to build a mode or cipher suite based on preferences or requirements. At least at this point in time,
all the provided options are seen as a suitable variant.

Constructs:

* AES-256-GCM (Key: 32 byte, Nonce: 12 byte, Tag: 16 byte) ``"AES256_GCM"``
* NaCl Secretbox (XSalsa20-Poly1305; Key: 32 byte, Nonce: 24 byte, Tag: 16 byte) ``"NACL_SECRETBOX"``
* XChaCha20-Poly1305 (Key: 32 byte, Nonce: 24 byte, Tag: 16 byte) ``"XCHACHA20_POLY1305"``

Key Derivation Functions:

* PBKDF2 with SHA-256 and 10000 iterations ``"PBKDF2_SHA256_I10K"``
* Argon2id with time=1, memory=64*1024 and threads=4 ``"ARGON2ID_T1_M65536_C4"``
* scrypt with N=32768, r=8 and p=1 ``"SCRYPT_N32768_R8_P1"``

Salts:

* Random 8 bytes ``"R8B"``
* Random 16 bytes ``"R16B"``

Encodings:

* Hex ``"HEX"``
* Base32 ``"BASE32"``
* Base62 ``"BASE62"`` (Base64 without the special characters)
* Base64 ``"BASE64"``
* Base64 URL ``"BASE64URL"`` (Base64 URL safe variant)

### Comparison

The combination of algorithms in use can be compared to other well-established tools:

* [openssl-enc](https://www.openssl.org/docs/man1.1.1/man1/enc.html) symmetric cipher routines as of version 1.1.1:
    * algorithms:
        * Various ciphers supported including legacy ciphers, but no AEAD by choice
        * e.g. AES-256-CTR (in some other implementations also GCM; but there the authentication tag is discarded)
        * If enabled PBKDF2 can be used, defaults use SHA-256 and 10000 iterations with a salt (8 bytes) to derive the key (32 bytes) and an IV (16 bytes)
    * differences: no authenticated encryption; selected cipher suite or version information is not included in the output
    * result structure, binary or all together Base64 encoded:
        * a prefix: `Salted__`
        * the salt
        * the ciphertext
* [Ansible Vault](https://docs.ansible.com/ansible/2.10/user_guide/vault.html#format-of-files-encrypted-with-ansible-vault) payload format 1.2:
    * algorithms:
        * AES-256-CTR
        * PBKDF2 with SHA-256 using 10000 iterations with a salt (32 bytes) to derive the AES key (32 bytes), HMAC key (32 bytes) and an IV (16 bytes)
    * differences: plaintext is padded to block size before encryption; separate HMAC is used as AES mode is CTR instead of GCM
    * result structure, separated by newlines:
        * a header (e.g. `$ANSIBLE_VAULT;1.2;AES256;vault-id-label`)
        * hex encoded salt
        * hex encoded HMAC of the ciphertext
        * hex encoded ciphertext (newlines after 80 characters)
