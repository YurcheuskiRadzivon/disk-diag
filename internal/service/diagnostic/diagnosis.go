package diagnostic

import (
	"fmt"
	"math"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

var nvmeStatusDescriptions = map[int]string{
	0x00: "Generic Command Status",
	0x01: "Invalid Command Opcode",
	0x02: "Invalid Field in Command",
	0x03: "Command ID Conflict",
	0x04: "Data Transfer Error",
	0x05: "Aborted Command",
	0x06: "Internal Error",
}

func (s *service) AnalyzeManual(data *models.SmartInfo) (*models.HealthReport, error) {
	stats := data.NVMeSmartHealthInformationLog

	tbWritten := float64(stats.DataUnitsWritten) * 512000 / (1024 * 1024 * 1024 * 1024)

	tempC := stats.Temperature - 273
	if tempC < -100 {
		tempC = stats.Temperature
	}

	lifeRemaining := 100 - stats.PercentageUsed
	if lifeRemaining < 0 {
		lifeRemaining = 0
	}

	report := &models.HealthReport{
		Method:      "Manual Algorithm",
		Status:      "Healthy",
		HealthScore: 100,
		Problems:    []string{},
		Metrics: models.KeyStats{
			Temperature:     tempC,
			LifeRemaining:   lifeRemaining,
			DataWrittenTB:   math.Round(tbWritten*100) / 100,
			PowerOnHours:    stats.PowerOnHours,
			MediaErrors:     stats.MediaErrors,
			UnsafeShutdowns: stats.UnsafeShutdowns,
		},
	}

	totalErrorCount := 0

	if !data.SmartStatus.Passed {
		report.Status = "Critical"
		report.HealthScore = 0
		report.Problems = append(report.Problems, "Критический сбой: Диск сообщает о неисправности (SMART Fail).")
	}

	if stats.CriticalWarning > 0 {
		report.Status = "Warning"
		report.HealthScore -= 30
		report.Problems = append(report.Problems, fmt.Sprintf("Критические предупреждения контроллера (код: %d).", stats.CriticalWarning))
	}

	if stats.MediaErrors > 0 {
		report.Status = "Critical"
		report.HealthScore -= 50
		report.Problems = append(report.Problems, fmt.Sprintf("Обнаружены ошибки целостности данных (Media Errors: %d). Возможна потеря файлов.", stats.MediaErrors))
	}

	if len(data.NVMeErrorInformationLog.Table) > 0 {

		for _, errEntry := range data.NVMeErrorInformationLog.Table {
			totalErrorCount += errEntry.ErrorCount

			statusName := errEntry.StatusField.String
			if statusName == "" {
				if name, ok := nvmeStatusDescriptions[errEntry.StatusField.StatusCode]; ok {
					statusName = name
				} else {
					statusName = fmt.Sprintf("Unknown Error (Code: %d)", errEntry.StatusField.StatusCode)
				}
			}

			errDetail := fmt.Sprintf("Тип ошибки: '%s' (Код: %d). Количество вхождений: %d.",
				statusName,
				errEntry.StatusField.StatusCode,
				errEntry.ErrorCount)

			if errEntry.StatusField.StatusCode == 2 {
				errDetail += fmt.Sprintf(" Проблема, вероятно, связана с драйвером NVMe или некорректными запросами от ОС. Расположение ошибки (ParmErrorLocation): %d.", errEntry.ParmErrorLocation)
			}

			report.Problems = append(report.Problems, errDetail)
		}
		if totalErrorCount > 10000 {
			report.HealthScore -= 20
			if report.Status == "Healthy" {
				report.Status = "Warning"
			}
		} else if totalErrorCount > 1000 {
			report.HealthScore -= 10
			if report.Status == "Healthy" {
				report.Status = "Warning"
			}
		}
	}

	if stats.UnsafeShutdowns > 0 {
		pointsDeduced := 0
		if stats.UnsafeShutdowns > 50 {
			pointsDeduced = 10
			if report.Status == "Healthy" {
				report.Status = "Warning"
			}
		} else {
			pointsDeduced = 2
		}
		report.HealthScore -= pointsDeduced
		report.Problems = append(report.Problems, fmt.Sprintf("Небезопасных отключений питания: %d. Это может повредить метаданные.", stats.UnsafeShutdowns))
	}

	critTemp := data.NVMeCompositeTemperatureThreshold.Critical - 273
	if critTemp <= 0 {
		critTemp = 80
	}

	if report.Metrics.Temperature >= critTemp {
		report.Status = "Critical"
		report.HealthScore -= 30
		report.Problems = append(report.Problems, fmt.Sprintf("Критический перегрев: %d°C.", report.Metrics.Temperature))
	} else if report.Metrics.Temperature >= critTemp-15 {
		report.Status = "Warning"
		report.HealthScore -= 10
		report.Problems = append(report.Problems, "Высокая рабочая температура.")
	}

	if stats.PercentageUsed >= 100 {
		report.Status = "Warning"
		report.HealthScore -= 20
		report.Problems = append(report.Problems, "Ресурс перезаписи исчерпан (100%+).")
	} else if stats.PercentageUsed > 90 {
		report.HealthScore -= 5
		report.Problems = append(report.Problems, "Высокий износ флеш-памяти (>90%).")
	}

	if report.HealthScore < 0 {
		report.HealthScore = 0
	}

	if report.Status == "Healthy" {
		report.Summary = fmt.Sprintf("Диск в отличном состоянии. Оценка здоровья %d/100. Износ всего %d%%.", report.HealthScore, stats.PercentageUsed)
	} else if report.Status == "Warning" {
		report.Summary = fmt.Sprintf("Диск требует внимания (Оценка %d/100). Обнаружены проблемы, связанные с %s, например, ошибки интерфейса (%d шт) или небезопасные отключения.", report.HealthScore, data.ModelName, totalErrorCount)
	} else {
		report.Summary = "КРИТИЧЕСКОЕ СОСТОЯНИЕ. Рекомендуется немедленно сделать резервную копию и заменить диск."
	}

	return report, nil
}
