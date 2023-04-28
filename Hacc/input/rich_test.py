
from rich.console import Console
from rich.prompt import Prompt

import time

from rich.align import Align
from rich.text import Text
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress
from rich.padding import Padding


console = Console()

def do_work():
    for i in range(3):
        console.print(f"checking vault component {i}")
        time.sleep(2)
    return

with console.status("Starting up...", spinner="point"):
    do_work()


console.rule()

for n in range(4):

    table = Table()

    table.add_column("#", justify="right", style="cyan", no_wrap=True)
    table.add_column("Service Name", style="magenta")

    table.add_row("1", "test")
    table.add_row("2", "testlongername")
    table.add_row("3", "test")
    table.add_row("4", "test")
    table.add_row("5", "test")
    table.add_row("6", "test")
    table.add_row("3", "test")
    table.add_row("4", "test")
    table.add_row("5", "test")
    table.add_row("6", "test")

    if n == 0:
        console.print(table)
    else:
        console.print(Padding(table, (2, 0, 0, 0)))
    name = Prompt.ask("Enter a service/#", default="more...")


console.rule(style="color(5)")

console.print([1, 2, 3])
console.print("[blue underline]Looks like a link")
console.print("FOO", style="white on blue")

console.input("What is [i]your[/i] [bold red]name[/]? :smiley: ")

name = Prompt.ask("Enter your name", choices=["Paul", "Jessica", "Duncan"], default="Paul")




with Progress() as progress:

    task1 = progress.add_task("[red]Downloading...", total=None)

    i = 0
    while i < 3:
        i += 1
        time.sleep(1)
        progress.update(task1, advance=1)
       
    progress.stop_task(task1)

    task2 = progress.add_task("[green]Processing...", total=5)
    

    while not progress.finished:
        progress.update(task2, advance=1.3)
        time.sleep(0.5)

console.print('done')