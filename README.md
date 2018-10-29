terrafmt-blocks
==================

Ruby script for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.

Formatting all hcl blocks in markdown files:
```shell
find . | egrep "markdown$" | while read f; do cat $f | ruby terrafmt-blocks.rb > tmp; mv tmp $f; done
``` 
