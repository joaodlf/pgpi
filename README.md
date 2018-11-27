# pgpi

**Note:** With the release of PostgreSQL 11, this tool is *mostly* not needed. 

> Partitions may have their own indexes, constraints and default values, distinct from those of other partitions. Indexes must be created separately for each partition.

-[PostgreSQL 10 documentation](https://www.postgresql.org/docs/10/static/ddl-partitioning.html)

**P**ost**g**res **P**artition **I**ndex is a CLI tool to facilitate creating an INDEX across multiple 
PostgreSQL 10 partitions/tables.

## Installation

A direct install is provided for macOS, Linux, and OpenBSD:

```
curl https://raw.githubusercontent.com/joaodlf/pgpi/master/install.sh | sh
```

Binaries (including Windows) can also be [downloaded](https://github.com/joaodlf/pgpi/releases).

## Usage

<p align="center">
  <img src="http://i.imgur.com/imSHZPj.jpg" alt="image"/>
</p>




