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
commands:
 # build and run commands written here will overwrite default ones so you can set your own :).
  # - name: build
  #   description: "echo test"
  #   command: echo test
  - name: hello
    description: "prints hello, World!" # this description is unused, this is only for documentation and shown only in command list.
    command: echo hello world, %1! # add args with %1, %2, etc.
    #shell: cmd > optional, default is cmd in Windows, sh in Linux. Set this to the shell you want to use.
    #shellargs: /C > optional, default is /C in Windows, -c in Linux. Set this to the arg needed to make the shell execute your command.
    env: # optional, set environment variables here.
      - name: TEST
        value: test
```

- List commands in growl.yaml using `growl l # or growl list`
- Run them with `growl [command] [args]`

- Run/build your project with `growl [b/build/r/run] [args]`
- Cross compile with  
`growl cross --os [os] --arch [arch] --ldflags [ldflags]`  
or `growl cross -o [os] -a [arch] -ld [ldflags]`
- Print availbable os/architectures with `growl cross list`