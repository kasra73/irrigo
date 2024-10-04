# Irrigo Framework

## Overview

Irrigo is a flexible and extensible framework for building HTTP servers with a configurable pipeline of filters. This framework allows you to define global and local configurations for your filters and easily manage them through a configuration file.

## Project Structure

```txt
pkg/
  config/
    config.go
    loader.go
  server/
    handler.go
    pipeline.go
  utils/
    merge.go
go.mod
go.sum
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- A configuration file (e.g., `config.yaml`)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/irrigo.git
    cd irrigo
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Configuration

Create a configuration file (e.g., config.yaml) with the necessary settings for your filters and server. The configuration file should follow the structure expected by the `Global` configuration type.

### Example `config.yaml`

```yaml
app:
  name: simple-server
  logging:
    level: INFO
    output: stdout
  metrics:
    statsd_host: localhost
    statsd_port: 8125
  server:
    bind_ip: 127.0.0.1
    bind_port: 8080
  custom:
    custom_setting_1: value1
    custom_setting_2: value2

pipeline:
  - "health-check"
  - "rate-limit"
  - "simple-filter"

filters:
  health-check:
    use: health-check
    type: built-in
    options: {}

  rate-limit:
    use: rate-limit
    type: plugin
    plugin_path: /path/to/rate-limit.so
    options:
      limit: 100
      window: 60
    overrides:
      logging:
        level: DEBUG
  
  simple-filter:
    use: my-app
    type: built-in
    options:
      option1: option_value1
      option2: option_value2
    overrides:
      logging:
        level: WARN
      custom:
        custom_setting_1: override_value1

```

### Example `main.go`

Below is an example of how to use the Irrigo framework in your `main.go`

 file:

```go
package main

import (
 "irrigo/pkg/config"
 "irrigo/pkg/server"
)

func main() {
  var globalConfig config.Global

  // Create a new ConfigLoader
  cl := config.NewConfigLoader(globalConfig)

  // Load the configuration
  configFile := "config.yaml" // Change this to your config file path
  err := cl.LoadConfig(configFile)
  if err != nil {
    panic(err)
  }

  // Provide a FilterFactory implementation
  filterFactory := MyFilterFactory{}

  // Create the pipeline
  p, err := server.CreatePipeline(globalConfig, filterFactory)
  if err != nil {
    panic(err)
  }

  // Starting server
  err = p.Start()
  if err != nil {
    panic(err)
  }
}
```

### Implementing a FilterFactory

You need to implement the `FilterFactory` interface to create your custom filters. Below is an example implementation:

```go
package main

import (
 "irrigo/pkg/config"
 "irrigo/pkg/server"
)

type MyFilterFactory struct{}

func (f MyFilterFactory) NewFilter(
  name string,
  globalConfig config.Global,
  mergedConfig config.GeneralSettings,
  localConfig config.Filter,
  next server.Handler) server.Filter {
  // Implement your filter creation logic here
  return &MyFilter{
    next: next,
  }
}

type MyFilter struct {
  next server.Handler
}

func (f *MyFilter) Handle(req *http.Request, start server.StartResponse) (io.Reader, error) {
  // Implement your filter logic here
  return f.next.Handle(req, start)
}
```

### Running the Server

To run the server, simply execute:

```sh
go run main.go
```

## Contributing

Feel free to submit issues, fork the repository, and send pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
