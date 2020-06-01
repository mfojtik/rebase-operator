package operator

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v32/github"
	slackgo "github.com/slack-go/slack"
	"golang.org/x/oauth2"
	"k8s.io/klog"

	"github.com/mfojtik/rebase-operator/pkg/operator/config"
	"github.com/mfojtik/rebase-operator/pkg/operator/fork"
	"github.com/mfojtik/rebase-operator/pkg/slack"
	"github.com/mfojtik/rebase-operator/pkg/slacker"
)

const bugzillaEndpoint = "https://bugzilla.redhat.com"

func Run(ctx context.Context, cfg config.OperatorConfig) error {
	slackClient := slackgo.New(cfg.Credentials.DecodedSlackToken(), slackgo.OptionDebug(true))

	// This slack client is used for debugging
	slackDebugClient := slack.NewChannelClient(slackClient, cfg.SlackAdminChannel, true)

	recorder := slack.NewRecorder(slackDebugClient, "RebaseOperator")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Credentials.GithubAPIKey},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	slackerInstance := slacker.NewSlacker(slackClient, slacker.Options{
		ListenAddress:     "0.0.0.0:3000",
		VerificationToken: cfg.Credentials.DecodedSlackVerificationToken(),
	})

	slackerInstance.Command("say <message>", &slacker.CommandDefinition{
		Description: "Say something.",
		Handler: func(req slacker.Request, w slacker.ResponseWriter) {
			msg := req.StringParam("message", "")
			w.Reply(msg)
		},
	})
	slackerInstance.DefaultCommand(func(req slacker.Request, w slacker.ResponseWriter) {
		w.Reply("Unknown command")
	})

	recorder.Eventf("OperatorStarted", "Rebase Operator Started\n\n```\n%s\n```\n", spew.Sdump(cfg.Anonymize()))

	// report command allow to manually trigger a reporter to run out of its normal schedule
	slackerInstance.Command("rebase prepare-branch <tag> <branch>", &slacker.CommandDefinition{
		Description: "Prepare a branch in openshift/kubernetes based on tag name from kubernetes/kubernetes.",
		Handler: auth(cfg, func(req slacker.Request, w slacker.ResponseWriter) {
			tag := req.StringParam("tag", "")
			branch := req.StringParam("branch", "")
			go func(tagName, branchName string) {
				if err := fork.PrepareBranch(ctx, githubClient, tagName, branchName); err != nil {
					_, _, _, err := w.Client().SendMessage(req.Event().Channel,
						slackgo.MsgOptionPostEphemeral(req.Event().User),
						slackgo.MsgOptionText(fmt.Sprintf(":warning: %s", err), false))
					if err != nil {
						klog.Error(err)
					}
					return
				}
				_, _, _, err := w.Client().SendMessage(req.Event().Channel,
					slackgo.MsgOptionPostEphemeral(req.Event().User),
					slackgo.MsgOptionText(fmt.Sprintf("Branch %q based on %q was created: https://github.com/openshift/kubernetes/tree/%s", branch, tag, branch), false))
				klog.Error(err)
			}(tag, branch)
			_, _, _, err := w.Client().SendMessage(req.Event().Channel,
				slackgo.MsgOptionPostEphemeral(req.Event().User),
				slackgo.MsgOptionText(fmt.Sprintf("Creating branch %q based on %q ...", branch, tag), false))
			if err != nil {
				klog.Error(err)
			}
		}, "group:admins"),
	})

	go slackerInstance.Run(ctx)

	<-ctx.Done()
	return nil
}
