terrafmt
==================

[![CircleCI](https://circleci.com/gh/katbyte/terrafmt/tree/master.svg?style=svg)](https://circleci.com/gh/katbyte/terrafmt/tree/master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/e80a8023626d4ecfa551cc75f88ae89f)](https://www.codacy.com/app/katbyte/terrafmt?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=katbyte/terrafmt&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/aaade40b149e1be650a8/maintainability)](https://codeclimate.com/github/katbyte/terrafmt/maintainability)

Ruby script for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

**PLEASE NOTE: this is a work in progress** 

First see what will be updated:
```shell
find . | egrep "markdown$" | sort | while read f; do ruby terrafmt.rb diff $f; done
``` 

Now format the terraform
```shell
find . | egrep "markdown$" | sort | while read f; do ruby terrafmt.rb fmt $f; done
``` 

if no file is specified stdin is used

```shell
cat FILE | ./terrafmt.rb diff
```

(todo proper examples with input & output)