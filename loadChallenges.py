import psycopg2 as driver
import json
import configparser
config = configparser.ConfigParser()
config.read('config.ini')

challenges = json.load(open('challenges.json'))
conn = driver.connect("dbname={} user={} password={}".format(
    config['system']['db'], config['system']['dbuser'], config['system']['dbpassword'])
)
cursor = conn.cursor()
for challenge in challenges:
    cursor.execute(
        """
            INSERT INTO challenges (name,category,flag,points,description,skip) VALUES (%s,%s,%s,%s,%s,%s)
            ON CONFLICT (name)
            DO
              UPDATE
                SET category=%s,flag=%s,points=%s,description=%s,skip=%s;
        """,
        (
            challenge, challenges[challenge]['category'], challenges[challenge]['flag'],
            challenges[challenge]['points'], challenges[challenge]['description'],
            challenges[challenge]['skip'], challenges[challenge]['category'],
            challenges[challenge]['flag'], challenges[challenge]['points'],
            challenges[challenge]['description'], challenges[challenge]['skip']
        )
    )
conn.commit()
cursor.close()
conn.close()
