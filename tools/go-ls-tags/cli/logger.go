/*
Copyright 2022 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"knative.dev/test-infra/pkg/logging"
)

func contextWithLogger() context.Context {
	ctx := context.Background()
	fallback := logging.FromContext(ctx)
	return logging.WithLogger(ctx, newLogger(fallback))
}

func newLogger(fallback *zap.SugaredLogger) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(logLevelFromEnv(fallback))
	cfg.OutputPaths = []string{"stderr"}
	log, err := cfg.Build()
	if err != nil {
		fallback.Fatal(err)
	}
	return log.Sugar()
}

func logLevelFromEnv(fallback *zap.SugaredLogger) zapcore.Level {
	defaultLevel := zapcore.WarnLevel.String()
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = defaultLevel
	}
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		fallback.Fatalf("Invalid log level: %q", level)
	}
	return l
}
