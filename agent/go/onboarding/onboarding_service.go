package main

import (
	logger "agentless/infra/log"
	"net/http"

	connector "agentless/infra/model/common"
	operations "agentless/infra/model/operations"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type OnboardingService struct {
	connectorID string
	accountID   string
}

func NewOnboardingService(connectorID string) *OnboardingService {
	accountId, err := getCallerAccount()
	if err != nil {
		panic("Failed to get caller account")
	}

	return &OnboardingService{
		connectorID: connectorID,
		accountID:   accountId,
	}
}

func (s *OnboardingService) GetV1OperationsInternalConfigMetadata(c *gin.Context) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) PostV1OperationsInternalInternalConfig(c *gin.Context) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) GetV1OperationsHealth(c *gin.Context) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) GetV1OperationsMetrics(c *gin.Context) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) GetV1OperationsEnvInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connector-id":   s.connectorID,
		"connector-type": connector.CloudAws,
		"info": gin.H{
			"general": nil,
		},
	})
}

func (s *OnboardingService) GetV1OperationsEnvUnitsList(c *gin.Context, params operations.GetV1OperationsEnvUnitsListParams) {
	c.JSON(http.StatusOK, gin.H{
		"id":     s.accountID,
		"name":   s.accountID,
		"parent": "None",
	})
}

func (s *OnboardingService) GetV1OperationsLogDownload(c *gin.Context, params operations.GetV1OperationsLogDownloadParams) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) PostV1OperationsLogStart(c *gin.Context) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) GetV1OperationsLogStatus(c *gin.Context, params operations.GetV1OperationsLogStatusParams) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) GetV1OperationsLogStop(c *gin.Context, params operations.GetV1OperationsLogStopParams) {
	logger.Log.Warningf("API Not implemented: %s", c.Request.URL.String())
	c.JSON(http.StatusOK, gin.H{})
}

func (s *OnboardingService) PostV1OperationsOnboarding(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connector-type": connector.CloudAws,
		"component-id":   s.connectorID,
		"general": gin.H{
			"operation-mode": "visibility",
			"auto-discovery": false,
		},
		"steps": []gin.H{
			{
				"name": "AwsOnboarding",
				"action": gin.H{
					"accounts": gin.H{
						"scope":  "accounts",
						"values": []string{s.accountID},
					},
				},
			},
		},
	})
}

func (s *OnboardingService) PostV1OperationsStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"detailed-status": []gin.H{
			{
				s.accountID: operations.Completed,
			},
		},
		"overall-status": operations.Completed,
		"steps":          nil,
	})
}

func (s *OnboardingService) PostVersionHandshake(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":          "1.0.0",
		"contract-version": "1.0.0",
		"component-id":     "1.0.0",
	})
}

func getCallerAccount() (string, error) {
	sess, err := session.NewSession(&aws.Config{
		//Region: aws.String("us-west-2"),
	})
	if err != nil {
		logger.Log.Fatalf("failed to create session, %v", err)
		return "", err
	}

	// Get account ID using STS
	stsSvc := sts.New(sess)
	stsResult, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		logger.Log.Fatalf("failed to get caller identity, %v", err)
		return "", err
	}

	return *stsResult.Account, nil
}
