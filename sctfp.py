from flask import Flask, render_template, url_for
import configparser
import arrow
import psycopg2 as driver
from werkzeug.security import safe_str_cmp

app = Flask(__name__)
config = configparser.ConfigParser()
config.read('config.ini')
app.config['SECRET_KEY'] = config['system']['secret_key']

# Init some global vars
start_time_un = arrow.get(config['system']['start_time']).to("utc")
end_time_un = arrow.get(config['system']['end_time']).to("utc")
start_time = start_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
end_time = end_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
team_mode = config.getboolean('system','teams')

# Ensure database is up and running
conn = driver.connect("dbname={} user={} password={}".format(
    config['system']['db'], config['system']['dbuser'], config['system']['dbpassword'])
)
cursor = conn.cursor()
cursor.execute("SELECT to_regclass('public.users')")
if cursor.fetchone()[0] == None: # Table doesn't exist
    cursor.execute("CREATE TABLE users (id serial PRIMARY KEY, username varchar, email varchar, password varchar, jwt varchar, solved varchar, points int, team int)")
cursor.execute("SELECT to_regclass('public.challenges')")
if cursor.fetchone()[0] == None:
    cursor.execute("CREATE TABLE challenges (id serial PRIMARY KEY, name varchar, flag varchar, points smallint, description varchar, url varchar, hint varchar, solved int)")
if team_mode: # Team mode enabled
    cursor.execute("SELECT to_regclass('public.teams')")
    if cursor.fetchone()[0] == None:
        cursor.execute("CREATE TABLE teams (id serial PRIMARY KEY, name varchar, token varchar(9), points int, solved smallint, members varchar)")
conn.commit()
cursor.close()
conn.close()

# Use the same values for all renders
def render_p(page, **kwargs):
    return render_template(
        page,
        page=page,
        passed = arrow.utcnow().timestamp > end_time_un.timestamp,
        prefs=config['prefs'],
        locale=config['locale'],
        system=config['system'],
        **kwargs
    )

# Views

@app.route('/about')
def about():
    return render_p('about.html')

@app.route('/')
def index():
    future = start_time if arrow.utcnow().timestamp < start_time_un.timestamp else end_time
    return render_p('index.html', start_time=start_time, end_time=end_time,
    future=future, start_time_un=start_time_un, end_time_un=end_time_un)

@app.errorhandler(404)
def page_not_found(e):
    return render_p('404.html')
