package goutil

import (
	"margo.sh/mg"
	"margo.sh/mgutil"
	"path/filepath"
)

const (
	ModEnvVar = "GO111MODULE"
)

// ModEnabled returns true of Go modules are enabled in srcDir
func ModEnabled(mx *mg.Ctx, srcDir string) bool {
	// - Inside GOPATH — defaults to old 1.10 behavior (ignoring modules)
	// - Outside GOPATH while inside a file tree with a go.mod — defaults to modules behavior
	// - GO111MODULE environment variable:
	//     unset or auto — default behavior above
	//     on — force module support on regardless of directory location
	//     off — force module support off regardless of directory location
	switch mx.Env.Getenv(ModEnvVar, "") {
	case "on":
		return true
	case "off":
		return false
	}

	bctx := BuildContext(mx)
	match := func(p string) bool {
		p = filepath.Join(p, "src")
		return p != "" && (mgutil.IsParentDir(p, srcDir) || srcDir == p)
	}
	if match(bctx.GOROOT) {
		return false
	}
	for _, gp := range PathList(bctx.GOPATH) {
		if match(gp) {
			return false
		}
	}
	return true
}