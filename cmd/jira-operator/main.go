// Copyright 2018 Jira Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"os"
	"runtime"

	stub "github.com/jmckind/jira-operator/pkg/stub"
	"github.com/jmckind/jira-operator/version"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	log "github.com/sirupsen/logrus"
)

func main() {
	configureLogging()
	printVersion()

	resource := "jira.atlassian.com/v1alpha1"
	kind := "Jira"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Fatalf("Failed to get watch namespace: %v", err)
	}
	resyncPeriod := 5
	log.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(stub.NewJiraHandler())
	sdk.Run(context.TODO())
}

func printVersion() {
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Infof("operator-sdk Version: %v", sdkVersion.Version)
	log.Infof("jira-operator Version: %v", version.Version)
}

func configureLogging() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.SetOutput(os.Stdout)

	// Allow the log level to be set using an environment variable
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}
