# hacker-news-cli


# Build Dependencies

Check [dep](https://github.com/golang/dep) for dependency management for install instruction.

For install dependencies, just run:

    dep ensure

For update a dependency version, just run:

    vim Gopkg.toml # Update the contraintrelated to the package version
    dep ensure

# Build

	make build

# Running

Execute the following to run and print available option

	./build/hacker-news-cli -h

# Design decision

* Uses Hacker News api to access the page stories as it is more reliable.
* Limits the number of stories to 255 or the limit api return limit. Whatever come first.
* To speedupt the stories load, it uses goroutines to get it. This provided a significant 
performance inprovement.
* The csv output option is separated from the save in file option. It does not look like would
make sense to put both in same parameter.
* The mocking fot the http client is a customized one instead of use go-mock. Go-mock is powerfull
but not flexible. It bind mocking and asserting in the same module.

# Things to improve

* Test coverage far from ideal with only 38% of coverage
* Although part of the code would be reusable, add a different source would require a 
considerable ammout of new code, especially for the API.
* Test cases doesnt validate error cases.
* Timeout not configurable.
