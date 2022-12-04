module github.com/situ-vault/situ-vault/cmd

go 1.19

require (
	github.com/situ-vault/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/nicksnyder/basen v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
)

replace github.com/situ-vault/situ-vault/pkg v0.0.0 => ./../pkg
