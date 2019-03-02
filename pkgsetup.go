// Copyright 2019 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goki/gopy/bind"
)

// 1 = pkg name, 2 = -user, 3 = version 4 = author, 5 = email, 6 = desc, 7 = url
const (
	setupTempl = `import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="%[1]s%[2]s",
    version="%[3]s",
    author="%[4]s",
    author_email="%[5]s",
    description="%[6]s",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="%[7]s",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: BSD License",
        "Operating System :: OS Independent",
    ],
    include_package_data=True,
)
`

	manifestTempl = `global-include *.so *.py
global-exclude build.py
`

	// 1 = pkg name
	bsdLicense = `BSD 3-Clause License

Copyright (c) 2018, The %[1]s Authors
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// 1 = pkg name, 2 = desc
	readmeTempl = `# %[1]s

%[2]s

`

	// 1 = pkg name, 2 = cmd, 3 = gencmd, 4 = vm (exe only)
	makefileTempl = `# Makefile for gopy pkg generation of python bindings to %[1]s
# File is generated by gopy (will not be overwritten though)
# %[2]s

PYTHON=%[4]s
PIP=$(PYTHON) -m pip

all: gen

gen:
	%[3]s

build:
	$(MAKE) -C %[1]s build

install:
	# this does a local install of the package, building the sdist and then directly installing it
	rm -rf dist build */*.egg-info *.egg-info
	$(PYTHON) setup.py sdist
	$(PIP) install dist/*.tar.gz

install-exe:
	# install executable into /usr/local/bin
	cp %[1]s/%[1]s /usr/local/bin

`
)

// GenPyPkgSetup generates python package setup files
func GenPyPkgSetup(odir, pkgname, cmdstr, user, version, author, email, desc, url, vm string) error {
	os.Chdir(odir)

	dashUser := user
	if user != "" {
		dashUser = "-" + user
	}

	sf, err := os.Create(filepath.Join(odir, "setup.py"))
	if err != nil {
		return err
	}
	fmt.Fprintf(sf, setupTempl, pkgname, dashUser, version, author, email, desc, url)
	sf.Close()

	mi, err := os.Create(filepath.Join(odir, "MANIFEST.in"))
	if err != nil {
		return err
	}
	fmt.Fprintf(mi, manifestTempl)
	mi.Close()

	lf, err := os.Create(filepath.Join(odir, "LICENSE"))
	if err != nil {
		return err
	}
	fmt.Fprintf(lf, bsdLicense, pkgname)
	lf.Close()

	rf, err := os.Create(filepath.Join(odir, "README.md"))
	if err != nil {
		return err
	}
	fmt.Fprintf(rf, readmeTempl, pkgname, desc)
	rf.Close()

	_, pyonly := filepath.Split(vm)
	gencmd := bind.CmdStrToMakefile(cmdstr)

	mf, err := os.Create(filepath.Join(odir, "Makefile"))
	if err != nil {
		return err
	}
	fmt.Fprintf(mf, makefileTempl, pkgname, cmdstr, gencmd, pyonly)
	mf.Close()

	return err
}
