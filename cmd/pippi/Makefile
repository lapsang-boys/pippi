all: pippi

pippi:
	cd frontend && npm install
	wails build

run:
	wails serve

.PHONY: all pippi clean

clean:
	$(RM) -r -v frontend/dist
	$(RM) -r -v frontend/node_modules
