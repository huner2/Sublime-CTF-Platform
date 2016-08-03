#!/usr/bin/env python

"""server.py -- the main flask server module"""

import simplejson as json
import random, string
import time
import re

from base64 import b64decode
from functools import wraps

from flask import Flask
from flask import jsonify
from flask import make_response
from flask import redirect
from flask import render_template
from flask import request
from flask import session
from flask import url_for
from flask.ext.seasurf import SeaSurf
from werkzeug.security import generate_password_hash
from flask_sqlalchemy import SQLAlchemy as sql
from sql_structure import * # Namespace pollution is best pollution

app = Flask(__name__, static_folder='static', static_url_path='')

# Load Config
config_str = open("config.json", "rb").read()
config = json.loads(config_str)

app.config['SQLALCHEMY_DATABASE_URI'] = config['db']
db = sql(app)
csrf = SeaSurf(app)

lang = None
sponsors = False
minpasslength = None

def get_user():
	return (False, None)

@app.route("/about")
def about():
	"""Displays the about page"""
	
	login, user = get_user()
	
	# Render the page
	render = render_template('frame.html', lang=lang, sponsored=sponsors, page='about.html', login=login, user=user)
	return make_response(render)

@app.route("/sponsors")
def sponsors():
	"""Displays the sponsors page if there are sponsors"""
	
	login, user = get_user()
	
	# Render the page
	if sponsors:
		render = render_template('frame.html', lang=lang, page='sponsors.html', login=login, user=user, sponsored=sponsors)
		return make_response(render)
	else:
		return redirect('/error/404') # In case it is manually entered with no sponsors to display.

@app.errorhandler(404)
def page_not_found(error):
	return redirect('/error/404')
	
@app.errorhandler(500)
def server_error(error):
	return redirect('/error/500')
	
@app.errorhandler(405)
def method_not_allowed(error):
	return redirect('/error/405')

@app.route("/error/<error>")
def error(error):
	"""Displays the error page with the corresponding error message"""
	
	login, user = get_user()
	
	if error not in lang['error']:
		error = "Unknown"
		
	# Render the page
	render = render_template('frame.html', lang=lang, sponsored=sponsors, page='error.html', login=login, user=user, error=error)
	return make_response(render)
	
@app.route("/register")
def register():
	"""Displays the register page"""
	
	login, user = get_user()
	
	if login:
		return redirect("/error/404") # If they are already logged in, don't let them register again.
	
	# Render the page
	render = render_template('frame.html', lang=lang, sponsored=sponsors, page='register.html', login=login, user=user,
		minplength=config["minimum_password_length"],maxplength=config["maximum_password_length"], minulength=config["minimum_username_length"],
		maxulength=config["maximum_username_length"],mtlength=config["maximum_team_name_length"],mtmembers=config["maximum_players_per_team"])
	return make_response(render)
	
@app.route("/check", methods = ["POST"])
def check_info():
	"""Check info currently entered on register form"""
	
	login, user = get_user
	
	if login:
		return redirect("/error/405")

@app.route("/register", methods = ["POST"])
def register_submit():
	"""Attempt to register user"""
	
	login, user = get_user()
	
	if login:
		return redirect("/error/405")
	
	firstname = request.form["firstname"]
	lastname = request.form["lastname"]
	email = request.form["email"]
	username = request.form["username"]
	password = request.form["password"]
	cteamname = request.form["cteam-name"]
	jteamname = request.form["jteam-name"]
	jteamcode = request.form["jteam-code"]
	
	# This is really inefficent and may be cut down later.
	
	if len(username) < config["minimum_username_length"] or len(username) > config["maximum_username_length"] or re.match('^[\w-]+$', username) is None:
		return redirect("/register")
	if len(password) < config["minimum_password_length"] or len(password) > config["maximum_password_length"]:
		return redirect("/register")
	if len(cteamname) > config["maximum_team_name_length"]:
		return redirect("/register")
	if cteamname != "" and (jteamname != "" or jteamcode != ""):
		return redirect("/register")
	if cteamname == "" and jteamname == "":
		return redirect("/register")
	if not firstname.isalpha() or not lastname.isalpha():
		return redirect("/register")
	if not re.match("[^@]+@[^@]+\.[^@]+", email):
		return redirect("/register")
	if not User.query.filter_by(username=username).first() == None:
		return redirect("/register")
	if not User.query.filter_by(email=email).first() == None:
		return redirect("/register")
		
	if cteamname != "":
		if not Team.query.filter_by(name=cteamname).first() == None:
			return redirect("/register")
		team = cteamname
		newTeam = Team(cteamname, str(random.randint(10000,99999)) + ''.join(random.choice(string.lowercase + string.uppercase) for i in range(5)))
		db.session.add(newTeam)
		
	if jteamname != "":
		if Team.query.filter_by(name=jteamname).first() == None:
			return redirect("/register")
		team = Team.query.filter_by(name=jteamname).first()
		if team.code != jteamcode:
			return redirect("/register")
		if len(User.query.filter_by(team=jteamname).all()) == 4:
			return redirect("/register")
		team = jteamname
	
	newUser = User(username,generate_password_hash(password),firstname,lastname,email,team)
	db.session.add(newUser)
	db.session.commit()
	
	# Database sessions close automatically, so don't complain.
	
	return redirect("/team")	

@app.route("/")
def index():
	"""Displays the index page"""
	
	login, user = get_user() 
	
	# Render the page
	render = render_template('frame.html', lang=lang, sponsored=sponsors, page='index.html', login=login, user=user)
	return make_response(render)

if __name__ == "__main__":
	"""Initializes variables and starts the server"""

	# Configure Security

	app.secret_key = config["secret_key"]

	# Load Language

	lang_str = open(config["language_file"], "rb").read()
	lang = json.loads(lang_str)
	lang = lang[config["language"]]
	
	# Sponsors
	
	sponsors = config["sponsors"]
	
	# Initialize Database
	db.create_all()
	
	# Run Server
	
	app.run(host=config['host'], port=config['port'], debug=config['debug'], threaded=True)
