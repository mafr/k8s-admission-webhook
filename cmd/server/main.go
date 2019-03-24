package main

import (
    "flag"
    log "github.com/sirupsen/logrus"
    "github.com/mafr/k8s-admission-webhook/pkg/server"
    "github.com/mafr/k8s-admission-webhook/pkg/validator"
)


func main() {
    listenAddress := flag.String("l", ":8443", "the address to listen on")
    certFile := flag.String("c", "/etc/webhook/certs/cert.pem", "server certificate in PEM format")
    keyFile := flag.String("k", "/etc/webhook/certs/key.pem", "server private key in PEM format")
    plainHttp := flag.Bool("p", false, "serve on plain HTTP for testing")

    flag.Parse()

    log.Infof("listening on %s", *listenAddress)
    if *plainHttp {
        log.Warn("running in plain HTTP mode (will NOT work in Kubernetes!)")
    }

    val := validator.ValidatorConfig{}
    val.Add(validator.CpuValidator{Max: "1000m"})
    val.Add(validator.MemValidator{Guaranteed: true})
    val.Add(validator.ReplicasValidator{Max: 3})

    httpServer := server.NewServer(*listenAddress, val)

    if *plainHttp {
        // This is for testing only, Kubernetes won't accept plain HTTP webhooks.
        log.Fatal(httpServer.ListenAndServe())
    } else {
        log.Fatal(httpServer.ListenAndServeTLS(*certFile, *keyFile))
    }
}
