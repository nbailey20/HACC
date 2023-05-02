
from rich.console import Console
from rich.prompt import Prompt

import time

from rich.align import Align
from rich.text import Text
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress
from rich.progress import SpinnerColumn
from rich.progress import BarColumn
from rich.progress import TextColumn
from rich.progress import TaskProgressColumn
from rich.padding import Padding


console = Console()
console.rule()


def do_work():
    for i in range(1):
        # console.print(f"checking vault component {i}")
        time.sleep(1)
    return

# with console.status("Starting up...", spinner="point"):
#     do_work()


# console.rule()

# for n in range(4):

#     table = Table()

#     table.add_column("#", justify="right", style="cyan", no_wrap=True)
#     table.add_column("Service Name", style="magenta")

#     table.add_row("1", "test")
#     table.add_row("2", "testlongername")
#     table.add_row("3", "test")
#     table.add_row("4", "test")
#     table.add_row("5", "test")
#     table.add_row("6", "test")
#     table.add_row("3", "test")
#     table.add_row("4", "test")
#     table.add_row("5", "test")
#     table.add_row("6", "test")

#     if n == 0:
#         console.print(table)
#     else:
#         console.print(Padding(table, (2, 0, 0, 0)))
#     name = Prompt.ask("Enter a service/#", default="more...")


# console.rule(style="color(5)")

# console.print([1, 2, 3])
# console.print("[blue underline]Looks like a link")
# console.print("FOO", style="white on blue")

# console.input("What is [i]your[/i] [bold red]name[/]? :smiley: ")

# name = Prompt.ask("Enter your name", choices=["Paul", "Jessica", "Duncan"], default="Paul")



# with console.status("Starting up...", spinner="point"):
#     console.print("[red]Downloading...")
#     do_work()
#     console.print("[green]Processing...")
#     do_work()
#     console.print("[blue]Fucking...")
#     do_work()



with Progress(
    SpinnerColumn(spinner_name="point"),
    TextColumn("[progress.description]{task.description}"),
    BarColumn(),
    TaskProgressColumn(),
    transient=True
) as progress:

    task1 = progress.add_task("[red]Downloading...", total=1)
    task2 = progress.add_task("[green]Processing...", total=1)
    task3 = progress.add_task("[blue]Fucking...", total=1)
    #task3 = progress.add_task()

    do_work()
    progress.update(task1, advance=1)
    
    do_work()
    progress.update(task2, advance=1)
    #task2 = progress.add_task("[green]Processing...", total=5)
    
    do_work()
    progress.update(task3, advance=1)


console.print(Panel('done'))