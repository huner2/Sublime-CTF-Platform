#!/usr/bin/env python

"""server.py -- the main flask server module"""

import dataset
import simplejson as json
import random
import time

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
	render = render_template('frame.html', lang=lang, page='register.html', login=login, user=user, mplength=minpasslength)
	return make_response(render)
	
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

	minpasslength = config["minimum_password_length"]
	app.secret_key = config["secret_key"]

	# Load Language

	lang_str = open(config["language_file"], "rb").read()
	lang = json.loads(lang_str)
	lang = lang[config["language"]]
	
	# Sponsors
	
	sponsors = config["sponsors"]
	

	# Run Server

	app.run(host=config['host'], port=config['port'], debug=config['debug'], threaded=True)
