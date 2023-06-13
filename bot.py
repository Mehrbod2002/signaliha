import re,os,time,requests
import mysql.connector
from dotenv import load_dotenv

token = 0
load_dotenv()

def get_db_connection():
    db = mysql.connector.connect(
        host=os.getenv("DB_HOST"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        database=os.getenv("DB_NAME")
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
    last_message_id = get_last_message_id(db)
    timestamp = time.mktime(time.strptime(message["sendTime"], "%Y/%m/%d | %H:%M"))
    detail = message['text'].lower()
    message_id = message['id']
    if ('سر به سر' in detail or 'exit' in detail or 'خارج' in detail) and message["relatedMessage"] is not None:
        cursor = db.cursor()
        query = "UPDATE messages SET exit = TRUE WHERE message_id = %s"
        values = (message_id,)
        cursor.execute(query, values)
        db.commit()
        cursor.close()
    elif ('risk' in detail or 'ریسک' in detail) and message["relatedMessage"] is not None:
        cursor = db.cursor()
        query = "UPDATE messages SET risk = TRUE WHERE message_id = %s"
        values = (message_id,)
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
            if index == 2:
                base_currency = d
                continue
            if d == "lvg":
                leverage = enums[index + 1]
                continue
            if d == "entry":
                detecting_entry = True
                detecting_tp, detecting_sl, detecting_margin = False, False, False
                continue
            if d == "tp":
                detecting_tp = True
                detecting_sl, detecting_entry, detecting_margin = False, False, False
                continue
            if d == "sl":
                detecting_sl = True
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
                t = re.findall(r"[-+]?\d*\.?\d+", d)
                if t is not None:
                    entry_list.extend(t)
            if detecting_sl:
                t = re.findall(r"[-+]?\d*\.?\d+", d)
                if t is not None:
                    entry_sl.extend(t)
            if detecting_tp:
                t = re.findall(r"[-+]?\d*\.?\d+", d)
                if t is not None:
                    entry_tp.extend(t)
            if "cross" in d.lower():
                cross_or_isolate = "cross"
            if ("sell" in d.lower() or "short" in d.lower()) and (index == 4 or index == 5):
                buy_or_sell = "sell"
        
        if coin != "" and message_id > last_message_id:
            cursor = db.cursor()
            query = "INSERT INTO messages (message_id, coin, base_currency, platform, leverage, side, entries, margin, sl, timestamp, exit, risk) " \
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
            values = (message_id, coin, base_currency, cross_or_isolate, leverage, buy_or_sell, entry_list, entry_margin, entry_list, timestamp, False, False)
            cursor.execute(query, values)
            db.commit()
            cursor.close()

def main():
    get_token()
    time.sleep(0.5)
    while True:
        if token != 0:
            try:
                db = get_db_connection()
                get_message_url = 'https://club.caronlineofficial.com/api/Chat/GetMessages'
                headers = {"Authorization": "Bearer " + token}
                response = requests.get(get_message_url, headers=headers)
                
                if response.status_code == 200:
                    messages = response.json()[::-1]
                    for message in messages:
                        process_message(message, db)
                elif response.status_code == 403:
                    get_token()
                else:
                    continue
            except mysql.connector.Error as error:
                db.close()
                continue
            except requests.RequestException as request_error:
                continue
        
        time.sleep(10)

if __name__ == "__main__":
    create_messages_table_if_not_exists()
    main()
