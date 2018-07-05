# i3-timer

A simple timer for the [i3](https://i3wm.org/) window manager's status bar built for [i3blocks](https://github.com/vivien/i3blocks).

![screenshot](screenshot.png)

## Installation

```shell
go get github.com/claudiodangelis/i3-timer
```


Add these lines to the configuration file, which is usually located at `~/.i3blocks.conf`:

```ini
[i3timer]
command=~/go/bin/i3-timer -alarm-command="notify-send 'i3-timer' 'Alarm Elapsed!'; play /usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga"
interval=10
```

Note: higher values for `interval` lead to poor timer accuracy.

### Colors

The `-colors` flag enables colored output. The config parameter`markup=pango` is required, for example:

```ini
[i3timer]
command=~/go/bin/i3-timer...
markup=pango
interval=10
```

#### Color Codes

| Color | Code |
| ----- | ---- |
| Green | Timer is in its first half |
| Yellow | Timer is in its second half |
| Red | Timer is in its last quarter |


**Note**: if pango markup is not rendered, for example the output is something like `<span color="red">Timer: 5m0s</span>`, you will need to set `font pango:Monospace 10` to the `bar` section of the i3's configuration file.

### Arguments

| Argument | Type | Description |
| -------- | ---- | ----------- |
| `alarm-command` | String | Command(s) to be executed when the  alarm fires |
| `-colors` | Boolean | Prints colorized timer |
| `-debug` | Boolean | Prints debug information |

## Usage



You can interact with the timer by clicking and/or using the scrollwheel.

| Event | Outcome |
| ----- | ------- |
| Left Click | _(does nothing for now, will toggle elapsed/remaining time)_
| Middle Click | Starts the timer |
| Right Click | Stops and resets the timer |
| Scroll Up | Adjusts the timer to add one minute |
| Scroll Down | Adjusts the timer to remove one minute |


If passed, the value of `-alarm-command` will be executed when the timer runs out. In the example above the value is `notify-send 'i3-timer' 'Alarm Elapsed!'; play /usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga`, which will pop up a notification on the top right corner of the screen and play a sound.

You can execute whatever program(s) you like.


**Note**: `i3-timer` stores timer's data into a hidden JSON file, `.i3-timer.json` stored at your home directory.


## Credits

The original author and current maintainer is [Claudio d'Angelis](https://github.com/claudiodangelis). Contributions are always welcome.


## License 

MIT
