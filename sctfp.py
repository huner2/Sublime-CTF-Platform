from flask import Flask, render_template, url_for
import configparser
import arrow

app = Flask(__name__)
config = configparser.ConfigParser()
config.read('config.ini')

# Init some global vars
start_time_un = arrow.get(config['system']['start_time']).to("utc")
end_time_un = arrow.get(config['system']['end_time']).to("utc")
start_time = start_time_un.format('MMM D, YYYY - h:mm A') + " UTC"
end_time = end_time_un.format('MMM D, YYYY - h:mm A') + " UTC"

# Use the same values for all renders
def render_p(page, **kwargs):
    return render_template(
        page,
        page=page,
        prefs=config['prefs'],
        locale=config['locale'],
        **kwargs
    )

@app.route('/')
def index():
    future = start_time if arrow.utcnow().timestamp < start_time_un.timestamp else end_time
    return render_p('index.html', start_time=start_time, end_time=end_time,
    future=future, start_time_un=start_time_un, end_time_un=end_time_un)
