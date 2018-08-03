rm -rf *.pem
rm -rf *.srl

echo "create a new certificate authority using this configuration"
openssl req -new -x509 -days 9999 -config ca.cnf -keyout ca-key.pem -out ca-crt.pem

echo "now that we have our certificate authority in ca-key.pem and ca-crt.pem, let’s generate a private key for the server"
openssl genrsa -out server-key.pem 4096

echo "generate a certificate signing request"
openssl req -new -config server.cnf -key server-key.pem -out server-csr.pem

echo "sign the request"
openssl x509 -req -extfile server.cnf -days 999 -passin "pass:password" -in server-csr.pem -CA ca-crt.pem -CAkey ca-key.pem -CAcreateserial -out server-crt.pem

echo "You can use server-key.pem, server-crt.pem, ca-crt.pem to set up https server. Don’t forget to tell your operating system to trust our newly created certificate authority by installing ca-crt.pem and marking it as trusted"

