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
	succ,fail,healthy := 0,0,false

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
			log.Println("3 consequetive successes. Healthy!")
			healthy = true
		}

		if fail >= 3 {
			log.Println("3 consequetive failues. Unhealthy. :(")
			healthy = false
		}

		if healthy {
			sendRipMsg(healthy_rm)
		} else {
			sendRipMsg(unhealthy_rm)
		}

		time.Sleep(advInt)
	}
}
