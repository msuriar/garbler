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


type Announcer interface {
	Announce(chan health) ()
}

type RIPAnnouncer struct {
	adv_interval time.Duration
	prefix string
	healthy_metric uint32
	unhealthy_metric uint32
}

func (a *RIPAnnouncer) Announce(ch chan health) (){
	go func() {
		healthy_rm := newHealthyRipMsg(a.prefix)
		unhealthy_rm := newUnhealthyRipMsg(a.prefix)
		ticker := time.Tick(a.adv_interval)
		status := <- ch
		for {
			select {
			case status = <- ch:
				log.Println("Received update.")
			case <- ticker:
				switch status {
				case Healthy:
					sendRipMsg(healthy_rm)
					log.Println("Sent healthy update.")
				case Unhealthy:
					sendRipMsg(unhealthy_rm)
					log.Println("Sent unhealthy update.")
				}
			}
		}
	}()
}
