# Request

## Parts of a request

A request to loxwebhook looks like this:

`https://your.domain.com/dvi/garage_door/open?k=11111111-1111-1111-1111-111111111111`

or more abstract

`https://domain/control_type/control_name/control_action?k=SecretKey`

| Part             | Description |
| ---------        | ---------------------------------- |
| domain           | The domain where the server that runs loxwebhook is reachable |
| control_type     | The type of the control we are accessing. Currently only `dvi` for "Digital virtual input" is supported |
| control_name     | The name of the control. It must exactly match the name we used in the [controls file](controls_files.md). |
| control_action | The action we want to send to the control. The action must be allowed in the [controls file](controls_files.md). |
| SecretKey        | A secret key configured in the [controls file](controls_files.md). Please read and understand the [Security Q&A](security_qa.md) before you choose a key. |

## Additional parameters

| Parameter   | Descriptions |
|-------------|--------------|
| simulate    | Prevents loxwebhook from sending requests to the Loxone Miniserver and returns config details |
