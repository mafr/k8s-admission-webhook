package validator

import (
	adm "k8s.io/api/admission/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Validator interface {
	Validate(req *adm.AdmissionRequest) *adm.AdmissionResponse
}

type ValidatorFunc func(req *adm.AdmissionRequest) *adm.AdmissionResponse

// An Adapter to make sure we can use a ValidatorFunc as a Validator.
func (f ValidatorFunc) Validate(req *adm.AdmissionRequest) *adm.AdmissionResponse {
	return f(req)
}

type ValidatorConfig struct {
	validators []Validator
}

func (r *ValidatorConfig) Add(val Validator) {
	r.validators = append(r.validators, val)
}

func (r *ValidatorConfig) AddFunc(f ValidatorFunc) {
	r.validators = append(r.validators, f)
}

// Call all registered Validators and stop with the first negative one.
func (r *ValidatorConfig) Validate(request *adm.AdmissionRequest) *adm.AdmissionResponse {
	for _, val := range r.validators {
		resp := val.Validate(request)
		if !resp.Allowed {
			return resp
		}
	}

	return NewResponse(true, "success")
}

// Return a Validator that always returns the given response, ignoring any input.
// Useful for testing.
func StaticValidator(allowed bool, msg string) Validator {
	return ValidatorFunc(func(req *adm.AdmissionRequest) *adm.AdmissionResponse {
		return NewResponse(allowed, msg)
	})
}

func NewResponse(allowed bool, msg string) *adm.AdmissionResponse {
	return &adm.AdmissionResponse{
		Allowed: allowed,
		Result: &meta.Status{
			Message: msg,
		},
	}
}
