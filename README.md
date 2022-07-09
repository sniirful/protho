# Protho
Protho is a proxy server application made to be easy to use and built for both Linux and Windows.

# Downloads
- [Linux x64](https://github.com/nSneerfulBike/protho/releases/download/0.1.0/protho-linux-amd64.zip)
- [Windows x64](https://github.com/nSneerfulBike/protho/releases/download/0.1.0/protho-windows-amd64.zip)

# Usage
The following samples will be shown using the Linux bash style. The arguments work the same way on all operating systems.

### Show help page
Shows the help page for this product.
```bash
./protho
```

### Port forwarding
Let's say you have a service listening on port 8080 which is not exposed, the following command will open port 80 and will forward all the data to port 8080.
```bash
./protho --in-port 80 --out-port 8080
```

### Port forwarding + filtering
Same case as before, but now you want to protect your service from attacks using some regex rules. You can setup a config.json and give it to the command like this.
```bash
./protho --in-port 80 --out-port 8080 -c ./config.json
```

### Port forwarding to external server
In this case you might have an external service on the web or in a private network (for instance in a virtual machine) and want to open this proxy server to act as a middleware.
```bash
./protho --in-port 12345 --out-server 192.168.122.12 --out-port 12345
```