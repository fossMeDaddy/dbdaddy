# DBDaddy [draft]

A simple database management CLI tool.

## Installation
not available on windows cuz fuck you, thats why. use wsl or something...
```
curl https://raw.githubusercontent.com/fossMeDaddy/dbdaddy/main/installer/install.sh | bash
```

## Note from the author
Look prisma is a great tool from devex pov but applications dont run on developers, they run on a machine.

Their auto-generated migrations feature is really the only reason I personally use it, else I rawdawg them SQL right in the DB.

But I thought, well there has to be an alternative to it, you can't just recommend prisma to your CTO right!? how will you back your decision?
> sorry but i am a fullstack (frontend) developer and i don't know how to write SQL & definitely don't wanna learn it anyway,
> so here's an automagic migrations gen tool & with it comes a bloated ORM architecture deal with it.

With this tool, I aim to give you an alternative to prisma minus the ORM (write your own fuckin SQL)

Just define a schema, this tool will handle migrations between schema changes.

Along with that, some other useful functionality like executing SQL, backup & restore, branches, etc.

supported databases:
- [x] Postgres (cuz' i fuckin love it)
- [ ] MySQL (in progress, facing too many problems that shouldn't even have been there in pg... so that's gonna have to wait a little)
- [ ] SQlite (haven't got a chance to touch it, but it won't take time I promise)

## Features
Here are some features the CLI covers (for PG)
- checking out to database branches, cloning & then checking out into database branches (kinda like `git`)
- inspecting database/table schema
- deleting said databases
- dumping databases
- restoring databases from dump
- executing raw SQL queries & exporting query
- AND LASTLY, MY FAVOURITE, a custom diffing & migrations engine, to detect changes & generate necessary SQL for it

In Progress:
- project structure & schema definition
- switch migrations to manual mode because you are a competent developer & would like some control over your own database

## Migrations
it's actually simpler than most people would think... so simple in fact that I would NOT recommend it for production use (for now).

indexes, triggers & custom types are not supported.

constraints, columns, tables & views are diffed and tracked for changes but there is no modification support as of now for any of them,
which means that, let's say if you make a field nullable in SQL, this tool will generate migrations to remove the column
from your database & then re-create it with SQL column definition such that the field is now nullable. DATA WILL BE LOST.

---
> DISCLAIMER:
> you might draw a conclusion that these docs are mean to read, and you're not wrong! but that's the whole satire :)
