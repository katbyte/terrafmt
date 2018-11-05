tfformatter
==================

[![CircleCI](https://circleci.com/gh/katbyte/chef-crowd/tree/master.svg?style=svg)](https://circleci.com/gh/katbyte/chef-crowd/tree/master)
[![Maintainability](https://api.codeclimate.com/v1/badges/e6dbb8dfc1fe75929d16/maintainability)](https://codeclimate.com/github/katbyte/tfformatter/maintainability)

Ruby script for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

**PLEASE NOTE: this is a work in progress** 

First see what will be updated:
```shell
find . | egrep "markdown$" | sort | while read f; do ruby tfformatter.rb diff $f; done
``` 

Now format the terraform
```shell
find . | egrep "markdown$" | sort | while read f; do ruby tfformatter.rb fmt $f; done
``` 

if no file is specified stdin is used

```shell
cat FILE | ./tfformatter.rb diff
```

(todo proper examples with input & output)