# joncalhoun-dl 🔥⬇

Downloads Go tutorial videos from joncalhoun.io

> **Before you proceed, note that you must be a paid user for the paid content to download**

Kindly create your account [here](https://courses.calhoun.io/signup?). Jon is a great teacher, consider buying his premium courses if you want to.

## How to use

+ Run the following commands

```bash
    $ git clone git@github.com:timolinn/joncalhoun-dl.git
    $ cd joncalhoun-dl
    $ go run main.go --email=jon@doe.com --password=12345 --course=gophercises
     [joncalhoun.io]: fetching video urls for gophercises
     [joncalhoun.io]: fetching data from https://courses.calhoun.io/courses/cor_gophercises...
```

### Command [OPTIONS]

+ `--email` : Your account email address. Sign up [here](https://courses.calhoun.io/signup?)
+ `--password` : Your account password. _Unlike the unix password prompt, this will not hide your password by default, you'll have to keep an eye over your shoulder 😉_
+ `--course` : This is the name of the course you want to download. **Defaults to `"gophercises"`**

### Supported courses

+ [x] [gophercises](https://courses.calhoun.io/courses/cor_gophercises)
+ [x] [testwithgo](https://courses.calhoun.io/courses/cor_test)
+ [ ] [algorithms](https://courses.calhoun.io/courses/cor_algo)

### Contributing

There is still a couple features to implement, check the TODO list below and send a pull request.

## TODO

+ [x] Add caching for requests
+ [x] Add default output directory
+ [ ] Add unit tests
+ [ ] provide packaged release and semver
+ [ ] check for authentication error
+ [ ] Add output directoy flag
+ [ ] prevent signin when using cache
+ [ ] choose video quality

### Authors

+ Timothy Onyiuke _([twitter](https://twitter.com/timolinn_))_

If you find this repository to be of any help, please consider giving it Star! 🔥
