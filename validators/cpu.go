package validators

import (
    "fmt"
    apps "k8s.io/api/apps/v1"
    core "k8s.io/api/core/v1"
    res "k8s.io/apimachinery/pkg/api/resource"
)


type CpuValidator struct {
    Max string
}


func (v CpuValidator) Validate(dep apps.Deployment) (bool, error) {
    maxCpu := parseCpu(v.Max)
    containers := dep.Spec.Template.Spec.Containers

    var sumCpuReq int64 = 0
    for _, c := range containers {
        req := c.Resources.Requests

        if getCpu(req) <= 0 {
            return false, fmt.Errorf("%s: cpu request not set", c.Name)
        }

        sumCpuReq += getCpu(req)
    }

    if maxCpu >= 0 && sumCpuReq > maxCpu {
        return false, fmt.Errorf("cpu request too high")
    }

    return true, nil
}

func parseCpu(val string) int64 {
    cpu := res.MustParse(val)
    return cpu.MilliValue()
}

func getCpu(res core.ResourceList) int64 {
    return res.Cpu().MilliValue()
}
