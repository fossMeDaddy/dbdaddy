# DBDaddy [draft]

A simple database management CLI tool.

safe for production use when version = `1.0.0`

stable enough for development DB use, NOW.

## Installation
only for linux or unix-based OS.

no windows.

why windows?

dont use windows.

LATEST VERSION:
```
curl https://raw.githubusercontent.com/fossMeDaddy/dbdaddy/main/installer/install.sh | bash
```

To install a specific version from the past add a `VERSION` shell variable before the bash command. Look for available releases in the releases tab.
```
curl https://raw.githubusercontent.com/fossMeDaddy/dbdaddy/main/installer/install.sh | VERSION=v0.5.1 bash
```

having doubts about installing a random binary from the internet? here's a remark from the author that might help:
> I used to be a frontend-heavy MERN-Stack specialized Javascript developer for God's sake!
> I am barely capable of writing software that works,
> you really think I can bypass your anti-virus, steal your information AND get away with it?
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

not only prisma, but almost every ORM suffers from problems like not being performant enough, not being comprehensive enough or being a black box that external users know very little about.

throw in the impatient gippidy syntax searches of the modern age and you've got yourself a flaming red hot ball of garbage queries eating away too much compute & memory on either the database or on your $2500 nodejs/deno/bun/whatever-the-fuck-next-shiny-runtime-is-gonna-be k8s cluster.

zooming out for a second...

what do you really need when working with databases?
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
- **schema definition in dbdaddy projects**
- **push/pull commands and merging migrations**

Planned:
- **studio UI for my fellow soydevs**, because these days nobody with a $1000 clerk subscription seems to use a CLI for more than 10s
- **db visualize & diagram**, aim is to be similar to dbdiagram.io in UI but use SQL for schema definition
- **hybrid migrations** allowing you to interleave your own sql between autogenerated sql
- **testing suite**, everyday i push to main, i push with fear, i dont want fear

## Migrations [unstable]
it's actually simpler than most people would think... so simple in fact that I would NOT recommend it for production use (for now).

indexes, triggers & custom types are not supported.

constraints, columns, tables & views are diffed and tracked for changes but there is no modification support as of now for any of them,
which means that, let's say if you make a field nullable in SQL, this tool will generate migrations to remove the column
from your database & then re-create it with SQL column definition such that the field is now nullable. DATA WILL BE LOST.

## Quickstart Guide

The CLI requires a config file: `dbdaddy.config.json` to connect with your database. It has connection credentials like host, port, params, user, password, etc.

When you first install & run `dbdaddy` it asks you for a connection uri, if not provided, default PostgreSQL credentials are used that can be changed later at any point in time.

`dbdaddy` handles multiple databases on your database server as "branches", there is always a "current branch"
on which you perform read/write operations like `inspect`ing the schema, `exec`uting SQL statements, etc.

You can `checkout` (change current branch) to different branches with the below command.

```
dbdaddy checkout postgres2
```

To create a new branch, copy the contents of the current branch into the new branch and checkout into it, use `checkout` with `-n`/`--new` option:

```
dbdaddy checkout -n my_new_db
```

Creating a new empty branch independent of the current branch and checking out into it looks like this (using `-c`/`--clean` option combined with `-n` option):

```
dbdaddy checkout -nc my_fresh_new_db
```

Let's take an example scenario:

Your friend gave you [this pg dump file](https://gist.github.com/fossMeDaddy/60c45d0b595d9167a3bd7556c1c31332), you'd like to create a table and give it back to them, let's see how you can do this with `dbdaddy`

Create a new database & checkout into it
```
dbdaddy checkout -nc friends_with_dbs
```

> NOTE:
> In case, `dbdaddy` hasn't been set up on your machine, you will be asked to input your local/dev database connection uri before your command is processed.

Run the restore command with a file option
```
dbdaddy restore --file /path/to/dump
```
This will execute & restore the dump in the current branch.

This should've restored your newly created database with your friend's dump (if I didn't fuck up).

Now that you both have the same DB state, let's inspect it
```
dbdaddy inspect --all
```
> NOTE:
> `--all` option prints everything, without this option, you're provided a searchable prompt to pick out a table/view and inspect its schema

After inspecting through all the tables, you can't see the table `person` which is very important for very valid reasons in this sample e-commerce db schema and not just for this tutorial.

Let's write a SQL script!
```sql
-- FILE: ./create_person.sql
CREATE TABLE person (
  id SERIAL PRIMARY KEY,
  name text NOT NULL,
  age integer CHECK (age > 18)
);
```

Executing the sql script, this should run successfully (again, if I didn't fuck up)
```
dbdaddy exec ./create_person.sql
```

Running `inspect --all` again, (hopefully) there you have it! a newly created `person` table.

Now run `dumpme`
```
dbdaddy dumpme
```
This will (hopefully) take a backup of your current branch and store it at a central location for all dumps near the JSON config file for future easy searching/finding.

The absolute path to this dump is also printed to the console, so you can find it and send it to your friend.

HAPPY DATABASING!

> NOTE:
> the tutorial in the quickstart guide covers only a subset of all the important commands, for more detailed info, please run `dbdaddy <COMMAND> --help` and `dbdaddy help <COMMAND>`

## Contributing Guide

contributors are advised to kindly turn off their fucking copilot :)

---

security audit will start after `1.0.0`, expect common & obvious security bugs to be fixed by `1.1.0`
