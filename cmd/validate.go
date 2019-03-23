package main

import (
    "encoding/json"
    "io"
    "log"
    adm "k8s.io/api/admission/v1beta1"
    apps "k8s.io/api/apps/v1"
    meta "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/util/yaml"
    "os"
    "mafr.de/admission-policy/validators"
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


func review(req *adm.AdmissionRequest, vals []validators.DeploymentValidator) *adm.AdmissionResponse {
    switch req.Kind.Kind {
	case "Deployment":
        var deployment apps.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
            log.Printf("failed to unmarshal deployment")
            return validators.NewAdmissionResponse(false, err.Error())
		}

        return validators.Validate(deployment, vals)
    default:
        log.Printf("unknown resource type: %v", req.Kind.Kind)
    }

    return validators.NewAdmissionResponse(true, "success")
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

    vals := []validators.DeploymentValidator{
        validators.CpuValidator{Max: "2000m"},
        validators.MemValidator{Guaranteed: true},
        validators.ReplicasValidator{Max: 3},
    }

    resp := review(rev.Request, vals)
    printJSON(os.Stdout, resp)
}
