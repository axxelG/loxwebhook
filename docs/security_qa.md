# Security Q & A

## Is it secure to use loxwebhook and expose my Loxone Miniserver to the public internet?

Yes, beside the fact that it is never 100 % secure to connect something to the internet, you can consider the usage as secure. Loxwebhook has implemented measurements to mitigate the risks:

1. **Key authentication** A request is only forwarded to the Loxone Miniserver if a matching authentication key is provided. Please read the question "How can I be sure my authentication keys are secure?" on this page.

1. **Transport layer encryption** Data that is transferred over the public internet is [TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) encrypted. Only [https](https://en.wikipedia.org/wiki/HTTPS) connections are allowed. A widely trusted [CA](https://en.wikipedia.org/wiki/Certificate_authority) ([Let's Encrypt](https://letsencrypt.org/)) is used. This keeps the data save while it is on the public internet.

1. **Rate limiting for Requests** Loxwebhook does not accept more than ~1 request per second. This makes [brute force attacks](https://en.wikipedia.org/wiki/Brute-force_attack) on the secret keys nearly impossible and prevents the Loxone Miniserver from being overloaded.

Beside all security measurements provided by Loxwebhook you need to be aware that every use case for loxwebhook involves someone who sends requests. You need to trust this second party.

## How can I be sure my authentication keys are secure?

Everybody who knows a key can access the assigned control(s) on you Loxone Miniserver. That's why you must keep them secret. You can (and should) use any ASCII-Character (A-Z upper and lower case), numbers, hyphens (-) and underscores (_).

- Use hard to guess and long keys. It's obvious that a key like "lamp" is not suitable. UUIDs are a good choice. You can easily create them. Use `cat /proc/sys/kernel/random/uuid` on Linux or `[guid]::NewGuid()` in Windows Powershell.

- Create a unique authentication key for every purpose.