# Queue added!

## Remember to run `golem generate` to update the generated code.

`Queues` are used by `Workers`. For the `Worker` to use the `Queue`, you must include it
in the `add` command, like below:

```bash
golem add worker [name] --queue [queue name]
```

For existing workers, you can add a queue by editing the `.golem/workers/[name].yaml` file.

```yaml
queue: name
```
