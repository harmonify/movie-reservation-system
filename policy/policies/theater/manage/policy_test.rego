package policies.theater.manage_test

import rego.v1

test_allow_admin if {
	data.policies.theater.manage.allow with input as allow_admin_test_data
}

test_deny_user if {
	not data.policies.theater.manage.allow with input as deny_user_test_data
}

deny_user_test_data := {"subject": {
	"first_name": "alice",
	"roles": ["user"],
}}

allow_admin_test_data := {"subject": {
	"first_name": "Alice",
	"roles": ["admin"],
}}
