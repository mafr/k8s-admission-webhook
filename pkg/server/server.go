package server

import (
	"bytes"
	"github.com/mafr/k8s-admission-webhook/pkg/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	adm "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
	"time"
)

type Server struct {
	decoder   runtime.Decoder
	validator validator.ValidatorConfig
}

func NewServer(listenAddress string, validator validator.ValidatorConfig) *http.Server {
	server := &Server{
		decoder:   createDecoder(),
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	review := adm.AdmissionReview{}
	_, _, err = s.decoder.Decode(data, nil, &review)
	if err != nil {
		log.WithError(err).Error("failed to decode admission request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	printJSON(&buf, &review)
	log.Debugf("input: %s", buf.String())

	if isResponsible(review) {
		review.Response = s.validator.Validate(review.Request)
	} else {
		review.Response = validator.NewResponse(true, "out of scope")
	}

	logReview(review)
	monitorResponse(review.Response)
	printJSON(w, review)
}

// TODO: implement inclusion/exclusion mechanism to define the scope
func isResponsible(review adm.AdmissionReview) bool {
	return review.Request.Namespace != "kube-system"
}

func createDecoder() runtime.Decoder {
	f := serializer.NewCodecFactory(runtime.NewScheme())
	return f.UniversalDeserializer()
}

func logReview(review adm.AdmissionReview) {
	log.WithFields(log.Fields{
		"uid":       review.Request.UID,
		"kind":      review.Request.Kind.Kind,
		"group":     review.Request.Resource.Group,
		"version":   review.Request.Resource.Version,
		"resource":  review.Request.Resource.Resource,
		"name":      review.Request.Name,
		"namespace": review.Request.Namespace,
		"allowed":   review.Response.Allowed,
		"message":   review.Response.Result.Message,
	}).Info("processed admission request")
}
