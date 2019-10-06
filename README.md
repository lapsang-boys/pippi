
[![Pippi & Lilla Gubben](https://github.com/lapsang-boys/pippi/blob/gh-pages/inc/img/pippi.png)](https://github.com/lapsang-boys/pippi/blob/gh-pages/inc/img/pippi.png)

# Pippi

[![CircleCI](https://circleci.com/gh/lapsang-boys/pippi/tree/master.svg?style=svg)](https://circleci.com/gh/lapsang-boys/pippi/tree/master)

---

An exploratory and modular reverse engineering platform.

## Build

```bash
make
```

## Run

### Back-end

Install `forego` dependency (or run the commands listed in the [Procfile](Procfile)):
```bash
go get -u github.com/ddollar/forego
```

Terminal 1:
```bash
make run_backend
```

### Front-end

Terminal 2:
```bash
make run_frontend
```

## Ports

* `upload`:  1100
* `bin`:     1200
* `disasm`:  1300
   - `disasm-objdump`:  1310
* `strings`: 1400
