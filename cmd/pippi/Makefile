all: pippi

pippi:
	cd frontend && npm install
	wails build

.SHELL:
run:
	wails serve &
	sleep 5; cd frontend && npx ng serve --poll=2000

.PHONY: all pippi clean

clean:
	$(RM) -r -v frontend/dist
	$(RM) -r -v frontend/node_modules
