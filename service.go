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
	"time"
)

func StartService(cmd, prefix string, advInt, cmd_to time.Duration) {
	log.Println("Health check command:", cmd)
	log.Println("Prefix:", prefix)
	succ,fail,status := 0,0,Unknown

	healthy_rm := newHealthyRipMsg(prefix)
	unhealthy_rm := newUnhealthyRipMsg(prefix)

	for {
		if runCheck(cmd, cmd_to) {
			succ += 1
			fail = 0
		} else {
			succ = 0
			fail += 1
		}

		if succ >= 3 {
			log.Println("3 consecutive successes. Healthy!")
			status = Healthy
		}

		if fail >= 3 {
			log.Println("3 consecutive failures. Unhealthy. :(")
			status = Unhealthy
		}

		switch status {
		case Healthy:
			sendRipMsg(healthy_rm)
		case Unhealthy:
			sendRipMsg(unhealthy_rm)
		case Unknown:
		default:
			log.Fatalln("Unknown status value:", status)
		}

		time.Sleep(advInt)
	}
}
