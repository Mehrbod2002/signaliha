# pm2 start bot.py --interpreter /usr/local/bin/python3.8 --watch
import re,os,time,requests,json
import pymysql
from dotenv import load_dotenv
from jdatetime import datetime as jdatetime_datetime

token = 0
load_dotenv(".env")

def get_db_connection():
    DB_HOST="localhost"
    DB_PORT=3306
    DB_USER="signal"
    DB_PASSWORD="DtrPuxeHW6wWQ#g^"
    DB_NAME="arz"
    db = pymysql.connect(
        host=DB_HOST,
        user=DB_USER,
        password=DB_PASSWORD,
        database=DB_NAME
    )
    return db

def create_messages_table_if_not_exists():
    db = get_db_connection()
    cursor = db.cursor()
    cursor.execute("SHOW TABLES LIKE 'messages'")
    table_exists = cursor.fetchone()
    
    if not table_exists:
        cursor.execute("""
            CREATE TABLE messages (
                id INT AUTO_INCREMENT PRIMARY KEY,
                message_id INT,
                coin VARCHAR(50),
                base_currency VARCHAR(50),
                platform VARCHAR(50),
                leverage VARCHAR(50),
                side VARCHAR(50),
                entries VARCHAR(255),
                margin VARCHAR(255),
                sl VARCHAR(255),
                timestamp INT,
                exit BOOL,
                risk BOOL
            )
        """)
        db.commit()
    
    cursor.close()

def get_token():
    global token
    response = requests.post("https://club.caronlineofficial.com/api/Token/Login", json={
        "password": "Best1234",
        "phoneNumber": "09124343432"
    })
    if response.status_code == 200:
        data = response.json()
        if data["isSuccess"]:
            token = data["data"]

def get_last_message_id(db):
    cursor = db.cursor()
    query = "SELECT MAX(message_id) FROM messages"
    cursor.execute(query)
    last_message_id = cursor.fetchone()[0]
    cursor.close()
    return last_message_id

