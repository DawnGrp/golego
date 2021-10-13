package auth

default allow = false

allow {
	input.age > 18
	input.method == "post"
	input.path == "doc"
}

hidFiled = x {
	x = [1, 2, 3]
}

# https://play.openpolicyagent.org/p/hk4TzZxXZG
# https://play.openpolicyagent.org/p/nT79bxm098
