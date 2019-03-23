package server

import (
    "bytes"
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "k8s.io/apimachinery/pkg/runtime/serializer"
    "k8s.io/apimachinery/pkg/runtime"
    adm "k8s.io/api/admission/v1beta1"
    "github.com/mafr/kubernetes-admission-webhook/validators"
)


type Server struct {
    Decoder runtime.Decoder
}

func NewServer(listenAddress string) *http.Server {
    server := &Server{
        Decoder: createDecoder(),
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/validate", server.HandleValidate)

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
        log.Printf("failed to read body: %v", err)
        return
    }

    review := adm.AdmissionReview{}
    _, _, err = s.Decoder.Decode(data, nil, &review)
    if err != nil {
        log.Printf("failed to decode admission request: %v", err)
        return
    }

    var buf bytes.Buffer
    printJSON(&buf, &review)
    log.Printf("input: %s", buf.String())

    // TODO: do something with the review

    review.Response = validators.NewAdmissionResponse(false, "not with me")

    printJSON(w, review)
}

func printJSON(w io.Writer, data interface{}) {
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    enc.Encode(&data)
}

func createDecoder() runtime.Decoder {
    f := serializer.NewCodecFactory(runtime.NewScheme())
    return f.UniversalDeserializer()
}
