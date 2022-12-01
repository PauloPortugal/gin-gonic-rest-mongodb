MAKEFLAGS += -j2

audit:
	go list -m all | nancy sleuth
