package server

import (
    "bytes"
    "io/ioutil"
    log "github.com/sirupsen/logrus"
    "net/http"
    "time"
    "k8s.io/apimachinery/pkg/runtime/serializer"
    "k8s.io/apimachinery/pkg/runtime"
    adm "k8s.io/api/admission/v1beta1"
    "github.com/mafr/k8s-admission-webhook/pkg/validator"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


type Server struct {
    decoder runtime.Decoder
    validator validator.ValidatorConfig
}

func NewServer(listenAddress string, validator validator.ValidatorConfig) *http.Server {
    server := &Server{
        decoder: createDecoder(),
        validator: validator,
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/validate", server.HandleValidate)
    mux.Handle("/metrics", promhttp.Handler())

    httpServer := &http.Server{
        Addr:           listenAddress,
        Handler:        mux,
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   5 * time.Second,
        MaxHeaderBytes: 1 << 10,
    }

    return httpServer
}


func (s *Server) HandleValidate(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.WithError(err).Error("failed to read request body")
        return
    }

    review := adm.AdmissionReview{}
    _, _, err = s.decoder.Decode(data, nil, &review)
    if err != nil {
        log.WithError(err).Error("failed to decode admission request")
        return
    }

    log.WithFields(log.Fields{
        "uid": review.Request.UID,
        "kind": review.Request.Kind.Kind,
        "group": review.Request.Resource.Group,
        "version": review.Request.Resource.Version,
        "resource": review.Request.Resource.Resource,
        "name": review.Request.Name,
        "namespace": review.Request.Namespace,
    }).Info("received admission request")

    var buf bytes.Buffer
    printJSON(&buf, &review)

    log.Debugf("input: %s", buf.String())

    review.Response = s.validator.Validate(review.Request)
    monitorResponse(review.Response)

    printJSON(w, review)
}

func createDecoder() runtime.Decoder {
    f := serializer.NewCodecFactory(runtime.NewScheme())
    return f.UniversalDeserializer()
}
