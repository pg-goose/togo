# Togo
> Stupid Simple todo-list app

Most to-do list apps are full-featured tools. I found them distracting so I needed something simpler, like paper and pen, but on my laptop. This tool replicates the simplicity of writing tasks in a notebook for a work session.

![togo ui](ui.png)

### **Install**
Download the latest binaries for [Linux](https://github.com/pg-goose/togo/releases/latest/download/togo-linux), [Windows](https://github.com/pg-goose/togo/releases/latest/download/togo-win) or [MacOS](https://github.com/pg-goose/togo/releases/latest/download/togo-darwin)  

> [!TIP]
> I recomend renaming the binary to `togo` and moving it to a directory added to PATH. The main reason for this was convenience, after all.
> ```
> sudo chmod +x togo-linux
> sudo mv togo-linux /usr/bin/togo 
> ```

### **Controls**  
- `Arrows` to move up and down
- `Enter` to create a task on the task input
- `Enter` or `Space` to mark a task as done
- `Backspace` to errase a task
- `Ctrl+C` to quit

### **Saved tasks**  
Tasks are saved at `$HOME/.config/togo/tasks.json`

### **Task format**  
```
{
    "task" : string,
    "completed" : bool
}
```

### Built with Bubbletea

[Bubbletea github](github.com/charmbracelet/bubbletea)
