# ztgrep

Search for file names and contents through nested archives.

Useful for locating data lost within many levels of compressed archives.

Supports the following compression formats for **both archives and files**:
- gzip
- bzip2
- xz (requires [xz-utils](https://tukaani.org/xz/) with `xz` CLI on `$PATH`)
- zstd (requires [zstd](https://github.com/facebook/zstd) with `zstd` CLI on `$PATH`)
- uncompressed

As well as the following archive formats:
- Tar (V7, USTAR, PAX, GNU, STAR)
- [ZIP](https://en.wikipedia.org/wiki/ZIP_(file_format)) (with size limitation)

Nested archives and compressed files must have a recognizable file extension to be searched.

If multiple paths are specified, they are searched in parallel with nondeterministic output order.
However, output order is deterministic for any single path.

Nested ZIP files must be read into memory to be searched.
By default, ZIP files larger 10 MB are not searched.
The `-z` option may be used to specify this size limit.

```
Usage:
  ztgrep [OPTIONS] regexp paths...

Search Options:
  -b, --skip-body     Skip file bodies
  -n, --skip-name     Skip file names inside of tarballs
  -z, --max-zip-size= Maximum zip file size to search (default: 10 MB)

Help Options:
  -h, --help          Show this help message
```
