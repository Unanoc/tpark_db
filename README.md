# Forum API 
Main project for Database course at Technopark Mail.ru (https://park.mail.ru).

# API Documentation
https://tech-db-forum.bozaro.ru/

## Functional testing
The functionality of the API will be verified by automated functional testing.

Testing method:

* Docker container is collected from storage;
* Docker container is running;
* Go runs a script that will conduct testing;
* Docker container stops.

To create a Go script locally, simply run the command:
```
go get -u -v github.com/bozaro/tech-db-forum
go build github.com/bozaro/tech-db-forum
```
After that, the executable file `tech-db-forum` will be created in the current directory.

### Launch functional testing

To start functional testing, you need to run the following command:
```
./tech-db-forum func -u http://localhost:5000/api -r report.html
```

Parametr                              | Description
---                                   | ---
-h, --help                            | List of parameters
-u, --url[=http://localhost:5000/api] | URL of the application
-k, --keep                            | Continue testing after the first failed test
-t, --tests[=.*]                      | Mask of run tests (regular expression)
-r, --report[=report.html]            | File Name for Detailed Functional Test Report
