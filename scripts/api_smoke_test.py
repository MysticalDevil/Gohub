#!/usr/bin/env python3
"""Run real HTTP smoke tests against the Gohub API.

Usage:
  python3 scripts/api_smoke_test.py --base-url http://127.0.0.1:3000

The script assumes the API server is already running with a clean test
database. It creates its own users/categories/topics through public HTTP
endpoints and verifies every route registered in routes/api.go.
"""

from __future__ import annotations

import argparse
import json
import mimetypes
import struct
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
import zlib
from dataclasses import dataclass
from typing import Any


def png_1x1_red() -> bytes:
    def chunk(kind: bytes, data: bytes) -> bytes:
        checksum = zlib.crc32(kind + data) & 0xFFFFFFFF
        return struct.pack(">I", len(data)) + kind + data + struct.pack(">I", checksum)

    signature = b"\x89PNG\r\n\x1a\n"
    ihdr = struct.pack(">IIBBBBB", 1, 1, 8, 2, 0, 0, 0)
    idat = zlib.compress(b"\x00\xff\x00\x00")
    return signature + chunk(b"IHDR", ihdr) + chunk(b"IDAT", idat) + chunk(b"IEND", b"")


@dataclass
class APIResponse:
    status: int
    body: dict[str, Any]


class APISmoke:
    def __init__(self, base_url: str) -> None:
        self.base_url = base_url.rstrip("/")
        self.results: list[str] = []

    def run(self) -> None:
        self.wait_until_ready()

        self.expect("GET /api/v1/user unauthorized", "GET", "/api/v1/user", 401)
        self.expect("GET /api/v1/links", "GET", "/api/v1/links", 200)

        self.expect(
            "POST /api/v1/auth/verify-codes/captcha",
            "POST",
            "/api/v1/auth/verify-codes/captcha",
            200,
            {},
            require_data_keys=("captcha_id", "captcha_image"),
        )
        self.expect(
            "POST /api/v1/auth/verify-codes/phone",
            "POST",
            "/api/v1/auth/verify-codes/phone",
            200,
            {
                "phone": "00012345678",
                "captcha_id": "captcha_skip_test",
                "captcha_answer": "123456",
            },
        )
        self.expect(
            "POST /api/v1/auth/verify-codes/email",
            "POST",
            "/api/v1/auth/verify-codes/email",
            200,
            {
                "email": "apiowner@testing.com",
                "captcha_id": "captcha_skip_test",
                "captcha_answer": "123456",
            },
        )

        owner_token = self.signup_with_phone()
        other_token = self.signup_with_email()

        self.expect(
            "POST /api/v1/auth/signup/phone/exist",
            "POST",
            "/api/v1/auth/signup/phone/exist",
            200,
            {"phone": "00012345678"},
            data_contains={"exist": True},
        )
        self.expect(
            "POST /api/v1/auth/signup/email/exist",
            "POST",
            "/api/v1/auth/signup/email/exist",
            200,
            {"email": "apiother@testing.com"},
            data_contains={"exist": True},
        )

        login_token = self.login_and_refresh()
        self.expect(
            "GET /api/v1/user authorized",
            "GET",
            "/api/v1/user",
            200,
            headers=self.auth_headers(login_token),
            require_data_keys=("id", "name"),
        )

        self.user_routes(login_token)
        category_id = self.category_routes(owner_token, other_token)
        self.topic_routes(owner_token, other_token, category_id)

    def signup_with_phone(self) -> str:
        response = self.expect(
            "POST /api/v1/auth/signup/using-phone",
            "POST",
            "/api/v1/auth/signup/using-phone",
            201,
            {
                "phone": "00012345678",
                "verify_code": "123456",
                "name": "apiowner",
                "password": "password123",
                "password_confirm": "password123",
            },
            require_data_keys=("token", "user"),
        )
        return self.object_data(response)["token"]

    def signup_with_email(self) -> str:
        response = self.expect(
            "POST /api/v1/auth/signup/using-email",
            "POST",
            "/api/v1/auth/signup/using-email",
            201,
            {
                "email": "apiother@testing.com",
                "verify_code": "123456",
                "name": "apiother",
                "password": "password123",
                "password_confirm": "password123",
            },
            require_data_keys=("token", "user"),
        )
        return self.object_data(response)["token"]

    def login_and_refresh(self) -> str:
        self.expect(
            "POST /api/v1/auth/login/using-phone",
            "POST",
            "/api/v1/auth/login/using-phone",
            200,
            {"phone": "00012345678", "verify_code": "123456"},
            require_data_keys=("token",),
        )
        response = self.expect(
            "POST /api/v1/auth/login/using-password",
            "POST",
            "/api/v1/auth/login/using-password",
            200,
            {
                "login_id": "apiowner",
                "password": "password123",
                "captcha_id": "captcha_skip_test",
                "captcha_answer": "123456",
            },
            require_data_keys=("token",),
        )
        token = self.object_data(response)["token"]
        refresh = self.expect(
            "POST /api/v1/auth/login/refresh-token",
            "POST",
            "/api/v1/auth/login/refresh-token",
            200,
            headers=self.auth_headers(token),
            require_data_keys=("token",),
        )
        return self.object_data(refresh)["token"]

    def user_routes(self, token: str) -> None:
        headers = self.auth_headers(token)
        self.expect("GET /api/v1/users", "GET", "/api/v1/users", 200)
        self.expect(
            "PUT /api/v1/users",
            "PUT",
            "/api/v1/users",
            200,
            {"name": "apiowner", "city": "beijing", "introduction": "hello world"},
            headers=headers,
        )
        self.expect(
            "PUT /api/v1/users/email",
            "PUT",
            "/api/v1/users/email",
            200,
            {"email": "apiowner2@testing.com", "verify_code": "123456"},
            headers=headers,
        )
        self.expect(
            "PUT /api/v1/users/phone",
            "PUT",
            "/api/v1/users/phone",
            200,
            {"phone": "00012345679", "verify_code": "123456"},
            headers=headers,
        )
        self.expect(
            "PUT /api/v1/users/password",
            "PUT",
            "/api/v1/users/password",
            200,
            {
                "password": "password123",
                "new_password": "password456",
                "new_password_confirm": "password456",
            },
            headers=headers,
        )
        self.expect_multipart(
            "PUT /api/v1/users/avatar",
            "PUT",
            "/api/v1/users/avatar",
            200,
            "avatar",
            "avatar.png",
            png_1x1_red(),
            headers=headers,
        )
        self.expect(
            "POST /api/v1/auth/password-reset/using-phone",
            "POST",
            "/api/v1/auth/password-reset/using-phone",
            200,
            {"phone": "00012345679", "verify_code": "123456", "password": "password789"},
        )
        self.expect(
            "POST /api/v1/auth/password-reset/using-email",
            "POST",
            "/api/v1/auth/password-reset/using-email",
            200,
            {
                "email": "apiowner2@testing.com",
                "verify_code": "123456",
                "password": "password999",
            },
        )

    def category_routes(self, owner_token: str, other_token: str) -> str:
        self.expect("GET /api/v1/categories", "GET", "/api/v1/categories", 200)
        created = self.expect(
            "POST /api/v1/categories",
            "POST",
            "/api/v1/categories",
            201,
            {"name": "catone", "description": "category description"},
            headers=self.auth_headers(owner_token),
            require_data_keys=("id", "name"),
        )
        category_id = str(self.object_data(created)["id"])
        self.expect(
            "PUT /api/v1/categories/:id forbidden",
            "PUT",
            f"/api/v1/categories/{category_id}",
            403,
            {"name": "catbad", "description": "category description"},
            headers=self.auth_headers(other_token),
        )
        self.expect(
            "PUT /api/v1/categories/:id",
            "PUT",
            f"/api/v1/categories/{category_id}",
            200,
            {"name": "cattwo", "description": "updated description"},
            headers=self.auth_headers(owner_token),
        )

        delete_target = self.expect(
            "POST /api/v1/categories for delete",
            "POST",
            "/api/v1/categories",
            201,
            {"name": "catdel", "description": "delete category"},
            headers=self.auth_headers(owner_token),
        )
        delete_id = str(self.object_data(delete_target)["id"])
        self.expect(
            "DELETE /api/v1/categories/:id forbidden",
            "DELETE",
            f"/api/v1/categories/{delete_id}",
            403,
            headers=self.auth_headers(other_token),
        )
        self.expect(
            "DELETE /api/v1/categories/:id",
            "DELETE",
            f"/api/v1/categories/{delete_id}",
            200,
            headers=self.auth_headers(owner_token),
        )
        self.expect(
            "DELETE /api/v1/categories/:id not found",
            "DELETE",
            "/api/v1/categories/999999",
            404,
            headers=self.auth_headers(owner_token),
        )
        return category_id

    def topic_routes(self, owner_token: str, other_token: str, category_id: str) -> None:
        self.expect("GET /api/v1/topics", "GET", "/api/v1/topics", 200)
        created = self.expect(
            "POST /api/v1/topics",
            "POST",
            "/api/v1/topics",
            201,
            {
                "title": "new topic",
                "body": "this is a new topic body",
                "category_id": category_id,
            },
            headers=self.auth_headers(owner_token),
            require_data_keys=("id", "title"),
        )
        topic_id = str(self.object_data(created)["id"])
        self.expect("GET /api/v1/topics/:id", "GET", f"/api/v1/topics/{topic_id}", 200)
        self.expect(
            "PUT /api/v1/topics/:id forbidden",
            "PUT",
            f"/api/v1/topics/{topic_id}",
            403,
            {
                "title": "bad topic",
                "body": "this is a bad topic body",
                "category_id": category_id,
            },
            headers=self.auth_headers(other_token),
        )
        self.expect(
            "PUT /api/v1/topics/:id",
            "PUT",
            f"/api/v1/topics/{topic_id}",
            200,
            {
                "title": "updated topic",
                "body": "this is an updated topic body",
                "category_id": category_id,
            },
            headers=self.auth_headers(owner_token),
        )
        self.expect(
            "POST /api/v1/topics missing category",
            "POST",
            "/api/v1/topics",
            422,
            {
                "title": "missing cat",
                "body": "this topic points at a missing category",
                "category_id": "999999",
            },
            headers=self.auth_headers(owner_token),
        )
        self.expect(
            "DELETE /api/v1/topics/:id",
            "DELETE",
            f"/api/v1/topics/{topic_id}",
            200,
            headers=self.auth_headers(owner_token),
        )
        self.expect(
            "DELETE /api/v1/topics/:id not found",
            "DELETE",
            "/api/v1/topics/999999",
            404,
            headers=self.auth_headers(owner_token),
        )

    def wait_until_ready(self) -> None:
        deadline = time.time() + 20
        last_error: Exception | None = None
        while time.time() < deadline:
            try:
                self.request("GET", "/api/v1/links")
                return
            except Exception as exc:
                last_error = exc
                time.sleep(0.25)
        raise AssertionError(f"API did not become ready: {last_error}")

    def expect(
        self,
        name: str,
        method: str,
        path: str,
        status: int,
        body: dict[str, Any] | None = None,
        headers: dict[str, str] | None = None,
        require_data_keys: tuple[str, ...] = (),
        data_contains: dict[str, Any] | None = None,
    ) -> APIResponse:
        response = self.request(method, path, body, headers)
        self.assert_status(name, response, status)
        if require_data_keys or data_contains:
            data = self.object_data(response)
            for key in require_data_keys:
                if key not in data or data[key] in ("", None):
                    raise AssertionError(f"{name}: missing data.{key}: {response.body}")
            for key, expected in (data_contains or {}).items():
                if data.get(key) != expected:
                    raise AssertionError(
                        f"{name}: expected data.{key}={expected!r}, got {data.get(key)!r}"
                    )
        self.results.append(name)
        print(f"ok {name}")
        return response

    def expect_multipart(
        self,
        name: str,
        method: str,
        path: str,
        status: int,
        field_name: str,
        filename: str,
        content: bytes,
        headers: dict[str, str] | None = None,
    ) -> APIResponse:
        boundary = "gohub-api-smoke-boundary"
        content_type = mimetypes.guess_type(filename)[0] or "application/octet-stream"
        payload = (
            f"--{boundary}\r\n"
            f'Content-Disposition: form-data; name="{field_name}"; filename="{filename}"\r\n'
            f"Content-Type: {content_type}\r\n\r\n"
        ).encode() + content + f"\r\n--{boundary}--\r\n".encode()
        request_headers = dict(headers or {})
        request_headers["Content-Type"] = f"multipart/form-data; boundary={boundary}"
        response = self.raw_request(method, path, payload, request_headers)
        self.assert_status(name, response, status)
        self.results.append(name)
        print(f"ok {name}")
        return response

    def request(
        self,
        method: str,
        path: str,
        body: dict[str, Any] | None = None,
        headers: dict[str, str] | None = None,
    ) -> APIResponse:
        payload = None if body is None else json.dumps(body).encode()
        request_headers = dict(headers or {})
        if payload is not None:
            request_headers["Content-Type"] = "application/json"
        return self.raw_request(method, path, payload, request_headers)

    def raw_request(
        self,
        method: str,
        path: str,
        payload: bytes | None,
        headers: dict[str, str] | None = None,
    ) -> APIResponse:
        request_headers = {"User-Agent": "gohub-api-smoke"}
        request_headers.update(headers or {})
        req = urllib.request.Request(
            self.url(path),
            data=payload,
            headers=request_headers,
            method=method,
        )
        try:
            with urllib.request.urlopen(req, timeout=10) as res:
                return self.parse_response(res.status, res.read())
        except urllib.error.HTTPError as err:
            return self.parse_response(err.code, err.read())

    def parse_response(self, status: int, raw_body: bytes) -> APIResponse:
        if not raw_body:
            return APIResponse(status=status, body={})
        try:
            parsed = json.loads(raw_body.decode())
        except json.JSONDecodeError as exc:
            raise AssertionError(f"response was not JSON: status={status}, body={raw_body!r}") from exc
        if not isinstance(parsed, dict):
            raise AssertionError(f"response JSON was not an object: {parsed!r}")
        return APIResponse(status=status, body=parsed)

    def assert_status(self, name: str, response: APIResponse, expected: int) -> None:
        if response.status != expected:
            raise AssertionError(
                f"{name}: expected HTTP {expected}, got {response.status}: {response.body}"
            )

    def object_data(self, response: APIResponse) -> dict[str, Any]:
        data = response.body.get("data")
        if data is None:
            return {}
        if not isinstance(data, dict):
            raise AssertionError(f"expected object data, got: {data!r}")
        return data

    def auth_headers(self, token: str) -> dict[str, str]:
        return {"Authorization": f"Bearer {token}"}

    def url(self, path: str) -> str:
        return urllib.parse.urljoin(self.base_url + "/", path.lstrip("/"))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", default="http://127.0.0.1:3000")
    args = parser.parse_args()

    suite = APISmoke(args.base_url)
    suite.run()
    print(f"passed {len(suite.results)} API smoke checks")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"api smoke failed: {exc}", file=sys.stderr)
        raise SystemExit(1)
