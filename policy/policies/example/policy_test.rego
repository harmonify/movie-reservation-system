package example_test

import data.example
import rego.v1

test_allow_user if {
	example.allow with input as allow_user_test_data
}

test_allow_admin if {
	example.allow with input as allow_admin_test_data
}

test_deny_user if {
	not example.allow with input as deny_user_test_data
}

allow_user_test_data := {
	"method": "GET",
	"path": ["salary", "alice"],
	"subject": {
		"user": "alice",
		"groups": ["user"],
	},
}

allow_admin_test_data := {
	"method": "GET",
	"path": ["salary", "bob"],
	"subject": {
		"user": "alice",
		"groups": ["admin"],
	},
}

deny_user_test_data := {
	"method": "GET",
	"path": ["salary", "bob"],
	"subject": {
		"user": "alice",
		"groups": ["user"],
	},
}
