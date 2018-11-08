# Makroud benchmark

## Requirements

At least **Go 1.11.0**.

## Setup

First of all, install latest version of others library:

```bash
go get -u -v "github.com/go-gorp/gorp"
go get -u -v "github.com/go-xorm/xorm"
go get -u -v "github.com/jinzhu/gorm"
go get -u -v "github.com/jmoiron/sqlx"
```

> **NOTE:** This benchmark doesn't include **[SQLBoiler](https://github.com/volatiletech/sqlboiler)**
and **[Kallax](https://github.com/src-d/go-kallax)** since they rely on code generation and not reflection.

## Execute

```bash
go test -run=XXX -bench=SelectAll -benchmem -benchtime=10s
go test -run=XXX -bench=SelectSubset -benchmem -benchtime=10s
go test -run=XXX -bench=SelectComplex -benchmem -benchtime=10s
go test -run=XXX -bench=Insert -benchmem -benchtime=10s
go test -run=XXX -bench=Update -benchmem -benchtime=10s
go test -run=XXX -bench=Delete -benchmem -benchtime=10s
```

## Benchmark graph

```bash
cd graph && python3 graph.py
```

## Acknowledgement

**SQLX** queries in this benchmark could be optimized.

However, since **Makroud** is built on top of **SQLX**, this benchmark highlight the difference of
performance using the same instructions.

**In every case**, using **SQLX** with a optimized workflow is more efficient than **Makroud**.

## Test Machine

```
OS:     Archlinux x86_64 Linux-4.18.8
CPU:    Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
Memory: 16GB
Go:     go version go1.11.2 linux/amd64
```

## Credits

* [SQLBoiler Benchmark repository](https://github.com/volatiletech/boilbench)
