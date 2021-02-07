module github.com/situ-vault/situ-vault/gui

go 1.16

require (
	fyne.io/fyne v1.4.3
	github.com/situ-vault/situ-vault/pkg v0.0.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/situ-vault/situ-vault/pkg v0.0.0 => ./../pkg
