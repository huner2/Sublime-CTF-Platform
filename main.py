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
minpasslength = None

def rendering(page, login, user):
	return render_template('frame.html', lang=lang, page=page, login=login, user=user)

def get_user():
	return (False, None)

@app.route("/")
def index():
	"""Displays the index page"""
	
	login, user = get_user() 
	
	# Render the page
	render = rendering('index.html', login, user)
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

	# Run Server

	app.run(host=config['host'], port=config['port'], debug=config['debug'], threaded=True)
