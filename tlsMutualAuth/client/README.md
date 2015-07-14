# gen

cfssl gencert -ca ../certsNkey/ca.pem -ca-key ../certsNkeys/ca-key.pem csr_client.json 
| cfssljson -bare client


# curl query:

 curl --cacert ../certsNkeys/ca.pem -E client.both.pem https://localhost:8080/hello

# p12 export

openssl pkcs12 -export -out client.p12 -inkey client-key.pem -in client.pem -certfile ../certsNkeys/ca.pem

