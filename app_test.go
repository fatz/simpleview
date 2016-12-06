package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {

	assert := assert.New(t)

	var simpleview Simpleview
	var jsonString = `{"foo":"bar","baz":[{"name": "boo"}]}`

	object, err := simpleview.readJSON(jsonString)
	assert.Nil(err, "Should generate no error")

	str, err := object.String("foo")
	assert.Equal(str, "bar", "foo should equal bar")
	assert.Nil(err, "Should generate no error")

	slice, err := object.Array("baz")

	assert.Nil(err, "Should generate no error")

	assert.Len(slice, 1, "Should have one entity")
}

func TestHost(t *testing.T) {

	assert := assert.New(t)

	var simpleview Simpleview
	var jsonString = `{
"results": [{
"attrs": {
"__name": "api01.application.example.com",
"acknowledgement": 0,
"acknowledgement_expiry": 0,
"action_url": "",
"active": true,
"address": "10.11.12.13",
"address6": "",
"check_attempt": 1,
"check_command": "hostalive",
"check_interval": 60,
"check_period": "",
"check_timeout": null,
"command_endpoint": "",
"display_name": "api01.application.example.com",
"enable_active_checks": true,
"enable_event_handler": true,
"enable_flapping": false,
"enable_notifications": true,
"enable_passive_checks": true,
"enable_perfdata": true,
"event_command": "",
"flapping": false,
"flapping_last_change": 1470408695.671307,
"flapping_negative": 2295,
"flapping_positive": 0,
"flapping_threshold": 30,
"force_next_check": false,
"force_next_notification": false,
"groups": [
"datacenter-nl",
"linux"
],
"ha_mode": 0,
"icon_image": "",
"icon_image_alt": "",
"last_check": 1470408695.663835,
"last_check_result": {
"active": true,
"check_source": "node01.icinga2.example.com",
"command": [
"/usr/lib64/nagios/plugins/check_ping",
"-H",
"10.11.12.13",
"-c",
"5000,100%",
"-w",
"3000,80%"
],
"execution_end": 1470408695.663742,
"execution_start": 1470408691.650308,
"exit_status": 0,
"output": "PING OK - Packet loss = 0%, RTA = 0.49 ms",
"performance_data": [
"rta=0.491000ms;3000.000000;5000.000000;0.000000",
"pl=0%;80;100;0"
],
"schedule_end": 1470408695.663835,
"schedule_start": 1470408691.65,
"state": 0,
"type": "CheckResult",
"vars_after": {
"attempt": 1,
"reachable": true,
"state": 0,
"state_type": 1
},
"vars_before": {
"attempt": 1,
"reachable": true,
"state": 0,
"state_type": 1
}
},
"last_hard_state": 0,
"last_hard_state_change": 1470055092.451354,
"last_in_downtime": false,
"last_reachable": true,
"last_state": 0,
"last_state_change": 1470055092.451354,
"last_state_down": 0,
"last_state_type": 1,
"last_state_unreachable": 0,
"last_state_up": 1470408695.671302,
"max_check_attempts": 3,
"name": "api01.application.example.com",
"next_check": 1470408751.6499999,
"notes": "",
"notes_url": "",
"original_attributes": null,
"package": "_cluster",
"paused": true,
"retry_interval": 30,
"state": 0,
"state_type": 1,
"templates": [
"api01.application.example.com",
"generic-host"
],
"type": "Host",
"vars": {
"application": "email-marketing",
"data_filesystem": "50 GB",
"dc": "NULL",
"environment": "production",
"manufacturer": "NULL",
"nodetype": "api",
"notification": {
"mail": {
"groups": [
"linux"
]
}
},
"os": "Linux",
"platform": "centos",
"platform_family": "rhel",
"platform_version": "6",
"rack": "NULL",
"row": "NULL",
"rtid": "1373",
"type": "VM"
},
"version": 0,
"volatile": false,
"zone": "nl"
},
"joins": {},
"meta": {},
"name": "api01.application.example.com",
"type": "Host"
}]}`

	object, err := simpleview.readJSON(jsonString)
	assert.Nil(err, "Should generate no error")

	nobject, err := object.ArrayOfObjects("results")
	assert.Nil(err, "Should generate no error")

	assert.Len(nobject, 1, "should have one entry")

	assert.Equal(nobject[0]["name"], "api01.application.example.com", "has name")

	jsonstr, err := json.Marshal(nobject)
	assert.Nil(err, "Should generate no error")
	assert.Equal(strings.Contains(string(jsonstr), "command"), true, "should include command")
}
