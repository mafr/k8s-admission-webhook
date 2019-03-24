package server

import (
    "encoding/json"
    "io"
    log "github.com/sirupsen/logrus"
)


func MustInitLogger(name string) {
    level, err := log.ParseLevel(name)
    if err != nil {
        log.WithField("value", name).Fatal("invalid log level")
    }

    log.SetLevel(level)
    log.AddHook(NewPrometheusHook())
}

func printJSON(w io.Writer, data interface{}) {
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    enc.Encode(&data)
}
