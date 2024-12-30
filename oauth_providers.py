from flask_dance.consumer import OAuth2ConsumerBlueprint
from functools import partial
from flask.globals import LocalProxy
from flask import current_app
import os

def make_yandex_blueprint(
    client_id=None,
    client_secret=None,
    scope=None,
    redirect_url=None,
    redirect_to=None,
    login_url=None,
    authorized_url=None,
    session_class=None,
    storage=None,
):
    """
    Создает blueprint для аутентификации через Яндекс OAuth.
    """
    scope = scope or ["login:info", "login:email"]
    yandex_bp = OAuth2ConsumerBlueprint(
        "yandex",
        __name__,
        client_id=os.getenv("YANDEX_CLIENT_ID"),
        client_secret=os.getenv("YANDEX_CLIENT_SECRET"),
        scope=scope,
        base_url="https://login.yandex.ru/info",
        authorization_url="https://oauth.yandex.ru/authorize",
        token_url="https://oauth.yandex.ru/token",
        redirect_url=redirect_url,
        redirect_to=redirect_to,
        login_url=login_url,
        authorized_url=authorized_url,
        session_class=session_class,
        storage=storage,
    )

    @yandex_bp.before_app_request
    def set_applocal_session():
        ctx = current_app.extensions.setdefault("yandex", {})
        ctx["session"] = yandex_bp.session

    return yandex_bp

yandex = LocalProxy(partial(
    lambda: current_app.extensions["yandex"]["session"]
))
