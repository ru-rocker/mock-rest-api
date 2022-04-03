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

For instance:

    # URL based 
    export MOCK_CONFIG_FILE=https://raw.githubusercontent.com/ru-rocker/mock-rest-api/main/config/mock.yaml

    # File based
    export MOCK_CONFIG_FILE=/tmp/mock.yaml

For configuratioin file references, please visit `config/mock.yaml`.

### Delay
To make a more real situation, you can put delay parameter. The delay will generate a random number between the min and max values. Leave it empty if you do not want to set delay. Or set min == max for consistent delay value.

### Conditional
Additional notes: the routes support conditional return for each specific cases. 
Worth noting that the conditions take the first matching rule.
The precedences are how you locate your route condition in the yaml file. Smaller index always wins.

There are 4 type of conditional rules:
* request_header
* request_param
* query_param
* request_body

# Run
To run directly from the source

    go get .
    go run .

# Dockerized
To build the docker file - owner only :)

    docker build -t rurocker/mock-rest-api:0.3 -t rurocker/mock-rest-api:latest .


This application is already on the docker hub.

    docker pull rurocker/mock-rest-api:latest

Run the container

    docker run -p8080:8080 --rm -it \
       -v /your-config-folder/your-config.yaml:/tmp/rest-mock.yaml \
       -e MOCK_CONFIG_FILE=/tmp/rest-mock.yaml \
       rurocker/mock-rest-api:latest

# More Explanation
Visit https://www.ru-rocker.com/2022/04/03/mocking-rest-api-to-speed-up-development-time/ for detail explanation.

# Support the creator
If you want to show a little appreciation, you can put the donation to the creator :).

Here are the creator wallets:

##### ETH
<img src="https://www.ru-rocker.com/wp-content/uploads/2022/04/ETH.jpeg" width="250">

##### BNB
<img src="https://www.ru-rocker.com/wp-content/uploads/2022/04/BNB.jpeg" width="250">