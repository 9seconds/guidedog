// Package options defines common options set for the guide-dog app.
package options

import (
	"fmt"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	lockfile "github.com/9seconds/guide-dog/lockfile"
)

// Options is just a storage of the possible options with some interpretations.
type Options struct {
	ConfigFormat    ConfigFormat
	ConfigPath      string
	Debug           bool
	Envs            map[string]string
	GracefulTimeout time.Duration
	LockFile        *lockfile.Lock
	PathsToTrack    []string
	PTY             bool
	Signal          syscall.Signal
	Supervisor      SupervisorMode
}

func (opt *Options) String() string {
	return fmt.Sprintf("<Options(configFormat='%v', configPath='%v', pathsToTrack='%v', debug='%t', envs='%v', gracefulTimeout='%d', lockFile='%v', signal='%v', supervisor='%v')>",
		opt.ConfigFormat,
		opt.ConfigPath,
		opt.PathsToTrack,
		opt.Debug,
		opt.Envs,
		opt.GracefulTimeout,
		opt.LockFile,
		opt.Signal,
		opt.Supervisor)
}

// NewOptions builds new Options struct based on the given parameter list
func NewOptions(debug bool,
	signal string,
	envs []string,
	gracefulTimeout time.Duration,
	configFormat string,
	configPath string,
	pathsToTrack []string,
	lockFile string,
	pty bool,
	supervise bool,
	restartOnConfigChanges bool) (options *Options, err error) {
	convertedConfigFormat, err := parseConfigFormat(configFormat)
	if err != nil {
		log.WithFields(log.Fields{
			"configFormat": configFormat,
			"error":        err,
		}).Errorf("Cannot convert configFormat.")
		return
	}

	convertedSignal, err := parseSignalName(signal)
	if err != nil {
		log.WithFields(log.Fields{
			"signal": signal,
			"error":  err,
		}).Errorf("Cannot convert signal.")
		return
	}

	convertedEnvs := parseEnvs(envs)

	supervisorMode := SupervisorModeNone
	if supervise {
		supervisorMode |= SupervisorModeSimple
	}
	if restartOnConfigChanges {
		supervisorMode |= SupervisorModeRestarting
	}

	var convertedLockFile *lockfile.Lock
	if lockFile != "" {
		convertedLockFile = lockfile.NewLock(lockFile)
	}

	options = &Options{
		ConfigFormat:    convertedConfigFormat,
		ConfigPath:      configPath,
		Debug:           debug,
		Envs:            convertedEnvs,
		GracefulTimeout: gracefulTimeout,
		LockFile:        convertedLockFile,
		PathsToTrack:    pathsToTrack,
		PTY:             pty,
		Signal:          convertedSignal,
		Supervisor:      supervisorMode,
	}

	return
}
