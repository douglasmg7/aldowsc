aldowsc
=======
Retrive aldo products

configuration
-----

## Set ZUNKA envoriment variable to define work path.
ZUNKA=~/.local/share/zunka

## Set ZUNKAENV envoriment variable to override configuration file.
ZUNKAENV=production
ZUNKAENV=development

Configuration file example (config.toml)
-------
Must be placed at location referenced by $ZUNKAPATH

```toml
# Zunka configuration file.

[all]
env = "development"
logDir = "log"
dbDir = "db"
listDir = "list"
xmlDir = "xml"

[zunkasrv]
logFileName = "zunkasrv.log"
dbFileName = "zunkasrv.db"
port = "8080"

[aldowsc]
logFileName = "aldowsc.log"
dbFileName = "aldowsc.db"
minPrice = 2000
maxPrice = 100000
```
