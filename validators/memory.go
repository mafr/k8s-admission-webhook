package validators

import (
    "fmt"
    apps "k8s.io/api/apps/v1"
    core "k8s.io/api/core/v1"
)

type MemValidator struct {}


func (v MemValidator) Validate(dep apps.Deployment) (bool, error) {
    containers := dep.Spec.Template.Spec.Containers

    for _, c := range containers {
        req := c.Resources.Requests
        lim := c.Resources.Limits

        if getMem(req) <= 0 {
            return false, fmt.Errorf("%s: memory requests not set", c.Name)
        }

        if getMem(lim) <= 0 {
            return false, fmt.Errorf("%s: memory limit not set", c.Name)
        }

        if getMem(req) != getMem(lim) {
            return false, fmt.Errorf("%s: memory request and limit must be equal", c.Name)
        }
    }

    return true, nil
}

func getMem(res core.ResourceList) int64 {
    return res.Memory().Value()
}
