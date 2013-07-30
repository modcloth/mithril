mithril
=======

[![Build Status](https://travis-ci.org/modcloth-labs/mithril.png?branch=master)](https://travis-ci.org/modcloth-labs/mithril)

HTTP -> AMQP proxy layer

## Development

The build/test cycle uses `make` and should not require specification of
a specific target:

``` bash
make
```

This will execute the `all` target, which cleans, builds, runs Go tests,
and runs the black-box <q>golden</q> suite.  Once some of the more
expensive one-time prerequisites are out of the way, the cycle should be
between 3-7 seconds.  If your build cycle is taking considerably longer,
please file an issue.

### Troubleshooting

#### `command not found: gvm`
If the build fails because `gvm` cannot be found, you can silence
the error by creating the following symlink (although use of
[gvm](https://github.com/moovweb/gvm) is encouraged, as this is what is
used on Travis CI):

``` bash
pwd
# => /Users/jdoe/workspace/mithril
ln -s $PWD ${GOPATH%%:*}/src/mithril
```

#### `command not found: psql`, or other failures related to PostgreSQL

In order to test the PostgreSQL integration, you will also have to have
a PostgreSQL server available.  The default URI used by the tests is the
following, which requires the presence of a `mithril_test` databasa:

```
postgres://localhost/mithril_test?sslmode=disable
```

You may specify your own URI via the `$MITHRIL_PG_URI` env var, e.g.:

``` bash
export MITHRIL_PG_URI='postgres://postgres@localhost?sslmode=disable'
```

#### `no such file to load -- minitest/spec (LoadError)` or `syntax error, unexpected ':', expecting ')'`

Your Ruby interpreter is too old :smile_cat:.  The test suite uses
`minitest` from the Ruby standard library as of version 1.9.  If you
install either `rvm` or `rbenv`, they should detect the presence of the
`.ruby-version` file and instruct you accordingly (or explosively).
