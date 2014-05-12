package main

import (
	"github.com/qiniu/log"
	"github.com/shxsun/go-sh"
	"github.com/shxsun/goyaml"
	"io/ioutil"
	"path/filepath"
	"strings"
)
import beeutils "github.com/astaxie/beego/utils"

// download src
func (b *Builder) get() (err error) {
	log.Debug("start get src to:", b.srcDir)
	exists := beeutils.FileExists(b.srcDir)
	b.sh.Command("go", "version").Run()
	if !exists {
		err = b.sh.Command("go", "get", "-v", "-d", b.project).Run()
		if err != nil {
			return
		}
	}
	b.sh.SetDir(b.srcDir)
	if b.ref == "-" {
		b.ref = "master"
	}
	// get code from remote
	if err = b.sh.Command("git", "fetch", "origin").Run(); err != nil {
		return
	}
	// change branch
	if err = b.sh.Command("git", "checkout", "-q", b.ref).Run(); err != nil {
		return
	}
	// update code
	if err = b.sh.Command("git", "merge", "origin/"+b.ref).Run(); err != nil {
		log.Warn("git merge error:", err)
		//return
	}
	// get sha
	out, err := sh.Command("git", "rev-parse", "HEAD", sh.Dir(b.srcDir)).Output()
	if err != nil {
		return
	}
	b.sha = strings.TrimSpace(string(out))

	// parse .gobuild
	b.rc = new(Assembly)
	rcfile := "public/gobuildrc"
	if b.sh.Test("f", ".gobuild") {
		rcfile = filepath.Join(b.srcDir, ".gobuild")
	}
	data, err := ioutil.ReadFile(rcfile)
	if err != nil {
		return
	}
	err = goyaml.Unmarshal(data, b.rc)
	return
}