def process_message(message, db):
    try:
        try:
            gregorian_date = jdatetime_datetime.strptime(message["sendTime"], "%Y/%m/%d  |  %H:%M").togregorian()
            timestamp = time.mktime(gregorian_date.timetuple())
        except:
            return
        detail = message['text'].lower()
        message_id = message['id']
        if ('سر به سر' in detail or 'exit' in detail or 'خارج' in detail or "کنسل" in detail or "cancel" in detail) and message["relatedMessage"] is not None:
            cursor = db.cursor()
            query = "UPDATE messages SET `exit` = TRUE WHERE message_id = %s"
            values = (message["relatedMessage"]["id"] ,)
            cursor.execute(query, values)
            db.commit()
            cursor.close()
        elif ('risk' in detail or 'ریسک' in detail) and message["relatedMessage"] is not None:
            cursor = db.cursor()
            query = "UPDATE messages SET risk = TRUE WHERE message_id = %s"
            values = (message["relatedMessage"]["id"] ,)
            cursor.execute(query, values)
            db.commit()
            cursor.close()
        elif 'entry' in detail:
            detecting_entry, detecting_tp, detecting_sl, detect_coin = False, False, False, False
            entry_list, entry_tp, entry_sl, entry_margin = [], [], [], []
            enums = detail.split(" ")
            leverage, cross_or_isolate, buy_or_sell, base_currency, coin = '20x', 'isolated', "buy", "usdt", ""
            for index, d in enumerate(enums):
                if d == "#":
                    detect_coin = True
                    continue
                if detect_coin:
                    if d == '' or d == "#":
                        coin = enums[index + 1]
                    else:
                        coin = d
                    detect_coin = False
                    continue
                if index == 1 and d == "usdt" and coin == '':
                    coins = re.findall(r'\b\w+\b', enums[0])
                    if coins != []:
                        coin = coins[0]
                        detect_coin = False
                if index == 2:
                    if coin == '':
                        coin = enums[index - 1]
                    base_currency = d
                    if base_currency == "" or d == "/":
                        base_currency = "usdt"
                    continue
                if index == 3 and d == "usdt":
                    if coin == "":
                        coin = enums[index - 1]
                if "⛔sl" in d:
                    t = re.findall(r"[-+]?\d*\.?\d+", d)
                    if t is not None and t != []:
                        detecting_tp = False
                        entry_sl.extend(sls)
                    continue
                if d == "lvg":
                    leverage = enums[index + 1]
                    continue
                if d == "entry":
                    detecting_entry = True
                    detecting_tp, detecting_sl, detecting_margin = False, False, False
                    continue
                if d == "tp" or ('tp' in d ):
                    detecting_tp = True
                    detecting_sl, detecting_entry, detecting_margin = False, False, False
                    t = re.findall(r"[-+]?\d*\.?\d+", d)
                    if t is not None and t != []:
                        if "sl" in d:
                            sls = re.findall(r"[-+]?\d*\.?\d+",enums[index+1])
                            if sls != []:
                                entry_sl.extend(sls)
                        entry_tp.extend(t)
                        detecting_tp = False
                    continue
                if d == "sl":
                    detecting_sl = True
                    if enums[index+1] == ":":
                        sls = re.findall(r"[-+]?\d*\.?\d+",enums[index+2])
                        if sls != []:
                            entry_sl.extend(sls)
                            detecting_sl = False
                    detecting_entry, detecting_tp, detecting_margin = False, False, False
                    continue
                if d == "margin":
                    detecting_entry, detecting_tp, detecting_sl = False, False, False
                    val = enums[index+1]
                    if val == ":":
                        val = enums[index+2].replace("%","")
                    else:
                        val = val.replace(":","").replace("%","")
                    entry_margin.append(val)
                    continue
                if detecting_entry:
                    x = d
                    if "tp" in x:
                        x = x.split("tp")[0]
                        detecting_sl = True
                        detecting_tp = True
                    t = re.findall(r"[-+]?\d*\.?\d+", x)
                    if t is not None and t != []:
                        entry_list.extend(t)
                        detecting_entry = False
                if detecting_sl:
                    if d == "margin":
                        t = re.findall(r"[-+]?\d*\.?\d+", re.findall(r"[-+]?\d*\.?\d+", d))
                        entry_sl.extend(t)
                    else:
                        t = re.findall(r"[-+]?\d*\.?\d+", d)
                        if t == []:
                            entry_sl.extend(re.findall(r"[-+]?\d*\.?\d+", enums[index+1]))
                        if t is not None:
                            entry_sl.extend(t)
                if detecting_tp:
                    t = re.findall(r"[-+]?\d*\.?\d+", d)
                    if t is not None:
                        if "sl" in d:
                            sls = re.findall(r"[-+]?\d*\.?\d+",enums[index+1])
                            if sls != []:
                                entry_sl.extend(sls)
                        entry_tp.extend(t)
                if "cross" in d.lower():
                    cross_or_isolate = "cross"
                if ("sell" in d.lower() or "short" in d.lower()) and (index == 4 or index == 5):
                    buy_or_sell = "sell"
            
            if coin != "":
                cursor = db.cursor()
                query = "INSERT INTO messages (message_id, coin, base_currency, platform, leverage, side, entries, tp, margin, sl, timestamp, `exit`, risk) " \
                        "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
                entry_list_json = json.dumps(entry_list)
                entry_margin_json = json.dumps(entry_margin)
                entry_sl_json = json.dumps(entry_sl)
                entry_tp_json = json.dumps(entry_tp)
                values = (message_id, coin, base_currency, cross_or_isolate, leverage, buy_or_sell, entry_list_json, entry_tp_json, entry_margin_json, entry_sl_json, timestamp, False, False)
                cursor.execute(query, values)
                db.commit()
                cursor.close()
    except Exception as e:
        return
def main():
    get_token()
    time.sleep(0.5)
    db = get_db_connection()
    last_message_id = get_last_message_id(db)
    if last_message_id is None:
        last_message_id = 0
    while True:
        if token != 0:
            try:
                get_message_url = 'https://club.caronlineofficial.com/api/Chat/GetMessages'
                headers = {"Authorization": "Bearer " + token}
                response = requests.get(get_message_url, headers=headers)
                if response.status_code == 200:
                    messages = response.json()
                    for message in messages:
                        if last_message_id < message["id"]:
                            process_message(message, db)
                            last_message_id = message['id']
                elif response.status_code == 403:
                    get_token()
                else:
                    continue
            except pymysql.connector.Error as error:
                db.close()
                continue
            except requests.RequestException as request_error:
                continue
            except Exception as e:
                continue
        time.sleep(10)

if __name__ == "__main__":
    create_messages_table_if_not_exists()
    main()