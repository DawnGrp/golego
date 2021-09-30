package user

default allow = false

allow {
	input.age > 18
}

allow {
	data.name = "zeta"
}

user_role {
	input.age > 10
}

user_role {
	data.name = "zeta1"
}
