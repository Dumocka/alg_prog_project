from app import app, db
from app import User, Survey, Question, Option, Response

def upgrade_database():
    with app.app_context():
        # Пересоздаем базу данных
        db.drop_all()
        db.create_all()
        print("Database recreated successfully!")

if __name__ == '__main__':
    upgrade_database()
