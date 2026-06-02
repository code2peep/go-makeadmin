#!/usr/bin/env python3
"""P1 HTTP smoke runner for a disposable go-makeadmin database.

The script intentionally requires P1_SMOKE_ALLOW_WRITE=1 because it exercises
create/update/delete APIs. Point it at a local disposable P1 database-backed
API server, not a shared or production-like database.
"""

from __future__ import annotations

import base64
import json
import os
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Any, Callable


SMOKE_MATRIX = [
    ("auth", "POST /system/login", "read", "login returns a token"),
    ("auth", "GET /system/admin/self", "read", "token resolves current admin"),
    ("auth", "GET /system/menu/route", "read", "token resolves route menus"),
    ("common", "GET /common/index/config", "read", "public config resolves from ma_setting"),
    ("common", "GET /common/index/console", "read", "console resolves from ma_setting"),
    ("log", "GET /system/log/login", "read", "login logs can be queried"),
    ("role", "POST /system/role/add", "write", "role can be created"),
    ("log", "GET /system/log/operate", "read", "operate logs can be queried after a write action"),
    ("role", "POST /system/role/edit", "write", "role can be edited"),
    ("role", "POST /system/role/del", "write", "role can be deleted after admin cleanup"),
    ("admin", "POST /system/admin/add", "write", "admin can be created"),
    ("admin", "POST /system/admin/edit", "write", "admin can be edited"),
    ("admin", "POST /system/admin/disable", "write", "admin status can be toggled"),
    ("admin", "POST /system/admin/del", "write", "admin can be deleted"),
    ("menu", "POST /system/menu/add", "write", "menu action can be created"),
    ("menu", "POST /system/menu/edit", "write", "menu action can be edited"),
    ("menu", "POST /system/menu/del", "write", "menu action can be deleted"),
    ("dict", "POST /setting/dict/type/add", "write", "dict type can be created"),
    ("dict", "POST /setting/dict/data/add", "write", "dict item can be created"),
    ("dict", "POST /setting/dict/data/edit", "write", "dict item can be edited"),
    ("dict", "POST /setting/dict/data/del", "write", "dict item can be deleted"),
    ("dict", "POST /setting/dict/type/del", "write", "dict type can be deleted"),
    ("file", "POST /common/album/cateAdd", "write", "file category can be created"),
    ("file", "POST /common/upload/image", "write", "image upload creates ma_file metadata"),
    ("file", "POST /common/album/albumDel", "write", "uploaded file metadata can be deleted"),
    ("file", "POST /common/album/cateDel", "write", "file category can be deleted"),
    ("codegen", "POST /gen/importTable", "write", "table metadata can be imported"),
    ("codegen", "POST /gen/syncTable", "write", "imported metadata can be synced"),
    ("codegen", "GET /gen/previewCode", "read", "imported metadata can render templates"),
    ("codegen", "GET /gen/downloadCode", "read", "imported metadata can render zip"),
    ("codegen", "POST /gen/delTable", "write", "imported metadata can be deleted"),
]


class SmokeError(RuntimeError):
    pass


@dataclass
class Client:
    base_url: str
    token: str = ""

    def api_json(
        self,
        method: str,
        path: str,
        *,
        params: dict[str, Any] | None = None,
        body: dict[str, Any] | None = None,
        label: str,
    ) -> dict[str, Any]:
        raw = self._request(method, path, params=params, body=body)
        try:
            payload = json.loads(raw.decode("utf-8"))
        except json.JSONDecodeError as exc:
            raise SmokeError(f"{label}: response is not JSON: {raw[:300]!r}") from exc
        code = payload.get("code")
        if code != 200:
            raise SmokeError(f"{label}: code={code}, msg={payload.get('msg')!r}, data={payload.get('data')!r}")
        print(f"OK: {label}")
        return payload

    def api_bytes(self, method: str, path: str, *, params: dict[str, Any] | None = None, label: str) -> bytes:
        raw = self._request(method, path, params=params)
        if not raw:
            raise SmokeError(f"{label}: empty response")
        print(f"OK: {label}")
        return raw

    def upload_png(self, path: str, *, cid: int, label: str) -> dict[str, Any]:
        boundary = f"----go-makeadmin-p1-smoke-{int(time.time())}"
        png = base64.b64decode(
            "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII="
        )
        parts: list[bytes] = []
        parts.append(_multipart_field(boundary, "cid", str(cid).encode("utf-8")))
        parts.append(_multipart_file(boundary, "file", "p1-smoke.png", "image/png", png))
        parts.append(f"--{boundary}--\r\n".encode("utf-8"))
        raw = self._request(
            "POST",
            path,
            data=b"".join(parts),
            content_type=f"multipart/form-data; boundary={boundary}",
        )
        payload = json.loads(raw.decode("utf-8"))
        if payload.get("code") != 200:
            raise SmokeError(f"{label}: code={payload.get('code')}, msg={payload.get('msg')!r}, data={payload.get('data')!r}")
        print(f"OK: {label}")
        return payload

    def _request(
        self,
        method: str,
        path: str,
        *,
        params: dict[str, Any] | None = None,
        body: dict[str, Any] | None = None,
        data: bytes | None = None,
        content_type: str | None = None,
    ) -> bytes:
        query = urllib.parse.urlencode(params or {}, doseq=True)
        url = self.base_url + path
        if query:
            url = f"{url}?{query}"
        headers = {"Accept": "application/json"}
        if self.token:
            headers["token"] = self.token
        if body is not None:
            data = json.dumps(body, separators=(",", ":")).encode("utf-8")
            content_type = "application/json"
        if content_type:
            headers["Content-Type"] = content_type
        req = urllib.request.Request(url, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=15) as resp:
                return resp.read()
        except urllib.error.HTTPError as exc:
            raw = exc.read()
            raise SmokeError(f"{method} {url}: HTTP {exc.code}: {raw[:500]!r}") from exc
        except urllib.error.URLError as exc:
            raise SmokeError(f"{method} {url}: {exc}") from exc


