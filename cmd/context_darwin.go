//go:build darwin

package cmd

import (
	"github.com/senzing/senzing-tools/cmdhelper"
)

var ContextBoolsForOsArch = []cmdhelper.ContextBool{}

var ContextIntsForForOsArch = []cmdhelper.ContextInt{}

var ContextStringsForOsArch = []cmdhelper.ContextString{
	{
		Default: cmdhelper.OsLookupEnvString("SENZING_TOOLS_SENZING_DIRECTORY", ""),
		Envar:   "SENZING_TOOLS_SENZING_DIRECTORY",
		Help:    "Path to the SenzingAPI installation directory [%s]",
		Option:  "senzing-directory",
	},
	{
		Default: cmdhelper.OsLookupEnvString("SENZING_TOOLS_CONFIG_PATH", ""),
		Envar:   "SENZING_TOOLS_CONFIG_PATH",
		Help:    "Path to SenzingAPI's configuration directory [%s]",
		Option:  "config-path",
	},
	{
		Default: cmdhelper.OsLookupEnvString("SENZING_TOOLS_RESOURCE_PATH", ""),
		Envar:   "SENZING_TOOLS_RESOURCE_PATH",
		Help:    "Path to SenzingAPI's config, schema, and templates directory [%s]",
		Option:  "resource-path",
	},
	{
		Default: cmdhelper.OsLookupEnvString("SENZING_TOOLS_SUPPORT_PATH", ""),
		Envar:   "SENZING_TOOLS_SUPPORT_PATH",
		Help:    "Path to SenzingAPI's data directory [%s]",
		Option:  "support-path",
	},
}
