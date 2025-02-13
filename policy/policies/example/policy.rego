package example

import rego.v1

default allow := false

allow if {
	input.method == "GET"
	input.path == ["salary", input.subject.user]
}

allow if is_admin

is_admin if "admin" in input.subject.groups
