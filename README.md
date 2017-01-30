# secret2value

Update enviroment variable to expand symbolic referrence to Docker Swarm secret with the actual value.
By convention, enviroment values with value `"secret:<name>"` will have their value replaced
by the content of file `/run/secrets/<name>`. The path to secrets may be overriden by setting
`$SECRETS_PATH`

The intention is to run this as a wrapper for another command which already expects the values
to be present in enviroment variable. For example:

```bash
secret2value confd -onetime -backend env
```
