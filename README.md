# terrafmt

==================



Go CLI for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

TODO: Depends on:


**PLEASE NOTE: this is a work in progress** 

First see what will be updated:
```shell
find . | egrep ".html.markdown" | sort | while read f; do ./terrafmt diff $f; done
``` 

Now format the terraform
```shell
find . | egrep ".html.markdown" | sort | while read f; do ./terrafmt $f; done
``` 

if no file is specified stdin is used

```shell
cat FILE | ./terrafmt diff
```

When working with provider acceptance tests with unquoted format placeholders you can use sed to make the blocks valid:



(todo proper examples with input & output)

