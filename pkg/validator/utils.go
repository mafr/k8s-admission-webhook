package validator

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
    adm "k8s.io/api/admission/v1beta1"
    apps "k8s.io/api/apps/v1"
)

// Convenience function to extract a Deployment from the request, if it exists.
func GetDeployment(req *adm.AdmissionRequest) (*apps.Deployment, bool) {
    if req.Kind.Kind != "Deployment" {
        return nil, false
    }

    // TODO: req.Object.Raw is the object *before* default values have been applied.
    var deployment apps.Deployment
    if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
        log.WithError(err).Error("failed to unmarshal deployment")
        return nil, false
    }

    return &deployment, true
}
