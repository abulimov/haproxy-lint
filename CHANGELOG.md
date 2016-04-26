## v0.6.2 [2016-04-26]

- More fixes for detection of unreachable rules with `or`

## v0.6.1 [2016-04-25]

- Fixed false positive for unreachable rules checks with `or`
- Fixed missing detection for multiple unreachable rules

## v0.6.0 [2016-04-22]

- Refactoring - moved all generic functions from checks to lib/acl and lib/backend
- New ACL type to keep ACLs negation and other attributes
- Added unreachable rules check

## v0.5.0 [2016-04-19]

- Added ability to exclude some config lines based on regexp for non-native checks
- Added more filtering for HAProxy output
- Refactored main func to improve readability

## v0.4.1 [2016-04-18]

- Added ability to exclude some config lines based on regexp

## v0.4.0 [2016-04-08]

- Execute checks in parallel
- More tests

## v0.3.0 [2016-04-08]

- Added deprecated rules check
- Fixed parser GetUsage for 'no option'
- Added ability to run HAProxy binary in check mode and return parsed warnings

## v0.2.1 [2016-04-07]

- Fixed duplicated rules recognition
- Small fixes and style changes

## v0.2.0 [2016-04-06]

- Added duplicate lines check
- Added rules precedence check

## v0.1.0 [2016-04-05]

- First versioned release
- Contains 4 basic rules
