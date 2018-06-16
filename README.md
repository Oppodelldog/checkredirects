[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/checkredirects)](https://goreportcard.com/report/github.com/Oppodelldog/checkredirects) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/Oppodelldog/checkredirects/master/LICENSE) [![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/checkredirects/status.svg)](http://nulldog.de:12080/Oppodelldog/checkredirects) [![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/checkredirects/badge.svg?branch=master)](https://coveralls.io/github/Oppodelldog/checkredirects?branch=master)

# check redirects
> this tool checks a list of redirects


## redirects file
Redirects file has to be named **redirects** and must be accissible next to checkredirects command.

File is csv, must be separated by ```\t```

**Column1:** *checked url*

**Column2:** *expected target url*


## concurrent requests
Set Parameter ```-c``` to define the number concurrent requests.

**Default=1**
