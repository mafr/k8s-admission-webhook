package main

import (
    "encoding/json"
    "io"
    apps "k8s.io/api/apps/v1"
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


func main() {
    dep, err := parseDeployment(os.Stdin)
    if err != nil {
        panic(err)
    }

    vals := []validators.DeploymentValidator{
        validators.CpuValidator{Max: "2000m"},
        validators.MemValidator{Guaranteed: true},
    }

    resp := validators.Validate(dep, vals)
    printJSON(os.Stdout, resp)
}
