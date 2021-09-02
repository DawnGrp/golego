package user

default allow = false

allow {
	input.age > 18
	input.method == "post"
	input.path == "doc"
}
