import html
import flask

app = flask.Flask(__name__)
i = 0


@app.route("/")
def index():
    global i
    i += 1
    print(i)
    return "Hello, World!"


@app.route("/ip")
def ip():
    ip = flask.request.remote_addr
    actual_ip = flask.request.headers.get("X-Real-IP")

    return f"Remote IP: {html.escape(ip)}\nActual IP: {html.escape(actual_ip)}"


if __name__ == "__main__":
    app.run(port=9812)
