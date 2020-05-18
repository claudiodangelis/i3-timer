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
| (default) | Timer is idle  |
| Green | Timer is in its first half |
| Yellow | Timer is in its second half |
| Red | Timer is in its last quarter |


**Note**: if pango markup is not rendered, for example the output is something like `<span color="red">Timer: 5m0s</span>`, you will need to set `font pango:Monospace 10` to the `bar` section of the i3's configuration file.

#### Status label

| Label | Status |
| ----- | ------ |
| `[I]` | Timer is idle |
| `[R]` | Timer is showing remaining time |
| `[E]` | Timer is showing elapsed time |

### Arguments

| Argument | Type | Description |
| -------- | ---- | ----------- |
| `-alarm-command` | String | Command(s) to be executed when the  alarm fires |
| `-colors` | Boolean | Prints colorized timer |
| `-debug` | Boolean | Prints debug information |
| `-recurring` | Boolean | Restarts timer after the alarm fires |
| `-duration` | Integer | Sets the default duration, in minutes, of the timer |
| `-autostart` | Boolean | Starts the timer as soon as i3 starts |

## Usage



You can interact with the timer by clicking and/or using the scrollwheel.

| Event | Outcome |
| ----- | ------- |
| Left Click | Toggles between elapsed/remaining time |
| Middle Click | Starts the timer |
| Right Click | Stops and resets the timer |
| Scroll Up | Adjusts the timer to add one minute |
| Scroll Down | Adjusts the timer to remove one minute |



If passed, the value of `-alarm-command` will be executed when the timer runs out. In the example above the value is `notify-send 'i3-timer' 'Alarm Elapsed!'; play /usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga`, which will pop up a notification on the top right corner of the screen and play a sound.

You can execute whatever program(s) you like.

## Key Bindings

You can setup key bindings to start and stop the timer by using i3blocks' [signaling](https://github.com/vivien/i3blocks#signal) feature. What you should do is to create 2 additional blocks with no `interval` property but set the proper `signal` value for the blocks, then set i3 configuration to bind the emission of the given signals to a keyboard shortcut.

The following example will bind the start of timer to **Mod1+Shift+Control+k** (which will emit, for example, signal `10`) and the stop to **Mod1+Shift+Control+l** (which will emit, for example, signal `11`). After emitting signal `10` or `11`, both key shortcuts will emit signal `12` that will refresh the timer's "gui".


1. Create the main block to configure the timer in i3blocks configuration file:

    ```ini
    [i3timer]
    command=~/go/bin/i3-timer -alarm-command="notify-send 'i3-timer' 'Alarm Elapsed!'; play /usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga"
    interval=10
    signal=12
    ```

2. Create the "start timer" block  in i3blocks configuration file

    ```ini
    [i3timer]
    command=~/go/bin/i3-timer -exec-start
    signal=10
    ```

3. Create the "stop timer" block  in i3blocks configuration file

    ```ini
    [i3timer]
    command=~/go/bin/i3-timer -exec-stop
    signal=11
    ```

4. Create the keyboard shortcut binding in i3 configuration file
    ```sh
    # start timer
    bindsym --release Mod1+Shift+Control+k exec bash -c "pkill -SIGRTMIN+10 i3blocks && pkill -SIGRTMIN+12 i3blocks"
    # stop timer
    bindsym --release Mod1+Shift+Control+l exec bash -c "pkill -SIGRTMIN+11 i3blocks && pkill -SIGRTMIN+12 i3blocks"
    ```

5. Restart i3 and enjoy

**Note**: `i3-timer` stores timer's data into a hidden JSON file, `.i3-timer.json` stored at your home directory.


## Credits

The original author and current maintainer is [Claudio d'Angelis](https://github.com/claudiodangelis). Refer to [contributors page](https://github.com/claudiodangelis/i3-timer/graphs/contributors) for a full list of people who have contributed to this project. Contributions are always welcome.


## License 

MIT
