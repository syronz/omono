# OMONO

[![BuildStatus](https://api.travis-ci.org/syronz/omono.svg?branch=master)](http://travis-ci.org/syronz/omono) 
[![ReportCard](https://goreportcard.com/badge/github.com/syronz/omono)](https://goreportcard.com/report/github.com/syronz/omono) 
[![codecov](https://codecov.io/gh/syronz/omono/branch/master/graph/badge.svg)](https://codecov.io/gh/syronz/omono)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6938819425f94f6f9d8046b4fdfdcbc1)](https://www.codacy.com/manual/syronz/omono?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=syronz/omono&amp;utm_campaign=Badge_Grade)
[![Coverage Status](https://coveralls.io/repos/github/syronz/omono/badge.svg?branch=master)](https://coveralls.io/github/syronz/omono?branch=master)
[![codebeat badge](https://codebeat.co/badges/f7ed90cf-4793-4b82-acd3-00fecf4e3817)](https://codebeat.co/projects/github-com-syronz-omono-master)
[![Maintainability](https://api.codeclimate.com/v1/badges/129904e9ab5aca417faa/maintainability)](https://codeclimate.com/github/syronz/omono/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/129904e9ab5aca417faa/test_coverage)](https://codeclimate.com/github/syronz/omono/test_coverage)
[![GolangCI](https://golangci.com/badges/github.com/gojek/darkroom.svg)](https://golangci.com/r/github.com/syronz/omono)
[![GoDoc](https://godoc.org/github.com/syronz/omono?status.png)](https://godoc.org/github.com/syronz/omono)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

refactoring omega for a monolith project

## Run
in the main directory

```bash
source cmd/restapi/restapi.sample.env
reflex -r '\.go' -s -- sh -c 'go run cmd/restapi/main.go'
```

## Logrus levels

```go
plog.ServerLog.Trace(err.Error())
plog.ServerLog.Debug(err.Error())
plog.ServerLog.Info(err.Error())
plog.ServerLog.Warn(err.Error())
plog.ServerLog.Error(err.Error())
plog.ServerLog.Fatal(err.Error())
plog.ServerLog.Panic(err.Error())
```

## Docker Requirement
run mysql
```bash
docker run --rm --name db-mysql -d -v mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=Qaz1@345 -e TZ='Asia/Baghdad' -p 3306:3306 mysql
```

## sed command for mass update
```
for i in $(grep -rl gorm);do sed -i 's/github.com\/jinzhu/gorm.io/' $i ;done
```
