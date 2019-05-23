### Foreman's sync API
Server: sjc01-c01-pds10.c01.ringcentral.com

API: https://sjc01-c01-pds10.c01.ringcentral.com:8086/api/v1/

Web interface link: https://sjc01-c01-pds10.c01.ringcentral.com:8086/ 

On server API running as a systemd unit **gofsync**

Dependency management tool for third-party packages used [dep](https://github.com/golang/dep)

For compile project you need installed **dep**, then in project folder:
```bash
dep ensure
make
```

#### API endpoints
- Auth - POST body: ``{"username":"user","password":"pass","remember_me":false}``

| URI           | Method           |   |
| ------------- |:-------------:| -----:|
| /signin      | POST | Set jwt in cookie and user info in localstorage |
| /refreshjwt  | POST | Triggered after each request  |
- Host groups - POST body: ``{"source_host":"name","target_host":"name","target_hg_id":id,"source_hg_id":id}``

| URI                    | Method           |   |
| -----------------------|:-------------:| -----:|
| /hg                    |   GET         | Return all hostgroups from local base      |
| /hg/{host}             |   GET         | Return all hostgroups from local base for {host}      |
| /hg/{host}/{hgId}    |   GET         | Return hostgroup from local base    |
| /hg/overrides/{hgName}    |   GET         | Return all overrides for hostgroup from local base     |
| /hg/foreman/get/{host}/{hgName}    |   GET         | Return hostgroup from remote foreman     |
| /hg/foreman/update/{host}/{hgName}    |   GET         | Update host group in base     |
| /hg/foreman/check/{host}/{hgName}    |   GET         | Check host group exist in target foreman     |
| /hg/upload    |   POST         | Upload host group to target foreman     |
| /hg/check    |   POST         | Check host group exist in local base    |
| /hg/update    |   POST         | Update host group to target foreman (remove and upload)       |

- Locatios - POST body: ``NONE``

| URI                    | Method           |   |
| -----------------------|:-------------:| -----:|
| /loc                    |   GET         | Return all locations      |
| /loc/{host}             | POST | Update locations in base for {host} |
| /loc/overrides/{locName}                  |   GET         | Return all overrides for {locName}    |

- Environments - POST body: ``NONE``

| URI                    | Method           |   |
| -----------------------|:-------------:| -----:|
| /env/{host}                |   GET         | Return all environments for {host}      | 
| /env/{host}                |   POST         | Update environments in base for {host}      | 

- Other

| URI                    | Method           |   |
| -----------------------|:-------------:| -----:|
| /hosts                    |   GET         | Return all serving hosts      |
| /pc/{host}                    |   GET         | Return all puppet classes for host      |
| /sc/{sc_id}                    |   GET         | Return smart classe by ID      |


#### Cmd arguments
```bash
  -action string
        If specified run update for one of env|loc|pc|sc|hg|pcu
  -conf string
        Config file, TOML
  -file string
        File contain hosts divide by new line
  -host string
        Foreman FQDN
  -server
        Run as web server daemon
```
