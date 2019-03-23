package server

import (
    "encoding/json"
    "github.com/mafr/kubernetes-admission-webhook/pkg/validators"
    "log"
    adm "k8s.io/api/admission/v1beta1"
    apps "k8s.io/api/apps/v1"
)

func (s *Server) Review(req *adm.AdmissionRequest) *adm.AdmissionResponse {
    switch req.Kind.Kind {
	case "Deployment":
        var deployment apps.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
            log.Printf("failed to unmarshal deployment")
            return validators.NewAdmissionResponse(false, err.Error())
		}

        return validators.Validate(deployment, s.deploymentValidators)
    default:
        log.Printf("unknown resource type: %v", req.Kind.Kind)
    }

    return validators.NewAdmissionResponse(true, "success")
}
