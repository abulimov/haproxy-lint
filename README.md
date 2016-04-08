[![Build Status](https://travis-ci.org/abulimov/haproxy-lint.svg?branch=master)](https://travis-ci.org/abulimov/haproxy-lint)

# haproxy-lint

Simple linter for [HAProxy](http://haproxy.org) configuration.

## Editor plugins

There is official [Atom](http://atom.io) plugin - [linter-haproxy](https://atom.io/packages/linter-haproxy).

## Installation

Grab latest release on [releases page](https://github.com/abulimov/haproxy-lint/releases),
or build from source.

**To get more warnings you need a local HAProxy executable.**
Install it with [Homebrew](http://brew.sh) on OS X or package manager of your choice on Linux.


### Building from source

You need working Go compiler.
Tested against Go 1.5+

On Linux/OSX:

```console
# set GOPATH to some valid path
export GOPATH=~/go && mkdir -p ~/go
go get github.com/abulimov/haproxy-lint
```

Now you have *haproxy-lint* binary.


## Usage

You can specify switch to JSON output
with `--json`flag (useful for editor plugins integration).

Also you can manually disable running local HAProxy binary in check mode with
`--run-haproxy=false` flag.

```console
haproxy-lint /etc/haproxy/haproxy.cfg
24:0:warning: ACL h_some declared but not used
18:0:warning: backend unused-servers declared but not used
25:0:critical: backend undefined-servers used but not declared


haproxy-lint --json /etc/haproxy/haproxy.cfg
[
  {
    "col": 0,
    "line": 24,
    "severity": "warning",
    "message": "ACL h_some declared but not used"
  },
  {
    "col": 0,
    "line": 18,
    "severity": "warning",
    "message": "backend unused-servers declared but not used"
  },
  {
    "col": 0,
    "line": 25,
    "severity": "critical",
    "message": "backend undefined-servers used but not declared"
  }
]
```

## HAProxy check mode

In case if you have locally installed HAProxy,
it gets run with `-c` flag to check specified file,
and it's output gets parsed and returned as a linter warning.

If locally installed HAProxy is found, some of Native rules does not get
executed, as they just duplicate HAProxy's own checks.

## Native Rules

| #   | Severity | Rule                                          | Runs when local HAProxy found |
|-----|----------|-----------------------------------------------|-------------------------------|
| 001 | critical | backend used but not declared                 | yes                           |
| 002 | warning  | backend declared but not used                 | yes                           |
| 003 | warning  | acl declared but not used                     | yes                           |
| 004 | critical | acl used but not declared                     | no                            |
| 005 | warning  | rule order masking real evaluation precedence | no                            |
| 006 | warning  | duplicate rules found                         | yes                           |
| 007 | warning  | deprecated keywords found                     | no                            |

## License

Licensed under the [MIT License](http://opensource.org/licenses/MIT),
see **LICENSE**.
