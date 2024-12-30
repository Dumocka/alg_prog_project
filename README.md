# Survey Platform

A web-based survey platform built with Flask that allows users to create, take, and manage surveys.

## Features

- User authentication (register/login)
- Create custom surveys
- Multiple question types (text, single choice, multiple choice)
- Take surveys
- View survey results
- Responsive design

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd pzhendthis
```

2. Create a virtual environment and activate it:
```bash
python -m venv venv
.\venv\Scripts\activate  # On Windows
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Run the application:
```bash
python app.py
```

The application will be available at `http://localhost:5000`

## Usage

1. Register a new account or login with existing credentials
2. Create a new survey by clicking "Create New Survey"
3. Add questions and customize their types (text, single choice, multiple choice)
4. Save the survey
5. Share the survey link with participants
6. View responses in the survey results section

## Project Structure

```
pzhendthis/
├── app.py              # Main application file
├── requirements.txt    # Python dependencies
├── static/            # Static files
│   ├── style.css     # CSS styles
│   └── script.js     # JavaScript code
└── templates/         # HTML templates
    ├── base.html     # Base template
    ├── index.html    # Home page
    ├── login.html    # Login page
    ├── register.html # Registration page
    └── ...           # Other templates
```

## License

This project is licensed under the MIT License.
