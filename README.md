# ztgrep
Recursively search through nested compressed tarballs and files.

Supports the following compression formats:
- gzip
- bzip2
- xz (requires `xz` CLI)
- zstd (requires `zstd` CLI)
- uncompressed

Nested tarballs/files must have a recognizable file extension.

```
Usage:
  ztgrep [OPTIONS] regexp paths...

Search Options:
  -b, --skip-body  Skip file bodies
  -n, --skip-name  Skip file names inside of tarballs

Help Options:
  -h, --help       Show this help message
```
