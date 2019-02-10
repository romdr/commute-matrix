# commute-matrix
A command line utility to display a matrix of commute times using google maps with real traffic at specific times

## Features

* Displays commute time matrices from/to a set of addresses, using with-traffic estimates from the Google Distance Matrix API
* Commute times can be obtained immediately, or can be configured and obtained at specific times
* Yaml configuration (rename config.yaml.example to config.yaml)
* Requires a Google Distance Matrix [API key](https://developers.google.com/maps/documentation/distance-matrix/get-api-key)

## Usage

```
λ go run main.go
Commute matrix scheduled...

λ go run main.go -now

               TO 01/20/2019 21:18              | 2111 7TH AVE, SEATTLE, WA 98121 | 85 PIKE ST, SEATTLE, WA 98101
+-----------------------------------------------+---------------------------------+-------------------------------+
  12309 SE 23rd Pl, Bellevue, WA 98005          |                            18.0 |                          17.7
  1005 8th St, Kirkland, WA 98033               |                            19.4 |                          20.6
  600 NW Richmond Beach Rd, Shoreline, WA 98177 |                            19.9 |                          21.0
  5700 24th Ave NW, Seattle, WA 98107           |                            13.2 |                          14.6
  2740 61st Ave SW, Seattle, WA 98116           |                            16.4 |                          17.4

              FROM 01/20/2019 21:18             | 2111 7TH AVE, SEATTLE, WA 98121 | 85 PIKE ST, SEATTLE, WA 98101
+-----------------------------------------------+---------------------------------+-------------------------------+
  12309 SE 23rd Pl, Bellevue, WA 98005          |                            18.1 |                          19.8
  1005 8th St, Kirkland, WA 98033               |                            20.2 |                          23.1
  600 NW Richmond Beach Rd, Shoreline, WA 98177 |                            20.4 |                          23.4
  5700 24th Ave NW, Seattle, WA 98107           |                            13.3 |                          14.3
  2740 61st Ave SW, Seattle, WA 98116           |                            16.5 |                          18.3

λ go run main.go -h
Usage of commute-matrix:
  -now
        Print the commute matrix now (instead of scheduling it)
```

## Dependencies

```
go get github.com/jasonlvhit/gocron
go get github.com/olekukonko/tablewriter
go get googlemaps.github.io/maps
go get gopkg.in/yaml.v2
```

## License

[MIT License](https://github.com/shazbits/commute-matrix/blob/master/LICENSE)
