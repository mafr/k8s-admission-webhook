package validators

import (
    adm "k8s.io/api/admission/v1beta1"
    apps "k8s.io/api/apps/v1"
    meta  "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentValidator interface {
    Validate(apps.Deployment) (bool, error)
}


func Validate(dep apps.Deployment, vals []DeploymentValidator) *adm.AdmissionResponse {
    val := ComposedValidator{Validators: vals}
    allowed := true
    msg := "success"

    if ok, err := val.Validate(dep); !ok {
        allowed = false
        msg = err.Error()
    }

    return NewAdmissionResponse(allowed, msg)
}


type ComposedValidator struct {
    Validators []DeploymentValidator
}


func (v ComposedValidator) Validate(dep apps.Deployment) (bool, error) {
    for _, val := range v.Validators {
        ok, err := val.Validate(dep)
        if !ok {
            return ok, err
        }
    }

    return true, nil
}


func NewAdmissionResponse(allowed bool, msg string) *adm.AdmissionResponse {
    return &adm.AdmissionResponse{
        Allowed: allowed,
        Result: &meta.Status{
            Message: msg,
        },
    }
}
