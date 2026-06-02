#!/usr/bin/env python3
"""Validate example module manifests."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
EXAMPLES = ROOT / "examples"
HTTP_METHODS = {"GET", "POST", "PUT", "PATCH", "DELETE"}
CODE_RE = re.compile(r"^[a-z][a-z0-9_]*(?::[a-z][a-z0-9_]+){1,2}$")


class ManifestError(RuntimeError):
    pass


def require_str(data: dict[str, Any], key: str) -> str:
    value = data.get(key)
    if not isinstance(value, str) or not value.strip():
        raise ManifestError(f"{key} must be a non-empty string")
    return value


def require_bool(data: dict[str, Any], key: str) -> bool:
    value = data.get(key)
    if not isinstance(value, bool):
        raise ManifestError(f"{key} must be a boolean")
    return value


def validate_permission_code(code: str) -> None:
    if not CODE_RE.match(code):
        raise ManifestError(f"invalid permission code: {code}")


def validate_manifest(path: Path) -> None:
    data = json.loads(path.read_text())
    if data.get("version") != 1:
        raise ManifestError("version must be 1")
    module = require_str(data, "module")
    require_str(data, "entity")
    require_str(data, "table")
    require_str(data, "backendPackage")
    require_bool(data, "runtimeRegistered")
    require_bool(data, "requiresSchema")

    permissions = data.get("permissions")
    if not isinstance(permissions, list) or not permissions:
        raise ManifestError("permissions must be a non-empty list")
    permission_codes: set[str] = set()
    for permission in permissions:
        if not isinstance(permission, dict):
            raise ManifestError("permission entries must be objects")
        code = require_str(permission, "code")
        validate_permission_code(code)
        if code in permission_codes:
            raise ManifestError(f"duplicate permission code: {code}")
        permission_codes.add(code)
        for key in ("name", "module", "resource", "action"):
            require_str(permission, key)

    backend = data.get("backend")
    if not isinstance(backend, dict):
        raise ManifestError("backend must be an object")
    routes = backend.get("routes")
    if not isinstance(routes, list) or not routes:
        raise ManifestError("backend.routes must be a non-empty list")
    seen_routes: set[tuple[str, str]] = set()
    for route in routes:
        if not isinstance(route, dict):
            raise ManifestError("backend route entries must be objects")
        method = require_str(route, "method").upper()
        route["method"] = method
        if method not in HTTP_METHODS:
            raise ManifestError(f"unsupported route method: {method}")
        route_path = require_str(route, "path")
        if not route_path.startswith("/"):
            raise ManifestError(f"route path must start with /: {route_path}")
        route_key = (method, route_path)
        if route_key in seen_routes:
            raise ManifestError(f"duplicate route: {method} {route_path}")
        seen_routes.add(route_key)
        permission = require_str(route, "permission")
        if permission not in permission_codes:
            raise ManifestError(f"route permission is not declared: {permission}")

    frontend = data.get("frontend")
    if not isinstance(frontend, dict):
        raise ManifestError("frontend must be an object")
    api_path = require_str(frontend, "api")
    if not api_path.startswith("admin/src/api/"):
        raise ManifestError(f"frontend.api must live under admin/src/api: {api_path}")
    views = frontend.get("views")
    if not isinstance(views, list) or not views:
        raise ManifestError("frontend.views must be a non-empty list")
    for view in views:
        if not isinstance(view, str) or not view.startswith("admin/src/views/"):
            raise ManifestError(f"frontend view must live under admin/src/views: {view}")

    menu = data.get("menu")
    if not isinstance(menu, dict):
        raise ManifestError("menu must be an object")
    for key in ("code", "parent", "type", "name", "routePath", "routeName", "component", "permission"):
        require_str(menu, key)
    if menu["type"] not in {"catalog", "page", "button"}:
        raise ManifestError(f"unsupported menu.type: {menu['type']}")
    if menu["permission"] not in permission_codes:
        raise ManifestError(f"menu permission is not declared: {menu['permission']}")
    if not isinstance(menu.get("visible"), bool):
        raise ManifestError("menu.visible must be a boolean")
    if not isinstance(menu.get("sort"), int):
        raise ManifestError("menu.sort must be an integer")
    if module not in menu["component"] and module not in menu["routeName"]:
        raise ManifestError("menu should reference the module in component or routeName")


def main() -> int:
    manifests = sorted(EXAMPLES.glob("*/manifest.json"))
    if not manifests:
        raise ManifestError("no module manifests found")
    for manifest in manifests:
        try:
            validate_manifest(manifest)
        except (json.JSONDecodeError, ManifestError) as exc:
            print(f"FAIL: {manifest.relative_to(ROOT)}: {exc}", file=sys.stderr)
            return 1
        print(f"OK: {manifest.relative_to(ROOT)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
