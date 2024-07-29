# Bug reproduction for go-gst #94

https://github.com/go-gst/go-gst/issues/94

## Description

Heavily concurrent pipelines sometimes corrupt memory while detaching signal handlers.

## Requirements

- Go 1.22 installed
- pkg-config installed
- gstreamer and glib installed where pkg-config can find them

## Steps to reproduce

```bash
go run .
```

This may compile a while (because of cgo) and then runs the program. This crashes pretty quickly on my machine.

After the first compilation, the cgo part is cached and it will compile faster.

## Crash output

```
(issue94:80736): GLib-GObject-CRITICAL **: 13:35:53.239: g_closure_ref: assertion 'closure->ref_count > 0' failed

(issue94:80736): GLib-GObject-CRITICAL **: 13:35:53.239: g_closure_unref: assertion 'closure->ref_count > 0' failed
.corrupted size vs. prev_size
SIGABRT: abort
PC=0x7fed51605e44 m=3 sigcode=18446744073709551610
signal arrived during cgo execution
```

```
(issue94:83391): GLib-GObject-CRITICAL **: 13:45:17.790: ../glib/gobject/gsignal.c:2685: instance '0x7bfd841e70f0' has no handler with id '3639'
```

## crashing Code

The unreffing/invalidation of the closures is pretty complicated and a bit obfuscated by the fact that CGo calls the unrefs on our side. On closure invalidation (either by [disconnecting the handler](https://github.com/go-gst/go-glib/blob/main/glib/gobject.go#L323) or other ways) cgo calls the [attached finalizer](https://github.com/go-gst/go-glib/blob/main/glib/connect.go#L32) (defined [in the C part](https://github.com/go-gst/go-glib/blob/main/glib/glib.go.h#L153-L157)) that [removes the closure from our map](https://github.com/go-gst/go-glib/blob/main/glib/connect.go#L114).