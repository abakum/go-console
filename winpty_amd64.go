//go:build windows && amd64
// +build windows,amd64

package console

import (
	"embed"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

// go:embed winpty/386/*
//var winpty_deps embed.FS

//encrypted:embed winpty/amd64/*
var winpty_deps encryptedfs.FS

//go:embed winpty/amd64/*.enc
var winpty_depsEnc embed.FS

func init() {
	winpty_deps = encryptedfs.InitFS(winpty_depsEnc, key)
}
