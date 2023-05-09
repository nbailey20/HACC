import sys

try:
    from rich.console import Console
except:
    print('Python module "rich" required for HACC client. Install (pip install rich) and try again.')
    sys.exit()

NUM_CHOICES_PER_TABLE = 10

console = Console()
