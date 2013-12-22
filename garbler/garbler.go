/*
 * Copyright 2013 Murali Suriar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/msuriar/garbler"
	"os"
	"time"
)

var (
	prefix,cmd string
	advInt time.Duration
	successes, failures, cmd_timout int
)

func init() {
	flag.StringVar(&prefix, "prefix", "", "Prefix to advertise.")
	flag.StringVar(&cmd, "cmd", "", "Healthcheck command to run.")
	flag.DurationVar(&advInt, "advInt", 30 * time.Second, `Interval between
	updates.`)
	flag.IntVar(&successes, "successes", 3, `Number of consecutive
	successes before garbler considers the service healthy.`)
	flag.IntVar(&failures, "failures", 3, `Number of consecutive
	failures before garbler considers the service unhealthy.`)
	flag.IntVar(&cmd_timout, "timeout", 30, `Time in seconds before garbler
	considers the healthcheck command to have failed, and kills it.`)
}

func main() {
	flag.Parse()

	err := false

	if prefix == "" {
		fmt.Println("prefix is a required flag.")
		err = true
	}

	if cmd == "" {
		fmt.Println("cmd is a required flag.")
		err = true
	}

	if err { os.Exit(1) }
	
	garbler.StartService(cmd, prefix, advInt)
}
