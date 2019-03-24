package validator

import (
    "fmt"
    adm "k8s.io/api/admission/v1beta1"
    core "k8s.io/api/core/v1"
)

type MemValidator struct {
    Guaranteed bool
}


func (v MemValidator) Validate(req *adm.AdmissionRequest) *adm.AdmissionResponse {
    dep, ok := GetDeployment(req)
    if !ok {
        return NewResponse(true, "ok")
    }

    containers := dep.Spec.Template.Spec.Containers

    for _, c := range containers {
        req := c.Resources.Requests
        lim := c.Resources.Limits

        if getMem(req) <= 0 {
            return NewResponse(false, fmt.Sprintf("%s: memory requests not set", c.Name))
        }

        if getMem(lim) <= 0 {
            return NewResponse(false, fmt.Sprintf("%s: memory limit not set", c.Name))
        }

        if v.Guaranteed && getMem(req) != getMem(lim) {
            return NewResponse(false, fmt.Sprintf("%s: memory request and limit must be equal", c.Name))
        }
    }

    return NewResponse(true, "ok")
}

func getMem(res core.ResourceList) int64 {
    return res.Memory().Value()
}
