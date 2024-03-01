# Migration DDL Checker

Usually in our release flow, we will want to do database migration before deploying source code.
Sometimes it's dangerous but sometimes it's not, depends on what migration code to execute.

This small tool analyse if migration DDL (i.e. CREATE TABLE, DROP INDEX, etc...) is hazardous or not.
If dangerous database operations contained, it reports those files.

# Usage

```
$ ./migration-ddl-checker --syntax [spanner|mysql|postgresql] --target-files [comma seperated file paths]
```

If syntax not specified, all target files will be reported (as hazardous files).
