/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cfg

import (
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// init
func init() {
	//
	// Searches for the environment.

	env := os.Getenv("APP_ENV")

	switch env {
	case "production":
		CurrentEnv = Prod
	case "development":
	default:
		CurrentEnv = Dev
	}

	//
	// Initialize Locales

	EnLocale = en.New()
	EnLocaleUni = ut.New(EnLocale, EnLocale)
	EnTrans, _ = EnLocaleUni.GetTranslator("en")

	//
	// Create a logger instance.

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if CurrentEnv == Prod {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		os.Stderr,
		zap.NewAtomicLevel(),
	)

	Log = zap.New(core)

	//
	// Feed `Cog`

	f, err := GetAssetPath("/cog.yaml")
	if err != nil {
		Log.Fatal(err.Error())
	}

	content, err := ioutil.ReadFile(f)
	if err != nil {
		Log.Fatal(err.Error())
	}

	// Replace environment values in config content.
	content = []byte(os.ExpandEnv(string(content)))

	err = yaml.Unmarshal(content, &Cog)
	if err != nil {
		Log.Fatal(err.Error())
	}

	// Provide default values for non-exsisting values.
	defaults.Set(&Cog)
}
