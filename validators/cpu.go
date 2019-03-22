package validators

import (
    "fmt"
    apps "k8s.io/api/apps/v1"
    core "k8s.io/api/core/v1"
)


type CpuValidator struct {
    Max int64
}


func (v CpuValidator) Validate(dep apps.Deployment) (bool, error) {
    containers := dep.Spec.Template.Spec.Containers

    var sumCpuReq int64 = 0
    for _, c := range containers {
        req := c.Resources.Requests

        if getCpu(req) <= 0 {
            return false, fmt.Errorf("%s: cpu request not set", c.Name)
        }

        sumCpuReq += getCpu(req)
    }

    if sumCpuReq > v.Max {
        return false, fmt.Errorf("cpu request too high")
    }

    return true, nil
}

func getCpu(res core.ResourceList) int64 {
    return res.Cpu().MilliValue()
}
