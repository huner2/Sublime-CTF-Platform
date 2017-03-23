from __main__ import app
from flask_sqlalchemy import SQLAlchemy as sql
import simplejson as json

db = sql(app)

# Load Config
config_str = open("config.json", "rb").read()
config = json.loads(config_str)

class User(db.Model):
	id = db.Column(db.Integer, primary_key=True)
	username = db.Column(db.String(config['maximum_username_length']), unique=True)
	password = db.Column(db.String(config['maximum_password_length']), unique=False)
	firstname = db.Column(db.String(80), unique=False)
	lastname = db.Column(db.String(80), unique=False)
	email = db.Column(db.String(120), unique=False)
	team = db.Column(db.String(config['maximum_team_name_length']), unique=False)
	
	def __init__(self, username, password, firstname, lastname, email, team):
		self.username = username
		self.password = password
		self.firstname = firstname
		self.lastname = lastname
		self.email = email
		self.team = team

class Team(db.Model):
	id = db.Column(db.Integer, primary_key=True)
	name = db.Column(db.String(config['maximum_team_name_length']), unique=True)
	code = db.Column(db.Integer, unique=False)
	
	def __init__(self, name, code):
		self.name = name
		self.code = code
		
class Challenge(db.Model):
	id = db.Column(db.Integer, primary_key=True)
	name = db.Column(db.String(80), unique=True)
	points = db.Column(db.Integer, unique=False)
	category = db.Column(db.String(80), unique=False)
	times_solved = db.Column(db.Integer, unique=False)
	
	def __init__(self, name, points, category, times_solved):
		self.name = name
		self.points = points
		self.category = category
		self.times_solved = times_solved
