## growl
A simple cli tool to make it easier to build and run your project, and automatize other tasks, by adding custom commands.  
This is pretty similar to NodeJS package.json commands.

## usage
First, make a growl.yaml file in your project. There's a example one in this repo.
```yaml
commands:
  - name: build # build and run commands written here will overwrite default ones :).
    description: "echo test"
    command: echo test
  - name: hello
    description: "prints hello, World!" # this description is unused, this is only for documentation and shown only in command list.
    command: echo "hello world, %1!" # add args with %1, %2, etc.
    shell: cmd # optional, default is cmd in Windows, bash in Linux. Set this to the shell you want to use.
    shellargs: /C # optional, default is /C in Windows, none in Linux. Set this to the arg needed to make the shell execute your command.
```

- List commands in growl.yaml using `growl l # or growl list`
- 