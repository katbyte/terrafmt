## 0.3.0 (2020-01-06)

- add support for `return '` blocks ([#39](https://github.com/katbyte/terrafmt/issues/39))
- use terraform-exec to manage the terraform executable ([#36](https://github.com/katbyte/terrafmt/issues/36))
- returns actionable error codes ([#33](https://github.com/katbyte/terrafmt/issues/33))
- suppot addtional go format verbs ([#31](https://github.com/katbyte/terrafmt/issues/31))
- add option to mask go format versd in `-blocks` output ([#29](https://github.com/katbyte/terrafmt/issues/29))
- the `blocks` command can now return blocks null-seperated ([#25](https://github.com/katbyte/terrafmt/issues/25))
- JSON output format for `blocks` command ([#23](https://github.com/katbyte/terrafmt/issues/23))
- tolerate whitespace at the beginning of the first line ([#12](https://github.com/katbyte/terrafmt/issues/12))
- include the terraform version in the output of `version` ([#8](https://github.com/katbyte/terrafmt/issues/8))

## 0.2.0 (2020-02-29)

- Replace `filename#linenumber` output with `filename:linenumber` ([#16](https://github.com/katbyte/terrafmt/issues/16))
- Support directory traversal with pattern matching in diff and fmt ([#14](https://github.com/katbyte/terrafmt/issues/14))

## 0.1.0 (2020-02-25)

Initial release!
