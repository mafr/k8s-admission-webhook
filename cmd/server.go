package main

import (
    "flag"
    "log"
    "github.com/mafr/kubernetes-admission-webhook/server"
)


func main() {
    listenAddress := flag.String("l", ":8443", "the address to listen on")
    certFile := flag.String("c", "/etc/webhook/certs/cert.pem", "server certificate in PEM format")
    keyFile := flag.String("k", "/etc/webhook/certs/key.pem", "server private key in PEM format")
    plainHttp := flag.Bool("p", false, "serve on plain HTTP for testing")

    flag.Parse()

    log.Printf("listening on %s", *listenAddress)
    if *plainHttp {
        log.Printf("running in plain HTTP mode")
    }

    httpServer := server.NewServer(*listenAddress)

    if *plainHttp {
        // This is for testing only, Kubernetes won't accept plain HTTP webhooks.
        log.Fatal(httpServer.ListenAndServe())
    } else {
        log.Fatal(httpServer.ListenAndServeTLS(*certFile, *keyFile))
    }
}
