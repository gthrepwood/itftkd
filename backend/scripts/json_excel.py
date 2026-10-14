#!/usr/bin/env python3
import json
import html
import os
import re
import sys
from datetime import datetime
from pathlib import Path
from collections import OrderedDict


def collect_fieldnames(data):
    fieldnames = OrderedDict()
    for item in data:
        for k in item.keys():
            fieldnames[k] = None
    return list(fieldnames.keys())


def json_to_html(json_path, html_path):
    with open(json_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    if not isinstance(data, list) or not data:
        raise ValueError("JSON root must be a non-empty list")

    fieldnames = collect_fieldnames(data)

    lines = []
    lines.append("<!DOCTYPE html>")
    lines.append('<html lang="en">')
    lines.append("<head>")
    lines.append('<meta charset="UTF-8">')
    lines.append("<title>Taekwondo Data</title>")
    lines.append("</head>")
    lines.append("<body>")
    lines.append("<table border='1'>")
    lines.append("<thead>")
    lines.append("<tr>")
    for k in fieldnames:
        lines.append(f"<th>{html.escape(str(k))}</th>")
    lines.append("</tr>")
    lines.append("</thead>")
    lines.append("<tbody>")
    for item in data:
        lines.append("<tr>")
        for k in fieldnames:
            val = item.get(k, "")
            if val is None:
                val = ""
            lines.append(f"<td>{html.escape(str(val))}</td>")
        lines.append("</tr>")
    lines.append("</tbody>")
    lines.append("</table>")
    lines.append("</body>")
    lines.append("</html>")

    with open(html_path, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))


def html_to_json(html_path, json_path):
    with open(html_path, "r", encoding="utf-8") as f:
        content = f.read()

    header_pattern = re.compile(r"<th\b[^>]*>(.*?)</th>", re.DOTALL | re.IGNORECASE)
    cell_pattern = re.compile(r"<td\b[^>]*>(.*?)</td>", re.DOTALL | re.IGNORECASE)

    headers = [html.unescape(m.group(1).strip()) for m in header_pattern.finditer(content)]
    if not headers:
        raise ValueError("No table headers found in HTML")

    data = []
    row_start = 0
    while True:
        row_match = re.search(r"<tr[^>]*>(.*?)</tr>", content[row_start:], re.DOTALL | re.IGNORECASE)
        if not row_match:
            break
        row_html = content[row_start:][row_match.start():row_match.end()]
        cells = [html.unescape(m.group(1).strip()) for m in cell_pattern.finditer(row_html)]
        if not cells:
            row_start += row_match.end()
            continue

        if len(cells) != len(headers):
            if len(cells) > len(headers):
                cells = cells[: len(headers)]
            else:
                cells = cells + [""] * (len(headers) - len(cells))

        item = {}
        for k, v in zip(headers, cells):
            if v == "":
                continue
            if k in ("number", "id", "steps", "version"):
                try:
                    item[k] = int(v)
                except ValueError:
                    item[k] = v
            elif k == "dupe":
                item[k] = v.lower() == "true"
            else:
                item[k] = v
        data.append(item)
        row_start += row_match.end()

    data = [d for d in data if any(v != "" for v in d.values())]

    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)


def main():
    if len(sys.argv) < 2:
        print("Usage:")
        print("  json_to_html  <input.json>  [output.html]")
        print("  html_to_json  <input.html>  [output.json]")
        sys.exit(1)

    mode = sys.argv[1].lower()
    if mode == "json_to_html":
        if len(sys.argv) < 3:
            print("Missing input JSON path")
            sys.exit(1)
        input_path = sys.argv[2]
        if len(sys.argv) >= 4:
            output_path = sys.argv[3]
        else:
            today = datetime.now().strftime("%Y%m%d")
            output_path = f"data_{today}.html"
        json_to_html(input_path, output_path)
        print(f"Converted {input_path} -> {output_path}")
    elif mode == "html_to_json":
        if len(sys.argv) < 3:
            print("Missing input HTML path")
            sys.exit(1)
        input_path = sys.argv[2]
        if len(sys.argv) >= 4:
            output_path = sys.argv[3]
        else:
            base = Path(input_path).stem
            m = re.match(r"data_(\d+)", base)
            if m:
                output_path = f"data_{m.group(1)}.json"
            else:
                output_path = str(Path(input_path).with_suffix(".json"))
        html_to_json(input_path, output_path)
        print(f"Converted {input_path} -> {output_path}")
    else:
        print(f"Unknown mode: {mode}")
        sys.exit(1)


if __name__ == "__main__":
    main()
