# `cancel`

`cancel` is a library which is made to just be a _simple_ form of `context.Context`.

What we wanted was a library to cancel tasks safely. `context.Context` does this, but it
provides several other features and complications we didn't ask for, such as
a (really) complex `WithValue` and tons of forward and back references.

It's not super optimized at the moment, but I wanted to make sure it was easy to use.

A `Canceller` is cancelled in one of three ways:
 - If its cancel function is called
 - If its deadline is reached
 - If its parent is cancelled

# What makes `cancel` better?

Well, I haven't written benchmarks yet. I'm not exactly a benchmark writer, so if someone
else wants to do that, that'd be awesome.

Cancel is much more readable and understandable as to how it works. If you look at the
implementation of `context`, it gets really complex really fast. I'm hoping it's
more memory and time efficient than `context`, although I'm not entirely sure.

Since the standard library uses `context` everywhere (for obvious reasons), there
is a function to convert a `Canceller` to a `Context`, `cancel.Context(canc)`

# Installation

`cancel` is a go module which is go-gettable. This means you will be able to fetch any version
of `cancel` that you want, and it will be compatible with the semantic versioning guarantee.

It is notable however that (as of me writing this README) the current major version is v0, so there
are no compatability guarantees until v1, which should be coming soon.

# Usage

You will find many similarities to `context`.

Root `Canceller`:
```
cancel.NeverCancel()
```

With a deadline:
```
canc := cancel.NeverCancel()
canc, f = cancel.WithDeadline(canc, time.Now().Add(5 * time.Second))
defer f()

// or more simply
canc, f = cancel.WithTimeLimit(canc, 5 * time.Second)
defer f()
```

With a cancel function:
```
canc, f = cancel.WithCanceller(canc)
// call f when you want to cancel
```


# Future features

I want to allow `WithDeadline` and `WithTimeLimit` to not return cancel functions. Currently
they do because creating them spins up a goroutine. Therefore the only two ways to
end the goroutine is to either reach the deadline or to cancel the parent. Perhaps this
is negligible, although I want to allow for efficient programming with this.

Perhaps instead of having `cancel.WithDeadline(canc, t)`, it should be `canc.WithDeadline(t)`.
It makes more intuitive sense for them to be methods rather than package-level functions.