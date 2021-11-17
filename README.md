# check redirects
> this tool checks a list of redirects


## redirects file
Redirects file has to be named **redirects** and must be accessible next to checkredirects command.
Use argument ```-f=anotherfile.txt``` to specify another file.

File is csv, must by default be separated by ```\t```.  
Use argument ```-d``` to specify another delimiter. Eg. ```-d=;```

| checked url | expected target url |
|------------------------------|----------------------------------|
| http://localhost:10099/test1 | http://localhost:10099/redirect1 |
| http://localhost:10099/test2 | http://localhost:10099/redirect1 |

**Column1:** *checked url*  
**Column2:** *expected target url*

## concurrent requests
Set Parameter ```-c``` to define the number concurrent requests.

**Default=1**
