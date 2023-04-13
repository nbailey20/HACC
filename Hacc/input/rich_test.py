from rich.layout import Layout
from rich.console import Console
from rich.live import Live
from rich.panel import Panel
from rich.prompt import Prompt
from rich.text import Text

import keyboard

refresh_string = [f"[green]{str(i)}" for i in range(9)]
refresh_string[0] = "[yellow]0"


console = Console()
layout = Layout(name="base")
layout.split(
    Layout(name="refresh", size=len(refresh_string)+2),
    Layout(name="consistent")
)
layout["refresh"].update(Panel("\n".join(refresh_string), expand=False))
consistent_text = Text("Make a selection or left/right arrows for more options...")
consistent_text.stylize("bold magenta")
layout["consistent"].update(consistent_text)

with Live(layout):
    idx = 0
    while True:
        event = keyboard.read_event()
        if event.event_type == keyboard.KEY_DOWN and event.name == 'down':
            refresh_string[idx] = f"[green]{str(idx)}"
            refresh_string[idx+1] = f"[yellow]{str(idx+1)}"
            layout["refresh"].update(Panel("\n".join(refresh_string), expand=False))
            idx += 1
        elif event.event_type == keyboard.KEY_DOWN and event.name == 'up':
            refresh_string[idx] = f"[green]{str(idx)}"
            refresh_string[idx-1] = f"[yellow]{str(idx-1)}"
            layout["refresh"].update(Panel("\n".join(refresh_string), expand=False))
            idx -= 1
        elif event.event_type == keyboard.KEY_DOWN and event.name in [str(i) for i in range(1,10)]:
            refresh_string[idx] = f"[green]{str(idx)}"
            idx = int(event.name)
            refresh_string[idx] = f"[yellow]{str(idx)}"
            layout["refresh"].update(Panel("\n".join(refresh_string), expand=False))

        # if idx < 0:
        #     idx = len(refresh_string)-1
        # elif idx > len(refresh_string)-2:
        #     idx = 0

