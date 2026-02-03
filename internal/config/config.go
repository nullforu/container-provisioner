package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	HTTPAddr        string
	ShutdownTimeout time.Duration

	Logging LoggingConfig
}

type LoggingConfig struct {
	Dir               string
	FilePrefix        string
	DiscordWebhookURL string
	SlackWebhookURL   string
	MaxBodyBytes      int
	WebhookQueueSize  int
	WebhookTimeout    time.Duration
	WebhookBatchSize  int
	WebhookBatchWait  time.Duration
	WebhookMaxChars   int
}

func Load() (Config, error) {
	var errs []error

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, fmt.Errorf("load .env: %w", err))
	}

	appEnv := getEnv("APP_ENV", "local")
	httpAddr := getEnv("HTTP_ADDR", ":8081")
	shutdownTimeout, err := getDuration("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		errs = append(errs, err)
	}

	logDir := getEnv("LOG_DIR", "logs")
	logPrefix := getEnv("LOG_FILE_PREFIX", "app")
	logMaxBodyBytes, err := getEnvInt("LOG_MAX_BODY_BYTES", 1024*1024)
	if err != nil {
		errs = append(errs, err)
	}

	logWebhookQueueSize, err := getEnvInt("LOG_WEBHOOK_QUEUE_SIZE", 1000)
	if err != nil {
		errs = append(errs, err)
	}

	logWebhookTimeout, err := getDuration("LOG_WEBHOOK_TIMEOUT", 5*time.Second)
	if err != nil {
		errs = append(errs, err)
	}

	logWebhookBatchSize, err := getEnvInt("LOG_WEBHOOK_BATCH_SIZE", 20)
	if err != nil {
		errs = append(errs, err)
	}

	logWebhookBatchWait, err := getDuration("LOG_WEBHOOK_BATCH_WAIT", 2*time.Second)
	if err != nil {
		errs = append(errs, err)
	}

	logWebhookMaxChars, err := getEnvInt("LOG_WEBHOOK_MAX_CHARS", 1800)
	if err != nil {
		errs = append(errs, err)
	}

	cfg := Config{
		AppEnv:          appEnv,
		HTTPAddr:        httpAddr,
		ShutdownTimeout: shutdownTimeout,
		Logging: LoggingConfig{
			Dir:               logDir,
			FilePrefix:        logPrefix,
			DiscordWebhookURL: getEnv("LOG_DISCORD_WEBHOOK_URL", ""),
			SlackWebhookURL:   getEnv("LOG_SLACK_WEBHOOK_URL", ""),
			MaxBodyBytes:      logMaxBodyBytes,
			WebhookQueueSize:  logWebhookQueueSize,
			WebhookTimeout:    logWebhookTimeout,
			WebhookBatchSize:  logWebhookBatchSize,
			WebhookBatchWait:  logWebhookBatchWait,
			WebhookMaxChars:   logWebhookMaxChars,
		},
	}

	if err := validateConfig(cfg); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return Config{}, errors.Join(errs...)
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}

func getEnvInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def, fmt.Errorf("%s must be an integer", key)
	}

	return n, nil
}

func getEnvBool(key string, def bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return def, fmt.Errorf("%s must be a boolean", key)
	}

	return b, nil
}

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return def, fmt.Errorf("%s must be a duration", key)
	}

	return d, nil
}

func validateConfig(cfg Config) error {
	var errs []error

	if cfg.HTTPAddr == "" {
		errs = append(errs, errors.New("HTTP_ADDR must not be empty"))
	}

	if cfg.Logging.Dir == "" {
		errs = append(errs, errors.New("LOG_DIR must not be empty"))
	}

	if cfg.Logging.FilePrefix == "" {
		errs = append(errs, errors.New("LOG_FILE_PREFIX must not be empty"))
	}

	if cfg.Logging.MaxBodyBytes <= 0 {
		errs = append(errs, errors.New("LOG_MAX_BODY_BYTES must be positive"))
	}

	if cfg.Logging.WebhookQueueSize <= 0 {
		errs = append(errs, errors.New("LOG_WEBHOOK_QUEUE_SIZE must be positive"))
	}

	if cfg.Logging.WebhookTimeout <= 0 {
		errs = append(errs, errors.New("LOG_WEBHOOK_TIMEOUT must be positive"))
	}

	if cfg.Logging.WebhookBatchSize <= 0 {
		errs = append(errs, errors.New("LOG_WEBHOOK_BATCH_SIZE must be positive"))
	}

	if cfg.Logging.WebhookBatchWait <= 0 {
		errs = append(errs, errors.New("LOG_WEBHOOK_BATCH_WAIT must be positive"))
	}

	if cfg.Logging.WebhookMaxChars <= 0 {
		errs = append(errs, errors.New("LOG_WEBHOOK_MAX_CHARS must be positive"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func Redact(cfg Config) Config {
	cfg.Logging.DiscordWebhookURL = redact(cfg.Logging.DiscordWebhookURL)
	cfg.Logging.SlackWebhookURL = redact(cfg.Logging.SlackWebhookURL)

	return cfg
}

func redact(value string) string {
	if value == "" {
		return ""
	}

	const (
		visiblePrefix = 2
		visibleSuffix = 2
	)
	if len(value) <= visiblePrefix+visibleSuffix {
		return "***"
	}

	return value[:visiblePrefix] + "***" + value[len(value)-visibleSuffix:]
}

func FormatForLog(cfg Config) string {
	cfg = Redact(cfg)
	var b strings.Builder
	fmt.Fprintf(&b, "AppEnv=%s\n", cfg.AppEnv)
	fmt.Fprintf(&b, "HTTPAddr=%s\n", cfg.HTTPAddr)
	fmt.Fprintf(&b, "ShutdownTimeout=%s\n", cfg.ShutdownTimeout)
	fmt.Fprintln(&b, "Logging:")
	fmt.Fprintf(&b, "  Dir=%s\n", cfg.Logging.Dir)
	fmt.Fprintf(&b, "  FilePrefix=%s\n", cfg.Logging.FilePrefix)
	fmt.Fprintf(&b, "  DiscordWebhookURL=%s\n", cfg.Logging.DiscordWebhookURL)
	fmt.Fprintf(&b, "  SlackWebhookURL=%s\n", cfg.Logging.SlackWebhookURL)
	fmt.Fprintf(&b, "  MaxBodyBytes=%d\n", cfg.Logging.MaxBodyBytes)
	fmt.Fprintf(&b, "  WebhookQueueSize=%d\n", cfg.Logging.WebhookQueueSize)
	fmt.Fprintf(&b, "  WebhookTimeout=%s\n", cfg.Logging.WebhookTimeout)
	fmt.Fprintf(&b, "  WebhookBatchSize=%d\n", cfg.Logging.WebhookBatchSize)
	fmt.Fprintf(&b, "  WebhookBatchWait=%s\n", cfg.Logging.WebhookBatchWait)
	fmt.Fprintf(&b, "  WebhookMaxChars=%d\n", cfg.Logging.WebhookMaxChars)

	return b.String()
}
