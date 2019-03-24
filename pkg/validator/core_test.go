package validator

import (
    adm "k8s.io/api/admission/v1beta1"
    "testing"
)

func TestStaticValidator(t *testing.T) {
    r := StaticValidator(true, "ok")

    if resp := r.Validate(nil); !resp.Allowed {
        t.Errorf("StaticValidator(nil) = %v; want 'allowed'", resp.Allowed)
    }

    if resp := r.Validate(nil); resp.Result.Message != "ok" {
        t.Errorf("StaticValidator(nil); message=%q", resp.Result.Message)
    }
}

func TestPerform(t *testing.T) {
    r := ValidatorConfig{}

    r.Add(StaticValidator(true, "ok"))
    r.AddFunc(myValidator)

    if resp := r.Validate(nil); resp.Allowed {
        t.Errorf("review.Perform(nil) = %v; want false", resp.Allowed)
    }
}

func myValidator(req *adm.AdmissionRequest) *adm.AdmissionResponse {
    return NewResponse(false, "nope")
}
