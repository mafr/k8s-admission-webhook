package validators

import (
    "fmt"
    apps "k8s.io/api/apps/v1"
)


type ReplicasValidator struct {
    Max int32
}


func (v ReplicasValidator) Validate(dep apps.Deployment) (bool, error) {
    replicas := dep.Spec.Replicas

    if *replicas > v.Max {
        return false, fmt.Errorf("replica count too high")
    }

    return true, nil
}
