import psycopg2 as driver
import configparser
config = configparser.ConfigParser()
config.read('config.ini')

conn = driver.connect("dbname={} user={} password={}".format(
    config['system']['db'], config['system']['dbuser'], config['system']['dbpassword'])
)
cursor = conn.cursor()
cursor.execute("CREATE TABLE IF NOT EXISTS public.users (id serial PRIMARY KEY, username varchar UNIQUE, email varchar, password varchar, team int);")
cursor.execute("CREATE TABLE IF NOT EXISTS public.teams (id serial PRIMARY KEY, name varchar UNIQUE, token varchar(9), points int);")
cursor.execute("CREATE TABLE IF NOT EXISTS public.challenges (id serial PRIMARY KEY, name varchar UNIQUE, category varchar, flag varchar, points smallint, description varchar, skip boolean);")
cursor.execute("CREATE TABLE IF NOT EXISTS public.solved (uid int REFERENCES users(id), cid smallint REFERENCES challenges(id), tid int REFERENCES teams(id));")
cursor.execute("CREATE TABLE IF NOT EXISTS public.sessions (uid int REFERENCES users(id), jwt varchar);")
conn.commit()
cursor.close()
conn.close()

print("Tables Created")
