package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/james-wukong/school-schedule/internal/config"
	infraPostgres "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"github.com/james-wukong/school-schedule/internal/interface/http/dto"
	schedulerUC "github.com/james-wukong/school-schedule/internal/usecase/schedule"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func StartScheduleConsumer(
	ctx context.Context, db *gorm.DB, cfg *config.Config, wLog *zerolog.Logger,
) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  cfg.Kafka.GroupID, // Identifies this worker set
		Topic:    cfg.Kafka.Topic,
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	wLog.Info().Msg("Worker started: Waiting for scheduling tasks...")

	for {
		// 1. Fetch the message from Kafka
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			wLog.Error().Err(err).Msgf("error fetching message: %v", err)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				continue
			case errors.Is(err, context.Canceled):
				return // graceful shutdown
			default:
				wLog.Error().Err(err).Msgf("fatal kafka error: %v", err)
				return
			}
		}

		// 2. Decode the Task
		var req dto.CreateScheduleRequest
		if err := json.Unmarshal(m.Value, &req); err != nil {
			wLog.Error().Err(err).Msgf("failed to unmarshal task: %v", err)
			continue
		}

		// 3. Run Scheduler and save schedules into database
		reqRepo := infraPostgres.NewRequirementRepository(db, wLog)
		roomRepo := infraPostgres.NewRoomRepository(db, wLog)
		tsRepo := infraPostgres.NewTimeslotRepository(db, wLog)
		schdRepo := infraPostgres.NewScheduleRepository(db, wLog)

		schdUC := schedulerUC.NewCreateScheduleUseCase(reqRepo, roomRepo, tsRepo, schdRepo)
		result, err := schdUC.Execute(ctx, req)
		if err != nil {
			wLog.Error().Err(err).Msgf("failed to execute task: %v", err)
			continue
		}

		// 4. Send a request to go-admin to generate export files
		if err := SchedulerAPIRequest(ctx, dto.CreateScheduleResponse{
			SchoolID:        req.SchoolID,
			SemesterID:      req.SemesterID,
			ScheduleVersion: result.ScheduleVersion,
		}, cfg, wLog); err != nil {
			wLog.Error().Err(err).Msgf("failed to make the exporter api call: %v", err)
		}
		// 5. Commit the message
		if err := reader.CommitMessages(ctx, m); err != nil {
			wLog.Error().Err(err).Msgf("failed to commit message: %v", err)
		}
		wLog.Info().Msg("mission complete...")
	}
}

func SchedulerAPIRequest(
	ctx context.Context, req dto.CreateScheduleResponse, cfg *config.Config, wLog *zerolog.Logger,
) error {
	wLog.Info().Msg("sending a request to Go-Admin to generate export files...")
	// Create a context that expires in 10 seconds
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // Always call cancel to release resources
	jsonData, err := json.Marshal(req)
	if err != nil {
		wLog.Error().Err(err).Msgf("failed to marshal response: %v", err)
		return err
	}
	apiURL, err := url.JoinPath(cfg.Scheduler.URI, "export", "reports")
	if err != nil {
		wLog.Error().Err(err).Msgf("failed to join url paths: %v", err)
		return err
	}
	apiReq, err := http.NewRequestWithContext(
		reqCtx,
		http.MethodPost,
		apiURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		wLog.Error().Err(err).Msgf("failed to create http request: %v", err)
		return err
	}
	apiReq.Header.Set("Content-Type", "application/json")

	for attempt := 0; attempt <= cfg.Scheduler.MaxRetries; attempt++ {
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(apiReq)
		// succeed, break loop, and return
		if err == nil && resp.StatusCode < 500 {
			var result dto.ExportAPIResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				wLog.Error().Err(err).Msgf("failed to decode api response: %v", err)
				continue
			}
			if result.Success {
				return nil
			}

		}
		if err != nil {
			wLog.Error().Err(err).Msgf("failed to get response from api: %v", err)
		}
		if attempt < cfg.Scheduler.MaxRetries {
			select {
			// Wait 1 second before retrying
			// Stop immediately if caller cancels request
			case <-time.After(1 * time.Second):
			case <-reqCtx.Done():
				wLog.Error().Err(err).Msgf("request cancelled: %+v", reqCtx.Err())
				return reqCtx.Err()
			}
		}
	}
	wLog.Error().Err(err).Msgf("retry http request: %v times", cfg.Scheduler.MaxRetries)
	return errors.New("max retries failed")
}
