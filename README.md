#### Arguments: 
  
  - count **string** - Pulled items (default "10")
  - file **string** - File contain hosts divided by new line
  - host **string** - Foreman FQDN
  - parallel **flag** - Parallel run
  - server **flag** - Run as web server daemon # not implemented now

---

#### Requirements
 - configuration solution [viper](https://github.com/spf13/viper)
 - terminal spinner/progress indicator [spinners](https://github.com/briandowns/spinner/)
 - sqlite3 driver [go-sqlite3](https://github.com/mattn/go-sqlite3/)
 
#### Config in TOML

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

#### Actions:
    // one thread 
    dbinit          - create base with $db_file name if not exist
    swefill         - create table with uniq SWEs
    swechec         - cross check SWE

    // may be multithreaded (-parallel)
    locations       - get locations
    swes            - get hostgroups with label ~ SWE/
    pclasses        - get puppet classes 
    sclasses        - get smart classes
    overridebase    - get smart clasees base
    overrideparams  - get smart clasees overrides


