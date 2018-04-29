from flask import Flask, render_template, url_for
import configparser
import arrow
import psycopg2 as driver

app = Flask(__name__)
config = configparser.ConfigParser()
config.read('config.ini')

# Init some global vars
start_time_un = arrow.get(config['system']['start_time']).to("utc")
end_time_un = arrow.get(config['system']['end_time']).to("utc")
start_time = start_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
end_time = end_time_un.format('MMM D, YYYY - h:mm A') + " UTC"

# Ensure database is up and running
conn = driver.connect("dbname={} user={} password={}".format(
    config['system']['db'], config['system']['dbuser'], config['system']['dbpassword'])
)
cursor = conn.cursor()
cursor.execute("SELECT to_regclass('public.users')")
if cursor.fetchone()[0] == None: # Table doesn't exist
    cursor.execute("CREATE TABLE users (id serial PRIMARY KEY, username varchar)")
cursor.execute("SELECT to_regclass('public.challenges')")
if cursor.fetchone()[0] == None:
    cursor.execute("CREATE TABLE challenges (id serial PRIMARY KEY, name varchar)")
# TODO: Continue setting up db
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
