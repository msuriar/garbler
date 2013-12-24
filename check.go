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

package garbler

import (
	"log"
	"os/exec"
	"time"
)

type health byte
const (
	Healthy health = iota
	Unhealthy
	Unknown
)

func runCheck(c string, cmd_to time.Duration) (result bool) {
	cmd := exec.Command("sh", "-c", c)
	done := make(chan error, 1)
	timeout := time.After(cmd_to)

	go func() {
		done <- cmd.Run()
	}()

	select {
	case <- timeout:
		log.Println("Command timed out.")
		if err := cmd.Process.Kill(); err != nil {
			log.Println("Failed to kill command after timeout.")
		} else {
			log.Println("Command killed successfully after timeout.")
		}
		<- done
		return false
	case err:= <- done:
		result := err == nil
		log.Println("Command completed. Succeeded: ", result)
		return result
	}
}
