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

	healthy_rm := newHealthyRipMsg(prefix)
	unhealthy_rm := newUnhealthyRipMsg(prefix)

	c := Command{cmd, advInt, cmd_to, 3, 3}
	p := c.Probe()
	ch := make(chan health, 1)
	go p(ch)

	ticker := time.Tick(advInt)
	status := <-ch
	for {
		select {
		case status = <- ch:
			log.Println("Received update.")
		case <- ticker:
			switch status {
			case Healthy:
				sendRipMsg(healthy_rm)
			case Unhealthy:
				sendRipMsg(unhealthy_rm)
			}
		}

	}
}
