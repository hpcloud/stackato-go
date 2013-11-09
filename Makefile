fmt:
	gofmt -w .

i:	fmt
	go install -v stackato/...