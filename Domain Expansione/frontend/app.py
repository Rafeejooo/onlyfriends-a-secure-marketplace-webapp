from flask import Flask, render_template, request, redirect, make_response, jsonify, request
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from dotenv import load_dotenv
import os
import requests, jwt


app = Flask(__name__, template_folder='templates', static_folder='static')


BACKEND_URL = "http://backend:8000"  # Docker service name

load_dotenv()  # Load environment variables from .env

JWT_SECRET = os.getenv("JWT_SECRET")
DB_USER = os.getenv("DB_USER")

limiter = Limiter(
    get_remote_address,  
    app=app,
    default_limits=["200 per day", "50 per hour"]
)

@app.route('/')
def index():
    error = request.args.get("error")
    return render_template("login.html", error=error)

@app.route('/login', methods=['POST'])
@limiter.limit("5 per minute") 
def login():
    email = request.form['email']
    password = request.form['pswd']

    try:
        # Send login credentials to Go backend
        res = requests.post(
            f"{BACKEND_URL}/login/submit",
            data={"email": email, "pswd": password},
            allow_redirects=False
        )

        if res.status_code == 303:
            jwt_token = res.cookies.get("token")
            resp = make_response(redirect("/main"))
            if jwt_token:
                resp.set_cookie("token", jwt_token, httponly=True)
            return resp
        else:
            return redirect("/?error=Invalid+credentials")

    except Exception as e:
        return redirect(f"/?error=Backend+Error:+{str(e)}")

@app.route('/register', methods=['POST'])
@limiter.limit("5 per minute") 
def register():
    email = request.form['email']
    password = request.form['pswd']
    username = request.form['txt']
    phone = request.form['broj']

    try:
        res = requests.post(
            f"{BACKEND_URL}/register",
            data={"email": email, "pswd": password, "txt": username, "broj": phone}
        )

        if res.status_code == 303:
            return redirect('/')
        else:
            return render_template("login.html", error="Registration failed")

    except Exception as e:
        return render_template("login.html", error=f"Error: {str(e)}")
    
# ORDER submission: sets status = pending, returns JSON for popup JS
@app.route('/order', methods=['POST'])
def submit_order():
    name = request.form.get("name")
    date = request.form.get("date")
    package = request.form.get("package")
    # talent_id = request.form.get("talent_id")
    token = request.cookies.get("token")

    print("🌐 Incoming order:", name, date, package)  # ← add this!

    try:
        print("📤 Sending to backend:", flush=True)
        print("name:", name, flush=True)
        print("date:", date, flush=True)
        print("package:", package, flush=True)
        # print("talent_id:", talent_id, flush=True)

        res = requests.post(
            f"{BACKEND_URL}/order/confirm",  
            data={
                "name": name,
                "date": date,
                "package": package,
                # "talent_id": talent_id
            },
            cookies={"token": token}
        )

        print("🔁 Backend response:", res.status_code, res.text)

        if res.status_code == 200:
            try:
                order_id = res.json().get("order_id")
                print("✅ Success! Order ID:", order_id)
                return jsonify({"success": True, "order_id": order_id})
            except Exception as e:
                print("❌ JSON parse error:", str(e))
                return jsonify({"success": False, "error": "Invalid JSON from backend"})
        else:
            return jsonify({"success": False, "error": res.text})

    except Exception as e:
        print("🔥 Exception:", str(e))
        return jsonify({"success": False, "error": str(e)})


# CONFIRM payment form, returns JSON for popup JS
@app.route('/confirm-payment', methods=['POST'])
def confirm_payment():
    order_id = request.form.get("order_id")

    try:
        res = requests.post(
            f"{BACKEND_URL}/payment/confirm",
            data={"order_id": order_id}
        )

        if res.status_code == 200:
            return jsonify({"success": True})
        else:
            return jsonify({"success": False, "error": "Failed to confirm payment"})

    except Exception as e:
        return jsonify({"success": False, "error": str(e)})

@app.route('/logout')
def logout():
    resp = make_response(redirect("/"))
    
    # Clear the token cookie (remove JWT from browser)
    resp.set_cookie("token", "", expires=0, httponly=True, samesite='Lax')

    return resp


JWT_ALGORITHM = "HS256"

@app.route('/main')
def main_page():
    token = request.cookies.get("token")
    if not token:
        return redirect("/?error=Login+required")
    
    try:
        payload = jwt.decode(token, JWT_SECRET, algorithms=[JWT_ALGORITHM])
        username = payload.get("username", "User")
    except Exception as e:
        print("JWT decode error:", str(e))
        return redirect("/?error=Invalid+token")

    return render_template("main.html", user_name=username)

    # Optionally, verify token by calling backend
    # res = requests.get(f"{BACKEND_URL}/validate", cookies={"token": token})

    # return render_template("main.html")

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)
