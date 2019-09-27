# Kaburaya AutoScaler

An autoscaler controller based on queuing theory.

- The controller automatically estimates the throughput of servers while running.
- It reduces the effects due to delay of load detection.
- And it also reduces the effects due to delay of increasing servers.

## Architecture

![image](https://user-images.githubusercontent.com/1845486/65784527-35289080-e18d-11e9-98eb-a155ed8967cc.png)

## Simulation

```sh
$ go run cmd/simulator/
$ go run cmd/simulator/main.go --step 500 \
                               --DT 0.00001 \
                               --lambda 20,100,5000,200,10000,300,15000,400,10000 \
                               --mu 1000 \
                               --in-delay 1.0 \
                               --delay 6.0 \
                               --out-delay 6.0
```

![image](https://user-images.githubusercontent.com/1845486/65784472-16c29500-e18d-11e9-9fee-718bac3bbdbd.png)

## License

[MIT](https://github.com/monochromegane/kaburaya-autoscaler/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)

