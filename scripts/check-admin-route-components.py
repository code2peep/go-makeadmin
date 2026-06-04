#!/usr/bin/env python3
import csv
import pathlib
import re
import sys


ROOT = pathlib.Path(__file__).resolve().parents[1]
SEED_SQL = ROOT / "sql" / "p1.seed.sql"
VIEWS_DIR = ROOT / "admin" / "src" / "views"


def parse_sql_tuple(line: str) -> list[str]:
    text = line.strip().rstrip(",;")
    if not text.startswith("(") or not text.endswith(")"):
        return []
    reader = csv.reader([text[1:-1]], quotechar="'", skipinitialspace=True)
    return [item.strip() for item in next(reader)]


def load_seed_menus() -> list[dict[str, str]]:
    content = SEED_SQL.read_text()
    match = re.search(
        r"INSERT INTO `ma_menu` .*? VALUES\n(?P<values>.*?);\n",
        content,
        flags=re.S,
    )
    if not match:
        raise SystemExit("FAIL: cannot find ma_menu seed block")

    menus = []
    for line in match.group("values").splitlines():
        values = parse_sql_tuple(line)
        if not values:
            continue
        if len(values) < 18:
            raise SystemExit(f"FAIL: malformed ma_menu seed row: {line}")
        menus.append(
            {
                "id": values[0],
                "parent_id": values[1],
                "menu_type": values[2],
                "name": values[3],
                "route_path": values[5],
                "component": values[7],
            }
        )
    return menus


def main() -> int:
    menus = load_seed_menus()
    ids = {menu["id"] for menu in menus}
    children_by_parent = {}
    failures = []

    for menu in menus:
        children_by_parent.setdefault(menu["parent_id"], []).append(menu)
        if menu["parent_id"] != "0" and menu["parent_id"] not in ids:
            failures.append(f"{menu['name']}: parent_id {menu['parent_id']} not found")

    for menu in menus:
        if menu["menu_type"] == "page":
            if not menu["route_path"]:
                failures.append(f"{menu['name']}: page route_path is empty")
            if not menu["component"]:
                failures.append(f"{menu['name']}: page component is empty")
                continue
            component_file = VIEWS_DIR / f"{menu['component']}.vue"
            if not component_file.is_file():
                failures.append(f"{menu['name']}: missing view {component_file.relative_to(ROOT)}")
        elif menu["menu_type"] == "catalog" and not children_by_parent.get(menu["id"]):
            failures.append(f"{menu['name']}: catalog has no children")

    if failures:
        print("FAIL: admin route component contract failed")
        for failure in failures:
            print(f"- {failure}")
        return 1

    page_count = sum(1 for menu in menus if menu["menu_type"] == "page")
    catalog_count = sum(1 for menu in menus if menu["menu_type"] == "catalog")
    print(f"OK: admin route component contract passed ({page_count} pages, {catalog_count} catalogs)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
