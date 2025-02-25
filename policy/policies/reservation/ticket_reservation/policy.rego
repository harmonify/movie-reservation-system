package policies.reservation.ticket_reservation

import rego.v1

default allow := false

allow if {
	input.user.is_email_verified
	input.user.is_phone_number_verified
}
