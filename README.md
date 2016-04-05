# haproxy-lint

Simple linter for [HAProxy](http://haproxy.org) configuration.

## Installation

Tested against Go 1.5+

On Linux/OSX:

```
# set GOPATH to some valid path
export GOPATH=~/go && mkdir -p ~/go
go get github.com/abulimov/haproxy-lint
```

Now you have *haproxy-lint* binary.


## Usage

You can specify switch to JSON output
with `--json`flag (useful for editor plugins integration).

```console
$ haproxy-lint /etc/haproxy/haproxy.cfg
24:0:warning: ACL h_some declared but not used
18:0:warning: backend unused-servers declared but not used
25:0:critical: backend undefined-servers used but not declared


$ haproxy-lint --json /etc/haproxy/haproxy.cfg
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

## Rules

| #   | Severity | Rule                          |
|-----|----------|-------------------------------|
| 001 | critical | backend used but not declared |
| 002 | warning  | backend declared but not used |
| 003 | warning  | acl declared but not used     |
| 004 | critical | acl used but not declared     |


## License

Licensed under the [MIT License](http://opensource.org/licenses/MIT),
see **LICENSE**.
