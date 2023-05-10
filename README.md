Introduction
--------------------
EPR (Efficient Programmed Router) is a software tool that allows requests targeted at a specific port to be redirected to any other local server. It is a binary software tool that is written in Rust and supports threading. EPR operates by copying a received request and resending it to a local server targeted by the binding, and then responding to the client with the content of the response of the binded server.

EPR is a highly efficient and scalable tool, making it an ideal solution for handling large volumes of traffic. It is also highly flexible, allowing for easy configuration via a YAML file.

Features
--------------------
EPR has several features that make it a powerful tool for handling network traffic, including:

1. Port Binding: EPR allows requests targeted at a specific port to be redirected to any other local server, making it easy to route traffic between servers.

2. Threading Support: EPR supports threading, making it highly scalable and able to handle large volumes of traffic.

3. Efficient Routing: EPR is highly efficient, using advanced algorithms to quickly route traffic between servers.

4. Flexible Configuration: EPR is easy to configure via a YAML file, allowing for quick and easy setup.

5. Support for Server-Sent Events binding (binary files binding is on it's way)

Configuration
--------------------
EPR is configured via a YAML file that follows the following architecture, that should be saved in the root of the executable as `config.epr.yaml`:

```
server:
    http_port:       int
    https_port:      (optionnal)int
    https_cert_path: (optionnal)string
    https_key_path:  (optionnal)string
bindings:
    [binding_to_port]:
        max_res_wait_time:  usize
        transfert_scheme:   (optionnal)string
```

The `binding_to_port` is the port number to which the request is targeted, and the `max_res_wait_time` is the maximum time EPR waits for a response from the local server before timing out.

Usage
--------------------
To use EPR, you first need to download and install it on your server. Once installed, you can configure it via the YAML file.

After configuration, you can start EPR by running the binary file. Once running, EPR will listen for incoming requests on the specified port and redirect them to the appropriate local server.

Conclusion
--------------------
EPR is a powerful tool for routing network traffic between servers. It is highly efficient and scalable, making it ideal for handling large volumes of traffic. Its flexible configuration via a YAML file makes it easy to set up and configure, and its threading support ensures that it can handle high loads without slowing down.