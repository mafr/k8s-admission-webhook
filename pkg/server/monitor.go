package server

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    adm "k8s.io/api/admission/v1beta1"
    "strconv"
)

var validationReqs *prometheus.CounterVec

func init() {
    validationReqs = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "validation_requests_total",
            Help: "Number of validation requests processed, by result",
        },
        []string{"allowed"},
    )
}

func monitorResponse(response *adm.AdmissionResponse) {
    allowed := strconv.FormatBool(response.Allowed)

    validationReqs.WithLabelValues(allowed).Inc()
}
