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
team_mode = config.getboolean('system','teams')
challenges = []
ltime = 0

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

# Get challenges
def queryChal():
    global challenges
    global ltime
    ctime = arrow.utcnow().timestamp
    if (ctime < (ltime + 6000)): # Only allow refreshing of challenges every 10 minutes to reduce server queries
        return
    ltime = ctime
    conn = driver.connect("dbname={} user={} password={}".format(
        config['system']['db'], config['system']['dbuser'], config['system']['dbpassword'])
    )
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM challenges;")
    lchal = cursor.fetchall()
    wchal = []
    if (lchal[0] == None): return
    for chal in lchal:
        wchal.append(
            {"name": chal[1], "category": chal[2], "flag": chal[3],
            "points": chal[4], "description": chal[5], "skip": chal[6]}
        )
    challenges = wchal

queryChal() # Query the challenges once at the beginning
# Views

@app.route('/challenges')
def challengepage():
    queryChal()
    return render_p('challenges.html')

@app.route('/about')
def about():
    return render_p('about.html')

@app.route('/')
def index():
    start_time = start_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
    end_time = end_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
    future = start_time if arrow.utcnow().timestamp < start_time_un.timestamp else end_time
    return render_p('index.html', start_time=start_time, end_time=end_time,
    future=future, start_time_un=start_time_un, end_time_un=end_time_un)

@app.errorhandler(404)
def page_not_found(e):
    return render_p('404.html')
