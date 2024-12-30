from flask import Flask, render_template, request, redirect, url_for, flash, jsonify
from flask_login import LoginManager, login_user, login_required, logout_user, current_user
from werkzeug.security import generate_password_hash, check_password_hash
from datetime import datetime, timedelta
from flask_dance.contrib.github import make_github_blueprint, github
from flask_dance.consumer import oauth_authorized
from oauth_providers import make_yandex_blueprint, yandex
from dotenv import load_dotenv
from database import mongo, User, init_db
from bson import ObjectId
import os
import jwt
from functools import wraps

# Загрузка переменных окружения
load_dotenv()

app = Flask(__name__)
app.config['SECRET_KEY'] = os.getenv('FLASK_SECRET_KEY')
app.config['MONGO_URI'] = 'mongodb://mongodb:27017/pzhendthis'

# Инициализация MongoDB
init_db(app)

def create_jwt_token(username):
    payload = {
        'username': username,
        'exp': datetime.utcnow() + timedelta(days=1)
    }
    return jwt.encode(payload, app.config['SECRET_KEY'], algorithm='HS256')

def verify_jwt_token(token):
    try:
        payload = jwt.decode(token, app.config['SECRET_KEY'], algorithms=['HS256'])
        return payload['username']
    except:
        return None

def jwt_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        token = request.cookies.get('jwt_token')
        if not token:
            return redirect(url_for('login'))
        
        username = verify_jwt_token(token)
        if not username:
            return redirect(url_for('login'))
            
        return f(*args, **kwargs)
    return decorated_function

# GitHub OAuth config
app.config['GITHUB_OAUTH_CLIENT_ID'] = os.getenv('GITHUB_OAUTH_CLIENT_ID')
app.config['GITHUB_OAUTH_CLIENT_SECRET'] = os.getenv('GITHUB_OAUTH_CLIENT_SECRET')
github_bp = make_github_blueprint()
app.register_blueprint(github_bp, url_prefix='/login/github')

# Yandex OAuth config
DOMAIN = os.getenv('DOMAIN', 'example.com')
yandex_bp = make_yandex_blueprint(
    client_id=os.getenv('YANDEX_OAUTH_CLIENT_ID'),
    client_secret=os.getenv('YANDEX_OAUTH_CLIENT_SECRET'),
    redirect_url=f"https://{DOMAIN}/login/yandex/authorized"
)
app.register_blueprint(yandex_bp, url_prefix='/login/yandex')

# Login manager setup
login_manager = LoginManager(app)
login_manager.login_view = 'login'

@login_manager.user_loader
def load_user(user_id):
    return User.get_by_id(user_id)

# OAuth handlers
@oauth_authorized.connect_via(github_bp)
def github_logged_in(blueprint, token):
    if not token:
        flash("Failed to log in with GitHub.", category="error")
        return False

    resp = blueprint.session.get("/user")
    if not resp.ok:
        flash("Failed to fetch user info from GitHub.", category="error")
        return False

    github_info = resp.json()
    github_user_id = str(github_info["id"])
    
    # Поиск пользователя по GitHub ID
    user_data = mongo.db.users.find_one({
        'oauth_provider': 'github',
        'oauth_id': github_user_id
    })

    if not user_data:
        # Создаем нового пользователя
        user_data = {
            'username': github_info["login"],
            'email': github_info.get("email"),
            'oauth_provider': 'github',
            'oauth_id': github_user_id,
            'created_at': datetime.utcnow()
        }
        result = mongo.db.users.insert_one(user_data)
        user_data['_id'] = result.inserted_id
    
    user = User(user_data)
    login_user(user)
    flash("Successfully signed in with GitHub.")
    return False

@oauth_authorized.connect_via(yandex_bp)
def yandex_logged_in(blueprint, token):
    if not token:
        flash("Failed to log in with Yandex.", category="error")
        return False

    resp = blueprint.session.get("/info")
    if not resp.ok:
        flash("Failed to fetch user info from Yandex.", category="error")
        return False

    yandex_info = resp.json()
    yandex_user_id = str(yandex_info["id"])
    
    # Поиск пользователя по Yandex ID
    user_data = mongo.db.users.find_one({
        'oauth_provider': 'yandex',
        'oauth_id': yandex_user_id
    })

    if not user_data:
        # Создаем нового пользователя
        user_data = {
            'username': yandex_info.get("login", f"yandex_user_{yandex_user_id}"),
            'email': yandex_info.get("default_email"),
            'oauth_provider': 'yandex',
            'oauth_id': yandex_user_id,
            'created_at': datetime.utcnow()
        }
        result = mongo.db.users.insert_one(user_data)
        user_data['_id'] = result.inserted_id
    
    user = User(user_data)
    login_user(user)
    flash("Successfully signed in with Yandex.")
    return False

