# Event added!

## Remember to run `golem generate` to update the generated code.

The implementation of this event is in `app/events_[name].go`. You will need to implement the
`on[Name]` function to handle the event. You can also add additional functions to this file to
help with the implementation.

`Sender` events can be used with the `app.Events.Send()` function to notify other services. `Receiver`
events are handled with `on[Name]` function mentioned above. The `on[Name]` function of `Receiver`
events with the `Proxy` setting are slightly different, they are intended to transform the incoming
event into an outgoing event. For more information, run:

> golem readme
