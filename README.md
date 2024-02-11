# gomiddleman

`gomiddleman` is a secure, and efficient proxy service written in Go, designed 
to facilitate secure communication between clients and backend services. 
Operating as a Man-in-the-Middle (MitM) proxy, gomiddleman specializes in 
handling mutual TLS (mTLS) connections, ensuring that all communications are 
authenticated and encrypted without compromising on performance.

See our development status from [Issue 1](https://github.com/hyperifyio/gomiddleman/issues/1).

## Quick Start Guide

To get started with `gomiddleman`, download the latest release for your operating
system and execute it with the necessary flags:

```bash
wget https://github.com/hyperifyio/gomiddleman/releases/download/v0.0.1/gomiddleman-v0.0.1-linux-amd64.zip
unzip ./gomiddleman-v0.0.1-linux-amd64.zip
cd gomiddleman-v0.0.1-linux-amd64
./gomiddleman
```

For global installation:

```
sudo cp ./gomiddleman /usr/local/bin/gomiddleman
```

## Setting Up for Development

To build `gomiddleman` from source:

```bash
git clone git@github.com:hyperifyio/gomiddleman.git
cd gomiddleman
make
./gomiddleman
```

## License

Copyright (c) Heusala Group Ltd. All rights reserved.

Each software release of `goeventd` is initially under the HG Evaluation and 
Non-Commercial License for the first two years. This allows use, modification, 
and distribution for non-commercial and evaluation purposes only. Post this 
period, the license transitions to the standard MIT license, permitting broader
usage, including commercial applications. For full details, refer to the 
[LICENSE.md](LICENSE.md) file. 

**Commercial usage licenses can be obtained under separate agreements.**

## Usage and Configuration

For usage and configuration options:

```bash
./gomiddleman --help
```

Output:

```
Usage of ./gomiddleman:
  -ca string
    	proxy ca as PEM file (default "ca.pem")
  -cert string
    	proxy certificate as PEM file (default "cert.pem")
  -client-ca string
    	CA file to use for dynamic client certificate generation (default "ca.pem")
  -client-ca-key string
    	CA key to use for dynamic client certificate generation (default "ca-key.pem")
  -key string
    	proxy key as PEM file (default "key.pem")
  -port string
    	port on which the proxy listens (default "8080")
  -target string
    	target where to proxy connections (default "http://localhost:3000")
  -type string
    	type of proxy (tcp or tls) (default "tls")
```
