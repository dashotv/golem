# Golem

- [ ] if Group Model matches a model name, use that model
- [ ] call out to dotenvx after generate?
- [ ] add UI generator? or just static package generator? and routes wiring?
- [ ] cache plugin to handle boilerplate instead of injecting into each app
- [ ] routes plugin to avoid the echo/v4 problem?
- [ ] migrate tool with vN -> vN scripts / code

## Done

- [x] switch to fae?
- [x] move app to internal? properly support output config
- [x] name files .gen.go?
- [x] update Dockerfile
- [x] update drone file
- [x] update air file

## Cancelled

- break up into packages? there's too much interdependency
- condense generated code into a single file? too complicated with the modify, would have to rewrite the generators
