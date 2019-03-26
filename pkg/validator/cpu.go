package validator

import (
    "fmt"
    adm "k8s.io/api/admission/v1beta1"
    core "k8s.io/api/core/v1"
    res "k8s.io/apimachinery/pkg/api/resource"
    "github.com/kelseyhightower/envconfig"
)


type CpuValidator struct {
    Max string
}

func NewCpuValidator() CpuValidator {
    v := CpuValidator{}
    envconfig.MustProcess("cpu", &v)
    return v
}

func (v CpuValidator) Validate(req *adm.AdmissionRequest) *adm.AdmissionResponse {
    if v.Max == "" {
        return NewResponse(true, "ok")
    }

    dep, ok := GetDeployment(req)
    if !ok {
        return NewResponse(true, "ok")
    }

    maxCpu := parseCpu(v.Max)
    containers := dep.Spec.Template.Spec.Containers

    var sumCpuReq int64 = 0
    for _, c := range containers {
        req := c.Resources.Requests

        if getCpu(req) <= 0 {
            return NewResponse(false, fmt.Sprintf("%s: cpu request not set", c.Name))
        }

        sumCpuReq += getCpu(req)
    }

    if maxCpu >= 0 && sumCpuReq > maxCpu {
        return NewResponse(false, fmt.Sprintf("total cpu request too high"))
    }

    return NewResponse(true, "ok")
}

func parseCpu(val string) int64 {
    cpu := res.MustParse(val)
    return cpu.MilliValue()
}

func getCpu(res core.ResourceList) int64 {
    return res.Cpu().MilliValue()
}
