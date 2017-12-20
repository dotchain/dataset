all:
	mkdir -p json/compact
	go run tools/gensplices/gen.go > json/compact/splices.json

