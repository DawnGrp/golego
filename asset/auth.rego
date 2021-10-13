package auth

default allow = false

allow {
	input.age > 18
	data.name = "zeta"
}

hidFiled = x {
	x = [1, 2, 3]
}

# input.params.age > 18
# input.params.method == "post"
# input.params.path == "doc"
# data.name = "zeta" 
