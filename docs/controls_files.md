# controls files

## Controls directory

You can use `controls.d/example.toml.disable` ([online version](https://github.com/axxelG/loxwebhook/blob/master/controls.d/example.toml.disabled)) as a good starting point to create your own controls file.

All filles ending with `.toml` in the controls directory will be imported.

The decision to keep everything in one file or use multiple files is up to you. All authentication keys and names of controls must be unique for all files. If you have configured an authentication key `Key1` in `file1.toml` you cannot configure `Key1` again in `file2.toml` but you can use `Key1` in a control definition in `file2.toml`.

## Controls files

Control files must be valid [TOML](https://github.com/toml-lang/toml)

### Section `[AuthKeys]`

List of key/value pairs. Format: `name = "key"`
You can use any ASCII-Character (A-Z upper and lower case), numbers, hyphens (-) and underscores (_). Authentication keys are case sensitive.

Example

```toml
[AuthKeys]
testOne   = "43b2c690-f281-42bb-af2d-979f5dbe9517"
testTwo   = "69b9a1ad-1224-4c93-8411-e88e65ebe582"
testThree = "84627dbd-bd68-476f-9e53-35522285783b"
```

### Section `[Controls]`

Table (dictionary) of control definitions.

### Control definition

The name of a control definition must only consist of ASCII-Characters (A-Z), numbers, hyphens (-) and underscores (_).

| Field    | Descriptions                                                  |
|----------|---------------------------------------------------------------|
| Category | Type of control. Currently only `dvi` for "digital virtual input" is supported  |
| ID | Miniserver internal ID number of the control. You can find the ID number in Loxone Config if you select the control and look at Property / Common / Connection |
| Allowed  | Array of allowed commands. You can find a list of allowed command on the [Loxone website](https://www.loxone.com/enen/kb/web-services/) |
| AuthKeys | Array of key names that can access this control. The names must exactly match a name configured in Section `[AuthKeys]`. You can use authentication keys defined in another controls file. |

Examples

```toml
[Controls.test1]
Category = "dvi"
ID = 7
Allowed = [
    "on",
]
AuthKeys = [
    "testOne",
]

[Controls.test2]
Category = "dvi"
ID = 6
Allowed = [
    "on",
    "off",
]
AuthKeys = [
    "testTwo",
    "testThree",
]
```