package user

default allow = false

allow {
	input.params.age > 18
	input.params.method == "post"
	input.params.path == "doc"
}
