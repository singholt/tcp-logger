# TCP Logger

Measure Firelens' log loss using this app that produces a high number of logs per second.

## Build the logger app

```
make docker_build_logger
```

## Logger arguments

The logger application accepts the following command-line arguments:

1. `-host` (string)
    - Description: Specifies the Fluent host address.
    - Default: "127.0.0.1"
    - Usage: `-host=<host_address>`

2. `-port` (int)
    - Description: Specifies the Fluent port number.
    - Default: "24224"
    - Usage: `-port=<port_number>`

3. `-steadyRate` (int)
    - Description: Sets the steady state rate.
    - Default: 10000
    - Usage: `-steadyRate=<rate>`

4. `-burst` (int)
    - Description: Sets the burst value.
    - Default: 0
    - Usage: `-burst=<value>`

5. `-time` (int)
    - Description: Specifies the time in minutes.
    - Default: 1
    - Usage: `-time=<minutes>`

6. `-async` (bool)
    - Description: Enables or disables Fluent Async mode.
    - Default: true
    - Usage: `-async=<true|false>`

7. `-subSecond` (bool)
    - Description: Enables or disables Fluent Sub-second precision.
    - Default: true
    - Usage: `-subSecond=<true|false>`

8. `-help` (bool)
    - Description: Prints the application usage information.
    - Default: false
    - Usage: `-help`

Note: All the above arguments are optional!

## Example usage
```
./logger -host="localhost" -port=24224 -steadyRate=1000 -burst=0 -time=2 -async=false -subSecond=true
```

To view the app usage information, run:
```
./logger -usage
```
