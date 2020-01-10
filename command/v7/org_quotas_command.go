package v7

import (
	"strconv"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/v7/shared"
	"code.cloudfoundry.org/cli/util/ui"
	"code.cloudfoundry.org/clock"
)

//go:generate counterfeiter . OrgQuotasActor

type OrgQuotasActor interface {
	//GetOrgQuotas(labelSelector string) ([]v7action.Buildpack, v7action.Warnings, error)
}

type OrgQuotasCommand struct {
	usage           interface{} `usage:"CF_NAME org-quotas"`
	relatedCommands interface{} `related_commands:"org-quota"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       OrgQuotasActor
}

func (cmd *OrgQuotasCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	sharedActor := sharedaction.NewActor(config)
	cmd.SharedActor = sharedActor

	ccClient, uaaClient, err := shared.GetNewClientsAndConnectToCF(config, ui, "")
	if err != nil {
		return err
	}
	cmd.Actor = v7action.NewActor(ccClient, config, sharedActor, uaaClient, clock.NewClock())

	return nil
}

func (cmd OrgQuotasCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(false, false)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	cmd.UI.DisplayTextWithFlavor("Getting org quotas as {{.Username}}...", map[string]interface{}{
		"Username": user.Name,
	})
	cmd.UI.DisplayNewline()

	// shamelessly copied over from buildpacks command

	//buildpacks, warnings, err := cmd.Actor.GetOrgQuotas(cmd.Labels)
	//cmd.UI.DisplayWarnings(warnings)
	//if err != nil {
	//	return err
	//}
	//
	//if len(buildpacks) == 0 {
	//	cmd.UI.DisplayTextWithFlavor("No buildpacks found")
	//} else {
	//	cmd.displayTable(buildpacks)
	//}
	return nil
}

func (cmd OrgQuotasCommand) displayTable(buildpacks []v7action.Buildpack) {
	if len(buildpacks) > 0 {
		var keyValueTable = [][]string{
			{"position", "name", "stack", "enabled", "locked", "filename"},
		}
		for _, buildpack := range buildpacks {
			keyValueTable = append(keyValueTable, []string{
				strconv.Itoa(buildpack.Position.Value),
				buildpack.Name,
				buildpack.Stack,
				strconv.FormatBool(buildpack.Enabled.Value),
				strconv.FormatBool(buildpack.Locked.Value),
				buildpack.Filename,
			})
		}

		cmd.UI.DisplayTableWithHeader("", keyValueTable, ui.DefaultTableSpacePadding)
	}
}
