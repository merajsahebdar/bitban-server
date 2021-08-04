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
	"io"
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// init
func init() {
	//
	// Search for the environment.

	env := os.Getenv("APP_ENV")

	switch env {
	case "production":
		CurrentEnv = Prod
	case "development":
		CurrentEnv = Dev
	default:
		CurrentEnv = Test
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

	var content []byte

	if f, err := GetEtcPath("/cog.yaml"); err != nil {
		if pak, err := pkger.Open("/configs/cog.yml"); err != nil {
			Log.Fatal(err.Error())
		} else {
			defer pak.Close()

			var c buffer.Buffer
			io.Copy(&c, pak)
			content = c.Bytes()
		}
	} else {
		if c, err := ioutil.ReadFile(f); err != nil {
			Log.Fatal(err.Error())
		} else {
			content = c
		}
	}

	// Replace environment values in config content.
	content = []byte(os.ExpandEnv(string(content)))

	if err := yaml.Unmarshal(content, &Cog); err != nil {
		Log.Fatal(err.Error())
	}

	// Provide default values for non-exsisting values.
	defaults.Set(&Cog)
}
