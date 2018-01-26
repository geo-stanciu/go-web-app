## Goals

This website is a demo (mostly a demo for golang and https://github.com/geo-stanciu/go-utils package).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## First Steps

As a first step after cloning this repository, you might need to run the following command (this downloads the needed dependencies):

```
go get -d
```

## TODO
- Password rules for register + change password
    - do not use common words

- Menu allocation on roles.

## Features
- Requests, roles and membership menu distribution initialized from initialize-database.go.
- Requests and controller + action are identified in the db.
  Calls to the propper action from a certain controller are made automatically.
- Auth. cookies are encrypted.
- Password rules:
    - do not use the last x passwords
    - use at least x characters
    - use at least x letters
    - use at least x numbers
    - use at least x uppercase
    - use at least x non alphanumerics
    - param for validity period
    - set the date when the password expires at password creation
    - do not use repetitive groups of letters
    - must not contain the username
    - redirect user to change his password if password is temporary
- Anti XRSF
- Router paths (Named here "Requests". See initialize-requests.go)
- Acces control (see initialize-database.go and initialize-access-rules.go)

## Connect Strings Examples

### PostgreSQL
```xml
<database>
  <db-type>postgres</db-type>
  <db-url>host=localhost port=5432 user=geo password=geo dbname=devel sslmode=disable options='--application_name=GoWebsiteExample --search_path=wmeter,public --client_encoding=UTF8'</db-url>
</database>
```

### MySql
```xml
<database>
  <db-type>mysql</db-type>
  <db-url>geo:geo@tcp(localhost:3306)/devel?parseTime=true&amp;sql_mode=%27ORACLE,TRADITIONAL%27</db-url>
</database>
```

### Oracle12c
```xml
<database>
  <db-type>oci8</db-type>
  <db-url>geo/geo@db1</db-url>
</database>
```
### Oracle11g
```xml
<database>
  <db-type>oracle11g</db-type>
  <db-url>geo/geo@db1</db-url>
</database>
```

### Sql Server
```xml
<database>
  <db-type>mssql</db-type>
  <db-url>server=localhost;database=devel;user id=geo;password=geo;port=1433;app name=GoWebsiteExample</db-url>
</database>
```

## Needs

Uses the following sql pachages:
- "github.com/denisenkom/go-mssqldb"
- "github.com/go-sql-driver/mysql"
- "github.com/lib/pq"
- "github.com/mattn/go-oci8"

If support is not needed for all of the above databases, remove some of the above imported packages.

## For Oracle driver:

1. Put oci8.pc path to your PKG_CONFIG_PATH environment variable.
2. You need:
   - go
   - oracle client or database
   - gcc from mingw64 - mine is installed in C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\bin
     and I put it in my path
   - pkg-config for Windows
     - copy pkg-config_0.26-1_win32.zip/bin/pkg-config.exe into
        C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\bin
     - copy gettext-runtime_0.18.1.1-2_win32.zip/bin/intl.dll into
        C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\bin
     - copy glib_2.28.8-1_win32.zip/bin/libglib-2.0-0.dll into
        C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\bin
3. go get github.com/mattn/go-oci8