# Routes
@app.route('/')
def index():
    surveys = list(mongo.db.surveys.find())
    return render_template('index.html', surveys=surveys)

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        username = request.form['username']
        password = request.form['password']
        
        if mongo.db.users.find_one({'username': username}):
            flash('Username already taken')
            return redirect(url_for('register'))
        
        user_data = {
            'username': username,
            'password': generate_password_hash(password),
            'created_at': datetime.utcnow()
        }
        result = mongo.db.users.insert_one(user_data)
        user_data['_id'] = result.inserted_id
        
        user = User(user_data)
        login_user(user)
        return redirect(url_for('index'))
    return render_template('register.html')

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form.get('username')
        password = request.form.get('password')
        
        user = mongo.db.users.find_one({'username': username})
        
        if user and check_password_hash(user['password'], password):
            user_obj = User(user)
            login_user(user_obj)
            
            # Создаем JWT токен
            token = create_jwt_token(username)
            response = redirect(url_for('index'))
            response.set_cookie('jwt_token', token)
            return response
            
        flash('Invalid username or password')
    return render_template('login.html')

@app.route('/logout')
@login_required
def logout():
    logout_user()
    return redirect(url_for('index'))

@app.route('/create_survey', methods=['GET', 'POST'])
@jwt_required
def create_survey():
    if request.method == 'POST':
        token = request.cookies.get('jwt_token')
        username = verify_jwt_token(token)
        user = mongo.db.users.find_one({'username': username})
        
        survey_data = {
            'title': request.form.get('title'),
            'questions': request.form.getlist('questions[]'),
            'user_id': str(user['_id']),
            'created_at': datetime.utcnow()
        }
        
        result = mongo.db.surveys.insert_one(survey_data)
        return redirect(url_for('index'))
        
    return render_template('create_survey.html')

@app.route('/edit_survey/<survey_id>', methods=['GET', 'POST'])
@login_required
def edit_survey(survey_id):
    survey = mongo.db.surveys.find_one({'_id': ObjectId(survey_id)})
    if not survey or survey['user_id'] != current_user.id:
        flash('Survey not found or access denied')
        return redirect(url_for('index'))
    
    if request.method == 'POST':
        data = request.get_json()
        mongo.db.surveys.update_one(
            {'_id': ObjectId(survey_id)},
            {'$set': {
                'title': data['title'],
                'description': data['description'],
                'questions': data['questions']
            }}
        )
        return jsonify({'success': True})
    
    return render_template('edit_survey.html', survey=survey)

@app.route('/take_survey/<survey_id>')
@login_required
def take_survey(survey_id):
    survey = mongo.db.surveys.find_one({'_id': ObjectId(survey_id)})
    if not survey:
        flash('Survey not found')
        return redirect(url_for('index'))
    return render_template('take_survey.html', survey=survey)

@app.route('/survey_results/<survey_id>')
@login_required
def survey_results(survey_id):
    survey = mongo.db.surveys.find_one({'_id': ObjectId(survey_id)})
    if not survey or survey['user_id'] != current_user.id:
        flash('Survey not found or access denied')
        return redirect(url_for('index'))
    
    answers = list(mongo.db.responses.find({'survey_id': ObjectId(survey_id)}))
    return render_template('survey_results.html', survey=survey, answers=answers)

@app.route('/delete_survey/<survey_id>', methods=['POST'])
@login_required
def delete_survey(survey_id):
    survey = mongo.db.surveys.find_one({'_id': ObjectId(survey_id)})
    if not survey or survey['user_id'] != current_user.id:
        flash('Survey not found or access denied')
        return redirect(url_for('index'))
    
    # Удаляем опрос и все связанные ответы
    mongo.db.surveys.delete_one({'_id': ObjectId(survey_id)})
    mongo.db.responses.delete_many({'survey_id': ObjectId(survey_id)})
    
    flash('Survey deleted successfully')
    return redirect(url_for('index'))

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
