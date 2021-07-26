package ssh

import (
	"go.uber.org/zap"
	"regeet.io/api/internal/conf"
)

// sshLog
type sshLog struct {
	*zap.Logger
}

// newLog
func newLog() *sshLog {
	return &sshLog{conf.Log.With(zap.String("pkg", "app.api.ssh"))}
}
