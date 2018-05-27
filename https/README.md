# Usage
1. Config tls/generate_cert.go, and run it to generate ``cert.pem`` and ``key.pem``. Note that, don't forget to add hostname to /etc/hots if using synthetic domain names

2. server/httpsServer.go will use ``cert.pem`` and ``key.pem`` to start server

3. request.go use ``cert.pem`` to connect, as self-signed certificate==ca certificate




Addition Reading

https://engineering.circle.com/https-authorized-certs-with-node-js-315e548354a2