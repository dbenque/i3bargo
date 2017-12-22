# i3bargo

Run i3status, decode the output, update the output with addition or midifying existing entries.

The implementation is used to retrieve the current context and namespace for Kubernetes and append the information in the i3bar on top of the report done by i3status. For the moment there is no configuration for i3bargo.

Then update the the i3 config file with:

```
# Start i3bar to display a workspace bar (plus the system information i3status
# finds out, if available)
bar {
        status_command /home/dbenque/bin/i3bargo
        #status_command i3status
        tray_output primary
}
```

For another way od doing the same kind of thing have a look at: https://github.com/nicklan/i3config

