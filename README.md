terrafmt
==================

[![CircleCI](https://circleci.com/gh/katbyte/terrafmt/tree/master.svg?style=svg)](https://circleci.com/gh/katbyte/terrafmt/tree/master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5aa82bb8d0de4c52b270c9030297eea9)](https://www.codacy.com/app/katbyte/tfformatter?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=katbyte/tfformatter&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/e6dbb8dfc1fe75929d16/maintainability)](https://codeclimate.com/github/katbyte/tfformatter/maintainability)

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