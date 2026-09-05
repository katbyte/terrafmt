# Security Policy

## Supported Versions

Only the [latest release](https://github.com/katbyte/terrafmt/releases/latest) is supported — please update before reporting an issue.

## Reporting a Vulnerability

Please **do not** open a public issue for security vulnerabilities.

Instead, report privately via [GitHub's private vulnerability reporting](https://github.com/katbyte/terrafmt/security/advisories/new).

I will do my best to acknowledge reports within 2 weeks and aim to release a fix or mitigation within 6 weeks for confirmed issues; timelines are best-effort.

## Scope

`terrafmt` parses and rewrites untrusted input — terraform blocks embedded in markdown documentation and Go test files. Issues where crafted input causes crashes, hangs, or writes outside the file being formatted are particularly relevant.
