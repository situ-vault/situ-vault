module github.com/situ-vault/situ-vault/cmd

go 1.16

require (
	github.com/situ-vault/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/situ-vault/situ-vault/pkg v0.0.0 => ./../pkg
