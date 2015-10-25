# gow Readmap

## 0.1

### Parse flag
Parse flag from `gow` command. 

- [x] Watch directory path
- [ ] Designate run command
- [ ] Watch file extension

#### Watch directory path
`--path` and alias `-p` flag

#### Designate run command
`--command` and alias `-c` flag

#### Watch file extension
`--file` and alias `-f` flag

## 0.2

- [ ] Parse grow DSL `json`
- [ ] Human Readable DSL for `Web Designer`, `Manager` and etc...
- [ ] `init` command for create init `json` file.

### Parse `gow` DSL `JSON` file

Support parse `json` file.

However, `gow` want easy to use for all managerial position.  
So, DSL design need to be Human Readable.  
e.g. `create`, `edit:`, `delete:` and etc.

```json
{
    "path" : "./",
    "extension" : [ "go", "Makefile" ],
    "event" : [
        {
            "create": [
                {
                  "command": "echo create",
                }
            ]
        },
        {
            "edit": [
                {
                  "command": "echo edit",
                }
            ]
        },
        {
            "delete": [
                {
                  "command": "echo Delete!!",
                }
            ]
        }
    ]
}
```

### Create init json file
Easy to create `gow` DSL `json` file.  
I think `gow init` command flag.

# Milstone
[Milstone](https://github.com/zchee/gow/milestones)

