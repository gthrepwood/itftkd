
import json
from pathlib import Path
from collections import OrderedDict

data_path = Path("/home/crown/dev/forgejo/PersonalRepo/questionary1/backend/data/data.json")
output_path = Path("/home/crown/dev/forgejo/PersonalRepo/questionary1/backend/data/DATA_TABLE.md")

with open(data_path, "r", encoding="utf-8") as f:
    data = json.load(f)

all_keys = OrderedDict()
for item in data:
    for k in item.keys():
        all_keys[k] = None

priority = ["number", "topic", "name", "id", "script", "pronunciation", "translation", "description", "grade", "belt_color", "diagram", "steps", "descr", "link", "version", "dupe"]
ordered_keys = []
for p in priority:
    if p in all_keys:
        ordered_keys.append(p)
        del all_keys[p]
for k in all_keys.keys():
    ordered_keys.append(k)

lines = []
lines.append("# Taekwondo Dictionary - Data Table")
lines.append("")
lines.append(f"Total items: {len(data)}")
lines.append("")

header_line = "| # | " + " | ".join(ordered_keys) + " |"
separator = "| " + " | ".join(["---"] * (len(ordered_keys) + 1)) + " |"

lines.append(header_line)
lines.append(separator)

for idx, item in enumerate(data, start=1):
    row = [str(idx)]
    for k in ordered_keys:
        val = item.get(k, "")
        if isinstance(val, bool):
            val = str(val).lower()
        else:
            val = str(val)
        row.append(val)
    lines.append("| " + " | ".join(row) + " |")

lines.append("")

output_path.parent.mkdir(parents=True, exist_ok=True)
with open(output_path, "w", encoding="utf-8") as f:
    f.write("\n".join(lines))

print(f"Done. Written to {output_path}")
