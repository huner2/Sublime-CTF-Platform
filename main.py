#!/usr/bin/env python

"""server.py -- the main flask server module"""

import dataset
import simplejson as json
import random
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

app = Flask(__name__, static_folder='static', static_url_path='')
csrf = SeaSurf(app)

lang = None
config = None
sponsors = False
minpasslength = None

def get_user():
	return (False, None)

@app.route("/about")
def about():
	"""Displays the about page"""
	
	login, user = get_user()
	
	# Render the page
	render = render_template('frame.html', lang=lang, page='about.html', login=login, user=user)
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
	render = render_template('frame.html', lang=lang, page='error.html', login=login, user=user, error=error)
	return make_response(render)
	
@app.route("/register")
def register():
	"""Displays the register page"""
	
	login, user = get_user()
	
	if login:
		return redirect("/") # If they are already logged in, don't let them register again.
	
	# Render the page
	render = render_template('frame.html', lang=lang, page='register.html', login=login, user=user,
		minplength=config["minimum_password_length"],maxplength=config["maximum_password_length"], minulength=config["minimum_username_length"],
		maxulength=config["maximum_username_length"],mtlength=config["maximum_team_name_length"],mtmembers=config["maximum_players_per_team"])
	return make_response(render)

@app.route("/register", methods = ["POST"])
def register_submit():
	"""Attempt to register user"""
	
	login, user = get_user()
	
	if login:
		return redirect("/error/500")
	
	username = request.form["username"]
	password = request.form["password"]
	cteamname = request.form["cteam-name"]
	jteamname = request.form["jteam-name"]
	jteamcode = request.form["jteam-code"]
	
	if len(username) < config["minimum_username_length"] or len(username) > config["maximum_username_length"] or re.match('^[\w-]+$', username) is None:
		return redirect("/register")
	if len(password) < config["minimum_password_length"] or len(password) > config["maximum_password_length"]:
		return redirect("/register")
	if len(cteamname) > config["maximum_team_name_length"]:
		return redirect("/register")
	if cteamname != "" and (jteamname != "" or jteamcode != ""):
		return redirect("/error/500") # Prevent accidents
	
	# Don't access database until it is needed
	
	
	return redirect("/")	

@app.route("/")
def index():
	"""Displays the index page"""
	
	login, user = get_user() 
	
	# Render the page
	render = render_template('frame.html', lang=lang, page='index.html', login=login, user=user)
	return make_response(render)

if __name__ == "__main__":
	"""Initializes variables and starts the server"""

	# Load Config

	config_str = open("config.json", "rb").read()
	config = json.loads(config_str)

	# Configure Security

	app.secret_key = config["secret_key"]

	# Load Language

	lang_str = open(config["language_file"], "rb").read()
	lang = json.loads(lang_str)
	lang = lang[config["language"]]
	
	# Sponsors
	
	sponsors = config["sponsors"]
	
	# Run Server
	
	app.run(host=config['host'], port=config['port'], debug=config['debug'], threaded=True)
