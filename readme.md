## growl
A simple cli tool to make it easier to build and run your project, and automatize other tasks, by adding custom commands.  
This is pretty similar to NodeJS package.json commands.

## installation
```sh
go install github.com/checkm4ted/growl
```

## usage
First, make a growl.yaml file in your project. There's a example one in this repo.
```yaml
shell: cmd /C # > optional, default is cmd /C in Windows, bash -c in Linux. Set this to the shell you want to use. (with the prefix to make it work too)
globalenv: # optional, set global (for all commands) environment variables like this.
  - name: TESTGLOBAL
    value: global var test

# commands:
commands:
  - name: b
    description: "build the app to a given os"
    command: growl cross -os %1 # add args with %1, %2, etc.
    
  - name: test
    description: "show environment variables"
    command: echo hello, %TESTGLOBAL% # in cmd you print env vars with %VAR%
    extra: # optional, set extra commands (ran after the main one) for a command like this.
      - echo hello, extra %TEST%
      - echo this is a nicer way to run multiple commands than using "&&"
    env: # optional, set environment variables for a command like this.
      - name: TEST
        value: test
```

- List commands in growl.yaml using `growl l # or growl list`
- Run them with `growl [command] [args]`

- Run/build your project with `growl [b/build/r/run] [args]`
- Cross compile with  
```bash
growl cross --os [os] --arch [arch] --ldflags "[ldflags]" [--static] [--light] [--cgo] 
#or
growl cross -o [os] -a [arch] -ld "[ldflags]" [-s] [-l] [-c] 
```
`--static` adds "-extldflags=-static" to ldflags and `--light` adds `-w -s`.

or `growl cross -o [os] -a [arch] -ld [ldflags]`
- Print availbable os/architectures with `growl cross list`