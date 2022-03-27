# Overview
HTTP REST Mock Endpoint.

Build under `golang` and `gin` framework.

For now, this only for `content-type: application/json` data only.

# Usage
For starter, all you need to do only update the `config/mock.yaml` files. This is a default mode.

Whenver you need more dynamic configuration files, set `MOCK_CONFIG_FILE` environment variable.
The application will look at the environment variable for loading configuration.

For config file references, please visit `config/mock.yaml`.

# Run

    go get .
    go run .