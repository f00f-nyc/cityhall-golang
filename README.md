This is the go library for City Hall Enterprise Settings Server

# ABOUT

 This project can be installed using:

```
go get github.com/f00f-nyc/cityhall-golang/
````

# USAGE

 The intention is to use the built-in City Hall web site for actual
 settings management, and then use this library for consuming those
 settings, in an application.  As such, there is really only one
 command to be familiar with:

```go
settings, _ := cityhall.NewSettingsFromUrl("http://path.to.server/api")
val, _ := settings.GetVal("/test/val1")
```

# LICENSE

This code is licensed under the MIT License

