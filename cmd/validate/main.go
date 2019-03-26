package main

import (
    "encoding/json"
    "io"
    adm "k8s.io/api/admission/v1beta1"
    apps "k8s.io/api/apps/v1"
    meta "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/util/yaml"
    "os"
    "github.com/mafr/k8s-admission-webhook/pkg/validator"
)


func parseDeployment(r io.Reader) (apps.Deployment, error) {
    dec := yaml.NewYAMLOrJSONDecoder(r, 1024)

    var dep apps.Deployment
    err := dec.Decode(&dep)

    return dep, err
}

func printJSON(w io.Writer, data interface{}) {
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    enc.Encode(&data)
}

func main() {
    dep, err := parseDeployment(os.Stdin)
    if err != nil {
        panic(err)
    }

    x, _ := json.Marshal(dep)
    rev := adm.AdmissionReview{
        TypeMeta: meta.TypeMeta{
            Kind: "AdmissionReview",
        },
        Request: &adm.AdmissionRequest{
            UID: "abc",
            Kind: meta.GroupVersionKind{
                Kind: "Deployment",
            },
            Object: runtime.RawExtension{
                Raw: x,
            },
        },
    }

    printJSON(os.Stdout, rev)

    val := validator.ValidatorConfig{}
    val.Add(validator.NewCpuValidator())
    val.Add(validator.NewMemValidator())
    val.Add(validator.NewReplicasValidator())

    resp := val.Validate(rev.Request)

    printJSON(os.Stdout, resp)
}
