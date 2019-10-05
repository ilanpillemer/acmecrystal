# Acmecrystal

You can build and run acmecrystal.

```
go build
```

Acmecrystal watches the acme log for any Crystal files being Put, and then reformats the file and reloads the window.
That's all it does. It's nice.

# Know Issues
It replaces the whole file on format, and leaves the cursor at the end. This needs to be fixed still.

