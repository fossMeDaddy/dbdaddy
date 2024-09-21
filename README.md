# DBDaddy [draft]

A simple database management CLI tool.

safe for production use when version = `v1.0.0`

stable enough for development DB use, NOW.

## Installation
only for linux or unix-based OS.

no windows.

why windows?

dont use windows.
```
curl https://raw.githubusercontent.com/fossMeDaddy/dbdaddy/main/installer/install.sh | bash
```

having doubts about installing a random binary from the internet? here's a remark from the author that might help:
> i am a recovering dev, used to write js for a living (soydev), barely capable of writing software that works, you really think I can bypass your anti-virus, steal your information AND get away with it?
>
> trust me bro.

## Usage Instructions
```
dbdaddy --help
```

## Note from the author

Prisma is great from devex pov but applications dont run on developers, they run on a machine.

I thought, well there has to be an alternative to it, you can't just recommend prisma to your CTO right? how will you pitch it?
> sorry but i am a full-stack (MERN with frontend specialisation) developer and i don't know how to write SQL & definitely don't wanna learn it. Anyway,
> here's an automagic migrations gen tool & with it comes a bloated ORM architecture deal with it.

jokes aside, some of the performance-related concerns were also raised by codedamn (company) in [this](https://codedamn.com/news/product/dont-use-prisma) article.

this problem arises when you create leaky abstractions, over things that shouldn't have been abstracted in the first place.

what do you need when working with databases?
1. a tool that gives you a good enough interface to do backups, migrations, custom one-off queries, etc. on the databases
2. efficient SQL to query the data
3. compile-time guard-rails when writing SQL

here's how you can do it WITHOUT AN ORM:
1. `DBDaddy` aka this tool
2. get good at SQL, bruh
3. amazing tools like `sqlc` for type-safe SQL queries

`DBDaddy` is not just a tool, it's an opinion.

## Features
supported databases:
- [x] Postgres
- [ ] MySQL
- [ ] SQlite

Here are some features the CLI covers (PG only, for now)
- using databases on a server like branches (like git but inferior)
- inspecting database/table schema
- deleting databases
- dumping databases
- restoring databases
- executing raw SQL queries & exporting queries as formatted readable text & CSV
- SQL migration file generation between schema changes

In Progress:
- **project structure & schema definition**
- **hybrid migrations allowing you to interleave your own sql between autogenerated sql**

Planned:
- **studio UI for my fellow soydevs**, because these days nobody with a $1000 clerk subscription seems to use a CLI for more than 10s
- **db visualize & diagram**, aim is to be similar to dbdiagram.io in UI but use SQL for schema definition
- **testing suite**, everyday i push to main, i push with fear, i dont want fear

## Migrations [unstable]
it's actually simpler than most people would think... so simple in fact that I would NOT recommend it for production use (for now).

indexes, triggers & custom types are not supported.

constraints, columns, tables & views are diffed and tracked for changes but there is no modification support as of now for any of them,
which means that, let's say if you make a field nullable in SQL, this tool will generate migrations to remove the column
from your database & then re-create it with SQL column definition such that the field is now nullable. DATA WILL BE LOST.
