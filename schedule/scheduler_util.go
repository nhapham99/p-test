package scheduler

import (
	"payment-module/internals/constants"
	"payment-module/internals/models"
	"payment-module/logger"
	utils "payment-module/utils"
	"time"

	"github.com/go-co-op/gocron"
)

const MAX_CONCCURENCY int = 8

var healthCheckScheduler *gocron.Scheduler
var versionScheduler *gocron.Scheduler
var cronJobScheduler *gocron.Scheduler

// convert types take an int and return a string value.
type handleVersionCheck func()

func initHealthCheckScheduler() {
	if healthCheckScheduler != nil {
		return
	}
	healthCheckScheduler = gocron.NewScheduler(time.UTC)
	healthCheckScheduler.SetMaxConcurrentJobs(MAX_CONCCURENCY, gocron.RescheduleMode)
}

func initVersionScheduler() {
	if versionScheduler != nil {
		return
	}
	versionScheduler = gocron.NewScheduler(time.UTC)
	versionScheduler.SetMaxConcurrentJobs(MAX_CONCCURENCY, gocron.RescheduleMode)
}

func initCronJobScheduler() {
	if cronJobScheduler != nil {
		return
	}
	cronJobScheduler = gocron.NewScheduler(time.UTC)
	cronJobScheduler.SetMaxConcurrentJobs(MAX_CONCCURENCY, gocron.RescheduleMode)
}

func NewCheckVersionScheduler(interval int, fn handleVersionCheck) {
	initVersionScheduler()
	_, err := versionScheduler.Every(interval).Seconds().Tag("version-scheduler").Do(func() {
		fn()
	})
	if err != nil {
		logger.Infof("Check version scheduler error:%s", err.Error())
	}

	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	versionScheduler.StartAsync()
}

func NewHealthCheckScheduler(systemCode string, service models.Service, telegram models.Telegram) {
	if len(systemCode) <= 0 || len(service.Code) <= 0 || service.Interval <= 0 {
		return
	}
	logger.Infof("NewHealthCheckScheduler system code:%s service:%s", systemCode, service)
	initHealthCheckScheduler()
	tag := systemCode + "_" + service.Code
	_, err := healthCheckScheduler.Every(service.Interval).Seconds().Tag(tag).Do(func() {
		//do something...
		result, err := utils.RunHealthCheckRequestGet(service.URL)
		if err != nil {
			logger.Warnf("Request health-check error:%s", err.Error())
			if service.WarningType == constants.WARNING_TYPE_TELEGRAM {
				utils.SendMessageTeleFail(telegram.Token, telegram.ChatId, service.Code, service.Name, service.MessageError, err.ErrorMessage())
			}
		} else {
			logger.Debugf("Request health-check result:%t", result)
			if service.NotifySuccess == 1 {
				if service.WarningType == constants.WARNING_TYPE_TELEGRAM {
					utils.SendMessageTeleSuccess(telegram.Token, telegram.ChatId, service.Code, service.Name)
				}
			}
		}

	})
	if err != nil {
		logger.Infof("Health-Check scheduler error:%s", err.Error())
	}

	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	healthCheckScheduler.StartAsync()
}

func GetAllJobVersionCheck() {
	if versionScheduler == nil {
		return
	}
	jobs := versionScheduler.Jobs()
	if jobs == nil {
		logger.Infof("Version-Check don't have job")
	} else {
		for i := 0; i < len(jobs); i++ {
			logger.Infof("Version-Check job with tag:%s", jobs[i].Tags())
		}
	}
}

func ClearAllVersionCheckScheduler() {
	if versionScheduler == nil {
		return
	}
	versionScheduler.Clear()
}

func GetAllJobHealthCheck() {
	if healthCheckScheduler == nil {
		return
	}
	jobs := healthCheckScheduler.Jobs()
	if jobs == nil {
		logger.Infof("Health-Check don't have job")
	} else {
		for i := 0; i < len(jobs); i++ {
			logger.Infof("Health-Check job with tag:%s", jobs[i].Tags())
		}
	}
}

func ClearAllHealthCheckScheduler() {
	if healthCheckScheduler == nil {
		return
	}
	healthCheckScheduler.Clear()
}

func NewCronJob() {
	initCronJobScheduler()
	// cron expressions supported
	//s.Cron("*/1 * * * *").Do(task) // every minute
}
