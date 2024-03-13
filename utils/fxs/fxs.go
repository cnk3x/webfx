package fxs

import (
	"fmt"
	"strings"

	"github.com/cnk3x/webfx/utils/log"

	"go.uber.org/fx/fxevent"
)

func Logger() fxevent.Logger {
	return fxPrintf(func(module, format string, v ...any) {
		log.Output(1, log.DEBUG, fmt.Sprintf("%-10s", module)+"  "+format, v...)
	})
}

type fxPrintf func(module, format string, v ...any)

const (
	HOOK     = "HOOK"
	PROVIDE  = "PROVIDE"
	INVOKE   = "INVOKE"
	ERROR    = "ERROR"
	SUPPLY   = "SUPPLY"
	REPLACE  = "REPLACE"
	DECORATE = "DECORATE"
	RUN      = "RUN"
	RUNNING  = "RUNNING"
	STOP     = "STOP"
	LOGGER   = "LOGGER"
)

func (lPrintf fxPrintf) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		lPrintf(HOOK, " OnStart %s executing (caller: %s)", HOOK, e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			lPrintf(HOOK, "OnStart %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			lPrintf(HOOK, "OnStart %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		lPrintf(HOOK, "OnStop %s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			lPrintf(HOOK, "OnStop %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			lPrintf(HOOK, "OnStop %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			lPrintf(ERROR, "Failed to supply %v: %+v", e.TypeName, e.Err)
		} else if e.ModuleName != "" {
			lPrintf(SUPPLY, "%s from module %q", e.TypeName, e.ModuleName)
		} else {
			lPrintf(SUPPLY, "%s", e.TypeName)
		}
	case *fxevent.Provided:
		var privateStr string
		if e.Private {
			privateStr = " (PRIVATE) "
		}
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				lPrintf(PROVIDE, "%s%s <= %s from module %q", privateStr, rtype, e.ConstructorName, e.ModuleName)
			} else {
				lPrintf(PROVIDE, "%s%s <= %s", privateStr, rtype, e.ConstructorName)
			}
		}
		if e.Err != nil {
			lPrintf(ERROR, "after options were applied: %+v", e.Err)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				lPrintf(REPLACE, "%v from module %q", rtype, e.ModuleName)
			} else {
				lPrintf(REPLACE, "%v", rtype)
			}
		}
		if e.Err != nil {
			lPrintf(ERROR, "Failed to replace: %+v", e.Err)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				lPrintf(DECORATE, "%v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				lPrintf(DECORATE, "%v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			lPrintf(ERROR, "after options were applied: %+v", e.Err)
		}
	case *fxevent.Run:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = fmt.Sprintf(" from module %q", e.ModuleName)
		}
		lPrintf(RUN, "%v: %v%v", e.Kind, e.Name, moduleStr)
		if e.Err != nil {
			lPrintf(ERROR, "returned: %+v", e.Err)
		}
	case *fxevent.Invoking:
		if e.ModuleName != "" {
			lPrintf(INVOKE, "%s from module %q", e.FunctionName, e.ModuleName)
		} else {
			lPrintf(INVOKE, e.FunctionName)
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			lPrintf(ERROR, "fx.Invoke(%v) called from:\n%+vFailed: %+v", e.FunctionName, e.Trace, e.Err)
		}
	case *fxevent.Stopping:
		lPrintf(STOP, "%v", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			lPrintf(ERROR, "Failed to stop cleanly: %+v", e.Err)
		}
	case *fxevent.RollingBack:
		lPrintf(ERROR, "Start failed, rolling back: %+v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			lPrintf(ERROR, "Couldn't roll back cleanly: %+v", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			lPrintf(ERROR, "Failed to start: %+v", e.Err)
		} else {
			lPrintf(RUNNING, "")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			lPrintf(ERROR, "Failed to initialize custom logger: %+v", e.Err)
		} else {
			lPrintf(LOGGER, "Initialized custom logger from %v", e.ConstructorName)
		}
	default:
		lPrintf(ERROR, "Unknown event type: %T", e)
	}
}
