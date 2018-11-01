terrafmt-blocks
==================

Ruby script for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

Formatting all hcl blocks in markdown files:
```shell
find . | egrep "markdown$" | while read f; do echo "formatting $f:"; cat $f | ./terrafmt-blocks.rb > tmp; mv tmp $f; done
``` 

Can also be used to show the change to make

```shell
./terrafmt-blocks.rb diff FILE
```

or count the blocks requiring formatting of total blocks

```shell
 ./terrafmt-blocks.rb count FILE
```

if no file is specified script reads from stdin:

```shell
cat FILE | ./terrafmt-blocks.rb diff
```

(todo proper examples with input & output)