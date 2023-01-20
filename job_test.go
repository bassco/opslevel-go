package opslevel_test

import (
	"testing"

	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2022"
)

func TestRunnerRegister(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/register")
	// Act
	result, err := client.RunnerRegister()
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, ol.ID("1234"), result.Id)
}

func TestRunnerGetScale(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/scale")
	// Act
	result, err := client.RunnerScale("1234567890", 2, 3)
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 6, result.RecommendedReplicaCount)
}

func TestRunnerGetPendingJobs(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/get_pending")
	// Act
	result, token, err := client.RunnerGetPendingJob("1234567890", "1234")
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "public.ecr.aws/opslevel/cli:v2022.02.25", result.Image)
	autopilot.Equals(t, "ls -al", result.Commands[1])
	autopilot.Equals(t, ol.ID("12344321"), token)
}

func TestRunnerAppendJobLog(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/append_log")
	// Act
	err := client.RunnerAppendJobLog(ol.RunnerAppendJobLogInput{
		RunnerId:    "1234",
		RunnerJobId: "5678",
		SentAt:      ol.NewISO8601Date("2022-07-01T01:00:00.000Z"),
		Logs:        []string{"Log1", "Log2"},
	})
	// Assert
	autopilot.Ok(t, err)
}

func TestRunnerReportJobOutcome(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/report_outcome")
	// Act
	err := client.RunnerReportJobOutcome(ol.RunnerReportJobOutcomeInput{
		RunnerId:    "1234567890",
		RunnerJobId: "Z2lkOi8vb3BzbGV2ZWwvUnVubmVyczo6Sm9iUnVuLzE",
		Outcome:     ol.RunnerJobOutcomeEnumExecutionTimeout,
	})
	// Assert
	autopilot.Ok(t, err)
}

func TestRunnerUnregister(t *testing.T) {
	// Arrange
	client := ATestClient(t, "job/unregister")
	// Act
	err := client.RunnerUnregister("1234")
	// Assert
	autopilot.Ok(t, err)
}

func TestRunnerJobNumber(t *testing.T) {
	// Arrange
	job := ol.RunnerJob{
		Id: "Z2lkOi8vb3BzbGV2ZWwvUnVubmVyczo6Sm9iUnVuLzIyNQ",
	}
	// Act
	jobNumber := job.Number()
	// Assert
	autopilot.Equals(t, "225", jobNumber)
}

func TestRunnerJobNumberFailure(t *testing.T) {
	// Arrange
	job := ol.RunnerJob{
		Id: "Z2lkOi8vb3BzbGV2ZWwvUnVubmVyczo6Sm9iU",
	}
	// Act
	jobNumber := job.Number()
	// Assert
	autopilot.Equals(t, "Z2lkOi8vb3BzbGV2ZWwvUnVubmVyczo6Sm9iU", jobNumber)
}
