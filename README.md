terrafmt-blocks
==================

Ruby script for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

First see what will be updated:
```shell
find . | egrep "markdown$" | sort | while read f; do ruby terrafmt-blocks.rb diff $f; done
``` 

Now format the terraform
```shell
find . | egrep "markdown$" | sort | while read f; do ruby terrafmt-blocks.rb fmt $f; done
``` 

if no file is specified stdin is used

```shell
cat FILE | ./terrafmt-blocks.rb diff
```

(todo proper examples with input & output)