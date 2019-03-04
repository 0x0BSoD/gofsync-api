Config in TOML

```toml
[API]
username = "USERNAME"
password = "PASS"

[DB]
db_file = "FIELNAME SQLITE3"

[WEB]
port = "PORT int"

[RUNNING]
actions = ["dbinit","locations", "swes", "sclasses", "overridebase"]
```

Actions:

    dbinit             - create base with $db_file name if not exist
    locations          - get locations
    swes               - get hostgroups with label ~ SWE/
    sclasses           - get puppet classes
    overridebase       - get smart classes
    overrideadditional - get smart clasees overrides
    swefill            - create table with uniq SWEs
    swecheck           - cross check SWE

