goem
====

A detailed description can be found [here](http://big-elephants.com/2013-09/goem-the-missing-go-extension-manager/)
This article explains the capabilities of goem with a detailed example.

# go extension manager

## What is this?

This is a little nifty tool to simplify the go development process.
The vanilla go only allows one go path at a time. If you want to setup
several development environments, you have to setup them all by yourself
and set the proper environmental variables.
This is where goem helps by setting up one Go environment for each project.
Right now goem only supports git repositories.

## Requirements

goem is written in pure go and only uses the core libraries.
All you need is a working go and git installation.

## Installation
First clone this repository.
To install goem you can run the build shell script.

```
bash build
```

It compiles goem and tries to install it to "/usr/local/bin".

Of course you can just run
```
go build main.go
```
and put the binary where you want.

## How does it work?

You need a so called "Gofile". This is json file, which specifies the dependencies of
your Go project. An example Gofile is provided in "example/Gofile".
Now you can run
```
goem bundle
```

This command creates a dot-go directory and uses this as your GOPATH.
Now run
```
goem build binary_name
```
to build your binary called "binary_name".

## Supported Actions

### Init

* goem init

Init creates and initializes a new Gofile in current working directory.

A new Gofile looks like

```json
        {
            "env": [
                {
                    "name": "development",
                    "packages": []
                }
            ],
            "testdir": "./test"
        }
```

If Gofile already exists init will not override it.

### List

* goem list

List lists all packages currently installed in your local gopath.

### Bundle

* goem bundle

Bundle reads your Gofile and fetches all packages in the desired version.

An example Gofile looks like that:

```json
        {
            "env" : [
                {
                    "name" : "development",
                    "packages": [
                        {
                            "name" : "github.com/adjust/goenv",
                            "branch" : "<cd3a33acbec38335c00b4ae252274827893d4e5b"
                        }
                    ]
                }
            ]
        }
```

* env

This is your GO_ENV which is read from GO_ENV environment variable.
If the environment variable is not set, goem sets it to 'development'.
If the Gofile provides only one env, this env is read.

* packages

Packages have a name and a branch. The name is the same as the one
you would use for a 'go get package'. The branch can be a branch,
a commit hash or a tag.
If branch is set to 'self', goem assumes you have your library *.go
files in a subfolder with the same name as your project and symlinks this
to the goem Gopath.

e.g.:

```
                For example the following directory structure:
                    project_name
                    |
                    |__project_name
                    |   |
                    |   |__
                    |
                    |__.go
                    |   |
                    |   |__src
                    |       |
                    |       |__your_account
                    |
                    |
                    |__main.go

```

Would create a symlink from project_name/projectname to
project_name/.go/src/your_account/project_name
Also you can use a relative or absolute path. This will also put a symlink to
the goem Gopath.

### Build

* goem build

Build your project with all the dependencies listed in your Gofile.
If no 'binname' is provided, the binary is called a.out.
'binname' can also be a path:

For example the following directory structure:

```
            project_root
            |
            |__bin
            |   |
            |   |__
            |
            |__main.go

```

To build the binary in bin/ called 'nifty_bin', you would call
'goem build bin/nifty_bin'

Absolute paths are also allowed.

### Test

* goem test

Run the tests as you would with go test.
As goem commands need to be run from the projects root,
you have to specify the test directory, if it is not the root directory.

For example the following directory structure:

```
            project_root
            |
            |__test_folder
            |   |
            |   |_nifty_test.go
            |
            |__main.go

```
In this case you need to add the 'testdir' key in your Gofile:

e.g.:

```json
                    {
                        "testdir": "test_folder"
                    }
```

Then run goem test

### Help

* goem help

Shows the avaiable commands.


## License

This Software is licensed under the MIT License.

Copyright (c) 2012 adjust GmbH,
http://www.adjust.com

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
