package policies.theater.manage

import rego.v1

default allow := false

allow if {
	"admin" in input.subject.roles
}
