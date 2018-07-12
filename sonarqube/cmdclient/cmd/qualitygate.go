package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alauda/bergamot/sonarqube"
	"github.com/alauda/bergamot/sonarqube/cmdclient/utils"
	"github.com/spf13/cobra"
)

var qualitygateCmd = &cobra.Command{
	Use:   "qualitygate",
	Short: "Sonarqube qualitygate sub command",
}

var qualitygateSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select sonarqube qualitygate",
	Run:   selectQualitygate,
}

var (
	gateID             int
	defaultProjectName string
	propertiesPath     string
)

func init() {
	qualitygateCmd.AddCommand(qualitygateSelectCmd)
	RootCmd.AddCommand(qualitygateCmd)

	qualitygateSelectCmd.Flags().IntVar(&gateID, "gate-id", 0, "gate id for qualitygate")
	qualitygateSelectCmd.Flags().StringVar(&propertiesPath, "properties", "./sonar-project.properties", "project settings file path")
	qualitygateSelectCmd.Flags().StringVar(&defaultProjectName, "name", "", "sonarqube project name")
}

func selectQualitygate(cmd *cobra.Command, args []string) {
	sonar, err := sonarqube.NewSonarQubeArgs(sonarHost, sonarToken)
	if err != nil {
		log.Printf("init sonarqube error: %v", err)
		os.Exit(ExitCodeProgramError)
	}

	if gateID == 0 {
		sonar.Logger.Errorf("select qualitygate need gate id")
		os.Exit(ExitCodeProgramError)
	}

	// read config content
	configFilePath := filepath.Join(workDir, propertiesPath)
	configMap, err := utils.ParseConfigContent(configFilePath)
	if err != nil {
		sonar.Logger.Errorf("read config file %s error: %v", configFilePath, err)
		os.Exit(ExitCodeProgramError)
	}

	var (
		projectKey  string
		projectName string
		ok          bool
	)
	projectKey, ok = configMap["sonar.projectKey"]
	if !ok {
		sonar.Logger.Errorf("need sonar.projectKey in file %s", configFilePath)
		os.Exit(ExitCodeNeedMoreInfo)
	}
	projectName, ok = configMap["sonar.projectName"]
	if !ok {
		sonar.Logger.Debugf("config file %s does not have sonar.projectName, will use pipeline name", configFilePath)
		projectName = defaultProjectName
	}
	if projectName == "" {
		sonar.Logger.Errorf("neigher config file %s nor command args have sonar.projectName", configFilePath)
		os.Exit(ExitCodeNeedMoreInfo)
	}

	// create project
	if err := sonar.CreateProject(projectName, projectKey); err != nil {
		sonar.Logger.Errorf("%v", err)
		os.Exit(ExitCodeProgramError)
	}

	if _, err := sonar.SelectQualityGates(gateID, "", projectKey); err != nil {
		sonar.Logger.Errorf("%v", err)
		os.Exit(ExitCodeProgramError)
	}

	sonar.Logger.Infof("set quality gate %d for %s success", gateID, projectKey)
}
