from flask_pymongo import PyMongo
from datetime import datetime
from bson import ObjectId
from flask_login import UserMixin

mongo = PyMongo()

class User(UserMixin):
    def __init__(self, user_data):
        self.user_data = user_data

    def get_id(self):
        return str(self.user_data['_id'])

    @property
    def id(self):
        return self.user_data['_id']

    @property
    def username(self):
        return self.user_data['username']

    @property
    def email(self):
        return self.user_data.get('email')

    @property
    def oauth_provider(self):
        return self.user_data.get('oauth_provider')

    @staticmethod
    def get_by_id(user_id):
        try:
            user_data = mongo.db.users.find_one({'_id': ObjectId(user_id)})
            return User(user_data) if user_data else None
        except:
            return None

    @staticmethod
    def get_by_username(username):
        user_data = mongo.db.users.find_one({'username': username})
        return User(user_data) if user_data else None

    @staticmethod
    def create_user(username, password_hash, email=None, oauth_provider=None, oauth_id=None):
        user_data = {
            'username': username,
            'password_hash': password_hash,
            'email': email,
            'oauth_provider': oauth_provider,
            'oauth_id': oauth_id,
            'created_at': datetime.utcnow()
        }
        result = mongo.db.users.insert_one(user_data)
        user_data['_id'] = result.inserted_id
        return User(user_data)

def init_db(app):
    mongo.init_app(app)
    
    # Создаем индексы
    with app.app_context():
        # Уникальный индекс для username
        mongo.db.users.create_index('username', unique=True)
        # Уникальный индекс для oauth_id в комбинации с oauth_provider
        mongo.db.users.create_index([('oauth_id', 1), ('oauth_provider', 1)], 
                                  unique=True, 
                                  sparse=True)
        # Индекс для поиска по email
        mongo.db.users.create_index('email', sparse=True)

        # Индексы для опросов
        mongo.db.surveys.create_index('user_id')
        mongo.db.surveys.create_index('created_at')
        
        # Индексы для ответов
        mongo.db.responses.create_index([('survey_id', 1), ('user_id', 1)])
