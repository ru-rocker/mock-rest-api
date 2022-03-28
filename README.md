# Overview
HTTP REST Mock Endpoint.

Build under `golang` and `gin` framework.

For now, this only for `content-type: application/json` data only.

# Usage
For starter, all you need to do only update the `config/mock.yaml` files. This is a default mode.

Whenver you need more dynamic configuration, set `MOCK_CONFIG_FILE` environment variable.
The variable can retrieve a file path or an URL from the configuration file.
The application will look at the environment variable then loading configuration.
If the variable does not exist in, then it will load the default one.

For configuratioin file references, please visit `config/mock.yaml`.

# Run

    go get .
    go run .

# Dockerized


    docker build -t rurocker/mock-rest-api:0.1 -t rurocker/mock-rest-api:latest .


    docker run -p8080:8080 --rm -it \
       -v /your-config-folder/your-config.yaml:/tmp/rest-mock.yaml \
       -e MOCK_CONFIG_FILE=/tmp/rest-mock.yaml \
       rurocker/mock-rest-api:latest