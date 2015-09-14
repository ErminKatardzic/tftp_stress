# TFTP stress test
A small tool that I quickly put together because I couldn't find a simple way to generate a small burst of GET requests.
It generates strings to be used with the GET request by appending a random 6 hexadecimal number on the OUI.

It is split into two versions: one which sends requests one after the other while waiting for a response before sending another request and a version which sends all the requests almost at once, without waiting for a response.

### Dependency

* https://github.com/pin/tftp

### Usage example
```
./tftp_stress_serial -address=10.0.0.1 -port=69 -oui=00aabb -num=200
```

You can also use the -help flag to see the explanations for each flag and their default values.

##### Many thanks to [Dmitri Popov](https://github.com/pin) for making the TFTP client that I'm using here.
