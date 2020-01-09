package v7

import (
	"context"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v7/shared"
	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/cli/actor/loggingaction"
)

//go:generate counterfeiter . LogsActor

type LogsActor interface {
	GetStreamingLogsForApplicationByNameAndSpace(appName string, spaceGUID string, client v7action.LogCacheClient) (<-chan v7action.LogMessage, <-chan error, context.CancelFunc, v7action.Warnings, error)
	GetRecentLogsForApplicationByNameAndSpace(appName string, spaceGUID string, client v7action.LogCacheClient) ([]v7action.LogMessage, v7action.Warnings, error)
	ScheduleTokenRefresh() (chan bool, error)
}

type LogsCommand struct {
	RequiredArgs    flag.AppName `positional-args:"yes"`
	Recent          bool         `long:"recent" description:"Dump recent logs instead of tailing"`
	usage           interface{}  `usage:"CF_NAME logs APP_NAME"`
	relatedCommands interface{}  `related_commands:"app, apps, ssh"`

	UI             command.UI
	Config         command.Config
	CC_Client      *ccv3.Client
	SharedActor    command.SharedActor
	Actor          LogsActor
	LogCacheClient v7action.LogCacheClient
}

func (cmd *LogsCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor(config)

	ccClient, uaaClient, err := shared.GetNewClientsAndConnectToCF(config, ui, "")
	if err != nil {
		return err
	}
	cmd.CC_Client = ccClient

	cmd.Actor = v7action.NewActor(ccClient, config, nil, uaaClient, clock.NewClock())
	cmd.LogCacheClient = shared.NewLogCacheClient(ccClient.Info.LogCache(), config, ui)
	return nil
}

func (cmd LogsCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(true, true)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	cmd.UI.DisplayTextWithFlavor("Retrieving logs for app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{
			"AppName":   cmd.RequiredArgs.AppName,
			"OrgName":   cmd.Config.TargetedOrganization().Name,
			"SpaceName": cmd.Config.TargetedSpace().Name,
			"Username":  user.Name,
		})
	cmd.UI.DisplayNewline()

	if cmd.Recent {
		return cmd.displayRecentLogs()
	}
	quitNowChannel, err := cmd.Actor.ScheduleTokenRefresh()
	if err != nil {
		return err
	}

	err = cmd.streamLogs()
	quitNowChannel <- true
	return err
}

func (cmd LogsCommand) displayRecentLogs() error {
	messages, warnings, err := cmd.Actor.GetRecentLogsForApplicationByNameAndSpace(
		cmd.RequiredArgs.AppName,
		cmd.Config.TargetedSpace().GUID,
		cmd.LogCacheClient,
	)

	for _, message := range messages {
		cmd.UI.DisplayLogMessage(loggingaction.LogMessage{
			Message:        message.Message(),
			MessageType:    message.Type(),
			Timestamp:      message.Timestamp(),
			SourceType:     message.SourceType(),
			SourceInstance: message.SourceInstance(),
		}, true)
	}

	cmd.UI.DisplayWarnings(warnings)
	return err
}

func (cmd LogsCommand) streamLogs() error {
	messages, logErrs, cancelFunc, warnings, err := cmd.Actor.GetStreamingLogsForApplicationByNameAndSpace(
		cmd.RequiredArgs.AppName,
		cmd.Config.TargetedSpace().GUID,
		cmd.LogCacheClient,
	)

	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return err
	}

	var messagesClosed, errLogsClosed bool
	for {
		select {
		case message, ok := <-messages:
			if !ok {
				messagesClosed = true
				break
			}

				cmd.UI.DisplayLogMessage(loggingaction.LogMessage{
					Message:        message.Message(),
					MessageType:    message.Type(),
					Timestamp:      message.Timestamp(),
					SourceType:     message.SourceType(),
					SourceInstance: message.SourceInstance(),
				}, true)
		case logErr, ok := <-logErrs:
			if !ok {
				errLogsClosed = true
				break
			}
			cancelFunc()
			return logErr
		}

		if messagesClosed && errLogsClosed {
			break
		}
	}

	return nil
}
