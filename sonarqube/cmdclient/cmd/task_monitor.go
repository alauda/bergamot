package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alauda/bergamot/sonarqube"
	"github.com/alauda/bergamot/sonarqube/cmdclient/utils"
	"github.com/spf13/cobra"
)

var taskMonitorCmd = &cobra.Command{
	Use:   "taskmonitor",
	Short: "Sonarqube analysis task monitor sub command",
	Run:   monitorTask,
}

// task status refer https://docs.sonarqube.org/display/SONAR/Background+Tasks
// You can filter Background Tasks according to their Status: Pending, Success, Failed or Canceled (upper!)
const (
	TaskStatusSuccess  = "SUCCESS"
	TaskStatusPenging  = "PENDING"
	TaskStatusFailed   = "FAILED"
	TaskStatusCanceled = "CANCELED"
	ProjectStatusOK    = "OK"
	TaskStatusWarn     = "WARN"
	TaskStatusError    = "ERROR"
	TaskStatusNone     = "NONE"
)

// exit code
const (
	ExitCodeNormal             int = 0
	ExitCodeProgramError       int = 131
	ExitCodeTimeout            int = 132
	ExitCodeQualityGateError   int = 133
	ExitCodeQualityGateWarn    int = 134
	ExitCodeQualityGateUnknown int = 135
	ExitCodeNeedMoreInfo       int = 136
)

func monitorTask(cmd *cobra.Command, args []string) {
	sonar, err := sonarqube.NewSonarQubeArgs(sonarHost, sonarToken)
	if err != nil {
		log.Printf("init sonarqube error: %v", err)
		os.Exit(ExitCodeProgramError)
	}

	// get and print taskData
	taskData, err := collectTaskData(workDir)
	if err != nil {
		sonar.Logger.Errorf("get task data from %s error: %v", workDir, err)
		os.Exit(ExitCodeProgramError)
	}
	sonar.Logger.Infof("task data are:")
	for k, v := range taskData {
		sonar.Logger.Infof("%s=%s", k, v)
	}

	// judge use taskID or workDir
	if ceTaskID == "" {
		ceTaskID, err = getTaskID(taskData)
		if err != nil {
			sonar.Logger.Errorf("get task id error: %v", err)
			os.Exit(ExitCodeProgramError)
		}
	}

	errCh := make(chan error)
	successCh := make(chan struct{})

	go func(successCh chan struct{}, errCh chan error, sonar *sonarqube.SonarQube, taskID string) {
		stop := false
		for !stop {
			status, err := sonar.GetAnalysisTaskStatus(taskID)
			if err != nil {
				errCh <- err
				break
			}

			switch status {
			case TaskStatusSuccess:
				successCh <- struct{}{}
				stop = true
			case TaskStatusFailed:
				sonar.Logger.Errorf("task %s is failed", taskID)
				errCh <- fmt.Errorf("task %s is failed", taskID)
				stop = true
			case TaskStatusCanceled:
				sonar.Logger.Errorf("task %s is canceled", taskID)
				successCh <- struct{}{}
				stop = true
			default:
				sonar.Logger.Debugf("task %s status is %s, will sleep %d seconds and try again",
					taskID, status, monitorInterval/time.Second)
				time.Sleep(monitorInterval)
			}
		}
	}(successCh, errCh, sonar, ceTaskID)

	select {
	case <-time.After(monitorTimeout):
		sonar.Logger.Errorf("monitor task %s timeout, exit\n", ceTaskID)
		os.Exit(ExitCodeTimeout)
	case err := <-errCh:
		sonar.Logger.Errorf("monitor task %s error: %v", ceTaskID, err)
		os.Exit(ExitCodeProgramError)
	case <-successCh:
		sonar.Logger.Infof("task of id %s is success, begin checking project analysis status...", ceTaskID)
		projectKey, ok := taskData["projectKey"]
		if !ok {
			sonar.Logger.Errorf("task data does not have projectKey")
			os.Exit(ExitCodeProgramError)
		}
		projectStatus, err := sonar.GetQualityGatesProjectStatus("", "", projectKey)
		if err != nil {
			sonar.Logger.Errorf("get quality gates project status error: %v", err)
			os.Exit(ExitCodeProgramError)
		}

		switch projectStatus {
		case ProjectStatusOK:
			sonar.Logger.Infof("quality gate status of project %s is OK", projectKey)
		case TaskStatusWarn:
			if warn {
				sonar.Logger.Errorf("quality gate status of project %s", projectKey)
				os.Exit(ExitCodeQualityGateWarn)
			}
			sonar.Logger.Infof(`quality gate status of project %s is Warn, treat as success.
if want treat warn as fail, use --warn flag`, projectKey)
		case TaskStatusError, TaskStatusNone:
			sonar.Logger.Errorf("quality gate status of project %s", projectKey)
			os.Exit(ExitCodeQualityGateError)
		default:
			sonar.Logger.Errorf("quality gate status of project %s is unknown %s",
				projectKey, projectStatus)
			os.Exit(ExitCodeQualityGateUnknown)
		}
	}
}

var (
	ceTaskID        string
	warn            bool
	monitorTimeout  time.Duration
	monitorInterval time.Duration
)

func init() {
	taskMonitorCmd.Flags().StringVar(&ceTaskID, "id", "", "sonar analysis task id")
	taskMonitorCmd.Flags().BoolVar(&warn, "warn", false, "whether set analysis status 'warn' to error")
	taskMonitorCmd.Flags().DurationVar(&monitorTimeout, "t", time.Minute*time.Duration(30), "sonar analysis timeout")
	taskMonitorCmd.Flags().DurationVar(&monitorInterval, "i", time.Second*time.Duration(4), "check sonar analysis task status interval")
}

func getTaskID(taskData map[string]string) (string, error) {
	taskID, ok := taskData["ceTaskId"]
	if !ok {
		return "", fmt.Errorf("task data does not have ceTaskId")
	}

	return taskID, nil
}

func collectTaskData(workDir string) (map[string]string, error) {
	taskDataFilePath := filepath.Join(workDir, ".scannerwork/report-task.txt")
	return utils.ParseConfigContent(taskDataFilePath)
}
