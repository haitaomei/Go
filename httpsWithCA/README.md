# Usage
1. Config ``tls/server.cnf``, and run ``generate.sh``, this script will genereate **certificate** and **key** for CA, and server

2. server/httpsServer.go will use ``server-crt.pem`` and ``server-key.pem`` to start server

3. request.go use ``ca-crt.pem`` as the caCert to connect. Note that, web browser can also import this certificate to trust our newly created certificate authority


Simple Node.js server
-----------
    var express = require('express');
    var fs = require('fs');
    var https = require('https');

    var server = express();
    server.get('/', function (req, res) {
        res.send("Hello World!");
    });


    var options = { 
        key: fs.readFileSync('localhost_cert/server-key.pem'), 
        cert: fs.readFileSync('localhost_cert/server-crt.pem'), 
        ca: fs.readFileSync('localhost_cert/ca-crt.pem'), 
    }; 

    https.createServer(options, server).listen(60443)  


Simple Node.js request
-----------
    var fs = require('fs'); 
    var https = require('https');

    var options = { 
        hostname: 'www.haitao.mei', 
        port: 60443, 
        path: '/', 
        method: 'GET', 
        ca: fs.readFileSync('ca-crt.pem')
    }; 
    var req = https.request(options, function(res) { 
        res.on('data', function(data) {
            console.log(JSON.parse(data));
        }); 
    }).on("error", (err) => {
        console.log("Error: " + err.message);
    }); 
    req.end();


Addition Reading
-----------
https://engineering.circle.com/https-authorized-certs-with-node-js-315e548354a2
