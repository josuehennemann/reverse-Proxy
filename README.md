# Reverse-Proxy

Reverse-Proxy is a simple and powerful reverse proxy written in [golang](https://golang.org/).
With it you can receive http requests in a single address and redirect them transparently to other web servers. All this in a simple and agile way.

#### Prerequisites:
     Go 1.10 (or >)

     go get -u https://github.com/josuehennemann/logger
     go get -u https://github.com/josuehennemann/conf


#### Installation and usage:
To start using the reverse-proxy, simply perform the following step-by-step: 
First you must download the project:

	git clone https://github.com/josuehennemann/reverse-proxy

Browse directory:
   
	cd reverse-proxy

Compile and start web server:

	go build
	
    ./reverse-proxy -config=conf/reverse-proxy.conf

Then go to the browser:

	localhost:8080/github/josuehennemann

Ready, easy and practical!!!

#### Tests:
  Reverse-Proxy has client for testing. To run it you must open another terminal and access the client_test directory

    cd client_test
   
   Start the test web server
   
   	go run main.go

Now go to the browser:

    localhost:8080/test-get

Or send a file using [postman](https://www.getpostman.com/)
    
    localhost:8080/test-file
Or execute curl:
```
curl -X POST \
  http://localhost:8080/test-file \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F file=@/path/your/file.ext
```
#### Proxy rules:

By default all rules are within the file:

	files/rules.json

If there is any change in the file, it is necessary to reload the rules in the reverse-proxy, for that it is enough to access the URL:
    
    localhost:8080/admin/reload-rules

#### Settings:
The reverse-proxy has only 4 configurations, they are:
- **Httplisten** => Web address in http (Example: my-reverseproxy.com:80)
- **Httpslisten** => Web address in https, server is only started in https if there is value in this entry. (Example: my-reverseproxy.com:443)
- **Httpscertificate** => directory with the certificate files for https
- **Filepath** => redirection rule file directory
- **Logfile** => Directory where the log file will be saved (If empty, it will be the directory where the reverse-proxy is running)

These settings are in the file: conf/reverse-proxy.conf.


#### Limitations:
Some web servers may refuse access via proxy or respond http code 302 (redirect). Tthis usually occurs when no header is sent in HTTP request. If this happens, you can inform the [issues](https://github.com/josuehennemann/reverse-proxy/issues) so that you can work on this correction.