def _multipart_field(boundary: str, name: str, value: bytes) -> bytes:
    return (
        f"--{boundary}\r\n"
        f'Content-Disposition: form-data; name="{name}"\r\n'
        "\r\n"
    ).encode("utf-8") + value + b"\r\n"


def _multipart_file(boundary: str, name: str, filename: str, content_type: str, value: bytes) -> bytes:
    return (
        f"--{boundary}\r\n"
        f'Content-Disposition: form-data; name="{name}"; filename="{filename}"\r\n'
        f"Content-Type: {content_type}\r\n"
        "\r\n"
    ).encode("utf-8") + value + b"\r\n"


def print_matrix() -> None:
    print("module\tmethod_path\tmode\tassertion")
    for module, method_path, mode, assertion in SMOKE_MATRIX:
        print(f"{module}\t{method_path}\t{mode}\t{assertion}")


def data_list(payload: dict[str, Any]) -> list[dict[str, Any]]:
    data = payload.get("data")
    if isinstance(data, dict) and isinstance(data.get("lists"), list):
        return data["lists"]
    if isinstance(data, list):
        return data
    raise SmokeError(f"unexpected list payload shape: {payload!r}")


def flatten_tree(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
    result: list[dict[str, Any]] = []
    for item in items:
        result.append(item)
        children = item.get("children")
        if isinstance(children, list):
            result.extend(flatten_tree(children))
    return result


def find_id(items: list[dict[str, Any]], field: str, value: Any) -> int:
    for item in items:
        if item.get(field) == value:
            item_id = item.get("id")
            if isinstance(item_id, int):
                return item_id
    raise SmokeError(f"cannot find item where {field}={value!r}")


def require_env(name: str) -> str:
    value = os.environ.get(name, "")
    if not value:
        raise SmokeError(f"{name} is required")
    return value


def run() -> None:
    if "--print-matrix" in sys.argv:
        print_matrix()
        return

    if os.environ.get("P1_SMOKE_ALLOW_WRITE") != "1":
        raise SmokeError("refusing to run write smoke without P1_SMOKE_ALLOW_WRITE=1")

    base_url = os.environ.get("P1_SMOKE_BASE_URL", "http://127.0.0.1:8000/api").rstrip("/")
    username = os.environ.get("P1_SMOKE_ADMIN_USERNAME", "admin")
    password = os.environ.get("P1_SMOKE_ADMIN_PASSWORD") or os.environ.get("ADMIN_PASSWORD")
    if not password:
        raise SmokeError("P1_SMOKE_ADMIN_PASSWORD or ADMIN_PASSWORD is required")

    suffix = os.environ.get("P1_SMOKE_SUFFIX") or time.strftime("%m%d%H%M%S")
    client = Client(base_url=base_url)
    cleanups: list[Callable[[], None]] = []

    print(f"==> P1 smoke target: {base_url}")
    print(f"==> P1 smoke suffix: {suffix}")
    try:
        login = client.api_json(
            "POST",
            "/system/login",
            body={"username": username, "password": password},
            label="login",
        )
        token = login.get("data", {}).get("token")
        if not token:
            raise SmokeError("login did not return data.token")
        client.token = str(token)

        client.api_json("GET", "/system/admin/self", label="admin self")
        client.api_json("GET", "/system/menu/route", label="menu route")
        client.api_json("GET", "/system/menu/detail", params={"id": 120}, label="menu detail seed")
        client.api_json("GET", "/common/index/config", label="common index config")
        client.api_json("GET", "/common/index/console", label="common index console")
        client.api_json("GET", "/system/log/login", params={"pageNo": 1, "pageSize": 20}, label="login log list")

        role_name = f"P1SmokeRole{suffix}"
        client.api_json(
            "POST",
            "/system/role/add",
            body={"name": role_name, "sort": 91, "isDisable": 0, "remark": "p1 smoke", "menuIds": "120"},
            label="role add",
        )
        client.api_json("GET", "/system/log/operate", params={"pageNo": 1, "pageSize": 20}, label="operate log list")
        role_items = data_list(client.api_json("GET", "/system/role/list", params={"pageNo": 1, "pageSize": 60}, label="role list"))
        role_id = find_id(role_items, "name", role_name)
        cleanups.append(lambda role_id=role_id: client.api_json("POST", "/system/role/del", body={"id": role_id}, label="cleanup role"))
        client.api_json("GET", "/system/role/detail", params={"id": role_id}, label="role detail")
        client.api_json(
            "POST",
            "/system/role/edit",
            body={"id": role_id, "name": role_name + "Edit", "sort": 92, "isDisable": 0, "remark": "p1 smoke edit", "menuIds": "120"},
            label="role edit",
        )

        admin_username = f"p1smoke{suffix[-8:]}"
        admin_nickname = f"烟测{suffix[-6:]}"
        client.api_json(
            "POST",
            "/system/admin/add",
            body={
                "deptId": 1,
                "postId": 1,
                "username": admin_username,
                "nickname": admin_nickname,
                "password": "P1smoke123!",
                "avatar": "/api/static/backend_avatar.png",
                "role": role_id,
                "sort": 1,
                "isDisable": 0,
                "isMultipoint": 1,
            },
            label="admin add",
        )
        admin_items = data_list(
            client.api_json(
                "GET",
                "/system/admin/list",
                params={"pageNo": 1, "pageSize": 20, "username": admin_username},
                label="admin list",
            )
        )
        admin_id = find_id(admin_items, "username", admin_username)
        cleanups.append(lambda admin_id=admin_id: client.api_json("POST", "/system/admin/del", body={"id": admin_id}, label="cleanup admin"))
        client.api_json("GET", "/system/admin/detail", params={"id": admin_id}, label="admin detail")
        client.api_json(
            "POST",
            "/system/admin/edit",
            body={
                "id": admin_id,
                "deptId": 1,
                "postId": 1,
                "username": admin_username,
                "nickname": admin_nickname + "改",
                "password": "",
                "avatar": "/api/static/backend_avatar.png",
                "role": role_id,
                "sort": 2,
                "isDisable": 0,
                "isMultipoint": 1,
            },
            label="admin edit",
        )
        client.api_json("POST", "/system/admin/disable", body={"id": admin_id}, label="admin disable")

        menu_name = f"P1SmokeMenu{suffix[-6:]}"
        menu_perm = f"system:p1smoke:{suffix[-6:]}"
        client.api_json(
            "POST",
            "/system/menu/add",
            body={
                "pid": 120,
                "menuType": "A",
                "menuName": menu_name,
                "menuIcon": "",
                "menuSort": 99,
                "perms": menu_perm,
                "paths": "",
                "component": "",
                "selected": "",
                "params": "",
                "isCache": 0,
                "isShow": 1,
                "isDisable": 0,
            },
            label="menu add",
        )
        menu_tree = data_list(client.api_json("GET", "/system/menu/list", label="menu list"))
        menu_id = find_id(flatten_tree(menu_tree), "menuName", menu_name)
        cleanups.append(lambda menu_id=menu_id: client.api_json("POST", "/system/menu/del", body={"id": menu_id}, label="cleanup menu"))
        client.api_json("GET", "/system/menu/detail", params={"id": menu_id}, label="menu detail")
        client.api_json(
            "POST",
            "/system/menu/edit",
            body={
                "id": menu_id,
                "pid": 120,
                "menuType": "A",
                "menuName": menu_name + "Edit",
                "menuIcon": "",
                "menuSort": 100,
                "perms": menu_perm,
                "paths": "",
                "component": "",
                "selected": "",
                "params": "",
                "isCache": 0,
                "isShow": 1,
                "isDisable": 1,
            },
            label="menu edit",
        )

        dict_code = f"p1_smoke_{suffix[-8:]}"
        dict_name = f"P1烟测{suffix[-6:]}"
        client.api_json(
            "POST",
            "/setting/dict/type/add",
            body={"dictName": dict_name, "dictType": dict_code, "dictRemark": "p1 smoke", "dictStatus": 1},
            label="dict type add",
        )
        dict_types = data_list(
            client.api_json(
                "GET",
                "/setting/dict/type/list",
                params={"pageNo": 1, "pageSize": 20, "dictType": dict_code},
                label="dict type list",
            )
        )
        dict_type_id = find_id(dict_types, "dictType", dict_code)
        cleanups.append(
            lambda dict_type_id=dict_type_id: client.api_json(
                "POST", "/setting/dict/type/del", body={"ids": [dict_type_id]}, label="cleanup dict type"
            )
        )
        client.api_json(
            "POST",
            "/setting/dict/type/edit",
            body={"id": dict_type_id, "dictName": dict_name + "改", "dictType": dict_code, "dictRemark": "p1 smoke edit", "dictStatus": 1},
            label="dict type edit",
        )
        dict_value = f"value_{suffix[-6:]}"
        client.api_json(
            "POST",
            "/setting/dict/data/add",
            body={"typeId": dict_type_id, "name": "烟测项", "value": dict_value, "remark": "p1 smoke", "sort": 1, "status": 1},
            label="dict data add",
        )
        dict_data = data_list(
            client.api_json(
                "GET",
                "/setting/dict/data/list",
                params={"pageNo": 1, "pageSize": 20, "dictType": dict_code, "value": dict_value},
                label="dict data list",
            )
        )
        dict_data_id = find_id(dict_data, "value", dict_value)
        cleanups.append(
            lambda dict_data_id=dict_data_id: client.api_json(
                "POST", "/setting/dict/data/del", body={"ids": [dict_data_id]}, label="cleanup dict data"
            )
        )
        client.api_json(
            "POST",
            "/setting/dict/data/edit",
            body={"id": dict_data_id, "typeId": dict_type_id, "name": "烟测项改", "value": dict_value, "remark": "p1 smoke edit", "sort": 2, "status": 1},
            label="dict data edit",
        )

        category_name = f"P1SmokeFile{suffix[-6:]}"
        client.api_json(
            "POST",
            "/common/album/cateAdd",
            body={"pid": 0, "type": 10, "name": category_name},
            label="file category add",
        )
        categories = flatten_tree(
            data_list(
                client.api_json(
                    "GET", "/common/album/cateList", params={"type": 10, "keyword": category_name}, label="file category list"
                )
            )
        )
        category_id = find_id(categories, "name", category_name)
        cleanups.append(
            lambda category_id=category_id: client.api_json(
                "POST", "/common/album/cateDel", body={"id": category_id}, label="cleanup file category"
            )
        )
        client.api_json(
            "POST",
            "/common/album/cateRename",
            body={"id": category_id, "name": category_name + "Edit"},
            label="file category rename",
        )
        upload = client.upload_png("/common/upload/image", cid=category_id, label="image upload")
        file_id = upload.get("data", {}).get("id")
        if not isinstance(file_id, int):
            raise SmokeError(f"image upload did not return numeric data.id: {upload!r}")
        cleanups.append(
            lambda file_id=file_id: client.api_json(
                "POST", "/common/album/albumDel", body={"ids": [file_id]}, label="cleanup uploaded file"
            )
        )
        client.api_json(
            "POST",
            "/common/album/albumRename",
            body={"id": file_id, "name": f"p1-smoke-{suffix[-6:]}.png"},
            label="file rename",
        )

        client.api_json("POST", "/gen/importTable", params={"tables": "ma_setting"}, label="codegen import")
        codegen_rows = data_list(
            client.api_json("GET", "/gen/list", params={"pageNo": 1, "pageSize": 20, "tableName": "ma_setting"}, label="codegen list")
        )
        codegen_id = find_id(codegen_rows, "tableName", "ma_setting")
        cleanups.append(
            lambda codegen_id=codegen_id: client.api_json(
                "POST", "/gen/delTable", body={"ids": [codegen_id]}, label="cleanup codegen table"
            )
        )
        client.api_json("GET", "/gen/detail", params={"id": codegen_id}, label="codegen detail")
        client.api_json("POST", "/gen/syncTable", params={"id": codegen_id}, label="codegen sync")
        client.api_json("GET", "/gen/previewCode", params={"id": codegen_id}, label="codegen preview")
        zip_body = client.api_bytes("GET", "/gen/downloadCode", params={"tables": "ma_setting"}, label="codegen download")
        if not zip_body.startswith(b"PK"):
            raise SmokeError("codegen download did not return a zip payload")

    finally:
        failed_cleanups = 0
        for cleanup in reversed(cleanups):
            try:
                cleanup()
            except Exception as exc:  # noqa: BLE001 - keep cleanup best-effort.
                failed_cleanups += 1
                print(f"WARN: cleanup failed: {exc}", file=sys.stderr)
        if failed_cleanups:
            print(f"WARN: {failed_cleanups} cleanup step(s) failed", file=sys.stderr)

    print("==> p1-smoke completed")


if __name__ == "__main__":
    try:
        run()
    except SmokeError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        sys.exit(1)
